// Package pagination to paginate news
package pagination

type PaginationErr struct {
	ErrPage     error
	ErrPageSize error
}

type Pagination struct {
	Page       int
	PageSize   int
	Toal       int
	TotalPages int
	HasPrev    bool
	HasNext    bool
	PrevPage   int
	NextPage   int
	PaginationErr
}

func NewPagination(page, pageSize, total int) Pagination {
	totalPages := (total + pageSize - 1) / pageSize

	if page < 1 {
		page = 1
	}
	if page > totalPages {
		page = totalPages
	}
	return Pagination{
		Page:       page,
		PageSize:   pageSize,
		Toal:       total,
		TotalPages: totalPages,
		HasPrev:    page > 1,
		HasNext:    page < totalPages,
		PrevPage:   page - 1,
		NextPage:   page + 1,
	}
}
