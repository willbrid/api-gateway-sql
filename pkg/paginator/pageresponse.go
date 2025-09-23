package paginator

type PageResponse struct {
	Data       []any `json:"data"`
	PageNum    int   `json:"page_num"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int64 `json:"total_pages"`
}

func NewPageResponse(data []any, total int64, pageReq PageRequest) *PageResponse {
	totalPages := (total + int64(pageReq.PageSize) - 1) / int64(pageReq.PageSize)

	return &PageResponse{data, pageReq.PageNum, pageReq.PageSize, total, totalPages}
}
