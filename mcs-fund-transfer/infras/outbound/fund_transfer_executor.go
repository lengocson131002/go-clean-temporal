package outbound

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/lengocson131002/mcs-fund-transfer/domain"
	"github.com/lengocson131002/mcs-fund-transfer/usecase/outbound"
)

const (
	PaymentExecutorEndpoint = "http://10.96.60.91:7800/paymentexecution/v1/rest/executePayment"
)

type ExecutePaymentInWrapper struct {
	ExecutePaymentIn ExecutePaymentIn `json:"executePayment_in"`
}

type ExecutePaymentOutWrapper struct {
	ExecutePaymentOut ExecutePaymentOut `json:"executePayment_out"`
}

type ExecutePaymentIn struct {
	TransactionInfo TransactionInfo `json:"transactionInfo"`
	TransferInfo    TransferInfo    `json:"transferInfo"`
	Payment         Payment         `json:"payment"`
	CardInfo        CardInfo        `json:"cardInfo"`
}

type BranchInfo struct {
	BranchCode string `json:"branchCode"`
}

type TransactionInfo struct {
	CRefNum    string     `json:"cRefNum"`
	UserId     string     `json:"userId"`
	ClientCode string     `json:"clientCode"`
	BranchInfo BranchInfo `json:"branchInfo"`
}

type TransferInfo struct {
	PaymentType string `json:"paymentType"`
}

type Payment struct {
	BookingDate        string      `json:"bookingDate"`
	DebitAccount       string      `json:"debitAccount"`
	CreditAccount      string      `json:"creditAccount"`
	CustomerId         string      `json:"customerId"`
	OriginalCustomerId string      `json:"originalCustomerId"`
	CRefNum            string      `json:"cRefNum"`
	Amount             int64       `json:"amount"`
	Currency           string      `json:"currency"`
	Remarks            string      `json:"remarks"`
	EbUserId           string      `json:"ebUserId"`
	MobilePhoneNumber  interface{} `json:"mobilePhoneNumber"`
}

type CardInfo struct {
	CardAccountNum string `json:"cardAccountNum"`
	CardNumber     string `json:"cardNumber"`
}

type ExecutePaymentOut struct {
	TransactionInfo TransactionOutInfo `json:"transactionInfo"`
}

type TransactionOutInfo struct {
	CoreRefNum               string     `json:"coreRefNum"`
	CoreRefNum2              string     `json:"coreRefNum2"`
	CRefNum                  string     `json:"cRefNum"`
	PRefNum                  string     `json:"prefNum"`
	UserId                   string     `json:"userId"`
	TransactionStartTime     string     `json:"transactionStartTime"`
	TransactionCompletedTime string     `json:"transactionCompletedTime"`
	TransactionErrorCode     string     `json:"transactionErrorCode"`
	TransactionErrorMsg      string     `json:"transactionErrorMsg"`
	TransactionReturn        int64      `json:"transactionReturn"`
	TransactionReturnMsg     string     `json:"transactionReturnMsg"`
	ClientCode               string     `json:"clientCode"`
	BranchInfo               BranchInfo `json:"branchInfo"`
	DetailsCode              string     `json:"detailsCode"`
	DetailMessage            string     `json:"detailMessage"`
}

type fundTransferExecutor struct {
}

func NewFundTransferExecutor() outbound.FundTransferExecutor {
	return &fundTransferExecutor{}
}

func (f *fundTransferExecutor) ExecutePayment(ctx context.Context, req *domain.FunTransferTransaction) (*outbound.ExecutePaymentResponse, error) {
	c := http.Client{
		Timeout:   time.Second * 10,
		Transport: http.DefaultTransport,
	}

	data := ExecutePaymentIn{
		TransactionInfo: TransactionInfo{
			CRefNum:    req.CRefNum,
			UserId:     "alo123",
			ClientCode: "OMNI",
			BranchInfo: BranchInfo{
				BranchCode: "internetbanking",
			},
		},
		TransferInfo: TransferInfo{
			PaymentType: "InternalPayment",
		},
		Payment: Payment{
			BookingDate:        "2023-11-01 10:55:16.525",
			DebitAccount:       req.FromAccount,
			CreditAccount:      req.ToAccount,
			CustomerId:         "7619297",
			OriginalCustomerId: "7619297",
			CRefNum:            req.CRefNum,
			Amount:             req.Amount,
			Currency:           "VND",
			Remarks:            "Test chuyen tien",
			EbUserId:           "alo123",
		},
		CardInfo: CardInfo{
			CardAccountNum: "",
			CardNumber:     "",
		},
	}

	dataByte, err := json.Marshal(ExecutePaymentInWrapper{
		ExecutePaymentIn: data,
	})

	if err != nil {
		return nil, err
	}

	res, err := c.Post(PaymentExecutorEndpoint, "application/json", bytes.NewBuffer(dataByte))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	resBodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode <= 399 && len(resBodyBytes) == 0 {
		return &outbound.ExecutePaymentResponse{
			Success: true,
		}, err
	}

	var resBody ExecutePaymentOutWrapper
	err = json.Unmarshal(resBodyBytes, &resBody)
	if err != nil {
		return nil, err
	}

	return &outbound.ExecutePaymentResponse{
		Success: resBody.ExecutePaymentOut.TransactionInfo.TransactionErrorCode == "SUCCESS",
		Detail:  resBody.ExecutePaymentOut.TransactionInfo.TransactionErrorMsg,
		TransNo: resBody.ExecutePaymentOut.TransactionInfo.CoreRefNum2,
	}, err

}
