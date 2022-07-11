package usecase_test

import (
	"context"
	"time"

	"testing"

	"github.com/cut4cut/avito-test-work/internal/entity"
	"github.com/cut4cut/avito-test-work/internal/usecase"
	"github.com/golang/mock/gomock"
)

func TestAccountUseCase_Create(t *testing.T) {
	type fields struct {
		ctx         context.Context
		accountRepo *MockAccountRepo
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		wantErr bool
	}{
		{
			name: "Case of correct work",
			prepare: func(f *fields) {
				f.accountRepo.EXPECT().Create(f.ctx).Return(entity.Account{Id: 1, Balance: 0.0, CreatedDt: time.Now()}, nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				ctx:         context.Background(),
				accountRepo: NewMockAccountRepo(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			uc := usecase.New(f.accountRepo)
			if acc, err := uc.Create(f.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Create() account=%v error = %v, wantErr %v", acc, err, tt.wantErr)
			}
		})
	}
}

func TestAccountUseCase_GetById(t *testing.T) {
	type fields struct {
		ctx         context.Context
		accountRepo *MockAccountRepo
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		arg     int64
		wantErr bool
	}{
		{
			name: "Case of correct work",
			prepare: func(f *fields) {
				f.accountRepo.EXPECT().GetById(f.ctx, int64(1)).Return(entity.Account{Id: 1, Balance: 0.0, CreatedDt: time.Now()}, nil)
			},
			arg:     1,
			wantErr: false,
		},
		{
			name:    "Case of incorrect work: ID is negative",
			prepare: func(f *fields) {},
			arg:     -4,
			wantErr: true,
		},
		{
			name:    "Case of incorrect work: ID is zero",
			prepare: func(f *fields) {},
			arg:     0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				ctx:         context.Background(),
				accountRepo: NewMockAccountRepo(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			uc := usecase.New(f.accountRepo)
			if acc, err := uc.GetById(f.ctx, tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("GetById() account=%v error = %v, wantErr %v", acc, err, tt.wantErr)
			}
		})
	}
}

func TestAccountUseCase_UpdBalance(t *testing.T) {
	type fields struct {
		ctx         context.Context
		accountRepo *MockAccountRepo
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		arg1    int64
		arg2    float64
		wantErr bool
	}{
		{
			name: "Case of correct work",
			prepare: func(f *fields) {
				f.accountRepo.EXPECT().UpdBalance(f.ctx, int64(1), int64(-999), float64(25.0)).Return(entity.Account{Id: 1, Balance: 25.0, CreatedDt: time.Now()}, nil)
			},
			arg1:    1,
			arg2:    25,
			wantErr: false,
		},
		{
			name:    "Case of incorrect work: amount is zero",
			prepare: func(f *fields) {},
			arg1:    1,
			arg2:    0,
			wantErr: true,
		},
		{
			name:    "Case of incorrect work: ID is negative",
			prepare: func(f *fields) {},
			arg1:    -89,
			arg2:    25,
			wantErr: true,
		},
		{
			name:    "Case of incorrect work: ID is zero",
			prepare: func(f *fields) {},
			arg1:    0,
			arg2:    25,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				ctx:         context.Background(),
				accountRepo: NewMockAccountRepo(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			uc := usecase.New(f.accountRepo)
			if acc, err := uc.UpdBalance(f.ctx, tt.arg1, tt.arg2); (err != nil) != tt.wantErr {
				t.Errorf("UpdBalance() account=%v error = %v, wantErr %v", acc, err, tt.wantErr)
			}
		})
	}
}

func TestAccountUseCase_TransferAmount(t *testing.T) {
	type fields struct {
		ctx         context.Context
		accountRepo *MockAccountRepo
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		arg1    int64
		arg2    int64
		arg3    float64
		wantErr bool
	}{
		{
			name: "Case of correct work",
			prepare: func(f *fields) {
				f.accountRepo.EXPECT().TransferAmount(f.ctx, int64(1), int64(2), float64(5.0)).Return(
					entity.Account{Id: 1, Balance: 25.0, CreatedDt: time.Now()},
					entity.Account{Id: 2, Balance: 30.0, CreatedDt: time.Now()},
					nil)
			},
			arg1:    1,
			arg2:    2,
			arg3:    5,
			wantErr: false,
		},
		{
			name:    "Case of incorrect work: amount is zero",
			prepare: func(f *fields) {},
			arg1:    1,
			arg2:    2,
			arg3:    0,
			wantErr: true,
		},
		{
			name:    "Case of incorrect work: redeem ID is negative",
			prepare: func(f *fields) {},
			arg1:    -1,
			arg2:    2,
			arg3:    5,
			wantErr: true,
		},
		{
			name:    "Case of incorrect work: redeem ID is zero",
			prepare: func(f *fields) {},
			arg1:    0,
			arg2:    2,
			arg3:    5,
			wantErr: true,
		},
		{
			name:    "Case of incorrect work: accrual ID is negative",
			prepare: func(f *fields) {},
			arg1:    1,
			arg2:    -2,
			arg3:    5,
			wantErr: true,
		},
		{
			name:    "Case of incorrect work: accrual ID is zero",
			prepare: func(f *fields) {},
			arg1:    1,
			arg2:    0,
			arg3:    5,
			wantErr: true,
		},
		{
			name:    "Case of incorrect work: IDs are the same",
			prepare: func(f *fields) {},
			arg1:    2,
			arg2:    2,
			arg3:    5,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				ctx:         context.Background(),
				accountRepo: NewMockAccountRepo(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			uc := usecase.New(f.accountRepo)
			if accrAcc, redeemAcc, err := uc.TransferAmount(f.ctx, tt.arg1, tt.arg2, tt.arg3); (err != nil) != tt.wantErr {
				t.Errorf("UpdBalance() accrual account=%v redeem account=%v error = %v, wantErr %v", accrAcc, redeemAcc, err, tt.wantErr)
			}
		})
	}
}

func TestAccountUseCase_GetHistory(t *testing.T) {
	type fields struct {
		ctx         context.Context
		accountRepo *MockAccountRepo
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		arg1    int64
		arg2    uint64
		arg3    uint64
		arg4    string
		arg5    string
		wantErr bool
	}{
		{
			name: "Case of correct work",
			prepare: func(f *fields) {
				f.accountRepo.EXPECT().GetHistory(f.ctx, int64(1), uint64(3), uint64(0), "trans_dt", true).Return(
					[]*entity.Transaction{
						&entity.Transaction{Id: 1, TransDt: time.Now(), AccountId: 1, DocNum: 2, Type: "redeem", Amount: 5},
						&entity.Transaction{Id: 2, TransDt: time.Now(), AccountId: 1, DocNum: 2, Type: "redeem", Amount: 5},
						&entity.Transaction{Id: 3, TransDt: time.Now(), AccountId: 1, DocNum: 2, Type: "redeem", Amount: 5},
					},
					nil)
			},
			arg1:    1,
			arg2:    3,
			arg3:    0,
			arg4:    "trans_dt",
			arg5:    "true",
			wantErr: false,
		},
		{
			name:    "Case of incorrect work: ID is negative",
			prepare: func(f *fields) {},
			arg1:    -1,
			arg2:    3,
			arg3:    0,
			arg4:    "trans_dt",
			arg5:    "true",
			wantErr: true,
		},
		{
			name:    "Case of incorrect work: ID is zero",
			prepare: func(f *fields) {},
			arg1:    0,
			arg2:    3,
			arg3:    0,
			arg4:    "trans_dt",
			arg5:    "true",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				ctx:         context.Background(),
				accountRepo: NewMockAccountRepo(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			uc := usecase.New(f.accountRepo)
			if trans, err := uc.GetHistory(f.ctx, tt.arg1, tt.arg2, tt.arg3, tt.arg4, tt.arg5); (err != nil) != tt.wantErr {
				t.Errorf("GetHistory() trans history=%v error = %v, wantErr %v", trans, err, tt.wantErr)
			}
		})
	}
}
