package dto

import (
	"github.com/google/uuid"
)

type BillingAddress struct {
	Street  string `json:"street" binding:"required"`
	City    string `json:"city" binding:"required"`
	State   string `json:"state" binding:"required,len=2"`
	Country string `json:"country" binding:"required"`
	Zipcode string `json:"zipcode" binding:"required, numeric,len=5"`
}

type BillingInfo struct {
	CreditCardNumber string         `json:"creditCardNumber" binding:"required, numeric"`
	CreditCardType   string         `json:"creditCardType" binding:"required,oneof=VISA MASTERCARD AMEX"`
	ExpirationDate   string         `json:"expirationDate" binding:"required, datetime=01/2006"`
	CVV              string         `json:"cvv" binding:"required,numeric,len=3"`
	BillingAddress   BillingAddress `json:"billingAddress" binding:"required, dive"`
}

type CreateAccountInput struct {
	Name        string      `json:"name" binding:"required, min=3,max=100"`
	Description string      `json:"description" binding:"min=3,max=500"`
	TenantID    uuid.UUID   `json:"tenantId" binding:"required, uuid4"`
	ParentID    uuid.UUID   `json:"parentId" binding:"required, uuid4"`
	BillingInfo BillingInfo `json:"billingInfo" binding:"required, dive"`
}
