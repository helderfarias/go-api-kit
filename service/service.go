package service

import "context"

type Service func(ctx context.Context) (interface{}, error)

// Middleware is a chainable behavior modifier for services.
type Middleware func(Service) Service
