package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/cut4cut/avito-test-work/internal/entity"
)

// AccountUseCase - use case with account.
type AccountUseCase struct {
	repo AccountRepo
}

// New - create new account use case.
func New(r AccountRepo) *AccountUseCase {
	return &AccountUseCase{
		repo: r,
	}
}

func (uc *AccountUseCase) amountValidation(amount float64) (err error) {
	if amount < 0 {
		err = ErrorAmountIsNegative
	} else if amount == 0 {
		err = ErrorAmountIsZero
	}

	return
}

func (uc *AccountUseCase) idValidation(id int64) (err error) {
	if id < 0 {
		err = ErrorIdIsNegative
	} else if id == 0 {
		err = ErrorIdIsZero
	}

	return
}

// Create - create new account with default values.
func (uc *AccountUseCase) Create(ctx context.Context) (acc entity.Account, err error) {
	acc, err = uc.repo.Create(ctx)
	if err != nil {
		return acc, fmt.Errorf("AccountUseCase - Create - uc.repo.Create: %w", err)
	}

	return
}

// GetById - get account's values by ID.
func (uc *AccountUseCase) GetById(ctx context.Context, id int64) (acc entity.Account, err error) {
	err = uc.idValidation(id)
	if err != nil {
		return acc, fmt.Errorf("AccountUseCase - GetById - uc.idValidation: %w", err)
	}

	acc, err = uc.repo.GetById(ctx, id)
	if err != nil {
		return acc, fmt.Errorf("AccountUseCase - GetById - uc.repo.GetById: %w", err)
	}

	return
}

// UpdBalance - update account's balance.
func (uc *AccountUseCase) UpdBalance(ctx context.Context, id int64, amount float64) (acc entity.Account, err error) {
	err = uc.idValidation(id)
	if err != nil {
		return acc, fmt.Errorf("AccountUseCase - UpdBalance - uc.idValidation: %w", err)
	}

	if amount == 0 {
		return acc, fmt.Errorf("AccountUseCase - UpdBalance - uc.idValidation: %w", ErrorAmountIsZero)
	}

	acc, err = uc.repo.UpdBalance(ctx, id, -999, amount)
	if err != nil {
		return acc, fmt.Errorf("AccountUseCase - UpdBalance - uc.repo.UpdBalance: %w", err)
	}

	return
}

// TransferAmount - transfer amount of money from redeem account to accrual account.
func (uc *AccountUseCase) TransferAmount(ctx context.Context, redeemId, accrId int64, amount float64) (accrAcc, redeemAcc entity.Account, err error) {
	if accrId == redeemId {
		return accrAcc, redeemAcc, fmt.Errorf("AccountUseCase - TransferAmount - validation: %w", ErrorSameRedeemAccrId)
	}

	err = uc.amountValidation(amount)
	if err != nil {
		return accrAcc, redeemAcc, fmt.Errorf("AccountUseCase - TransferAmount - uc.amountValidation: %w", err)
	}

	err = uc.idValidation(redeemId)
	if err != nil {
		return accrAcc, redeemAcc, fmt.Errorf("AccountUseCase - TransferAmount - uc.idValidation: %w", err)
	}

	err = uc.idValidation(accrId)
	if err != nil {
		return accrAcc, redeemAcc, fmt.Errorf("AccountUseCase - TransferAmount - uc.idValidation: %w", err)
	}

	accrAcc, redeemAcc, err = uc.repo.TransferAmount(ctx, redeemId, accrId, amount)
	if err != nil {
		return accrAcc, redeemAcc, fmt.Errorf("AccountUseCase - TransferAmount - uc.repo.TransferAmount: %w", err)
	}

	return
}

// GetHistory - get history of transaction.
func (uc *AccountUseCase) GetHistory(ctx context.Context, id int64, limit, offset uint64, sort, isDecreasingValue string) (trans []*entity.Transaction, err error) {
	err = uc.idValidation(id)
	if err != nil {
		return trans, fmt.Errorf("AccountUseCase - GetHistory - uc.idValidation: %w", err)
	}

	if sort == "" {
		sort = "trans_dt"
	}

	isDecreasing := false
	if strings.ToLower(isDecreasingValue) == "true" {
		isDecreasing = true
	}

	trans, err = uc.repo.GetHistory(ctx, id, limit, offset, sort, isDecreasing)
	if err != nil {
		return trans, fmt.Errorf("AccountUseCase - GetHistory - uc.repo.GetHistory: %w", err)
	}

	return
}
