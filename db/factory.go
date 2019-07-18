package db

import (
	"errors"

	wrapper "github.com/helderfarias/sqlx-wrapper/db"
	"github.com/jmoiron/sqlx"
)

type ConnectionFactory interface {
	NewConnection() wrapper.UnitOfWork

	NewConnectionWithTransaction() (wrapper.UnitOfWork, error)

	Delegate() interface{}

	Close() error
}

type connectionFactory struct {
	db *sqlx.DB
}

func NewConnectionFactory(delegate *sqlx.DB) (ConnectionFactory, error) {
	if err := delegate.Ping(); err != nil {
		return nil, err
	}

	return &connectionFactory{db: delegate}, nil
}

func NewPostgresConnectionFactory(ds string, poolMin, poolMax int) (ConnectionFactory, error) {
	conn := sqlx.MustOpen("postgres", ds)
	if err := conn.Ping(); err != nil {
		return nil, err
	}

	conn.SetMaxIdleConns(poolMin)
	conn.SetMaxOpenConns(poolMax)
	return &connectionFactory{db: conn}, nil
}

func (f *connectionFactory) NewConnection() wrapper.UnitOfWork {
	return wrapper.NewUnitOfWork(f.db, nil)
}

func (f *connectionFactory) NewConnectionWithTransaction() (wrapper.UnitOfWork, error) {
	tx := f.db.MustBegin()
	if tx == nil {
		return nil, errors.New("Could not start transaction.")
	}

	return wrapper.NewUnitOfWork(nil, tx), nil
}

func (f *connectionFactory) Delegate() interface{} {
	return f.db
}

func (f *connectionFactory) Close() error {
	if f.db != nil {
		return f.db.Close()
	}
	return nil
}
