package brewery

type MonthlyReport struct {
	Year int     `json:"year"`
	Type string  `json:"membership"`
	Jan  float64 `json:"jan"`
	Feb  float64 `json:"feb"`
	Mar  float64 `json:"mar"`
	Apr  float64 `json:"apr"`
	May  float64 `json:"may"`
	Jun  float64 `json:"jun"`
	Jul  float64 `json:"jul"`
	Aug  float64 `json:"aug"`
	Sep  float64 `json:"sep"`
	Oct  float64 `json:"oct"`
	Nov  float64 `json:"nov"`
	Dec  float64 `json:"dec"`
}

type ReportsStorage interface {
	Monthlies() ([]*MonthlyReport, error)
}
