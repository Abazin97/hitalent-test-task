package services

import (
	"context"
	"errors"
	"hitalent-test-task/internal/domain/models"
	"hitalent-test-task/internal/repository"
	"strings"
	"time"

	"gorm.io/gorm"
)

type DepartmentService struct {
	departments repository.DepartmentRepository
	employees   repository.EmployeeRepository
}

type DepartmentTree struct {
	Department *models.Department `json:"department"`

	Employees []models.Employee `json:"employees,omitempty"`

	Children []DepartmentTree `json:"children"`
}

func (s *DepartmentService) CreateDepartment(
	ctx context.Context,
	name string,
	parentID *int64,
) (*models.Department, error) {

	name = strings.TrimSpace(name)

	if len(name) == 0 || len(name) > 200 {
		return nil, ErrInvalidDepartmentName
	}

	if parentID != nil {
		exists, err := s.departments.ExistsByName(ctx, parentID, name)
		if err != nil {
			return nil, err
		}

		if exists {
			return nil, ErrDepartmentAlreadyExists
		}

		parentExists, err := s.departments.Exists(
			ctx,
			*parentID,
		)
		if err != nil {
			return nil, err
		}

		if !parentExists {
			return nil, ErrDepartmentNotFound
		}
	}

	department := &models.Department{
		Name:     name,
		ParentID: parentID,
	}

	if err := s.departments.Create(
		ctx,
		department,
	); err != nil {
		return nil, err
	}

	return department, nil
}

func (s *DepartmentService) CreateEmployee(
	ctx context.Context,
	departmentID int64,
	fullName string,
	position string,
	hiredAt *time.Time,
) (*models.Employee, error) {

	exists, err := s.departments.Exists(
		ctx,
		departmentID,
	)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrDepartmentNotFound
	}

	fullName = strings.TrimSpace(fullName)
	position = strings.TrimSpace(position)

	if len(fullName) == 0 || len(fullName) > 200 {
		return nil, ErrInvalidEmployeeName
	}

	if len(position) == 0 || len(position) > 200 {
		return nil, ErrInvalidPosition
	}

	employee := &models.Employee{
		DepartmentID: departmentID,
		FullName:     fullName,
		Position:     position,
		HiredAt:      hiredAt,
	}

	if err := s.employees.Create(
		ctx,
		employee,
	); err != nil {

		return nil, err
	}

	return employee, nil
}

func (s *DepartmentService) GetDepartmentTree(
	ctx context.Context,
	id int64,
	depth int,
	includeEmployees bool,
) (*DepartmentTree, error) {

	if depth <= 0 {
		depth = 1
	}

	if depth > 5 {
		depth = 5
	}

	return s.buildTree(ctx, id, depth, includeEmployees)
}

func (s *DepartmentService) buildTree(ctx context.Context, id int64, depth int, includeEmployees bool) (*DepartmentTree, error) {

	department, err := s.departments.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDepartmentNotFound
		}
		return nil, err
	}

	node := &DepartmentTree{
		Department: department,
	}

	if includeEmployees {
		employees, err := s.employees.GetByDepartmentID(ctx, id)
		if err != nil {
			return nil, err
		}

		node.Employees = employees
	}

	if depth == 0 {
		return node, nil
	}

	children, err := s.departments.
		GetChildren(ctx, id)

	if err != nil {
		return nil, err
	}

	for _, child := range children {

		subtree, err := s.buildTree(
			ctx,
			child.ID,
			depth-1,
			includeEmployees,
		)

		if err != nil {
			return nil, err
		}

		node.Children = append(
			node.Children,
			*subtree,
		)
	}

	return node, nil
}

func (s *DepartmentService) UpdateDepartment(
	ctx context.Context,
	id int64,
	name *string,
	parentID *int64,
) (*models.Department, error) {

	department, err := s.departments.GetByID(ctx, id)
	if err != nil {
		return nil, ErrDepartmentNotFound
	}

	if name != nil {

		n := strings.TrimSpace(*name)

		if len(n) == 0 || len(n) > 200 {
			return nil, ErrInvalidDepartmentName
		}

		department.Name = n
	}

	if parentID != nil {

		if *parentID == id {
			return nil, ErrDepartmentCycle
		}

		isDescendant, err := s.departments.IsDescendant(ctx, id, *parentID)

		if err != nil {
			return nil, err
		}

		if isDescendant {
			return nil, ErrDepartmentCycle
		}

		checkName := department.Name

		exists, err := s.departments.ExistsByName(
			ctx,
			parentID,
			checkName,
		)
		if err != nil {
			return nil, err
		}

		if exists {
			return nil, ErrDepartmentAlreadyExists
		}

		department.ParentID = parentID
	}

	if err := s.departments.Update(
		ctx,
		department,
	); err != nil {

		return nil, err
	}

	return department, nil
}

func (s *DepartmentService) DeleteCascade(
	ctx context.Context,
	id int64,
) error {

	return s.departments.Delete(ctx, id)
}

func (s *DepartmentService) DeleteReassign(
	ctx context.Context,
	departmentID int64,
	targetDepartmentID int64,
) error {

	exists, err := s.departments.Exists(
		ctx,
		targetDepartmentID,
	)
	if err != nil {
		return err
	}

	if !exists {
		return ErrDepartmentNotFound
	}

	if departmentID == targetDepartmentID {
		return ErrDepartmentCycle
	}

	if err := s.employees.ReassignDepartment(ctx, departmentID, targetDepartmentID); err != nil {
		return err
	}

	return s.departments.Delete(
		ctx,
		departmentID,
	)
}

func NewDepartmentService(
	departments repository.DepartmentRepository,
	employees repository.EmployeeRepository,
) *DepartmentService {

	return &DepartmentService{
		departments: departments,
		employees:   employees,
	}
}
