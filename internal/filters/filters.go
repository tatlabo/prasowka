// Package filters - filter the news
package filters

import (
	"time"

	"prasowka/internal/pagination"
)

type ProductFilter struct {
	Brand   string
	Model   string
	Color   string
	Year    int
	Status  string
	Sort    int
	Code    string
	Version int
	pagination.Pagination
}

type UserFilter struct {
	UserName  string
	FirstName string
	LastName  string
	Email     string
	Role      int
	Sort      int
	pagination.Pagination
}

type Article struct {
	Title     string
	Body      string
	CreatedAt time.Time
	pagination.Pagination
}

type ServiceFilter struct {
	UserID    int
	ProductID int
	Name      string
	UserName  string
	Status    string
	Price     string
	Detail    string
	CreatedAt time.Time
	UpdatedAt int
	pagination.Pagination
}
