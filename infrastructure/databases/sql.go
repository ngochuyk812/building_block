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

type IDatabase interface {
	GetWriteDB() *sql.DB
	GetReadDB() *sql.DB
}

var _ IDatabase = (*Database)(nil)

type Database struct {
	writeDB *sql.DB
	readDB  *sql.DB
}

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
func NewDatabases(dbType DatabaseType, readDSN, writeDSN, dbName string, logger *zap.Logger) (IDatabase, error) {
	rs := &Database{}
	err := rs.InitDatabases(dbType, readDSN, writeDSN, dbName, logger)
	if err != nil {
		return nil, err
	}
	return rs, nil
}

func (a *Database) InitDatabases(dbType DatabaseType, readDSN, writeDSN, dbName string, logger *zap.Logger) error {
	var err error

	a.writeDB, err = newConnectSql(dbType, dbName, writeDSN, logger)
	if err != nil {
		return fmt.Errorf("failed to connect write DB: %w", err)
	}

	if err := runMigration(dbType.String(), dbName, a.writeDB, logger); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}
	if readDSN != "" {
		a.readDB, err = newConnectSql(dbType, dbName, readDSN, logger)
		if err != nil {
			return fmt.Errorf("failed to connect read DB: %w", err)
		}
	}

	return nil
}
func (a *Database) GetWriteDB() *sql.DB {
	return a.writeDB
}
func (a *Database) GetReadDB() *sql.DB {
	if a.readDB == nil {
		return a.writeDB
	}
	return a.readDB
}
func newConnectSql(dbType DatabaseType, dbName, conenctString string, logger *zap.Logger) (db *sql.DB, err error) {
	logger.Info("Staring Connect SQL...")
	db, err = sql.Open(dbType.String(), conenctString)
	if err != nil {
		logger.Error(fmt.Sprintf("Connect SQL error...%v", err))
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(0)

	if err := db.Ping(); err != nil {
		logger.Error(fmt.Sprintf("Error pinging database: %v", err))
	}
	logger.Info(fmt.Sprintf("Connect SQL successfull...%v", err))

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
