package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lengocson131002/mcs-fund-transfer/domain"
	"github.com/lengocson131002/mcs-fund-transfer/presentation/http/handler"
)

type FundTransferController struct {
}

func NewFundTransferController() *FundTransferController {
	return &FundTransferController{}
}

func (c *FundTransferController) FundTransfer(ctx *fiber.Ctx) error {
	return handler.RequestHandler[*domain.StartFundTransferRequest, *domain.StartFundTransferResponse](ctx)
}

func (c *FundTransferController) VerifyOTP(ctx *fiber.Ctx) error {
	return handler.RequestHandler[*domain.VerifyFundTransferOTPRequest, *domain.VerifyFundTransferOTPResponse](ctx)
}

func (c *FundTransferController) QueryFundTransfer(ctx *fiber.Ctx) error {
	return handler.RequestHandler[*domain.QueryFundTransferRequest, *domain.QueryFundTransferResponse](ctx)
}
