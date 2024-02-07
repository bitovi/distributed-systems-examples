package schemas

type WorkflowInput struct {
	Total    int       `json:"total"`
	Customer Customer  `json:"customer"`
	Products []Product `json:"products"`
}

type Product struct {
	ProductCode string `json:"productCode"`
	Quantity    int    `json:"quantity"`
	Price       int    `json:"price"`
}

type Customer struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	City      string `json:"city"`
}
