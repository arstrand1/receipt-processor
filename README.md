# Receipt Processor Project
This is an API that meets the POST and GET methods defined in Fetch Reward's Receipt Processor Challenge (https://github.com/fetch-rewards/receipt-processor-challenge).
It is coded in Go.

## Running the Program
The program runs on "localhost:8080". Once downloaded, you can run this program by navigating to its folder and using the Go command:
```go run . ```

### Endpoint: Process Receipts
* Path: http://localhost:8080/receipts/process
* Method: POST
* Payload: Receipt JSON
* Response: JSON containing an id for the receipt.

Example: ``` { "id": "12a95c10-0445-48ea-ae0c-535fdee7b075" } ```

### Endpoint: Get Points
* Path: http://localhost:8080/receipts/{id}/points
* Method: GET
* Response: A JSON object containing the number of points awarded.

Example: ``` { "points": 109 } ```
