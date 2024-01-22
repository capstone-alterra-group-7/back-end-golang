package usecases

import (
	"back-end-golang/configs"
	"back-end-golang/dtos"
	"context"
	"strconv"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

var s snap.Client

var c coreapi.Client

func InitiateCoreApiClient() {
	c.New(configs.EnvMidtransServerKey(), midtrans.Sandbox)
}

func InitializeSnapClient() {
	s.New(configs.EnvMidtransServerKey(), midtrans.Sandbox)
}

func CreateUrlTransactionWithGateway(input dtos.MidtransInput) (string, error) {
	InitializeSnapClient()
	s.Options.SetContext(context.Background())

	resp, err := s.CreateTransactionUrl(GenerateSnapReq(input))
	if err != nil {
		return "", err
	}
	return resp, nil
}

func GenerateSnapReq(input dtos.MidtransInput) *snap.Request {

	// Initiate Customer address
	custAddress := &midtrans.CustomerAddress{
		FName:       input.CustomerAddress.FName,
		LName:       "- Tripease",
		Phone:       input.CustomerAddress.Phone,
		Address:     input.CustomerAddress.Address,
		City:        input.CustomerAddress.City,
		Postcode:    input.CustomerAddress.Postcode,
		CountryCode: "IDN",
	}

	// Initiate Snap Request
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  input.TransactionDetails.OrderID,
			GrossAmt: int64(input.TransactionDetails.GrossAmt),
		},
		Expiry: &snap.ExpiryDetails{
			Unit:     "minutes",
			Duration: 60,
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName:    input.CustomerDetail.FName,
			LName:    input.CustomerDetail.LName,
			Email:    input.CustomerDetail.Email,
			Phone:    input.CustomerDetail.Phone,
			BillAddr: custAddress,
			ShipAddr: custAddress,
		},
		EnabledPayments: snap.AllSnapPaymentType,
		Items: &[]midtrans.ItemDetails{
			{
				ID:    strconv.Itoa(input.Items.ID),
				Price: int64(input.Items.Price),
				Qty:   int32(input.Items.Qty),
				Name:  input.Items.Name,
			},
		},
	}
	return snapReq
}
