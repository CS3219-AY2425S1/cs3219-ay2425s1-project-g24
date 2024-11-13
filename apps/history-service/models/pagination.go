package models

type HistoriesResponse struct {
	TotalCount  int                 `json:"totalCount"`
	TotalPages  int                 `json:"totalPages"`
	CurrentPage int                 `json:"currentPage"`
	Limit       int                 `json:"limit"`
	HasNextPage bool                `json:"hasNextPage"`
	Questions   []SubmissionHistory `json:"histories"`
}

func PaginateResponse(limit, offset int, histories []SubmissionHistory) *HistoriesResponse {
	start := offset
	end := offset + limit

	var paginatedHistory []SubmissionHistory
	if start <= len(histories) {
		if end >= len(histories) {
			end = len(histories)
		}
	} else {
		start = 0
		offset = 0
		end = limit
	}
	paginatedHistory = histories[start:end]

	// Calculate pagination info
	totalCount := len(histories)
	totalPages := (totalCount + limit - 1) / limit
	currentPage := (offset / limit) + 1
	if len(paginatedHistory) == 0 {
		currentPage = 0
	}
	hasNextPage := totalPages > currentPage

	// Construct response
	return &HistoriesResponse{
		TotalCount:  totalCount,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
		Limit:       limit,
		HasNextPage: hasNextPage,
		Questions:   paginatedHistory,
	}
}
