package tests

import (
	"gorm.io/gorm"

	usersDomain "mono_pardo/internal/domain/users"
)

type UserFixture struct {
	Users []usersDomain.User
}

func (f *UserFixture) Setup(db *gorm.DB) error {
	return db.Create(&f.Users).Error
}

func (f *UserFixture) Teardown(db *gorm.DB) error {
	return db.Unscoped().Delete(&f.Users).Error
}
