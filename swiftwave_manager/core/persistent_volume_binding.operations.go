package core

import (
	"context"
	"gorm.io/gorm"
)

// This file contains the operations for the PersistentVolumeBinding model.
// This functions will perform necessary validation before doing the actual database operation.

// Each function's argument format should be (ctx context.Context, db gorm.DB, ...)
// context used to pass some data to the function e.g. user id, auth info, etc.

func FindPersistentVolumeBindingsByApplicationId(ctx context.Context, db gorm.DB, applicationId string) ([]*PersistentVolumeBinding, error) {
	var persistentVolumeBindings []*PersistentVolumeBinding
	tx := db.Where("application_id = ?", applicationId).Find(&persistentVolumeBindings)
	return persistentVolumeBindings, tx.Error
}
