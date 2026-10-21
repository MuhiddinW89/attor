package analytics

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/MuhiddinW89/attor/internal/identity"
	"github.com/MuhiddinW89/attor/pkg/database"
	"github.com/google/uuid"
)

type PostgresRepository struct {
	db database.Querier
}

func NewPostgresRepository(
	db database.Querier,
) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

// =========================================================
// Sales filter
// =========================================================
//
// Week is NOT ISO week of the year.
//
// It is a calendar week inside the selected month:
//
// W01 = first calendar week of month
// W02 = second calendar week
// ...
// W06 = sixth calendar week if it exists.
//
// Weeks start on Monday.
//
// Example:
//
// August 2026:
//
// Mon Tue Wed Thu Fri Sat Sun
//                     01  02  -> W01
// 03  04  05  06  07  08  09  -> W02
// 10  11  12  13  14  15  16  -> W03
// 17  18  19  20  21  22  23  -> W04
// 24  25  26  27  28  29  30  -> W05
// 31                              -> W06
//
// =========================================================

func buildSalesFilter(
	filters Filters,
	session *identity.Session,
) (string, []interface{}) {

	conditions := []string{
		"sales.store_id = $1",
		"($2 = 'owner' OR sales.created_by_user_id = $3)",
	}

	args := []interface{}{session.StoreID, session.Role, session.UserID}

	add := func(
		condition string,
		value interface{},
	) {
		args = append(args, value)

		conditions = append(
			conditions,
			fmt.Sprintf(
				condition,
				len(args),
			),
		)
	}

	// =========================================================
	// Year
	// =========================================================
	//
	// Use YEAR, not ISOYEAR.
	//
	// ISOYEAR can belong to another calendar year around
	// New Year. For analytics we want normal calendar years.
	// =========================================================

	if filters.Year != nil {
		add(
			"EXTRACT(YEAR FROM sale_date) = $%d",
			*filters.Year,
		)
	}

	// =========================================================
	// Month
	// =========================================================

	if filters.Month != nil {
		add(
			"EXTRACT(MONTH FROM sale_date) = $%d",
			*filters.Month,
		)
	}

	// =========================================================
	// Week inside month
	// =========================================================
	//
	// Formula:
	//
	// ((day_of_month + first_day_offset - 2) / 7) + 1
	//
	// EXTRACT(ISODOW ...) returns:
	//
	// Monday    = 1
	// Tuesday   = 2
	// ...
	// Sunday    = 7
	//
	// This gives calendar weeks starting Monday.
	// =========================================================

	if filters.Week != nil &&
		filters.Year != nil &&
		filters.Month != nil {

		add(
			`(
				(
					EXTRACT(DAY FROM sale_date)::int
					+
					EXTRACT(
						ISODOW FROM make_date($%d, $%d, 1)
					)::int
					- 2
				) / 7
			) + 1 = $%d`,
			*filters.Year,
		)

		/*
			The condition above needs three parameters:
				year
				month
				week

			The generic add() helper only accepts one value,
			so we cannot use it for this condition.

			The actual week condition is added below instead.
		*/

		// Remove the incorrectly added condition/value.
		conditions = conditions[:len(conditions)-1]
		args = args[:len(args)-1]

		args = append(
			args,
			*filters.Year,
			*filters.Month,
			*filters.Week,
		)

		conditions = append(
			conditions,
			fmt.Sprintf(
				`(
					(
						EXTRACT(DAY FROM sale_date)::int
						+
						EXTRACT(
							ISODOW FROM make_date($%d, $%d, 1)
						)::int
						- 2
					) / 7
				) + 1 = $%d`,
				len(args)-2,
				len(args)-1,
				len(args),
			),
		)
	}

	// =========================================================
	// Day
	// =========================================================

	if filters.Day != nil {
		add(
			"EXTRACT(DAY FROM sale_date) = $%d",
			*filters.Day,
		)
	}

	return strings.Join(
		conditions,
		" AND ",
	), args
}

func (r *PostgresRepository) Get(
	ctx context.Context,
	filters Filters,
) (*Analytics, error) {
	session, ok := identity.SessionFromContext(ctx)
	if !ok {
		return nil, identity.ErrInvalidSession
	}

	where, args := buildSalesFilter(filters, session)

	result := &Analytics{}

	// =========================================================
	// Main statistics
	// =========================================================

	const mainQueryTemplate = `
		SELECT
			COALESCE(SUM(price), 0),
			COUNT(*),
			COALESCE(AVG(price), 0)
		FROM sales
		WHERE %s
	`

	mainQuery := fmt.Sprintf(
		mainQueryTemplate,
		where,
	)

	err := r.db.QueryRow(
		ctx,
		mainQuery,
		args...,
	).Scan(
		&result.TotalRevenue,
		&result.TotalSales,
		&result.AveragePurchase,
	)

	if err != nil {
		return nil, err
	}

	// =========================================================
	// Clients + repeat rate
	// =========================================================

	const clientsQueryTemplate = `
		SELECT
			COUNT(*) AS total_clients,
			COUNT(*) FILTER (
				WHERE purchase_count >= 2
			) AS repeat_customers
		FROM (
			SELECT
				client_id,
				COUNT(*) AS purchase_count
			FROM sales
			WHERE %s
			GROUP BY client_id
		) grouped_clients
	`

	clientsQuery := fmt.Sprintf(
		clientsQueryTemplate,
		where,
	)

	var totalClients int
	var repeatCustomers int

	err = r.db.QueryRow(
		ctx,
		clientsQuery,
		args...,
	).Scan(
		&totalClients,
		&repeatCustomers,
	)

	if err != nil {
		return nil, err
	}

	result.TotalClients = totalClients
	result.RepeatCustomers = repeatCustomers

	if totalClients > 0 {
		result.RepeatRate =
			float64(repeatCustomers) /
				float64(totalClients) *
				100
	} else {
		result.RepeatRate = 0
	}

	// =========================================================
	// Revenue overview
	// =========================================================
	//
	// Year only:
	//     group by month
	//
	// Month / Week / Day:
	//     group by day
	// =========================================================

	groupExpression :=
		"date_trunc('month', sale_date)"

	if filters.Month != nil ||
		filters.Week != nil ||
		filters.Day != nil {

		groupExpression =
			"date_trunc('day', sale_date)"
	}

	revenueQuery := fmt.Sprintf(
		`
		SELECT
			%s AS bucket,
			COALESCE(SUM(price), 0),
			COUNT(*)
		FROM sales
		WHERE %s
		GROUP BY bucket
		ORDER BY bucket
		`,
		groupExpression,
		where,
	)

	rows, err := r.db.Query(
		ctx,
		revenueQuery,
		args...,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			bucket  time.Time
			revenue float64
			sales   int
		)

		if err := rows.Scan(
			&bucket,
			&revenue,
			&sales,
		); err != nil {
			rows.Close()
			return nil, err
		}

		result.RevenueOverview =
			append(
				result.RevenueOverview,
				RevenuePoint{
					Label:   bucket.Format("Jan 2"),
					Revenue: revenue,
					Sales:   sales,
				},
			)
	}

	rows.Close()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// =========================================================
	// Top perfumes
	// =========================================================

	perfumeQuery := fmt.Sprintf(
		`
		SELECT
			perfume_name,
			COUNT(*) AS sales,
			COALESCE(SUM(price), 0) AS revenue
		FROM sales
		WHERE %s
		GROUP BY perfume_name
		ORDER BY sales DESC, revenue DESC
		LIMIT 5
		`,
		where,
	)

	rows, err = r.db.Query(
		ctx,
		perfumeQuery,
		args...,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var item PerfumeStat

		if err := rows.Scan(
			&item.PerfumeName,
			&item.Sales,
			&item.Revenue,
		); err != nil {
			rows.Close()
			return nil, err
		}

		result.TopPerfumes =
			append(
				result.TopPerfumes,
				item,
			)
	}

	rows.Close()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// =========================================================
	// Top customers
	// =========================================================

	customerWhere :=
		strings.ReplaceAll(
			where,
			"sales.",
			"s.",
		)

	customerQuery := fmt.Sprintf(
		`
		SELECT
			c.id,
			c.full_name,
			COUNT(*) AS sales,
			COALESCE(SUM(s.price), 0) AS revenue
		FROM sales s
		JOIN clients c
			ON c.id = s.client_id AND c.store_id = s.store_id
		WHERE %s
		GROUP BY c.id, c.full_name
		ORDER BY revenue DESC, sales DESC
		LIMIT 5
		`,
		customerWhere,
	)

	rows, err = r.db.Query(
		ctx,
		customerQuery,
		args...,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			item     CustomerStat
			clientID uuid.UUID
		)

		if err := rows.Scan(
			&clientID,
			&item.FullName,
			&item.Sales,
			&item.Revenue,
		); err != nil {
			rows.Close()
			return nil, err
		}

		item.ClientID = clientID.String()

		result.TopCustomers =
			append(
				result.TopCustomers,
				item,
			)
	}

	rows.Close()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// =========================================================
	// Key insights
	// =========================================================

	insightQuery := fmt.Sprintf(
		`
		SELECT
			COALESCE((
				SELECT TO_CHAR(
					sale_date,
					'Day'
				)
				FROM sales
				WHERE %s
				GROUP BY TO_CHAR(
					sale_date,
					'Day'
				)
				ORDER BY COUNT(*) DESC
				LIMIT 1
			), ''),

			COALESCE((
				SELECT
					EXTRACT(
						HOUR FROM sale_date
					)::int
				FROM sales
				WHERE %s
				GROUP BY EXTRACT(
					HOUR FROM sale_date
				)
				ORDER BY COUNT(*) DESC
				LIMIT 1
			), -1)
		`,
		where,
		where,
	)

	var (
		bestDay  string
		peakHour int
	)

	err = r.db.QueryRow(
		ctx,
		insightQuery,
		args...,
	).Scan(
		&bestDay,
		&peakHour,
	)

	if err != nil {
		return nil, err
	}

	if bestDay != "" {
		result.KeyInsights =
			append(
				result.KeyInsights,
				fmt.Sprintf(
					"Best selling day: %s",
					strings.TrimSpace(bestDay),
				),
			)
	}

	if peakHour >= 0 {
		result.KeyInsights =
			append(
				result.KeyInsights,
				fmt.Sprintf(
					"Peak sales hour: %02d:00",
					peakHour,
				),
			)
	}

	return result, nil
}
