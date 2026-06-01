package http

import "time"

type CreateDepartmentRequest struct {
	Name     string `json:"name"`
	ParentID *int64 `json:"parent_id"`
}

type CreateEmployeeRequest struct {
	FullName string     `json:"full_name"`
	Position string     `json:"position"`
	HiredAt  *time.Time `json:"hired_at"`
}

type UpdateDepartmentRequest struct {
	Name     *string `json:"name"`
	ParentID *int64  `json:"parent_id"`
}
