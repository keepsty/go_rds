package forms

import "github.com/keepsty/go_rds/pkg/pagination"

type GetHistoryForm struct {
	PaginationQ pagination.Pagination
	Search      string `form:"search"`
}
