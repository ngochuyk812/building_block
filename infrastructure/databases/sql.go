package databases

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"go.uber.org/zap"
)

type DatabaseType string

const (
	MYSQL    = "mysql"
	POSTGRES = "postgres"
)

func (s DatabaseType) String() string {
	switch s {
	case MYSQL:
		return "mysql"
	case POSTGRES:
		return "postgres"
	default:
		return "unknown"
	}

}

func NewConnectSql(dbType DatabaseType, dbName, conenctString string, logger *zap.Logger) (db *sql.DB, err error) {
	logger.Info("Staring Connect SQL...")
	db, err = sql.Open("mysql", "root:123321huy@tcp(localhost:3316)/auth-service?tls=false")
	if err != nil {
		logger.Error(fmt.Sprintf("Connect SQL error...%v", err))
		return nil, err
	}
	defer db.Close()

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(0)

	if err := db.Ping(); err != nil {
		logger.Error(fmt.Sprintf("Error pinging database: %v", err))
	}
	logger.Info(fmt.Sprintf("Connect SQL successfull...%v", err))

	err = runMigration(dbType.String(), dbName, db, logger)
	return db, err
}

func runMigration(dbType, dbName string, db *sql.DB, logger *zap.Logger) error {
	var driver database.Driver
	var err error
	logger.Info(fmt.Sprintf("Migragion SQL...: %v, %v", dbType, dbName))

	switch dbType {
	case MYSQL:
		driver, err = mysql.WithInstance(db, &mysql.Config{})
	case POSTGRES:
		driver, err = postgres.WithInstance(db, &postgres.Config{})
	default:
		return fmt.Errorf("unsupported DB driver: %s", dbType)
	}
	if err != nil {
		return fmt.Errorf("migration driver error: %w", err)
	}

	sourceURL := fmt.Sprintf("file://migration")
	m, err := migrate.NewWithDatabaseInstance(sourceURL, dbName, driver)
	if err != nil {
		logger.Error(fmt.Sprintf("migration init error: %v", err))
		return fmt.Errorf("migration init error: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Error(fmt.Sprintf("migration run error: %v", err))
		return fmt.Errorf("migration run error: %w", err)
	}
	logger.Info(fmt.Sprintf("Migragion SQL successfull"))

	return nil
}
