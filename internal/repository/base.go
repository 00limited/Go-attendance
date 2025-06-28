package repository

import (
	gormrepository "github.com/yourname/payslip-system/internal/helper/gorm-repository"
	"gorm.io/gorm"
)

type baseRepository struct {
	gormrepository.TransactionRepository
	db *gorm.DB
}

type BaseRepositoryInterface interface {
	gormrepository.TransactionRepository
	GetDB() *gorm.DB
}

// NewBaseRepository creates a new instance of BaseRepository
func NewBaseRepository(db *gorm.DB) BaseRepositoryInterface {
	return &baseRepository{
		TransactionRepository: gormrepository.NewGormRepository(db),
		db:                    db,
	}
}

// GetDB returns the underlying database instance
func (br *baseRepository) GetDB() *gorm.DB {
	return br.db
}
