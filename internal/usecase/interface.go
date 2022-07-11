// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/cut4cut/avito-test-work/internal/entity"
)

//go:generate mockgen -destination=./mocks_test.go -package=usecase_test github.com/cut4cut/avito-test-work/internal/usecase AccountRepo

type (
	// AccountRepo -.
	AccountRepo interface {
		Create(context.Context) (entity.Account, error)
		GetById(context.Context, int64) (entity.Account, error)
		UpdBalance(context.Context, int64, int64, float64) (entity.Account, error)
		TransferAmount(context.Context, int64, int64, float64) (entity.Account, entity.Account, error)
		GetHistory(context.Context, int64, uint64, uint64, string, bool) ([]*entity.Transaction, error)
	}
)
