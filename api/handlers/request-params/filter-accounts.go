package requestparams

import "strconv"

type ListAccountsRequest struct {
	Limit      string `json:"limit"`
	Page       string `json:"page"`
	OwnerId    string `json:"owner_id"`
	CreateadAt string `json:"createad_at"`
}

func (r *ListAccountsRequest) Validate() bool {
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
