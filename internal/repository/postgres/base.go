package postgres

import (
	"context"
	"gorm.io/gorm"
	"wildwest/internal/errors"
	"wildwest/pkg/contextutils"
)

type BaseRepository struct {
	db *gorm.DB
}

func (r *BaseRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

func (r *BaseRepository) Get(ctx context.Context, tx *gorm.DB, tableName string, fieldName string, fieldValue interface{}, dest interface{}) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	result := db.WithContext(ctx).Table(tableName).Where(fieldName+" = ?", fieldValue).First(dest)
	if result.Error != nil {
		contextData := contextutils.ExtractContextData(ctx)
		return errors.RecordNotFoundError(contextData, tableName)
	}
	return nil
}

func (r *BaseRepository) Create(ctx context.Context, tx *gorm.DB, tableName string, data interface{}) (int, error) {
	db := r.db
	if tx != nil {
		db = tx
	}
	result := db.WithContext(ctx).Table(tableName).Create(data)
	if result.Error != nil {
		contextData := contextutils.ExtractContextData(ctx)
		return 0, errors.CreateError(contextData, tableName, result.Error)
	}
	return int(result.RowsAffected), nil
}

func (r *BaseRepository) Update(ctx context.Context, tx *gorm.DB, tableName string, fieldName string, fieldValue interface{}, data interface{}) (int, error) {
	db := r.db
	if tx != nil {
		db = tx
	}
	result := db.WithContext(ctx).Table(tableName).Where(fieldName+" = ?", fieldValue).Updates(data)
	contextData := contextutils.ExtractContextData(ctx)
	if result.Error != nil {
		return 0, errors.UpdateError(contextData, tableName, result.Error)
	}
	if result.RowsAffected == 0 {
		return 0, errors.RecordNotFoundError(contextData, tableName)
	}
	return int(result.RowsAffected), nil
}

func (r *BaseRepository) Delete(ctx context.Context, tx *gorm.DB, tableName string, fieldName string, id interface{}) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	result := db.WithContext(ctx).Table(tableName).Where(fieldName+" = ?", id).Delete(nil)
	contextData := contextutils.ExtractContextData(ctx)
	if result.Error != nil {
		return errors.DeleteError(contextData, tableName, result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.RecordNotFoundError(contextData, tableName)
	}
	return nil
}
