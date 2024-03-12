package cron

type KafkaAccountStatementWrapper struct {
	Table     string                 `json:"table"`
	OpType    string                 `json:"op_type"`
	OpTs      string                 `json:"op_ts"`
	CurrentTs string                 `json:"current_ts"`
	Before    *KafkaAccountStatement `json:"before"`
	After     *KafkaAccountStatement `json:"after"`
}

type KafkaAccountStatement struct {
	TransID      string  `json:"TRANS_ID"`
	DateTrans    string  `json:"DATE_TRANS"`
	Account      string  `json:"ACCOUNT"`
	TransNo      string  `json:"TRANS_NO"`
	Descriptions string  `json:"DESCRIPTIONS"`
	Amount       float64 `json:"AMOUNT"`
	Currency     string  `json:"CURRENCY"`
	InsertToSMS  int     `json:"INSERT_TO_SMS"`
	CIF          int     `json:"CIF"`
	IsDebit      int     `json:"IS_DEBIT"`
}
