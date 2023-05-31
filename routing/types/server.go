package types

type TestParam struct {
	Phone string `json:"phone"`
}

type TestRes struct {
	Name    string  `json:"name"`
	Phone   string  `json:"phone"`
	Balance float64 `json:"balance"`
	Card    string  `json:"card"`
	Date    string  `json:"date"`
}
