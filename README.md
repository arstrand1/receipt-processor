# Receipt Processor Project
This is an API that meets the POST and GET methods defined in Fetch Reward's Receipt Processor Challenge (https://github.com/fetch-rewards/receipt-processor-challenge).
This program is coded in Go.

This code runs on "localhost:8080", so the specified methods use these paths:

Process Receipts
POST: http://localhost:8080/receipts/process

Get Points
GET: http://localhost:8080/receipts/{id}/points

