package clients

type CreateClientRequest struct {
	FullName  string  `json:"fullName"`
	Phone     string  `json:"phone"`
	Instagram *string `json:"instagram"`
}

type ClientResponse struct {
	ID        string  `json:"id"`
	FullName  string  `json:"fullName"`
	Phone     string  `json:"phone"`
	Instagram *string `json:"instagram,omitempty"`
}

type ClientListItem struct {
	ID       string `json:"id"`
	FullName string `json:"fullName"`
	Phone    string `json:"phone"`
}

type ClientDetailsResponse struct {
	ID         string  `json:"id"`
	FullName   string  `json:"fullName"`
	Phone      string  `json:"phone"`
	Instagram  *string `json:"instagram,omitempty"`
	BirthDate  *string `json:"birthDate,omitempty"`
}

