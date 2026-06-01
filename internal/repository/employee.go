package repository

import (
	"context"
	"hitalent-test-task/internal/domain/models"

	"gorm.io/gorm"
)

type employeeRepository struct {
	db *gorm.DB
}

func (r *employeeRepository) Create(
	ctx context.Context,
	emp *models.Employee,
) error {

	return r.db.
		WithContext(ctx).
		Create(emp).
		Error
}

func (r *employeeRepository) GetByDepartmentID(
	ctx context.Context,
	departmentID int64,
) ([]models.Employee, error) {

	var employees []models.Employee

	err := r.db.
		WithContext(ctx).
		Where("department_id = ?", departmentID).
		Order("full_name ASC").
		Find(&employees).
		Error

	if err != nil {
		return nil, err
	}

	return employees, nil
}

func (r *employeeRepository) ReassignDepartment(
	ctx context.Context,
	fromDepartmentID int64,
	toDepartmentID int64,
) error {

	return r.db.
		WithContext(ctx).
		Model(&models.Employee{}).
		Where("department_id = ?", fromDepartmentID).
		Update("department_id", toDepartmentID).
		Error
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{
		db: db,
	}
}
