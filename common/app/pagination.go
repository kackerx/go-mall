package app

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/kackerx/go-mall/config"
)

type Pagination struct {
	Page     int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	Total    int `json:"total,omitempty"`
}

func NewPagination(c *gin.Context) *Pagination {
	page, _ := strconv.Atoi(c.Query("page"))
	if page <= 0 {
		page = 1
	}

	cnf := config.Conf.App.Pagination
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	if pageSize > cnf.MaxSize {
		pageSize = cnf.MaxSize
	}

	if pageSize <= 0 {
		pageSize = cnf.DefaultSize
	}

	return &Pagination{page, pageSize, 0}
}

func (p *Pagination) SetTotal(total int) *Pagination {
	p.Total = total
	return p
}

func (p *Pagination) GetPage() int {
	return p.Page
}

func (p *Pagination) GetPageSize() int {
	return p.PageSize
}

func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}
