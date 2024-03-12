package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lengocson131002/go-clean-core/logger"
	"github.com/lengocson131002/mcs-fund-transfer/domain"
	"github.com/lengocson131002/mcs-fund-transfer/usecase/data"
	"github.com/lengocson131002/mcs-fund-transfer/usecase/outbound"
)

type startFundTransferHandler struct {
	workflow outbound.FundTransferWorkflow
	fData    data.FundTransferData
	logger   logger.Logger
}

func NewStartFundTransferHandler(
	worflow outbound.FundTransferWorkflow,
	data data.FundTransferData,
	logger logger.Logger,
) domain.StartFundTransferHandler {
	return &startFundTransferHandler{
		workflow: worflow,
		logger:   logger,
		fData:    data,
	}
}

func (h *startFundTransferHandler) Handle(ctx context.Context, request *domain.StartFundTransferRequest) (*domain.StartFundTransferResponse, error) {
	trans := &domain.FunTransferTransaction{
		FromAccount: request.FromAccount,
		ToAccount:   request.ToAccount,
		Amount:      request.Amount,
		CRefNum:     uuid.New().String(),
		CreatedAt:   time.Now(),
	}

	workflowRes, err := h.workflow.StartFundTransferWorkflow(ctx, trans)
	if err != nil {
		h.logger.Errorf(ctx, "failed to start fund transfer workflow: %v", err)
		return nil, err
	}

	trans.WorflowId = workflowRes.WorkflowId

	h.logger.Infof(ctx, "started fund transfer workflow id: %v", workflowRes.WorkflowId)

	err = h.fData.SaveFundTransferTransaction(ctx, trans)
	if err != nil {
		h.logger.Errorf(ctx, "failed to save fund transfer transaction: %v", err)
		return nil, err
	}

	h.logger.Infof(ctx, "saved fund transfer transaction: %v", trans.CRefNum)

	return &domain.StartFundTransferResponse{
		CRefNum:    trans.CRefNum,
		WorkflowId: workflowRes.WorkflowId,
	}, nil
}
