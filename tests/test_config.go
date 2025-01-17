package tests

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	usersDomain "mono_pardo/internal/domain/users"
	wordsDomain "mono_pardo/internal/domain/words"
	"mono_pardo/pkg/config"
)

var uniqueDBName string

// TestDB represents test database configuration and connection
type TestDB struct {
	DB *gorm.DB
}

// TestEnv holds all test environment components
type TestEnv struct {
	DB     *TestDB
	Router *gin.Engine
}

// NewTestEnv creates a new test environment
func NewTestEnv(t *testing.T) (*TestEnv, config.Config) {
	t.Helper()

	config := loadTestConfig(t)
	db := setupTestDB(t, config)
	router := gin.New()

	return &TestEnv{
		DB:     db,
		Router: router,
	}, *config
}

// loadTestConfig loads test configuration from environment
func loadTestConfig(t *testing.T) *config.Config {
	t.Helper()

	conf, err := config.LoadConfig("../../.")
	if err != nil {
		t.Fatalf("🚀 Could not load environment variables: %v", err)
	}

	return &conf
}

// setupTestDB creates and configures test database
func setupTestDB(t *testing.T, config *config.Config) *TestDB {
	t.Helper()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for attempts := 0; attempts < 5; attempts++ {
		uniqueDBName = fmt.Sprintf(
			"%s_%d_%d", config.DBTestName, time.Now().UnixNano(), r.Intn(100000))

		// Connect to postgres for creating test db
		adminDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
			config.DBHost, config.DBPort, config.DBUsername, config.DBPassword)

		adminDB, err := gorm.Open(postgres.Open(adminDSN), &gorm.Config{})
		if err != nil {
			t.Fatalf("Failed to connect to postgres: %v", err)
		}

		// Create test db
		sqlDB, err := adminDB.DB()
		if err != nil {
			t.Fatalf("Failed to get underlying *sql.DB: %v", err)
		}

		_, err = sqlDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", config.DBTestName))
		if err != nil {
			t.Fatalf("Failed to drop test database: %v", err)
		}

		_, err = sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s", uniqueDBName))
		if err != nil {
			if strings.Contains(err.Error(), "pg_database_datname_index") {
				t.Logf("Database name conflict, retrying with new name: %v", err)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			t.Fatalf("Failed to create test database: %v", err)
		}

		// Connect to Test Database
		testDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config.DBHost, config.DBPort, config.DBUsername, config.DBPassword, uniqueDBName)

		testDB, err := gorm.Open(postgres.Open(testDSN), &gorm.Config{})
		if err != nil {
			t.Fatalf("Failed to connect to test database: %v", err)
		}

		config.DBTestName = uniqueDBName

		return &TestDB{DB: testDB}
	}

	t.Fatalf("Failed to create a unique test database after multiple attempts")
	return nil
}

// Fixture interface for creating fixtures
type Fixture interface {
	Setup(db *gorm.DB) error
	Teardown(db *gorm.DB) error
}

// WithFixture prepare fixture for test
func (env *TestEnv) WithFixture(t *testing.T, fixture Fixture) func() {
	t.Helper()

	if err := fixture.Setup(env.DB.DB); err != nil {
		t.Fatalf("Failed to setup fixture: %v", err)
	}

	return func() {
		if err := fixture.Teardown(env.DB.DB); err != nil {
			t.Errorf("Failed to teardown fixture: %v", err)
		}
	}
}

// RunMigrations run migrations from indicated directory
func (env *TestEnv) RunMigrations(t *testing.T) {
	t.Helper()

	// There should be logic for running migrations
	// We can use golang-migrate/migrate or any other library
	models := []interface{}{
		// Models
		&usersDomain.User{},
		&wordsDomain.Word{},
	}

	if err := env.DB.DB.AutoMigrate(models...); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
}

// Cleanup (clean test enviroment)
func (env *TestEnv) Cleanup(t *testing.T) {
	t.Helper()

	config := loadTestConfig(t)

	// Get *sql.DB
	sqlDB, err := env.DB.DB.DB()
	if err != nil {
		t.Errorf("Failed to get underlying *sql.DB: %v", err)
		return
	}

	// Close connection
	if err := sqlDB.Close(); err != nil {
		t.Errorf("Failed to close database connection: %v", err)
	}

	// Remove test db
	adminDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUsername, config.DBPassword)

	adminDB, err := gorm.Open(postgres.Open(adminDSN), &gorm.Config{})
	if err != nil {
		t.Errorf("Failed to connect to postgres for cleanup: %v", err)
		return
	}

	sqlDB, err = adminDB.DB()
	if err != nil {
		t.Errorf("Failed to get underlying *sql.DB for cleanup: %v", err)
		return
	}

	_, err = sqlDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", uniqueDBName))
	if err != nil {
		t.Errorf("Failed to drop test database: %v", err)
	}
}
