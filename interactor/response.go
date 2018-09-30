package interactor

// SlimLoadResponder defines methods around response
type SlimLoadResponder interface {
	Present(response []Response)
}

// Response is used to draw the graph
type Response struct {
	UserID   string  `json:"user_id"`
	Weight   float64 `json:"weight"`
	Distance float64 `json:"distance"`
	Date     string  `json:"date"`
}
