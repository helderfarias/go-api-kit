package endpoint

import (
	"context"

	"github.com/helderfarias/go-api-kit/constants"
	"github.com/helderfarias/go-api-kit/db"
	"github.com/sirupsen/logrus"
)

func Database(dbfactory db.ConnectionFactory, key constants.DatabaseContextValue) Middleware {
	return func(next Endpoint) Endpoint {
		return func(parent context.Context, request interface{}) (response EndpointResponse, err error) {
			conn := dbfactory.NewConnection()
			if err != nil {
				return nil, err
			}

			ctx := context.WithValue(parent, constants.DatabaseContextValue(key), conn)
			return next(ctx, request)
		}
	}
}

func DatabaseWithTx(dbfactory db.ConnectionFactory, key constants.DatabaseContextValue) Middleware {
	return func(next Endpoint) Endpoint {
		return func(parent context.Context, request interface{}) (response EndpointResponse, err error) {
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

			resp, err := next(ctx, request)
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
