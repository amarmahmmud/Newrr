package core

import (
	"context"
	"gorm.io/gorm"
)

// This file contains the operations for the ImageRegistryCredential model.
// This functions will perform necessary validation before doing the actual database operation.

// Each function's argument format should be (ctx context.Context, db gorm.DB, ...)
// context used to pass some data to the function e.g. user id, auth info, etc.

func FindAllImageRegistryCredentials(ctx context.Context, db gorm.DB) ([]*ImageRegistryCredential, error) {
	var imageRegistryCredentials []*ImageRegistryCredential
	tx := db.Find(&imageRegistryCredentials)
	return imageRegistryCredentials, tx.Error
}

func (imageRegistryCredential *ImageRegistryCredential) FindById(ctx context.Context, db gorm.DB, id int) error {
	tx := db.First(imageRegistryCredential, id)
	return tx.Error
}

func (imageRegistryCredential *ImageRegistryCredential) Create(ctx context.Context, db gorm.DB) error {
	tx := db.Create(imageRegistryCredential)
	return tx.Error
}

func (imageRegistryCredential *ImageRegistryCredential) Update(ctx context.Context, db gorm.DB) error {
	tx := db.Save(imageRegistryCredential)
	return tx.Error
}

func (imageRegistryCredential *ImageRegistryCredential) Delete(ctx context.Context, db gorm.DB) error {
	tx := db.Delete(imageRegistryCredential)
	return tx.Error
}
