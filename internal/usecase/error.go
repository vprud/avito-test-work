package usecase

import "errors"

var (
	ErrorAmountIsNegative error = errors.New("amount is negative")
	ErrorAmountIsZero     error = errors.New("amount is zero")
	ErrorIdIsNegative     error = errors.New("ID is negative")
	ErrorIdIsZero         error = errors.New("ID is zero")
	ErrorSameRedeemAccrId error = errors.New("redeem and accrual ID are the same")
)
