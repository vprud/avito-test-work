package repo

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/cut4cut/avito-test-work/internal/entity"
	"github.com/cut4cut/avito-test-work/pkg/postgres"
	"github.com/georgysavva/scany/pgxscan"
	pgx "github.com/jackc/pgx/v4"
)

const _defaultEntityCap = 64
const isoLevel = pgx.Serializable

// AccountRepo - repository with account.
type AccountRepo struct {
	*postgres.Postgres
}

// New - create new account repository.
func New(pg *postgres.Postgres) *AccountRepo {
	return &AccountRepo{pg}
}

func selectTransactionType(amount float64) (string, error) {
	if amount > 0 {
		return "accrual", nil
	} else if amount < 0 {
		return "redeem", nil
	}
	return "", errors.New("amount in transaction is zero")
}

// Create - create new account with default values.
func (r *AccountRepo) Create(ctx context.Context) (acc entity.Account, err error) {
	sql, _, err := r.Builder.
		Insert("account").
		Columns("id, balance, created_dt").
		Values(
			sq.Expr("DEFAULT"),
			sq.Expr("DEFAULT"),
			sq.Expr("DEFAULT")).
		Suffix("RETURNING \"id\", \"balance\", \"created_dt\"").
		ToSql()
	if err != nil {
		return acc, fmt.Errorf("AccountRepo - Create - r.Builder: %w", err)
	}

	err = r.Pool.QueryRow(ctx, sql).Scan(&acc.Id, &acc.Balance, &acc.CreatedDt)
	if err != nil {
		return acc, fmt.Errorf("AccountRepo - updBalance - tx.QueryRow: %w", err)
	}

	return
}

// GetByID - get account's values by ID.
func (r *AccountRepo) GetById(ctx context.Context, id int64) (acc entity.Account, err error) {
	sql, _, err := r.Builder.
		Select("id, balance, created_dt").
		From("account").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return acc, fmt.Errorf("AccountRepo - GetByID - r.Builder: %w", err)
	}

	err = r.Pool.QueryRow(context.Background(), sql, id).Scan(&acc.Id, &acc.Balance, &acc.CreatedDt)
	if err != nil {
		return acc, fmt.Errorf("AccountRepo - GetByID - r.Pool.QueryRow: %w", err)
	}

	return
}

// updBalance - helper function to update the balance.
func (r *AccountRepo) updBalance(ctx context.Context, tx *pgx.Tx, transType string, id, docNum int64, amount float64) (acc entity.Account, err error) {
	sql, _, err := r.Builder.
		Select("balance").
		From("account").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return acc, fmt.Errorf("AccountRepo - updBalance - r.Builder: %w", err)
	}

	balance := 0.0
	err = (*tx).QueryRow(ctx, sql, id).Scan(&balance)
	if err != nil {
		return acc, fmt.Errorf("AccountRepo - updBalance - tx.QueryRow: %w", err)
	}

	if balance+amount < 0 {
		return acc, fmt.Errorf("AccountRepo - updBalance - tx.QueryRow: %w", ErrNotEnoughMoney)
	}

	sqlUpd, _, err := r.Builder.
		Update("account").
		Set("balance", sq.Expr("balance + $2")).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING \"id\", \"balance\", \"created_dt\"").
		ToSql()
	if err != nil {
		return acc, fmt.Errorf("AccountRepo - updBalance - r.Builder: %w", err)
	}

	sqlIns, _, err := r.Builder.
		Insert("fct_transcation").
		Values(
			sq.Expr("DEFAULT"),
			sq.Expr("DEFAULT"),
			id,
			docNum,
			transType,
			amount).
		ToSql()
	if err != nil {
		return acc, fmt.Errorf("AccountRepo - updBalance - r.Builder: %w", err)
	}

	err = (*tx).QueryRow(ctx, sqlUpd, id, amount).Scan(&acc.Id, &acc.Balance, &acc.CreatedDt)
	if err != nil {
		return acc, fmt.Errorf("AccountRepo - updBalance - tx.QueryRow: %w", err)
	}

	_, err = (*tx).Exec(ctx, sqlIns, id, docNum, transType, amount)
	if err != nil {
		return acc, fmt.Errorf("AccountRepo - updBalance - tx.Exec: %w", err)
	}

	return
}

// UpdBalance - update account's balance.
func (r *AccountRepo) UpdBalance(ctx context.Context, id, docNum int64, amount float64) (acc entity.Account, err error) {
	transType, err := selectTransactionType(amount)
	if err != nil {
		return acc, fmt.Errorf("AccountRepo - UpdBalance - selectTransactionType: %w", err)
	}

	tx, err := r.Pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: isoLevel})
	if err != nil {
		return acc, err
	}
	defer tx.Rollback(ctx)

	acc, err = r.updBalance(ctx, &tx, transType, id, docNum, amount)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, pgx.ErrTxCommitRollback) {
			return acc, fmt.Errorf("AccountRepo - UpdBalance - r.updBalance TRUXAuxa: %w", err)
		}
		return acc, fmt.Errorf("AccountRepo - UpdBalance - r.updBalance: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, pgx.ErrTxCommitRollback) {
			return acc, fmt.Errorf("AccountRepo - UpdBalance - TRUXAuxa: %w", err)
		}
		return acc, fmt.Errorf("AccountRepo - UpdBalance - tx.Commit: %w", err)
	}

	return
}

// TransferAmount - transfer amount of money from redeem account to accrual account.
func (r *AccountRepo) TransferAmount(ctx context.Context, redeemId, accrId int64, amount float64) (accrAcc, redeemAcc entity.Account, err error) {
	tx, err := r.Pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: isoLevel})
	if err != nil {
		return accrAcc, redeemAcc, err
	}
	defer tx.Rollback(ctx)

	redeemAcc, err = r.updBalance(ctx, &tx, "redeem", redeemId, accrId, -amount)
	if err != nil {
		return accrAcc, redeemAcc, fmt.Errorf("AccountRepo - TransferAmount - r.updBalance: %w", err)
	}

	accrAcc, err = r.updBalance(ctx, &tx, "accrual", accrId, redeemId, amount)
	if err != nil {
		return accrAcc, redeemAcc, fmt.Errorf("AccountRepo - TransferAmount - r.updBalance: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return accrAcc, redeemAcc, fmt.Errorf("AccountRepo - TransferAmount - tx.Commit: %w", err)
	}

	return
}

// GetHistory - get history of transaction.
func (r *AccountRepo) GetHistory(ctx context.Context, id int64, limit, offset uint64, sort string, isDecreasing bool) (trns []*entity.Transaction, err error) {
	pred := fmt.Sprintf("%s ASC", sort)

	if isDecreasing {
		pred = fmt.Sprintf("%s DESC", sort)
	}

	sql, _, err := r.Builder.
		Select("id, trans_dt, account_id, doc_num, type, amount").
		From("fct_transcation").
		Where(sq.Eq{"account_id": id}).
		OrderBy(pred).
		Limit(limit).
		Offset(offset).
		ToSql()
	if err != nil {
		return trns, fmt.Errorf("AccountRepo - GetHistory - r.Builder: %w", err)
	}

	if err := pgxscan.Select(
		ctx, r.Pool, &trns, sql, id,
	); err != nil {
		return nil, fmt.Errorf("AccountRepo - GetHistory - pgxscan.Select: %w", err)
	}

	return
}
