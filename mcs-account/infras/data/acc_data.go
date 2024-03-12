package data

import (
	"context"
	"fmt"

	"github.com/lengocson131002/go-clean-core/es"
	"github.com/lengocson131002/go-clean-core/util"
	"github.com/lengocson131002/mcs-account/domain"
	"github.com/lengocson131002/mcs-account/usecase/data"
)

type accData struct {
	esClient es.ElasticSearchClient
}

type esAccountBalanceModel struct {
	OcbAccountNumber           string
	OcbBranchCode              string
	OcbCustomerNumber          string
	CustomerNumberJointProfile string
	Currency                   string
	AccountOpeningDate         string
	LastAccountStatusCode      string
	LastAccountStatusDate      string
	Category                   string
	AccountTitle               string
	ShortTitle                 string
	OpenActualBal              int64
	OnlineActualBal            int64
	WorkingBalance             int64
	AccountOfficer             string
	ConditionGroup             string
	CurrNo                     string
	Op_ts                      string
	Current_ts                 string
}

const (
	IndexAccountBalance = "t24v2.fbnk_account.transf.1"
)

// GetBalance implements data.AccountData.
func (a *accData) GetBalance(ctx context.Context, accNumber string) (*data.AccountBalanceResponse, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"ocbAccountNumber": accNumber,
			},
		},
	}

	output, err := a.esClient.Search(
		ctx,
		fmt.Sprintf("%s*", IndexAccountBalance),
		es.WithSearchQuery(query),
		es.WithSearchSort([]string{"op_ts:desc"}),
	)

	if err != nil {
		return nil, err
	}

	if len(output) == 0 {
		return nil, domain.ErrorAccountNotFound
	}

	var balRes esAccountBalanceModel
	err = util.MapStruct(output[0], &balRes)
	if err != nil {
		return nil, err
	}

	return &data.AccountBalanceResponse{
		Currency:        balRes.Currency,
		OpenActualBal:   balRes.OpenActualBal,
		WorkingBalance:  balRes.WorkingBalance,
		OnlineActualBal: balRes.OnlineActualBal,
	}, nil
}

func NewAccountData(esClient es.ElasticSearchClient) data.AccountData {
	return &accData{
		esClient: esClient,
	}
}
