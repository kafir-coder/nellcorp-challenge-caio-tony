package requestparams

type GetTransactionsHistoryRequest struct {
	Limit    string `json:"limit"`
	Page     string `json:"page"`
	Sort     string `json:"sort"`
	Date     string `json:"date"`
	FromDate string `json:"from_date"`
	Todate   string `json:"to_date"`
}
