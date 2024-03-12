package outbound

import (
	"context"
	"fmt"
	"time"

	"github.com/lengocson131002/mcs-fund-transfer/domain"
	"github.com/lengocson131002/mcs-fund-transfer/usecase/outbound"
	"go.temporal.io/sdk/client"
)

const (
	FundTransferOTPVerifiedSignalName = "VERIFY_OTP_CHANNEL"
	FundTransferResponseSignalName    = "CREATE_TRANSACTION_CHANNEL"
	FundTransferWorkflow              = "TransferWorkflow"
	FundTransferTaskQueue             = "TransferTaskQueue"
)

type fundTransferWorkflow struct {
	temClient client.Client
}

func NewFundTransferWorkflow(
	temClient client.Client,
) outbound.FundTransferWorkflow {
	return &fundTransferWorkflow{
		temClient: temClient,
	}
}

// StartFundTransferWorkflow implements outbound.FundTransferWorkflow.
func (w *fundTransferWorkflow) StartFundTransferWorkflow(ctx context.Context, req *domain.FunTransferTransaction) (*outbound.StartFundTransferWorkflowResponse, error) {
	options := client.StartWorkflowOptions{
		ID:        NewFundTransferWorkflowId(),
		TaskQueue: FundTransferTaskQueue,
	}

	we, err := w.temClient.ExecuteWorkflow(ctx, options, FundTransferWorkflow, req)
	if err != nil {
		return nil, err
	}

	return &outbound.StartFundTransferWorkflowResponse{
		WorkflowId: we.GetID(),
	}, nil
}

// SignalFundTransferVerifiedOTP implements outbound.FundTransferWorkflow.
func (w *fundTransferWorkflow) SignalFundTransferVerifiedOTP(ctx context.Context, trans *domain.FunTransferTransaction) error {
	err := w.temClient.SignalWorkflow(ctx, trans.WorflowId, "", FundTransferOTPVerifiedSignalName, trans)

	if err != nil {
		return err
	}

	return nil
}

// SignalFundTransferCompleted implements outbound.FundTransferWorkflow.
func (w *fundTransferWorkflow) SignalFundTransferCompleted(ctx context.Context, trans *domain.FunTransferTransaction) error {
	err := w.temClient.SignalWorkflow(ctx, trans.WorflowId, "", FundTransferResponseSignalName, trans)

	if err != nil {
		return err
	}

	return nil
}

func NewFundTransferWorkflowId() string {
	return FundTransferWorkflow + fmt.Sprintf("%d", time.Now().Unix())
}
