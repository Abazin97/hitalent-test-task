package mocks

import (
	"context"
	"hitalent-test-task/internal/domain/models"

	"github.com/stretchr/testify/mock"
)

type DepartmentRepository struct {
	mock.Mock
}

func (m *DepartmentRepository) Exists(ctx context.Context, id int64) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *DepartmentRepository) ExistsByName(ctx context.Context, parentID *int64, name string) (bool, error) {
	args := m.Called(ctx, parentID, name)
	return args.Bool(0), args.Error(1)
}

func (m *DepartmentRepository) Create(ctx context.Context, d *models.Department) error {
	return m.Called(ctx, d).Error(0)
}

func (m *DepartmentRepository) GetByID(ctx context.Context, id int64) (*models.Department, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Department), args.Error(1)
}

func (m *DepartmentRepository) GetChildren(ctx context.Context, id int64) ([]models.Department, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]models.Department), args.Error(1)
}

func (m *DepartmentRepository) IsDescendant(ctx context.Context, id, parentID int64) (bool, error) {
	args := m.Called(ctx, id, parentID)
	return args.Bool(0), args.Error(1)
}

func (m *DepartmentRepository) Update(ctx context.Context, d *models.Department) error {
	return m.Called(ctx, d).Error(0)
}

func (m *DepartmentRepository) Delete(ctx context.Context, id int64) error {
	return m.Called(ctx, id).Error(0)
}

type EmployeeRepository struct {
	mock.Mock
}

func (m *EmployeeRepository) Create(ctx context.Context, e *models.Employee) error {
	return m.Called(ctx, e).Error(0)
}

func (m *EmployeeRepository) GetByDepartmentID(ctx context.Context, id int64) ([]models.Employee, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]models.Employee), args.Error(1)
}

func (m *EmployeeRepository) ReassignDepartment(ctx context.Context, from, to int64) error {
	return m.Called(ctx, from, to).Error(0)
}
