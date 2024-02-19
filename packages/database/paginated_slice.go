package database

type PaginatedSlice[M any] struct {
	Page    int  `json:"page"`
	PerPage int  `json:"per_page"`
	Items   *[]M `json:"items"`
}

func NewPaginatedSlice[M any](page int, perPage int, items *[]M) *PaginatedSlice[M] {
	return &PaginatedSlice[M]{
		Page:    page,
		PerPage: perPage,
		Items:   items,
	}
}
