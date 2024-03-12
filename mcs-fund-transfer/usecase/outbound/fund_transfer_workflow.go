package outbound

import (
	"context"

	"github.com/lengocson131002/mcs-fund-transfer/domain"
)

type StartFundTransferWorkflowRequest struct {
	FromAccount string `json:"fromAccount"`
	ToAccount   string `json:"toAccount"`
	Amount      int64  `json:"amount"`
	CRefNum     string `json:"cRefNum"`
}

type StartFundTransferWorkflowResponse struct {
	WorkflowId string `json:"workflowId"`
}

type FundTransferWorkflow interface {
	StartFundTransferWorkflow(context.Context, *domain.FunTransferTransaction) (*StartFundTransferWorkflowResponse, error)
	SignalFundTransferVerifiedOTP(context.Context, *domain.FunTransferTransaction) error
	SignalFundTransferCompleted(context.Context, *domain.FunTransferTransaction) error
}
