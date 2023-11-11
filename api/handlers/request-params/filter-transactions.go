package requestparams

import "strconv"

type GetTransactionsHistoryRequest struct {
	Limit    string `json:"limit"`
	Page     string `json:"page"`
	Sort     string `json:"sort"`
	Date     string `json:"date"`
	FromDate string `json:"from_date"`
	Todate   string `json:"to_date"`
}

func (r *GetTransactionsHistoryRequest) Validate() bool {
	if r.Limit == "" {
		r.Limit = "10"

	}

	if r.Page == "" {
		r.Page = "1"
	}
	page, _ := strconv.ParseInt(r.Page, 10, 64)
	if page <= 0 {
		r.Page = "1"
	}

	return true
}
