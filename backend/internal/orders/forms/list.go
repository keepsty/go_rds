package forms

import (
	"github.com/keepsty/go_rds/pkg/pagination"
)

type GetOrderListForm struct {
	PaginationQ  pagination.Pagination
	OnlyMyOrders bool   `form:"only_my_orders"`
	Search       string `form:"search"`
	Progress     string `form:"progress" json:"progress"`
	Environment  int    `form:"environment" json:"environment" `
}
