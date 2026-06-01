package http

import (
	"hitalent-test-task/internal/domain/models"
	"hitalent-test-task/internal/services"
	"time"
)

type EmployeeResponse struct {
	ID           int64      `json:"id"`
	DepartmentID int64      `json:"department_id"`
	FullName     string     `json:"full_name"`
	Position     string     `json:"position"`
	HiredAt      *time.Time `json:"hired_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

type DepartmentResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	ParentID  *int64    `json:"parent_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type DepartmentTreeResponse struct {
	Department DepartmentResponse       `json:"department"`
	Employees  []EmployeeResponse       `json:"employees,omitempty"`
	Children   []DepartmentTreeResponse `json:"children,omitempty"`
}

func ToDepartmentTreeResponse(tree *services.DepartmentTree) DepartmentTreeResponse {
	resp := DepartmentTreeResponse{
		Department: DepartmentResponse{
			ID:        tree.Department.ID,
			Name:      tree.Department.Name,
			ParentID:  tree.Department.ParentID,
			CreatedAt: tree.Department.CreatedAt,
		},
	}

	for _, emp := range tree.Employees {
		resp.Employees = append(resp.Employees, EmployeeResponse{
			ID:           emp.ID,
			DepartmentID: emp.DepartmentID,
			FullName:     emp.FullName,
			Position:     emp.Position,
			HiredAt:      emp.HiredAt,
			CreatedAt:    emp.CreatedAt,
		})
	}

	for _, child := range tree.Children {
		resp.Children = append(
			resp.Children,
			ToDepartmentTreeResponse(&child),
		)
	}

	return resp
}

func ToEmployeeResponse(emp *models.Employee) EmployeeResponse {
	return EmployeeResponse{
		ID:           emp.ID,
		DepartmentID: emp.DepartmentID,
		FullName:     emp.FullName,
		Position:     emp.Position,
		HiredAt:      emp.HiredAt,
		CreatedAt:    emp.CreatedAt,
	}
}

func ToDepartmentResponse(d *models.Department) DepartmentResponse {
	return DepartmentResponse{
		ID:        d.ID,
		Name:      d.Name,
		ParentID:  d.ParentID,
		CreatedAt: d.CreatedAt,
	}
}
