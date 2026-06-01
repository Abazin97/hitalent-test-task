package services

import "errors"

var (
	ErrDepartmentAlreadyExists  = errors.New("department exists")
	ErrDepartmentNotFound       = errors.New("department not found")
	ErrParentDepartmentNotFound = errors.New("parent department not found")

	ErrDepartmentCycle = errors.New("department cycle detected")

	ErrInvalidDepartmentName = errors.New("invalid department name")

	ErrInvalidEmployeeName = errors.New("invalid employee name")

	ErrInvalidPosition = errors.New("invalid position")
	ErrSelfParent      = errors.New("self parent")
)
