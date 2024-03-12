package workflow

import (
	"fmt"
	"time"
)

func NewFundTransferWorkflowId() string {
	return "FUND_TRANSFER_" + fmt.Sprintf("%d", time.Now().Unix())
}
