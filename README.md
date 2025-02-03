# Introduction

The documented project below represents a first for myself. I have never written in Go before now. Given the opportunity, I figured why not. There are certainly some areas that could use improvement, project is structured in a way similar to what I am used to in .NET, not split out quite as much for the sake of simplicity.

The following was built in Go using Gin. Initially I was torn between Gin and Mux, but Gin felt a bit more intuitive for someone who has never programmed in Go before, so I chose and stuck with it. That being said, the only packages necessary for the project were Gin and Google's uuid.

# Run Locally

The following was built using the Go 1.23.4 [Install here](https://go.dev/doc/install)

```
$ git clone https://github.com/cbjorne/fetchapi.git
$ go run fetchapi/cmd/main.go
```

# Endpoints

## Process Receipt

Processes receipt data JSON and returns an object with a UUID key associated with the point value rewarded

```http
POST /receipts/process
```
Request
```javascript
{
  "retailer": string,
  "purchaseDate": string,
  "purchaseTime": string,
  "items": [
    {
      "shortDescription": string,
      "price": string
    }
  ],
  "total": string
}
```
Response
```javascript
{
    "id": uuid
}
```

### Status Codes

| Status Code | Description |
| :--- | :--- |
| 200 | `OK` |
| 404 | `BAD REQUEST` |

## Get Points

Gets the point value associated with the ID passed in query params

```http
GET /receipts/{id}/points
```

Request
| Parameter | Type | Description |
| :--- | :--- | :--- |
| `receipt_id` | `uuid` | **Required**. Receipt Id corresponding to points awarded |

Response
```javascript
{
    "points": int
}
```

### Status Codes

| Status Code | Description |
| :--- | :--- |
| 200 | `OK` |
| 404 | `PAGE NOT FOUND` |




