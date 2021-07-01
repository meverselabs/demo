package helloworld

import "errors"

// errors
var (
	ErrPolicyShouldBeSetupInApplication = errors.New("policy should be setup in application")
	ErrMinusInput                       = errors.New("minus input")
	ErrMinusPoint                       = errors.New("minus point")
	ErrIsNotFranchiseAccount            = errors.New("is not franchise account")
	ErrInvalidPaymentID                 = errors.New("invalid paymentID")
	ErrExpiredPayment                   = errors.New("expired payment")
	ErrNotExistExchangeRequest          = errors.New("not exist exchange request")
	ErrLackOfPoint                      = errors.New("lack of point")
	ErrInvalidPaymentStatus             = errors.New("invalid payment status")
	ErrInvalidInvoiceAddr               = errors.New("invalid ivoice address")
	ErrNotEnoughtPoint                  = errors.New("not enought point")
	ErrNotExistRefund                   = errors.New("not exist refund")
	ErrIsNotInvoiceTx                   = errors.New("isNot invoice tx")
	ErrIsNotPaymentTx                   = errors.New("isNot payment tx")
	ErrInvalidRequestTXID               = errors.New("invalid request TXID")
	ErrRequestAlreadyProcessed          = errors.New("request already processed")
	ErrNotExistFranchise                = errors.New("not exist franchise")
	ErrNotInitPolicy                    = errors.New("not init policy")
	ErrInvalidPaymentAmount             = errors.New("invalid payment amount")
	ErrAlreadyRefund                    = errors.New("already refund")
	ErrInvalidTagSize                   = errors.New("invalid tag size")
)
