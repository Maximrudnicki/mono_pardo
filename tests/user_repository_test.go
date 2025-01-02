package tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	domain "mono_pardo/internal/domain/users"
	repository "mono_pardo/internal/infrastructure/users"
	"mono_pardo/pkg/config"
)

func TestUserRepository(t *testing.T) {
	lc, err := config.LoadConfig("..")
	if err != nil {
		t.Fatalf("🚀 Could not load environment variables: %v", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		lc.DBHost, lc.DBPort, lc.DBUsername, lc.DBPassword, lc.DBTestName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal("Failed to connect to the database:", err)
	}
	db.AutoMigrate(&domain.User{})

	userRepository := repository.NewPostgresRepositoryImpl(db)

	testUser := domain.User{
		Id:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Password: "testpassword",
	}

	t.Run("Test Save User", func(t *testing.T) {
		err := userRepository.Save(testUser)
		assert.NoError(t, err, "Expected no error while saving the user")

		err = userRepository.Save(testUser)
		assert.Error(t, err, "Expected an error while saving a user with the same email")
	})

	t.Run("Test FindByEmail", func(t *testing.T) {
		foundUser, err := userRepository.FindByEmail(testUser.Email)
		assert.NoError(t, err, "Expected no error while finding the user by Email")
		assert.Equal(t, testUser.Email, foundUser.Email, "Expected the found user to have the same email as the test user")
	})

	t.Run("Test FindByEmail Fail", func(t *testing.T) {
		_, err := userRepository.FindByEmail("notvalid@email.com")
		assert.Error(t, err, "Expected error while finding the user by Email")
	})

	t.Run("Test FindAll", func(t *testing.T) {
		users := userRepository.FindAll()
		assert.NotEmpty(t, users)
	})

	t.Run("Test FindById", func(t *testing.T) {
		foundUser, err := userRepository.FindById(testUser.Id)
		assert.NoError(t, err, "Expected no error while finding the user by ID")
		assert.Equal(t, testUser, foundUser, "Expected the found user to be the same as the test user")
	})

	t.Run("Test Delete User", func(t *testing.T) {
		userRepository.Delete(testUser.Id)
		_, err := userRepository.FindByEmail(testUser.Email)
		assert.Error(t, err, "Expected error while finding the user by Email")
	})

	db.Migrator().DropTable(&domain.User{})
}
