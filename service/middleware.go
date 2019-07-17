package service

import (
	"context"

	"github.com/helderfarias/go-api-kit/constants"
	"github.com/helderfarias/go-api-kit/db"
	"github.com/sirupsen/logrus"
)

func Database(dbfactory db.ConnectionFactory, key constants.DatabaseContextValue) Middleware {
	return func(next Service) Service {
		return func(parent context.Context) (interface{}, error) {
			conn := dbfactory.NewConnection()
			if conn != nil {
				return nil, nil
			}

			ctx := context.WithValue(parent, constants.DatabaseContextValue(key), conn)
			return next(ctx)
		}
	}
}

func DatabaseWithTx(dbfactory db.ConnectionFactory, key constants.DatabaseContextValue) Middleware {
	return func(next Service) Service {
		return func(parent context.Context) (interface{}, error) {
			tx, err := dbfactory.NewConnectionWithTransaction()
			if err != nil {
				return nil, err
			}

			ctx := context.WithValue(parent, constants.DatabaseContextValue(key), tx)

			defer func() {
				if r := recover(); r != nil {
					if err := tx.Rollback(); err != nil {
						logrus.Error(err)
					} else {
						logrus.Error(err)
					}
				}
			}()

			resp, err := next(ctx)
			if err == nil {
				if err := tx.Commit(); err != nil {
					logrus.Error(err)
				}
			} else {
				if err := tx.Rollback(); err != nil {
					logrus.Error(err)
				}
			}

			return resp, err
		}
	}
}
