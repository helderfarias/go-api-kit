package db

import (
	"errors"

	wrapper "github.com/helderfarias/sqlx-wrapper/db"
	"github.com/jmoiron/sqlx"
)

type ConnectionFactory interface {
	NewConnection() wrapper.UnitOfWork

	NewConnectionWithTransaction() (wrapper.UnitOfWork, error)
}

type postgresConnectionFactory struct {
	db *sqlx.DB
}

func NewPostgresConnectionFactory(ds string, poolMin, poolMax int) (ConnectionFactory, error) {
	conn := sqlx.MustOpen("postgres", ds)
	if err := conn.Ping(); err != nil {
		return nil, err
	}

	conn.SetMaxIdleConns(poolMin)
	conn.SetMaxOpenConns(poolMax)
	return &postgresConnectionFactory{db: conn}, nil
}

func (f *postgresConnectionFactory) NewConnection() wrapper.UnitOfWork {
	return wrapper.NewUnitOfWork(f.db, nil)
}

func (f *postgresConnectionFactory) NewConnectionWithTransaction() (wrapper.UnitOfWork, error) {
	tx := f.db.MustBegin()
	if tx == nil {
		return nil, errors.New("Nenhuma transação foi iniciada.")
	}

	return wrapper.NewUnitOfWork(nil, tx), nil
}
