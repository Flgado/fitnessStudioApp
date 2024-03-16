package dbfactory

import (
	"fmt"
	"time"

	"github.com/Flgado/fitnessStudioApp/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DBfactory interface {
	GetDbContext() (*sqlx.DB, error)
}

const (
	maxOpenConns    = 60
	connMaxLifetime = 120
	maxIdleConns    = 30
	connMaxIdleTime = 20
)

// DBFactory is an interface for creating database connections.
type DBFactory interface {
	GetDbContext() (*sqlx.DB, error)
}

// DBFactoryImpl is an implementation of the DBFactory interface.
type dBFactoryImpl struct {
	config config.Config
}

// NewDBFactory creates a new instance of DBFactoryImpl.
func NewDBFactory(c config.Config) DBFactory {
	return &dBFactoryImpl{config: c}
}

func (df *dBFactoryImpl) GetDbContext() (*sqlx.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		df.config.Postgres.PostgresqlHost,
		df.config.Postgres.PostgresqlPort,
		df.config.Postgres.PostgresqlUser,
		df.config.Postgres.PostgresqlDbname,
		df.config.Postgres.PostgresqlPassword,
	)

	db, err := sqlx.Connect(df.config.Postgres.PgDriver, dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifetime * time.Second)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(connMaxIdleTime * time.Second)

	return db, nil
}
