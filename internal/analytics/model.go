package analytics

type Filters struct {
	Year  *int
	Month *int
	Week  *int
	Day   *int
}

type Analytics struct {
	TotalRevenue    float64          `json:"total_revenue"`
	TotalSales      int              `json:"total_sales"`
	AveragePurchase float64          `json:"average_purchase"`
	TotalClients    int              `json:"total_clients"`
	RepeatRate      float64          `json:"repeat_rate"`
	RepeatCustomers int              `json:"repeat_customers"`

	RevenueOverview []RevenuePoint `json:"revenue_overview"`
	TopPerfumes     []PerfumeStat  `json:"top_perfumes"`
	TopCustomers    []CustomerStat `json:"top_customers"`
	KeyInsights     []string       `json:"key_insights"`
}

type RevenuePoint struct {
	Label   string  `json:"label"`
	Revenue float64 `json:"revenue"`
	Sales   int     `json:"sales"`
}

type PerfumeStat struct {
	PerfumeName string  `json:"perfume_name"`
	Sales       int     `json:"sales"`
	Revenue     float64 `json:"revenue"`
}

type CustomerStat struct {
	ClientID    string  `json:"client_id"`
	FullName    string  `json:"full_name"`
	Sales       int     `json:"sales"`
	Revenue     float64 `json:"revenue"`
}