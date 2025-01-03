package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	domain "mono_pardo/internal/domain/users"
	repository "mono_pardo/internal/infrastructure/users"
)

func TestUserRepository(t *testing.T) {
	env := NewTestEnv(t)
	defer env.Cleanup(t)

	env.RunMigrations(t)

	userRepository := repository.NewPostgresRepositoryImpl(env.DB.DB)

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
}
