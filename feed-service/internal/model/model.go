package model

type OkPaginationResult struct {
	Status   string             `json:"status"`
	Data     any                `json:"data"`
	Metadata PaginationMetadata `json:"metadata"`
}

type PaginationMetadata struct {
	CreatedAt  string `json:"createdAt"`
	TotalItems int    `json:"totalItems"`
	Sort       string `json:"sort"`
}

type OkResult struct {
	Status   string   `json:"status"`
	Data     any      `json:"data"`
	Metadata Metadata `json:"metadata"`
}

type ErrorResult struct {
	Status   string   `json:"status"`
	Message  string   `json:"message"`
	Metadata Metadata `json:"metadata"`
}

type Metadata struct {
	CreatedAt string `json:"createdAt"`
}
