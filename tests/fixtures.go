package tests

import (
	"gorm.io/gorm"

	usersDomain "mono_pardo/internal/domain/users"
	wordsDomain "mono_pardo/internal/domain/words"
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

type WordFixture struct {
	Words []wordsDomain.Word
}

func (f *WordFixture) Setup(db *gorm.DB) error {
	return db.Create(&f.Words).Error
}

func (f *WordFixture) Teardown(db *gorm.DB) error {
	return db.Unscoped().Delete(&f.Words).Error
}
