package data

import (
	"context"
	"fmt"
	"time"

	"github.com/lengocson131002/go-clean-core/es"
	"github.com/lengocson131002/go-clean-core/util"
	"github.com/lengocson131002/mcs-fund-transfer/domain"
	"github.com/lengocson131002/mcs-fund-transfer/usecase/data"
)

const (
	IndexFundTransfer    = "new-mcs-go.banktransfer.transf.1"
	IndexFundTransferOTP = "new-mcs-go.banktransfer.otp.1"
)

type fundTransferData struct {
	esClient es.ElasticSearchClient
}

func (d *fundTransferData) GetFundTransferTransactionByTransNo(ctx context.Context, transNo string) (*domain.FunTransferTransaction, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"transNo": transNo,
			},
		},
	}

	output, err := d.esClient.Search(
		ctx,
		fmt.Sprintf("%s*", IndexFundTransfer),
		es.WithSearchQuery(query),
	)

	if err != nil {
		return nil, err
	}

	if len(output) == 0 {
		return nil, domain.ErrorTransactionNotFound
	}

	var trans domain.FunTransferTransaction
	for _, v := range output {
		err = util.MapStruct(v, &trans, util.WithDecodeTimeFormat(time.RFC3339Nano))
		if err != nil {
			return nil, err
		}
		return &trans, nil
	}

	return nil, nil
}

func (d *fundTransferData) GetFundTransferOTP(ctx context.Context, cRefNum string, otp string) (*domain.FundTransferOTP, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"cRefNum": cRefNum,
						},
					},
					{
						"match": map[string]interface{}{
							"otp": otp,
						},
					},
				},
			},
		},
	}

	output, err := d.esClient.Search(
		ctx,
		fmt.Sprintf("%s*", IndexFundTransferOTP),
		es.WithSearchQuery(query),
		es.WithSearchSort([]string{"createdAt:desc"}),
	)

	if err != nil {
		return nil, err
	}

	if len(output) == 0 {
		return nil, domain.ErrorInvalidOTP
	}

	var trans domain.FundTransferOTP
	for _, v := range output {
		err = util.MapStruct(v, &trans, util.WithDecodeTimeFormat(time.RFC3339Nano))
		if err != nil {
			return nil, err
		}
		return &trans, nil
	}

	return nil, nil
}

func (d *fundTransferData) GetFundTransferTransaction(ctx context.Context, cRefNum string) (*domain.FunTransferTransaction, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"cRefNum": cRefNum,
			},
		},
	}

	output, err := d.esClient.Search(
		ctx,
		fmt.Sprintf("%s*", IndexFundTransfer),
		es.WithSearchQuery(query),
	)

	if err != nil {
		return nil, err
	}

	if len(output) == 0 {
		return nil, domain.ErrorTransactionNotFound
	}

	var trans domain.FunTransferTransaction
	for _, v := range output {
		err = util.MapStruct(v, &trans, util.WithDecodeTimeFormat(time.RFC3339Nano))
		if err != nil {
			return nil, err
		}
		return &trans, nil
	}
	return nil, nil
}

func (d *fundTransferData) SaveFundTransferOTP(ctx context.Context, otp *domain.FundTransferOTP) error {
	idx := fmt.Sprintf("%s-%s", IndexFundTransferOTP, time.Now().Format("2006.01.02"))
	err := d.esClient.Index(ctx, idx, otp)
	if err != nil {
		return fmt.Errorf("failed to create index for fund tranfer transaction otp. %w", err)
	}

	return nil
}

func (d *fundTransferData) SaveFundTransferTransaction(ctx context.Context, trans *domain.FunTransferTransaction) error {
	idx := fmt.Sprintf("%s-%s", IndexFundTransfer, time.Now().Format("2006.01.02"))
	err := d.esClient.Index(
		ctx,
		idx,
		trans,
		es.WithDocumentId(trans.CRefNum),
	)
	if err != nil {
		return fmt.Errorf("failed to create index for fund tranfer transaction. %w", err)
	}

	return nil
}

func NewFundTransferData(esClient es.ElasticSearchClient) data.FundTransferData {
	return &fundTransferData{
		esClient: esClient,
	}
}
