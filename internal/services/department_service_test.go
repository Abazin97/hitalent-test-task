package services

import (
	"context"
	"hitalent-test-task/test/mocks"
	"testing"

	"hitalent-test-task/internal/domain/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateDepartment_Success(t *testing.T) {
	deptRepo := new(mocks.DepartmentRepository)
	empRepo := new(mocks.EmployeeRepository)

	svc := NewDepartmentService(deptRepo, empRepo)

	parentID := int64(1)

	deptRepo.On("ExistsByName", context.Background(), &parentID, "backend").Return(false, nil)
	deptRepo.On("Exists", context.Background(), parentID).Return(true, nil)
	deptRepo.On("Create", context.Background(), mock.AnythingOfType("*models.Department")).Return(nil)

	dept, err := svc.CreateDepartment(context.Background(), "backend", &parentID)

	assert.NoError(t, err)
	assert.Equal(t, "backend", dept.Name)
}

func TestCreateDepartment_AlreadyExists(t *testing.T) {
	deptRepo := new(mocks.DepartmentRepository)
	empRepo := new(mocks.EmployeeRepository)

	svc := NewDepartmentService(deptRepo, empRepo)

	parentID := int64(1)

	deptRepo.On("ExistsByName", context.Background(), &parentID, "backend").Return(true, nil)

	_, err := svc.CreateDepartment(context.Background(), "backend", &parentID)

	assert.Equal(t, ErrDepartmentAlreadyExists, err)
}

func TestCreateEmployee_DepartmentNotFound(t *testing.T) {
	deptRepo := new(mocks.DepartmentRepository)
	empRepo := new(mocks.EmployeeRepository)

	svc := NewDepartmentService(deptRepo, empRepo)

	deptRepo.On("Exists", context.Background(), int64(1)).Return(false, nil)

	_, err := svc.CreateEmployee(context.Background(), 1, "John", "Dev", nil)

	assert.Equal(t, ErrDepartmentNotFound, err)
}

func TestCreateEmployee_Success(t *testing.T) {
	deptRepo := new(mocks.DepartmentRepository)
	empRepo := new(mocks.EmployeeRepository)

	svc := NewDepartmentService(deptRepo, empRepo)

	deptRepo.On("Exists", context.Background(), int64(1)).Return(true, nil)
	empRepo.On("Create", context.Background(), mock.AnythingOfType("*models.Employee")).Return(nil)

	e, err := svc.CreateEmployee(context.Background(), 1, "John", "Dev", nil)

	assert.NoError(t, err)
	assert.Equal(t, "John", e.FullName)
}

func TestUpdateDepartment_CycleError(t *testing.T) {
	deptRepo := new(mocks.DepartmentRepository)
	empRepo := new(mocks.EmployeeRepository)

	svc := NewDepartmentService(deptRepo, empRepo)

	parentID := int64(1)

	deptRepo.On("GetByID", context.Background(), int64(1)).
		Return(&models.Department{ID: 1, Name: "root"}, nil)

	_, err := svc.UpdateDepartment(context.Background(), 1, nil, &parentID)

	assert.Equal(t, ErrDepartmentCycle, err)
}

func TestDeleteReassign_Success(t *testing.T) {
	deptRepo := new(mocks.DepartmentRepository)
	empRepo := new(mocks.EmployeeRepository)

	svc := NewDepartmentService(deptRepo, empRepo)

	deptRepo.On("Exists", context.Background(), int64(2)).Return(true, nil)
	empRepo.On("ReassignDepartment", context.Background(), int64(1), int64(2)).Return(nil)
	deptRepo.On("Delete", context.Background(), int64(1)).Return(nil)

	err := svc.DeleteReassign(context.Background(), 1, 2)

	assert.NoError(t, err)
}
