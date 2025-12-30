# Field Mappings Documentation

This document describes the field mappings between the Drumkit Load API and Turvo's Shipment API.

## Overview

The integration maps fields from our internal Load model to Turvo's Shipment model for creation, and maps Turvo's responses back to our Load model for listing and retrieval.

## Create Load (Drumkit → Turvo)

### Required Fields

#### Customer
- **`customer.externalTMSId`** (string) → `customerOrder.customer.id` (integer)
  - Must be a valid integer string
  - Required field
  
- **`customer.name`** (string) → `customerOrder.customer.name`
  - Optional but recommended

- **`customer.refNumber`** (string) → `customerOrder.externalIds[]` (type: "Reference Number", key: "1401")
  - Maps to Turvo's external IDs array

#### Pickup
- **`pickup.readyTime`** (datetime) OR **`pickup.apptTime`** (datetime) → `startDate.date`
  - At least one is required
  - `readyTime` is preferred if both are provided
  - Format: RFC3339 (e.g., "2025-01-25T08:00:00Z")
  
- **`pickup.timezone`** (string) → `startDate.timeZone`
  - Defaults to "America/New_York" if not provided
  - Examples: "America/New_York", "America/Chicago", "UTC"

- **`pickup.city`** + **`pickup.state`** OR **`pickup.name`** → `lane.start`
  - At least one combination is required
  - If city and state are provided: `"City, State"` format
  - Otherwise: uses `pickup.name`

#### Consignee
- **`consignee.apptTime`** (datetime) → `endDate.date`
  - Required field
  - Format: RFC3339 (e.g., "2025-01-26T14:00:00Z")
  
- **`consignee.timezone`** (string) → `endDate.timeZone`
  - Defaults to pickup timezone if not provided

- **`consignee.city`** + **`consignee.state`** OR **`consignee.name`** → `lane.end`
  - At least one combination is required
  - If city and state are provided: `"City, State"` format
  - Otherwise: uses `consignee.name`

### Optional Fields

#### Status
- **`status`** (string) → `status.code.key` and `status.code.value`
  - Maps using status mapper (see STATUS_MAPPING.md)
  - Examples: "tendered", "covered", "dispatched", "delivered"
  - If not provided, Turvo will use default status

#### Freight Load ID
- **`freightLoadID`** (string) → `customId`
  - Note: Currently not mapped in create request (Turvo sets this in response)
  - Stored in our system for reference

#### PO Numbers
- **`poNums`** (string, comma-separated) → `customerOrder.externalIds[]` (type: "Purchase shipment #", key: "1400")
  - Example: "PO-001, PO-002, PO-003"
  - Each PO number becomes a separate external ID entry

#### BillTo
- **`billTo.externalTMSId`** (string) → `party[].account.id` (integer)
  - Must be a valid integer string if provided
  
- **`billTo.name`** (string) → `party[].account.name`
  - Maps to Turvo's Party array

#### Carrier
- **`carrier.externalTMSId`** (string) → `carrierOrder.carrier.id` (integer)
  - Required if carrier is provided
  - Must be a valid integer string
  
- **`carrier.name`** (string) → `carrierOrder.carrier.name`
  - Required if carrier is provided

#### Equipment
- **`totalWeight`** (float64) → `equipment[].weight`
  - Units: pounds (lb)
  - Key: "1520", Value: "lb"

- **`specifications.minTempFahrenheit`** and/or **`specifications.maxTempFahrenheit`** → `equipment[].temp`
  - If both provided: uses average temperature
  - If only one: uses that value
  - Units: Fahrenheit (°F)
  - Key: "1510", Value: "°F"

#### Route Miles
- **`routeMiles`** (float64) → `skipDistanceCalculation`
  - If provided: `skipDistanceCalculation = false`
  - If not provided: `skipDistanceCalculation = true`

### Fields NOT Mapped to Turvo (Stored in Our System Only)

These fields are accepted by our API but are not sent to Turvo:

- `externalTMSLoadID` - Our internal load identifier
- `freightLoadID` - Currently not sent (Turvo generates customId)
- All address fields (`addressLine1`, `addressLine2`, `city`, `state`, `zipcode`, `country`) in:
  - Customer
  - BillTo
  - Pickup
  - Consignee
- All contact fields (`contact`, `phone`, `email`) in:
  - Customer
  - BillTo
  - Pickup
  - Consignee
- `pickup.businessHours`, `pickup.refNumber`, `pickup.warehouseId`
- `consignee.businessHours`, `consignee.refNumber`, `consignee.mustDeliver`, `consignee.warehouseId`
- `pickup.apptNote` - Note: Not mapped (would require globalRoute)
- `consignee.apptNote` - Note: Not mapped (would require globalRoute)
- All carrier fields except `name` and `externalTMSId`:
  - `mcNumber`, `dotNumber`, `phone`, `dispatcher`, `sealNumber`, `scac`
  - Driver information (`firstDriverName`, `firstDriverPhone`, `secondDriverName`, `secondDriverPhone`)
  - `email`, `dispatchCity`, `dispatchState`
  - `externalTMSTruckId`, `externalTMSTrailerId`
  - All timestamp fields (`confirmationSentTime`, `dispatchedTime`, `pickupStart`, etc.)
  - `signedBy`
- `rateData` - All rate data fields
- `specifications` service flags:
  - `liftgatePickup`, `liftgateDelivery`
  - `insidePickup`, `insideDelivery`
  - `tarps`, `oversized`, `hazmat`, `straps`, `permits`, `escorts`, `seal`, `customBonded`, `labor`
- `inPalletCount`, `outPalletCount`, `numCommodities`
- `billableWeight`
- `operator`

## List/Get Load (Turvo → Drumkit)

### Basic Fields (from List Endpoint)

- **`id`** (integer) → `externalTMSLoadID` (string)
- **`customId`** (string) → `freightLoadID`
- **`status.code.key`** + **`status.code.value`** → `status` (string)
  - Mapped using status mapper

### Detailed Fields (from Get Shipment Endpoint)

When `includeDetails=true` is used, additional fields are mapped:

#### Dates
- **`startDate.date`** → `pickup.readyTime` and `pickup.apptTime`
- **`startDate.timeZone`** → `pickup.timezone`
- **`endDate.date`** → `consignee.apptTime`
- **`endDate.timeZone`** → `consignee.timezone`

#### Lane
- **`lane.start`** → `pickup.name`, `pickup.city`, `pickup.state`
  - Parses "City, State" format if applicable
- **`lane.end`** → `consignee.name`, `consignee.city`, `consignee.state`
  - Parses "City, State" format if applicable

#### Global Route (if available)
- **`globalRoute[].name`** → `pickup.name` or `consignee.name` (based on stopType)
- **`globalRoute[].address.line1`** → `pickup.addressLine1` or `consignee.addressLine1`
- **`globalRoute[].address.city`** → `pickup.city` or `consignee.city`
- **`globalRoute[].address.state`** → `pickup.state` or `consignee.state`
- **`globalRoute[].address.zip`** → `pickup.zipcode` or `consignee.zipcode`
- **`globalRoute[].timezone`** → `pickup.timezone` or `consignee.timezone`
- **`globalRoute[].notes`** → `pickup.apptNote` or `consignee.apptNote`
- **`globalRoute[].appointment.date`** → `pickup.apptTime` or `consignee.apptTime`
- **`globalRoute[].poNumbers[]`** → `poNums` (comma-separated string)
- **`globalRoute[].contact.name`** → `pickup.contact` or `consignee.contact`

#### Customer Order Route (if available)
- **`customerOrder[].route[].address`** → `pickup` or `consignee` address fields
- **`customerOrder[].route[].phone`** → `pickup.phone` or `consignee.phone`
- **`customerOrder[].route[].email`** → `pickup.email` or `consignee.email`
- **`customerOrder[].route[].contact`** → `pickup.contact` or `consignee.contact`

#### Customer Order
- **`customerOrder[].customer.id`** → `customer.externalTMSId` (string)
- **`customerOrder[].customer.name`** → `customer.name`
- **`customerOrder[].externalIds[]`** → `poNums` and `customer.refNumber`
  - Type "1400" or "Purchase shipment #" → `poNums`
  - Type "1401" or "Reference Number" → `customer.refNumber`

#### Carrier Order
- **`carrierOrder[].carrier.id`** → `carrier.externalTMSId` (string)
- **`carrierOrder[].carrier.name`** → `carrier.name`
- **`carrierOrder[].drivers[0].context.name`** → `carrier.firstDriverName`
- **`carrierOrder[].drivers[0].phone.number`** → `carrier.firstDriverPhone`
- **`carrierOrder[].drivers[0].email.email`** → `carrier.email`
- **`carrierOrder[].drivers[1].context.name`** → `carrier.secondDriverName`
- **`carrierOrder[].drivers[1].phone.number`** → `carrier.secondDriverPhone`
- **`carrierOrder[].externalIds[]`** → `carrier.externalTMSTruckId`, `carrier.externalTMSTrailerId`, `carrier.sealNumber`
  - Based on external ID type

#### Party
- **`party[].account.id`** → `billTo.externalTMSId` (string)
- **`party[].account.name`** → `billTo.name`

#### Equipment
- **`equipment[].weight`** → `totalWeight`
- **`equipment[].temp`** → `specifications.minTempFahrenheit` and `specifications.maxTempFahrenheit`
  - If only one value, both min and max are set to the same value

#### Specifications (from Services)
- **`globalRoute[].services[]`** → `specifications` flags
  - "Liftgate" → `liftgatePickup` or `liftgateDelivery` (based on stop)
  - "Inside Pickup" → `insidePickup`
  - "Inside Delivery" → `insideDelivery`
  - "Tarps" → `tarps`
  - "Hazmat" → `hazmat`
  - "Straps" → `straps`
  - "Seal" → `seal`

#### Distance
- **`customerOrder[].totalMiles`** → `routeMiles`

## Status Mapping

See `STATUS_MAPPING.md` for complete status code mappings between Drumkit API and Turvo.

## Notes

1. **GlobalRoute**: We do not send globalRoute stops to Turvo in create requests (not supported in Turvo's POST API). However, we read globalRoute from Turvo responses when available.

2. **CustomId**: The `freightLoadID` field is not currently sent to Turvo during creation. Turvo generates its own `customId` which is returned in the response.

3. **Address Details**: While we accept full address details in our API, only city/state or name are sent to Turvo (for lane mapping). Full addresses are stored in our system and can be retrieved from Turvo responses when `includeDetails=true`.

4. **Specifications**: Service flags (liftgate, inside, tarps, etc.) are not sent to Turvo in create requests. They are only mapped back from Turvo responses when available in globalRoute services.

5. **Validation**: Required fields are validated before conversion. See `validation.go` for validation rules.

