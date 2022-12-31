package expense

type Expense struct {
	ID     int      `json:"id""`
	Title  string   `json:"title"`
	Amount float64  `json:"amount"`
	Note   string   `json:"note"`
	tags   []string `json:"tags"`
}

type Err struct {
	Message string `json:"message"`
}
