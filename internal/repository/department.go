package repository

import (
	"context"
	"hitalent-test-task/internal/domain/models"

	"gorm.io/gorm"
)

type departmentRepository struct {
	db *gorm.DB
}

func (r *departmentRepository) Create(ctx context.Context, dep *models.Department) error {
	return r.db.WithContext(ctx).Create(dep).Error
}

func (r *departmentRepository) ExistsByName(ctx context.Context, parentID *int64, name string) (bool, error) {

	var count int64

	query := r.db.WithContext(ctx).
		Model(&models.Department{}).
		Where("LOWER(TRIM(name)) = LOWER(TRIM(?))", name)

	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *departmentRepository) GetByID(
	ctx context.Context,
	id int64,
) (*models.Department, error) {

	var dep models.Department

	err := r.db.WithContext(ctx).First(&dep, id).Error

	if err != nil {
		return nil, err
	}

	return &dep, nil
}

func (r *departmentRepository) Exists(ctx context.Context, id int64) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&models.Department{}).Where("id = ?", id).Count(&count).Error

	return count > 0, err
}

func (r *departmentRepository) Update(ctx context.Context, dep *models.Department) error {
	return r.db.WithContext(ctx).Save(dep).Error
}

func (r *departmentRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.Department{}, id).Error
}

func (r *departmentRepository) IsDescendant(ctx context.Context, parentID int64, childID int64,
) (bool, error) {

	var exists bool

	query := `
	WITH RECURSIVE tree AS (
		SELECT id, parent_id
		FROM departments
		WHERE id = ?

		UNION ALL

		SELECT d.id, d.parent_id
		FROM departments d
		INNER JOIN tree t
			ON d.parent_id = t.id
	)
	SELECT EXISTS(
		SELECT 1
		FROM tree
		WHERE id = ?
	)
	`

	err := r.db.
		WithContext(ctx).
		Raw(query, parentID, childID).
		Scan(&exists).
		Error

	return exists, err
}

func (r *departmentRepository) GetChildren(
	ctx context.Context,
	parentID int64,
) ([]models.Department, error) {

	var deps []models.Department

	err := r.db.
		WithContext(ctx).
		Where("parent_id = ?", parentID).
		Order("name").
		Find(&deps).
		Error

	return deps, err
}

func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &departmentRepository{
		db: db,
	}
}
