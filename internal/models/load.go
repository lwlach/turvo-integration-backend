package models

import "time"

// Load represents the Drumkit load format (input/output of our REST API)
type Load struct {
	ExternalTMSLoadID string          `json:"externalTMSLoadID,omitempty"`
	FreightLoadID     string          `json:"freightLoadID,omitempty"`
	Status            string          `json:"status,omitempty"`
	Customer          *Customer       `json:"customer,omitempty"`
	BillTo            *BillTo         `json:"billTo,omitempty"`
	Pickup            *Pickup         `json:"pickup,omitempty"`
	Consignee         *Consignee      `json:"consignee,omitempty"`
	Carrier           *Carrier        `json:"carrier,omitempty"`
	RateData          *RateData       `json:"rateData,omitempty"`
	Specifications    *Specifications `json:"specifications,omitempty"`
	InPalletCount     *int            `json:"inPalletCount,omitempty"`
	OutPalletCount    *int            `json:"outPalletCount,omitempty"`
	NumCommodities    *int            `json:"numCommodities,omitempty"`
	TotalWeight       *float64        `json:"totalWeight,omitempty"`
	BillableWeight    *float64        `json:"billableWeight,omitempty"`
	PoNums            string          `json:"poNums,omitempty"`
	Operator          string          `json:"operator,omitempty"`
	RouteMiles        *float64        `json:"routeMiles,omitempty"`
}

// Customer represents the customer object in Drumkit load format
type Customer struct {
	ExternalTMSId string `json:"externalTMSId"`
	Name          string `json:"name,omitempty"`
	AddressLine1  string `json:"addressLine1,omitempty"`
	AddressLine2  string `json:"addressLine2,omitempty"`
	City          string `json:"city,omitempty"`
	State         string `json:"state,omitempty"`
	Zipcode       string `json:"zipcode,omitempty"`
	Country       string `json:"country,omitempty"`
	Contact       string `json:"contact,omitempty"`
	Phone         string `json:"phone,omitempty"`
	Email         string `json:"email,omitempty"`
	RefNumber     string `json:"refNumber,omitempty"`
}

// BillTo represents the billTo object in Drumkit load format
type BillTo struct {
	ExternalTMSId string `json:"externalTMSId,omitempty"`
	Name          string `json:"name,omitempty"`
	AddressLine1  string `json:"addressLine1,omitempty"`
	AddressLine2  string `json:"addressLine2,omitempty"`
	City          string `json:"city,omitempty"`
	State         string `json:"state,omitempty"`
	Zipcode       string `json:"zipcode,omitempty"`
	Country       string `json:"country,omitempty"`
	Contact       string `json:"contact,omitempty"`
	Phone         string `json:"phone,omitempty"`
	Email         string `json:"email,omitempty"`
}

// Pickup represents the pickup object in Drumkit load format
type Pickup struct {
	ExternalTMSId string     `json:"externalTMSId,omitempty"`
	Name          string     `json:"name,omitempty"`
	AddressLine1  string     `json:"addressLine1,omitempty"`
	AddressLine2  string     `json:"addressLine2,omitempty"`
	City          string     `json:"city,omitempty"`
	State         string     `json:"state,omitempty"`
	Zipcode       string     `json:"zipcode,omitempty"`
	Country       string     `json:"country,omitempty"`
	Contact       string     `json:"contact,omitempty"`
	Phone         string     `json:"phone,omitempty"`
	Email         string     `json:"email,omitempty"`
	BusinessHours string     `json:"businessHours,omitempty"`
	RefNumber     string     `json:"refNumber,omitempty"`
	ReadyTime     *time.Time `json:"readyTime,omitempty"`
	ApptTime      *time.Time `json:"apptTime,omitempty"`
	ApptNote      string     `json:"apptNote,omitempty"`
	Timezone      string     `json:"timezone,omitempty"`
	WarehouseId   string     `json:"warehouseId,omitempty"`
}

// Consignee represents the consignee object in Drumkit load format
type Consignee struct {
	ExternalTMSId string     `json:"externalTMSId,omitempty"`
	Name          string     `json:"name,omitempty"`
	AddressLine1  string     `json:"addressLine1,omitempty"`
	AddressLine2  string     `json:"addressLine2,omitempty"`
	City          string     `json:"city,omitempty"`
	State         string     `json:"state,omitempty"`
	Zipcode       string     `json:"zipcode,omitempty"`
	Country       string     `json:"country,omitempty"`
	Contact       string     `json:"contact,omitempty"`
	Phone         string     `json:"phone,omitempty"`
	Email         string     `json:"email,omitempty"`
	BusinessHours string     `json:"businessHours,omitempty"`
	RefNumber     string     `json:"refNumber,omitempty"`
	MustDeliver   string     `json:"mustDeliver,omitempty"`
	ApptTime      *time.Time `json:"apptTime,omitempty"`
	ApptNote      string     `json:"apptNote,omitempty"`
	Timezone      string     `json:"timezone,omitempty"`
	WarehouseId   string     `json:"warehouseId,omitempty"`
}

// Carrier represents the carrier object in Drumkit load format
type Carrier struct {
	MCNumber                 string     `json:"mcNumber,omitempty"`
	DOTNumber                string     `json:"dotNumber,omitempty"`
	Name                     string     `json:"name,omitempty"`
	Phone                    string     `json:"phone,omitempty"`
	Dispatcher               string     `json:"dispatcher,omitempty"`
	SealNumber               string     `json:"sealNumber,omitempty"`
	Scac                     string     `json:"scac,omitempty"`
	FirstDriverName          string     `json:"firstDriverName,omitempty"`
	FirstDriverPhone         string     `json:"firstDriverPhone,omitempty"`
	SecondDriverName         string     `json:"secondDriverName,omitempty"`
	SecondDriverPhone        string     `json:"secondDriverPhone,omitempty"`
	Email                    string     `json:"email,omitempty"`
	DispatchCity             string     `json:"dispatchCity,omitempty"`
	DispatchState            string     `json:"dispatchState,omitempty"`
	ExternalTMSTruckId       string     `json:"externalTMSTruckId,omitempty"`
	ExternalTMSTrailerId     string     `json:"externalTMSTrailerId,omitempty"`
	ConfirmationSentTime     *time.Time `json:"confirmationSentTime,omitempty"`
	ConfirmationReceivedTime *time.Time `json:"confirmationReceivedTime,omitempty"`
	DispatchedTime           *time.Time `json:"dispatchedTime,omitempty"`
	ExpectedPickupTime       *time.Time `json:"expectedPickupTime,omitempty"`
	PickupStart              *time.Time `json:"pickupStart,omitempty"`
	PickupEnd                *time.Time `json:"pickupEnd,omitempty"`
	ExpectedDeliveryTime     *time.Time `json:"expectedDeliveryTime,omitempty"`
	DeliveryStart            *time.Time `json:"deliveryStart,omitempty"`
	DeliveryEnd              *time.Time `json:"deliveryEnd,omitempty"`
	SignedBy                 string     `json:"signedBy,omitempty"`
	ExternalTMSId            string     `json:"externalTMSId,omitempty"`
}

// RateData represents the rateData object in Drumkit load format
type RateData struct {
	CustomerRateType  string   `json:"customerRateType,omitempty"`
	CustomerNumHours  *float64 `json:"customerNumHours,omitempty"`
	CustomerLhRateUsd *float64 `json:"customerLhRateUsd,omitempty"`
	FscPercent        *float64 `json:"fscPercent,omitempty"`
	FscPerMile        *float64 `json:"fscPerMile,omitempty"`
	CarrierRateType   string   `json:"carrierRateType,omitempty"`
	CarrierNumHours   *float64 `json:"carrierNumHours,omitempty"`
	CarrierLhRateUsd  *float64 `json:"carrierLhRateUsd,omitempty"`
	CarrierMaxRate    *float64 `json:"carrierMaxRate,omitempty"`
	NetProfitUsd      *float64 `json:"netProfitUsd,omitempty"`
	ProfitPercent     *float64 `json:"profitPercent,omitempty"`
}

// Specifications represents the specifications object in Drumkit load format
type Specifications struct {
	MinTempFahrenheit *float64 `json:"minTempFahrenheit,omitempty"`
	MaxTempFahrenheit *float64 `json:"maxTempFahrenheit,omitempty"`
	LiftgatePickup    *bool    `json:"liftgatePickup,omitempty"`
	LiftgateDelivery  *bool    `json:"liftgateDelivery,omitempty"`
	InsidePickup      *bool    `json:"insidePickup,omitempty"`
	InsideDelivery    *bool    `json:"insideDelivery,omitempty"`
	Tarps             *bool    `json:"tarps,omitempty"`
	Oversized         *bool    `json:"oversized,omitempty"`
	Hazmat            *bool    `json:"hazmat,omitempty"`
	Straps            *bool    `json:"straps,omitempty"`
	Permits           *bool    `json:"permits,omitempty"`
	Escorts           *bool    `json:"escorts,omitempty"`
	Seal              *bool    `json:"seal,omitempty"`
	CustomBonded      *bool    `json:"customBonded,omitempty"`
	Labor             *bool    `json:"labor,omitempty"`
}

// LoadFilters represents filter parameters for listing loads
type LoadFilters struct {
	Status         string
	CustomerID     string
	PickupDateFrom *time.Time
	PickupDateTo   *time.Time
	Page           int
	Limit          int
	IncludeDetails bool // If true, fetch detailed shipment information for each load
}

// LoadListResponse represents the paginated response for listing loads
type LoadListResponse struct {
	Data       []Load     `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// Pagination represents pagination metadata
type Pagination struct {
	Total int `json:"total"`
	Pages int `json:"pages"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// LoadCreateResponse represents the response from creating a load
type LoadCreateResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}
