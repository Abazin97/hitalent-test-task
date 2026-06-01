package repository

import (
	"context"
	"hitalent-test-task/internal/domain/models"
)

type DepartmentRepository interface {
	Create(ctx context.Context, dep *models.Department) error
	ExistsByName(ctx context.Context, parentID *int64, name string) (bool, error)
	GetByID(ctx context.Context, id int64) (*models.Department, error)

	Update(ctx context.Context, dep *models.Department) error

	Delete(ctx context.Context, id int64) error

	Exists(ctx context.Context, id int64) (bool, error)

	IsDescendant(ctx context.Context, parentID int64, childID int64) (bool, error)

	GetChildren(ctx context.Context, parentID int64) ([]models.Department, error)
}

type EmployeeRepository interface {
	Create(ctx context.Context, emp *models.Employee) error
	GetByDepartmentID(ctx context.Context, departmentID int64) ([]models.Employee, error)
	ReassignDepartment(ctx context.Context, fromDepartmentID int64, toDepartmentID int64) error
}
