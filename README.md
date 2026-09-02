# SSL Custom API

A Go REST API that provides customer aging data by querying SAP HANA databases(especially designed for Sarbottam Steels)

## Features

- Customer accounts receivable aging reports


## Setup

1. Copy `.env.example` to `.env` and configure:

```bash
HANA_PORT=30015
HANA_USER=your-user
HANA_PASSWORD=your-password
HANA_SCHEMA=your-schema
```

2. Run the API:

```bash
go run cmd/api/main.go
```

The API listens on port `8080`.

## API

### GET /api/aging

Returns customer aging data.

**Parameters:**
- `CardCode` (required) - Customer code
- `CompanyDB` (required) - Company database/schema name

**Example:**
```bash
curl "http://localhost:8080/api/aging?CardCode=C001&CompanyDB=SBODemoUS"
```

## Tech Stack

- Go
- SAP HANA (go-hdb driver)
- godotenv
