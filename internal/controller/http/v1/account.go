package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/cut4cut/avito-test-work/internal/entity"
	"github.com/cut4cut/avito-test-work/internal/usecase"
	"github.com/cut4cut/avito-test-work/pkg/logger"
)

type accountRoutes struct {
	u usecase.AccountUseCase
	l logger.Interface
}

func newAccountRoutes(handler *gin.RouterGroup, u usecase.AccountUseCase, l logger.Interface) {
	r := &accountRoutes{u, l}

	h := handler.Group("/account")
	{
		h.GET("/history/:id", r.getHistory)
		h.POST("/", r.create)
		h.GET("/:id", r.getById)
		h.PUT("/:id", r.updBalance)
		h.PUT("/amount/:redeemId/transfer/:accrId", r.transferAmount)
	}
}

type correctResponse struct {
	Data interface{} `json:"data"`
}

type transferAccountPair struct {
	AccrAcc   entity.Account `json:"accrualAccount"`
	RedeemAcc entity.Account `json:"redeemAccount"`
}

// @Summary     Create new account
// @Description Create a new account with default fields and return in the response
// @ID          create
// @Tags  	    account
// @Accept      json
// @Produce     json
// @Success     200 {object} correctResponse
// @Failure     500 {object} response
// @Router      /account [post]
func (r *accountRoutes) create(c *gin.Context) {
	account, err := r.u.Create(c.Request.Context())
	if err != nil {
		r.l.Error(err, "http - v1 - create")
		errorMassage := fmt.Sprint("internal Error: ", errors.Unwrap(err))
		errorResponse(c, http.StatusInternalServerError, errorMassage)

		return
	}

	c.JSON(http.StatusOK, correctResponse{account})
}

// @Summary     Get account by ID
// @Description Returns account fields by ID in the response
// @ID          getById
// @Tags  	    account
// @Accept      json
// @Produce     json
// @Param       id   path      int  true  "Account ID"
// @Success     200 {object} correctResponse
// @Failure     500 {object} response
// @Router      /account/{id} [get]
func (r *accountRoutes) getById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		r.l.Error(err, "http - v1 - getById")
		errorResponse(c, http.StatusBadRequest, "incorrect account ID")

		return
	}

	account, err := r.u.GetById(c.Request.Context(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - getById")
		errorMassage := fmt.Sprint("internal Error: ", errors.Unwrap(err))
		errorResponse(c, http.StatusInternalServerError, errorMassage)

		return
	}

	c.JSON(http.StatusOK, correctResponse{account})
}

// @Summary     Update balance
// @Description Changing the account balance to the amount passed in the parameter
// @ID          updBalance
// @Tags  	    account
// @Accept      json
// @Produce     json
// @Param       id   path      int  true  "Account ID"
// @Param       amount    query     number  true  "The value by which the balance changes"
// @Success     200 {object} correctResponse
// @Failure     500 {object} response
// @Router      /account/{id} [put]
func (r *accountRoutes) updBalance(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		r.l.Error(err, "http - v1 - updBalance")
		errorResponse(c, http.StatusBadRequest, "incorrect account ID")

		return
	}

	amount, err := strconv.ParseFloat(c.Request.URL.Query().Get("amount"), 10)
	if err != nil {
		r.l.Error(err, "http - v1 - updBalance")
		errorResponse(c, http.StatusBadRequest, "incorrect amount")

		return
	}

	account, err := r.u.UpdBalance(c.Request.Context(), id, amount)
	if err != nil {
		r.l.Error(err, "http - v1 - updBalance")
		errorMassage := fmt.Sprint("internal Error: ", errors.Unwrap(err))
		errorResponse(c, http.StatusInternalServerError, errorMassage)

		return
	}

	c.JSON(http.StatusOK, correctResponse{account})
}

// @Summary     Money transaction
// @Description Transferring amount of money between accounts
// @ID          transferAmount
// @Tags  	    account
// @Accept      json
// @Produce     json
// @Param       redeemId   path      int  true  "Account ID for redeem funds"
// @Param       accrId   path      int  true  "Account ID for accrual funds"
// @Param       amount    query     number  true  "Amount of money to transfer"
// @Success     200 {object} transferAccountPair
// @Failure     500 {object} response
// @Router      /account/amount/{redeemId}/transfer/{accrId} [put]
func (r *accountRoutes) transferAmount(c *gin.Context) {
	// redeemId, accrId int64, amount
	redeemId, err := strconv.ParseInt(c.Param("redeemId"), 10, 64)
	if err != nil {
		r.l.Error(err, "http - v1 - transferAmount")
		errorResponse(c, http.StatusBadRequest, "incorrect redeem's ID")

		return
	}

	accrId, err := strconv.ParseInt(c.Param("accrId"), 10, 64)
	if err != nil {
		r.l.Error(err, "http - v1 - transferAmount")
		errorResponse(c, http.StatusBadRequest, "incorrect accrual's ID")

		return
	}

	amount, err := strconv.ParseFloat(c.Request.URL.Query().Get("amount"), 10)
	if err != nil {
		r.l.Error(err, "http - v1 - transferAmount")
		errorResponse(c, http.StatusBadRequest, "incorrect amount")

		return
	}

	accrAcc, redeemAcc, err := r.u.TransferAmount(c.Request.Context(), redeemId, accrId, amount)
	if err != nil {
		r.l.Error(err, "http - v1 - transferAmount")
		errorMassage := fmt.Sprint("internal Error: ", errors.Unwrap(err))
		errorResponse(c, http.StatusInternalServerError, errorMassage)

		return
	}

	pair := transferAccountPair{AccrAcc: accrAcc, RedeemAcc: redeemAcc}

	c.JSON(http.StatusOK, correctResponse{pair})
}

// @Summary     Transaction history
// @Description Return history of all account's transactions
// @ID          history
// @Tags  	    account
// @Accept      json
// @Produce     json
// @Param       id   path      int  true  "Account ID"
// @Param       limit    query     int  true  "The value of limit in pagination"
// @Param       offset    query     int  true  "The value of offset in pagination"
// @Param       sort    query     string  false  "Column name to sort"
// @Param       isDecreasing    query     bool  false  "Descending sort flag"
// @Success     200 {object} correctResponse
// @Failure     500 {object} response
// @Router      /account/history/{id} [get]
func (r *accountRoutes) getHistory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		r.l.Error(err, "http - v1 - getHistory")
		errorResponse(c, http.StatusBadRequest, "incorrect account ID")

		return
	}

	limit, err := strconv.ParseUint(c.Request.URL.Query().Get("limit"), 10, 64)
	if err != nil {
		r.l.Error(err, "http - v1 - getHistory")
		errorResponse(c, http.StatusBadRequest, "incorrect limit value")

		return
	}

	offset, err := strconv.ParseUint(c.Request.URL.Query().Get("offset"), 10, 64)
	if err != nil {
		r.l.Error(err, "http - v1 - getHistory")
		errorResponse(c, http.StatusBadRequest, "incorrect offset value")

		return
	}

	isDecreasing := c.Request.URL.Query().Get("isDecreasing")
	sort := c.Request.URL.Query().Get("sort")

	transactions, err := r.u.GetHistory(c.Request.Context(), id, limit, offset, sort, isDecreasing)
	if err != nil {
		r.l.Error(err, "http - v1 - history")
		errorMassage := fmt.Sprint("internal Error: ", errors.Unwrap(err))
		errorResponse(c, http.StatusInternalServerError, errorMassage)

		return
	}

	c.JSON(http.StatusOK, correctResponse{transactions})
}
