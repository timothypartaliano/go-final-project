package repository

import "gorm.io/gorm"

type ProductRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepositoryImpl {
	return &ProductRepositoryImpl{db: db}
}