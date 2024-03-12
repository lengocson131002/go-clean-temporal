### 2. START MCS
#### START ACCOUNT MCS
```
cd ./mcs-account
go run ./cmd/main.go
```

#### START FUND TRANSFER MCS
``` 
cd ./mcs-fund-transfer
go run ./cmd/main.go
```

### 2. FUND TRANSFER FLOW

#### STEP 1: Start fund transfer workflow
``` json
MCS: MCS-FUND-TRANSFER 
PROTOCOL: HTTP REST API - POST
URL: http://localhost:8088/api/v1/fund-transfer/start
PAYLOAD: application/json
Request:
{
    "fromAccount": "0347885267",
    "toAccount": "12345",
    "amount": 10000
}

Response:
{
    "result": {
        "status": 200,
        "code": "0",
        "message": "Success",
        "details": null
    },
    "data": {
        "cRefNum": "3e2ac7b7-3b25-42ce-b65c-5dd0f6db78fd",
        "workflowId": "TransferWorkflow1710212345"
    }
}

Transaction is created and saved in Elastic search
Query transaction: {{eshost}}/new-mcs-go.banktransfer.otp.1*/_search
Body: 
{
  "query": {
    "match": {"cRefNum": "a7510947-9bb4-4425-800f-911fb136dc4a"}
  }
}

```

#### STEP 2: Check balance debit account
``` json
MCS: MCS-ACCOUNT
PROTOCOL: KAFKABROKER
TOPIC: Request topic: OCB.REQUEST.CHECK_BALANCE, Response topic: OCB.REPLY.CHECK_BALANCE 
PAYLOAD: application/json
Request:
{"account":"0347885267"}

Response:
{
	"result": {
		"status": 200,
		"code": "0",
		"message": "Success",
		"details": null
	},
	"data": {
		"balance": 99905527408,
		"currency": "VND"
	}
}
```

#### STEP 3: Generate fund transfer OTP
``` json
MCS: MCS-FUND-TRANSFER 
PROTOCOL: KAFKABROKER
TOPIC: Request topic: OCB.REQUEST.GENERATE_OTP, Response topic: OCB.REPLY.GENERATE_OTP 
PAYLOAD: application/json
Request:
{"cRefNum": "3e2ac7b7-3b25-42ce-b65c-5dd0f6db78fd"}

Response:
{
	"result": {
		"status": 200,
		"code": "0",
		"message": "Success",
		"details": null
	},
	"data": null
}

OTP is created and saved in Elastic search
Query OTP: {{eshost}}/new-mcs-go.banktransfer.otp.1*/_search
Body: 
{
  "query": {
    "match": {"cRefNum": "a7510947-9bb4-4425-800f-911fb136dc4a"}
  }
}
```

#### STEP 4: User verify OTP
``` json
MCS: MCS-FUND-TRANSFER 
PROTOCOL: HTTP REST API - POST
URL: http://localhost:8088/api/v1/fund-transfer/verify-otp
PAYLOAD: application/json
Request:
{
    "cRefNum": "a7510947-9bb4-4425-800f-911fb136dc4a",
    "otp": "413814"
}

Response:
{
    "result": {
        "status": 200,
        "code": "0",
        "message": "Success",
        "details": null
    },
    "data": {
        "success": true
    }
}

==> Signal OTP VERIFIED Channel: VERIFY_OTP_CHANNEL
		Signal Payload: {
                    "worflowId": "FUND_TRANSFER_1710212304",
                    "fromAccount": "0347885267",
                    "toAccount": "12345",
                    "amount": 10000,
                    "cRefNum": "7ad2f7d6-cde7-431a-8359-f464aa314779",
                    "createdAt": "2024-03-12T09:58:24.692483962+07:00",
                    "transferAt": null,
                    "transNo": "abcd",
                    "status": 0
                }

```

#### STEP 5: Execute fund transaction
``` json
MCS: MCS-FUND-TRANSFER 
PROTOCOL: KAFKABROKER
TOPIC: Request: OCB.REQUEST.FUND_TRANSFER, Response: OCB.REPLY.FUND_TRANSFER
PAYLOAD: application/json
Request:
{"cRefNum": "a7510947-9bb4-4425-800f-911fb136dc4a"}
Response:
{
	"result": {
		"status": 200,
		"code": "0",
		"message": "Success",
		"details": null
	},
	"data": {}
}

=> Transaction is processing, and run async
```

#### STEP 5: Signal when transaction completed
``` json

MCS: MCS-FUND-TRANSFER 
PROTOCOL: TEMPORAL WORKFLOW SIGNAL
CHANEL: CREATE_TRANSACTION_CHANNEL
PAYLOAD: application/json
{
    "worflowId": "FUND_TRANSFER_1710212304",
    "fromAccount": "0347885267",
    "toAccount": "12345",
    "amount": 10000,
    "cRefNum": "7ad2f7d6-cde7-431a-8359-f464aa314779",
    "createdAt": "2024-03-12T09:58:24.692483962+07:00",
    "transferAt": null,
    "transNo": "abcd"
    "status": 0
}

```

==============================
#### Query transfer:
``` json

MCS: MCS-FUND-TRANSFER 
PROTOCOL: HTTP REST API - GET
URL: http://localhost:8088/api/v1/fund-transfer/query
PAYLOAD: application/json
Request:
{
    "cRefNum": "0ba86baf-83ab-4c5f-9335-d90dc2c26643"
}

Response:
{
    "result": {
        "status": 200,
        "code": "0",
        "message": "Success",
        "details": null
    },
    "data": {
        "worflowId": "TransferWorkflow1710216905",
        "fromAccount": "0347885267",
        "toAccount": "12345",
        "amount": 10000,
        "cRefNum": "0ba86baf-83ab-4c5f-9335-d90dc2c26643",
        "createdAt": "0001-01-01T00:00:00Z",
        "transferAt": "2024-03-12T11:21:51.961139538+07:00",
        "status": 3,
        "transNo": ""
    }
}
```