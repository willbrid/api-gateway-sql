package paginator

const defaultPageSize int = 10

type PageRequest struct {
	PageNum  int
	PageSize int
}

func NewPageRequest(pageNum, pageSize int) *PageRequest {
	return &PageRequest{pageNum, pageSize}
}

func (p PageRequest) Offset() int {
	if p.PageNum <= 0 {
		return 0
	}

	return (p.PageNum - 1) * p.PageSize
}

func (p PageRequest) Limit() int {
	if p.PageSize <= 0 {
		return defaultPageSize
	}

	return p.PageSize
}
