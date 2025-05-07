package postgres_repository

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	commontest "github.com/customeros/customeros/packages/server/customer-os-common-module/test"
	"github.com/testcontainers/testcontainers-go"
	"gorm.io/gorm"
)

var (
	postgresContainer testcontainers.Container
	gormDB            *gorm.DB
	repositories      *Repositories
)

func TestMain(m *testing.M) {
	postgresContainer, gormDB, _ = commontest.InitTestDB()
	defer func(postgresContainer testcontainers.Container, ctx context.Context) {
		commontest.TerminatePostgres(postgresContainer, ctx)
	}(postgresContainer, context.Background())

	postgresDB := config.PostgresDB{
		GormDB: gormDB,
	}

	repositories = InitRepositories(&postgresDB)

	createPostgresTables(gormDB)

	os.Exit(m.Run())
}

func createPostgresTables(db *gorm.DB) {
	db.Exec("create schema if not exists derived")

	postgresDB := config.PostgresDB{
		GormDB: db,
	}

	err := InitRepositories(&postgresDB).MigrateOpenlineDB()
	if err != nil {
		panic(err)
	}
}

func tearDownTestCase() func(tb testing.TB) {
	return func(tb testing.TB) {
		tb.Logf("Teardown test %v, cleaning postgres DB", tb.Name())
		// Query all table names in the public schema
		var tableNames []string
		rows, err := gormDB.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Rows()
		if err != nil {
			tb.Fatalf("failed to fetch table names: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var tableName string
			if err := rows.Scan(&tableName); err != nil {
				tb.Fatalf("failed to scan table name: %v", err)
			}
			// Optionally, you might want to exclude tables like migration tables here.
			// Wrap table names in double quotes to account for any special characters.
			tableNames = append(tableNames, fmt.Sprintf(`"%s"`, tableName))
		}

		// If any tables were found, build and execute a TRUNCATE query.
		if len(tableNames) > 0 {
			truncateSQL := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE;", strings.Join(tableNames, ", "))
			if err := gormDB.Exec(truncateSQL).Error; err != nil {
				tb.Fatalf("failed to truncate tables: %v", err)
			}
		}
	}
}
