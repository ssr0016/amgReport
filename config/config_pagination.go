package config

import (
	"os"
	"strconv"
)

const (
	DefaultPage      = 1
	DefaultPageLimit = 20
)

type PaginationConfig struct {
	Page      int
	PageLimit int
}

func (cfg *Config) LoadPaginationConfig() {
	page, err := strconv.Atoi(os.Getenv("PAGE"))
	if err != nil || page <= 0 {
		page = DefaultPage
	}
	cfg.Pagination.Page = page

	pageLimit, err := strconv.Atoi(os.Getenv("PER_PAGE"))
	if err != nil || pageLimit <= 0 {
		pageLimit = DefaultPageLimit
	}
	cfg.Pagination.PageLimit = pageLimit
}
