package integration_test

import (
	"log"
	"net/http"
	"os"
	"sort"
	"testing"
	"time"

	. "github.com/Eun/go-hit"
	"github.com/cut4cut/avito-test-work/internal/entity"
	"github.com/stretchr/testify/require"
)

const (
	host     = "app:8080"
	attempts = 20
	basePath = "http://" + host + "/v1"
	requests = 10
)

func TestMain(m *testing.M) {
	err := healthCheck(attempts)
	if err != nil {
		log.Fatalf("Integration tests: host %s is not available: %s", host, err)
	}

	log.Printf("Integration tests: host %s is available", host)

	code := m.Run()
	os.Exit(code)
}

func healthCheck(attempts int) (err error) {
	path := basePath + "/account"

	for attempts > 0 {
		err = Do(Post(path), Expect().Status().Equal(http.StatusOK))
		if err == nil {
			return nil
		}

		log.Printf("Integration tests: url %s is not available, attempts left: %d", path, attempts)

		time.Sleep(time.Second)

		attempts--
	}

	return err
}

// HTTP POST: /account.
func TestHttp_Create(t *testing.T) {
	Test(t,
		Description("Create new account: сase of correct work"),
		Post(basePath+"/account"),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().String().Contains(`{"data":{"id":`),
	)
}

// HTTP GET: /account/:id.
func TestHttp_GetId(t *testing.T) {
	Test(t,
		Description("Get account by ID: case of correct work"),
		Get(basePath+"/account/2"),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().String().Contains(`{"data":{"id":2`),
	)
	Test(t,
		Description("Get account by ID: incorrect ID"),
		Get(basePath+"/account/ry1"),
		Expect().Status().Equal(http.StatusBadRequest),
		Expect().Body().String().Contains(`{"error"`),
	)
	Test(t,
		Description("Get account by ID: negative value of ID"),
		Get(basePath+"/account/-1"),
		Expect().Status().Equal(http.StatusInternalServerError),
		Expect().Body().String().Contains(`internal Error: ID is negative`),
	)
	Test(t,
		Description("Get account by ID: zero value of ID"),
		Get(basePath+"/account/0"),
		Expect().Status().Equal(http.StatusInternalServerError),
		Expect().Body().String().Contains(`internal Error: ID is zero`),
	)
	Test(t,
		Description("Get account by ID: not exists account ID"),
		Get(basePath+"/account/56784"),
		Expect().Status().Equal(http.StatusInternalServerError),
		Expect().Body().String().Contains(`no rows in result set`),
	)
}

// HTTP PUT: /account/:id?amount=.
func TestHttp_UpdateBalance(t *testing.T) {
	Test(t,
		Description("Update account's balance: increase balance for accout with ID=1"),
		Put(basePath+"/account/1?amount=35"),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().String().Contains(`{"data":{"id":1,"balance":35`),
	)
	Test(t,
		Description("Update account's balance: increase balance for accout with ID=2"),
		Put(basePath+"/account/2?amount=5"),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().String().Contains(`{"data":{"id":2,"balance":5`),
	)
	Test(t,
		Description("Update account's balance: decrease balance"),
		Put(basePath+"/account/1?amount=-5"),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().String().Contains(`{"data":{"id":1,"balance":30`),
	)
	Test(t,
		Description("Update account's balance: not enough money for decrease balance"),
		Put(basePath+"/account/1?amount=-105"),
		Expect().Status().Equal(http.StatusInternalServerError),
		Expect().Body().String().Contains(`not enough money`),
	)
}

// HTTP PUT:  /account/amount/:redeemId/transfer/:accrId?amount=.
func TestHttp_TransferAmount(t *testing.T) {
	var accrualId, redeemId int
	var accrualBalance, redeemBalance float64

	Test(t,
		Description("Transfer amount between accounts: case of correct work"),
		Put(basePath+"/account/amount/2/transfer/1?amount=1"),
		Expect().Status().Equal(http.StatusOK),
		Store().Response().Body().JSON().JQ(".data.accrualAccount.id").In(&accrualId),
		Store().Response().Body().JSON().JQ(".data.redeemAccount.id").In(&redeemId),
		Store().Response().Body().JSON().JQ(".data.accrualAccount.balance").In(&accrualBalance),
		Store().Response().Body().JSON().JQ(".data.redeemAccount.balance").In(&redeemBalance),
	)

	require.Equal(t, 1, accrualId)
	require.Equal(t, 2, redeemId)
	require.Equal(t, 31.0, accrualBalance)
	require.Equal(t, 4.0, redeemBalance)

	Test(t,
		Description("Transfer amount between accounts: case of not enough money"),
		Put(basePath+"/account/amount/2/transfer/1?amount=150"),
		Expect().Status().Equal(http.StatusInternalServerError),
		Expect().Body().String().Contains(`not enough money`),
	)

	Test(t,
		Description("Transfer amount between accounts: case of not exists ID"),
		Put(basePath+"/account/amount/2/transfer/150?amount=1"),
		Expect().Status().Equal(http.StatusInternalServerError),
		Expect().Body().String().Contains(`no rows in result set`),
	)
}

// HTTP GET:  /account/history/:id?limit=&offset=&sort=&isDecreasing=
func TestHttp_GetHistory(t *testing.T) {
	var transactions *[]entity.Transaction
	var expectedTransactions []entity.Transaction = []entity.Transaction{
		{Id: 1, AccountId: 1, DocNum: -999, Type: "accrual", Amount: 35},
		{Id: 3, AccountId: 1, DocNum: -999, Type: "redeem", Amount: -5},
		{Id: 5, AccountId: 1, DocNum: -999, Type: "accrual", Amount: 1},
	}

	Test(t,
		Description("Get transaction history: case of correct work"),
		Get(basePath+"/account/history/1?limit=5&offset=0"),
		Expect().Status().Equal(http.StatusOK),
		Store().Response().Body().JSON().JQ(".data").In(&transactions),
	)

	require.Equal(t, 3, len(*transactions))

	for i, transaction := range *transactions {
		require.Equal(t, expectedTransactions[i].Type, transaction.Type)
		require.Equal(t, expectedTransactions[i].Amount, transaction.Amount)
	}

	Test(t,
		Description("Get transaction history: case of correct work"),
		Get(basePath+"/account/history/1?limit=5&offset=0&sort=type&isDecreasing=false"),
		Expect().Status().Equal(http.StatusOK),
		Store().Response().Body().JSON().JQ(".data").In(&transactions),
	)

	require.Equal(t, 3, len(*transactions))

	sort.Slice(expectedTransactions, func(i, j int) bool {
		return expectedTransactions[i].Type < expectedTransactions[j].Type
	})

	for i, transaction := range *transactions {
		require.Equal(t, expectedTransactions[i].Type, transaction.Type)
	}

	Test(t,
		Description("Get transaction history: case of correct work"),
		Get(basePath+"/account/history/1?limit=5&offset=0&sort=amount&isDecreasing=false"),
		Expect().Status().Equal(http.StatusOK),
		Store().Response().Body().JSON().JQ(".data").In(&transactions),
	)

	require.Equal(t, 3, len(*transactions))

	sort.Slice(expectedTransactions, func(i, j int) bool {
		return expectedTransactions[i].Amount < expectedTransactions[j].Amount
	})

	for i, transaction := range *transactions {
		require.Equal(t, expectedTransactions[i].Amount, transaction.Amount)
	}

	Test(t,
		Description("Get transaction history: case of correct work"),
		Get(basePath+"/account/history/1?limit=5&offset=0&sort=amount&isDecreasing=true"),
		Expect().Status().Equal(http.StatusOK),
		Store().Response().Body().JSON().JQ(".data").In(&transactions),
	)

	require.Equal(t, 3, len(*transactions))

	sort.Slice(expectedTransactions, func(i, j int) bool {
		return expectedTransactions[i].Amount > expectedTransactions[j].Amount
	})

	for i, transaction := range *transactions {
		require.Equal(t, expectedTransactions[i].Amount, transaction.Amount)
	}

	Test(t,
		Description("Get transaction history: case of not exists ID"),
		Get(basePath+"/account/history/56784?limit=5&offset=0"),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().String().Contains(`{"data":null`),
	)
	Test(t,
		Description("Get transaction history: case of incorect limit"),
		Get(basePath+"/account/history/1?limit=-5&offset=0"),
		Expect().Status().Equal(http.StatusBadRequest),
		Expect().Body().String().Contains(`incorrect limit value`),
	)
}
