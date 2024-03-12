package outbound

import (
	"context"

	"github.com/lengocson131002/mcs-fund-transfer/domain"
	"github.com/lengocson131002/mcs-fund-transfer/pkg/workflow"
	"github.com/lengocson131002/mcs-fund-transfer/usecase/outbound"
	"go.temporal.io/sdk/client"
)

const (
	FundTransferOTPVerifiedSignalName = "FUND_STRANSFER_OTP_VERIFIED"
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
		ID:        workflow.NewFundTransferWorkflowId(),
		TaskQueue: workflow.FundTransferTaskQueue,
	}

	we, err := w.temClient.ExecuteWorkflow(ctx, options, workflow.FundTransferWorkflow, req)
	if err != nil {
		return nil, err
	}

	return &outbound.StartFundTransferWorkflowResponse{
		WorkflowId: we.GetID(),
	}, nil
}

// SignalFundTransferVerifiedOTP implements outbound.FundTransferWorkflow.
func (w *fundTransferWorkflow) SignalFundTransferVerifiedOTP(ctx context.Context, trans *domain.FunTransferTransaction) error {
	err := w.temClient.SignalWorkflow(ctx, trans.WorflowId, "", FundTransferOTPVerifiedSignalName, nil)

	if err != nil {
		return err
	}

	return nil
}
