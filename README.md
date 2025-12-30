# Turvo Integration Backend

A Go backend API that integrates with Turvo's TMS (Transportation Management System) to create and manage loads (shipments).

## Overview

This backend provides a REST API that acts as an intermediary between client applications and Turvo's API. It handles:
- Creating loads in Turvo's system
- Listing loads with filtering and pagination
- Mapping between our internal Load model and Turvo's Shipment model
- Status code translation between our API and Turvo

## Features

- **Create Loads**: Create new shipments in Turvo with comprehensive field support
- **List Loads**: Retrieve loads with filtering (status, customer, pickup dates) and pagination
- **Detailed Load Information**: Optionally fetch full details for each load using `includeDetails` flag
- **Field Mapping**: Automatic mapping between our API format and Turvo's format
- **Validation**: Comprehensive validation of required fields before creating loads
- **Status Mapping**: Translation between our status codes and Turvo's status codes

## API Endpoints

### Create Load

**POST** `/loads`

Creates a new load in Turvo's system.

**Request Body:** See `sample_create_load.json` for an example

**Response:** `201 Created`
```json
{
  "id": "1000306839",
  "createdAt": "2025-01-20T10:30:00Z"
}
```

**Required Fields:**
- `customer.externalTMSId` (string) - Must be a valid integer
- `pickup.readyTime` OR `pickup.apptTime` (datetime) - At least one required
- `pickup.city` + `pickup.state` OR `pickup.name` - At least one combination required
- `consignee.apptTime` (datetime) - Required
- `consignee.city` + `consignee.state` OR `consignee.name` - At least one combination required
- `carrier.externalTMSId` (string) - Required if carrier is provided, must be a valid integer

**Example Request:**
```json
{
  "status": "tendered",
  "customer": {
    "externalTMSId": "834099",
    "name": "Acme Corporation",
    "refNumber": "REF-12345"
  },
  "pickup": {
    "name": "Warehouse Distribution Center",
    "city": "Newark",
    "state": "NJ",
    "readyTime": "2025-01-27T08:00:00Z",
    "timezone": "America/New_York"
  },
  "consignee": {
    "name": "Retail Store Location",
    "city": "Philadelphia",
    "state": "PA",
    "apptTime": "2025-01-28T14:00:00Z",
    "timezone": "America/New_York"
  },
  "carrier": {
    "name": "ABC Transport Inc.",
    "externalTMSId": "834145"
  },
  "poNums": "PO-001, PO-002, PO-003",
  "totalWeight": 15000.5,
  "specifications": {
    "minTempFahrenheit": 32.0,
    "maxTempFahrenheit": 40.0
  },
  "routeMiles": 95.5
}
```

### List Loads

**GET** `/loads`

Retrieves a paginated list of loads from Turvo.

**Query Parameters:**
- `status` (string, optional) - Filter by status
- `customerId` (string, optional) - Filter by customer ID
- `pickupDateSearchFrom` (datetime, optional) - Filter loads picking up from this date (RFC3339)
- `pickupDateSearchTo` (datetime, optional) - Filter loads picking up to this date (RFC3339)
- `page` (integer, optional) - Page number (default: 1, min: 1)
- `limit` (integer, optional) - Results per page (default: 20, min: 1, max: 100)
- `includeDetails` (string, optional) - Set to "true" or "1" to fetch detailed information

**Response:** `200 OK`
```json
{
  "data": [
    {
      "externalTMSLoadID": "1000306839",
      "freightLoadID": "FREIGHT-67890",
      "status": "tendered",
      "customer": { ... },
      "pickup": { ... },
      "consignee": { ... },
      "carrier": { ... }
    }
  ],
  "pagination": {
    "total": 50,
    "pages": 3,
    "page": 1,
    "limit": 20
  }
}
```

**Example Request:**
```
GET /loads?status=tendered&page=1&limit=20&includeDetails=true
```

## Field Mappings

Only specific fields are mapped to Turvo's API. See `docs/FIELD_MAPPINGS.md` for complete details.

### Fields Mapped to Turvo:

**Customer:**
- `externalTMSId` → `customerOrder.customer.id`
- `name` → `customerOrder.customer.name`
- `refNumber` → `customerOrder.externalIds[]`

**Pickup:**
- `readyTime` or `apptTime` → `startDate`
- `timezone` → `startDate.timeZone`
- `city` + `state` or `name` → `lane.start`

**Consignee:**
- `apptTime` → `endDate`
- `timezone` → `endDate.timeZone`
- `city` + `state` or `name` → `lane.end`

**Carrier:**
- `name` → `carrierOrder.carrier.name`
- `externalTMSId` → `carrierOrder.carrier.id`

**Other:**
- `poNums` → `customerOrder.externalIds[]`
- `billTo` → `party[]`
- `totalWeight` → `equipment[].weight`
- `specifications.minTempFahrenheit/maxTempFahrenheit` → `equipment[].temp`
- `routeMiles` → `skipDistanceCalculation` flag
- `status` → `status.code`

### Fields NOT Mapped to Turvo:

These fields are accepted by our API but stored only in our system:
- `externalTMSLoadID`, `freightLoadID`
- Address details (addressLine1, addressLine2, zipcode, country)
- Contact information (contact, phone, email)
- Carrier details (MC number, DOT number, drivers, etc.)
- `rateData` fields
- Specification service flags (liftgate, inside, tarps, etc.)
- `inPalletCount`, `outPalletCount`, `numCommodities`, `billableWeight`, `operator`

## Status Codes

The API supports the following status values. See `STATUS_MAPPING.md` for complete mapping to Turvo status codes:

- `quote_active`, `tendered`, `covered`, `dispatched`
- `at_pickup`, `en_route`, `at_delivery`, `delivered`
- `ready_for_billing`, `processing`, `carrier_paid`, `customer_paid`
- `completed`, `canceled`, `quote_inactive`, `picked_up`
- `route_complete`, `tender_offered`, `tender_accepted`, `tender_rejected`
- `draft`, `shipment_ready`, `acquiring_location`, `customs_hold`
- `arrived`, `available`, `out_gated`, `in_gated`
- `arriving_to_port`, `berthing`, `unloading`, `ramped`
- `deramped`, `departed`, `held`, `out_for_delivery`
- `in_transshipment`, `on_hold`, `interline`

## Configuration

The backend requires Turvo API credentials configured via environment variables:

- `TURVO_CLIENT_ID` - Turvo API client ID
- `TURVO_CLIENT_SECRET` - Turvo API client secret
- `TURVO_BASE_URL` - Turvo API base URL (e.g., `https://my-sandbox-publicapi.turvo.com`)

## Running the Application

1. Install dependencies:
```bash
go mod download
```

2. Set environment variables:
```bash
export TURVO_CLIENT_ID="your-client-id"
export TURVO_CLIENT_SECRET="your-client-secret"
export TURVO_BASE_URL="https://my-sandbox-publicapi.turvo.com"
```

3. Run the server:
```bash
go run main.go
```

The API will be available at `http://localhost:8080` by default.

## Project Structure

```
.
├── docs/
│   ├── FIELD_MAPPINGS.md          # Complete field mapping documentation
│   └── FRONTEND_DEVELOPMENT_GUIDE.md  # Guide for frontend developers
├── examples/
│   └── create_load_complete.json  # Complete example with all fields
├── internal/
│   ├── handler/
│   │   └── load/
│   │       └── handler.go         # HTTP handlers
│   ├── models/
│   │   ├── load.go               # Load model definitions
│   │   └── turvo.go              # Turvo API models
│   ├── service/
│   │   └── load/
│   │       ├── service.go        # Business logic
│   │       ├── validation.go     # Validation rules
│   │       └── status_mapper.go   # Status code mapping
│   └── turvo/
│       └── client.go             # Turvo API client
├── sample_create_load.json       # Minimal example (only mapped fields)
├── STATUS_MAPPING.md            # Status code mappings
└── main.go                      # Application entry point
```

## Documentation

- **Field Mappings**: See `docs/FIELD_MAPPINGS.md` for detailed field mapping documentation
- **Frontend Guide**: See `docs/FRONTEND_DEVELOPMENT_GUIDE.md` for frontend development instructions
- **Status Mapping**: See `STATUS_MAPPING.md` for status code mappings
- **Examples**: See `sample_create_load.json` for a minimal example with only mapped fields

## Important Notes

1. **GlobalRoute**: We do not send globalRoute stops to Turvo in create requests (not supported in Turvo's POST API). However, we read globalRoute from Turvo responses when available.

2. **CustomId**: The `freightLoadID` field is not currently sent to Turvo during creation. Turvo generates its own `customId` which is returned in the response.

3. **Address Details**: While we accept full address details in our API, only city/state or name are sent to Turvo (for lane mapping). Full addresses are stored in our system and can be retrieved from Turvo responses when `includeDetails=true`.

4. **Specifications**: Service flags (liftgate, inside, tarps, etc.) are not sent to Turvo in create requests. They are only mapped back from Turvo responses when available in globalRoute services.

5. **Validation**: Required fields are validated before conversion. See `internal/service/load/validation.go` for validation rules.

6. **Concurrent Details Fetching**: When `includeDetails=true`, the API fetches detailed information for each load concurrently using goroutines for improved performance.

## Testing

See `TESTING.md` for testing instructions and examples.

## License

[Add your license here]
