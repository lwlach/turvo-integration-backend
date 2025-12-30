package models

// TurvoShipmentsListResponse represents the response from Turvo's GET shipments/list endpoint
type TurvoShipmentsListResponse struct {
	Status  string                `json:"Status"`
	Details TurvoShipmentsDetails `json:"details"`
}

type TurvoShipmentsDetails struct {
	Pagination TurvoPagination `json:"pagination"`
	Shipments  []TurvoShipment `json:"shipments"`
}

type TurvoPagination struct {
	Start              int     `json:"start"`
	PageSize           int     `json:"pageSize"`
	TotalRecordsInPage int     `json:"totalRecordsInPage"`
	MoreAvailable      bool    `json:"moreAvailable"`
	LastObjectKey      *string `json:"lastObjectKey,omitempty"` // Can be null
}

// TurvoShipmentFilters represents filter parameters for Turvo's list shipments API
type TurvoShipmentFilters struct {
	Status        string // Status code (e.g., "2101")
	CustomerID    string // Customer ID
	PickupDateGte string // Pickup date greater than or equal (RFC3339 format)
	PickupDateLte string // Pickup date less than or equal (RFC3339 format)
	Start         int    // Start index for pagination
	PageSize      int    // Page size for pagination
}

// TurvoShipment represents Turvo's shipment model from the list endpoint
type TurvoShipment struct {
	ID               int                  `json:"id"`
	CustomID         string               `json:"customId,omitempty"`
	NetCustomerCosts float64              `json:"netCustomerCosts,omitempty"`
	NetCarrierCosts  float64              `json:"netCarrierCosts,omitempty"`
	NetRevenue       float64              `json:"netRevenue,omitempty"`
	Status           TurvoShipmentStatus  `json:"status,omitempty"`
	CustomerOrder    []TurvoCustomerOrder `json:"customerOrder,omitempty"`
	CarrierOrder     []TurvoCarrierOrder  `json:"carrierOrder,omitempty"`
	Party            []TurvoPartyOrder    `json:"party,omitempty"`
	Created          string               `json:"created,omitempty"`
	Updated          string               `json:"updated,omitempty"`
	LastUpdatedOn    string               `json:"lastUpdatedOn,omitempty"`
	CreatedDate      string               `json:"createdDate,omitempty"`
}

type TurvoShipmentStatus struct {
	Code TurvoStatusCode `json:"code,omitempty"`
}

type TurvoStatusCode struct {
	Key   string `json:"key"` // Changed from int to string (e.g., "2117", "2101")
	Value string `json:"value"`
}

type TurvoCustomerOrder struct {
	ID       int          `json:"id"`
	Customer TurvoAccount `json:"customer,omitempty"`
	Deleted  bool         `json:"deleted"`
}

type TurvoCarrierOrder struct {
	ID      int          `json:"id"`
	Carrier TurvoAccount `json:"carrier,omitempty"`
	Deleted bool         `json:"deleted"`
}

type TurvoPartyOrder struct {
	ID      int          `json:"id"`
	Account TurvoAccount `json:"account,omitempty"`
	Deleted bool         `json:"deleted"`
}

type TurvoAccount struct {
	ID    int                 `json:"id"`
	Name  string              `json:"name,omitempty"`
	Owner *TurvoParentAccount `json:"owner,omitempty"`
}

type TurvoParentAccount struct {
	Name string `json:"name,omitempty"`
	ID   int    `json:"id"`
}

// Legacy models for backward compatibility and create operations
// These may be used for creating shipments or detailed shipment views

type TurvoLocation struct {
	Name      string  `json:"name,omitempty"`
	Address   string  `json:"address,omitempty"`
	City      string  `json:"city,omitempty"`
	State     string  `json:"state,omitempty"`
	ZipCode   string  `json:"zipCode,omitempty"`
	Country   string  `json:"country,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

type TurvoCarrier struct {
	Name      string `json:"name,omitempty"`
	MCNumber  string `json:"mcNumber,omitempty"`
	DOTNumber string `json:"dotNumber,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Email     string `json:"email,omitempty"`
}

type TurvoParty struct {
	Name    string `json:"name,omitempty"`
	Contact string `json:"contact,omitempty"`
	Phone   string `json:"phone,omitempty"`
	Email   string `json:"email,omitempty"`
}

// TurvoShipmentCreate represents the shipment model for creating shipments in Turvo API
type TurvoShipmentCreate struct {
	LtlShipment             bool                       `json:"ltlShipment"`
	StartDate               TurvoDateWithTimezone      `json:"startDate"`
	EndDate                 TurvoDateWithTimezone      `json:"endDate"`
	Status                  *TurvoCreateStatus         `json:"status,omitempty"`
	Groups                  []TurvoGroup               `json:"groups,omitempty"`
	Contributors            []TurvoContributor         `json:"contributors,omitempty"`
	Equipment               []TurvoEquipment           `json:"equipment,omitempty"`
	Lane                    TurvoLane                  `json:"lane"`
	GlobalRoute             []TurvoGlobalRouteStop     `json:"globalRoute,omitempty"`
	SkipDistanceCalculation bool                       `json:"skipDistanceCalculation,omitempty"`
	ModeInfo                []TurvoModeInfo            `json:"modeInfo,omitempty"`
	CustomerOrder           []TurvoCreateCustomerOrder `json:"customerOrder"`
	CarrierOrder            []TurvoCreateCarrierOrder  `json:"carrierOrder,omitempty"`
	Party                   []TurvoCreateParty         `json:"party,omitempty"`
	UseRoutingGuide         bool                       `json:"use_routing_guide,omitempty"`
}

type TurvoDateWithTimezone struct {
	Date     string `json:"date"`
	TimeZone string `json:"timeZone"`
}

type TurvoCreateStatus struct {
	Code        TurvoStatusCode `json:"code,omitempty"`
	Notes       string          `json:"notes,omitempty"`
	Description string          `json:"description,omitempty"`
	Category    string          `json:"category,omitempty"`
}

type TurvoGroup struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Operation int    `json:"_operation,omitempty"`
}

type TurvoContributor struct {
	Type            TurvoKeyValue      `json:"type,omitempty"`
	ContributorUser TurvoUserReference `json:"contributorUser,omitempty"`
	Operation       int                `json:"_operation,omitempty"`
}

type TurvoUserReference struct {
	ID int `json:"id"`
}

type TurvoEquipment struct {
	Operation      int           `json:"_operation,omitempty"`
	Type           TurvoKeyValue `json:"type,omitempty"`
	Size           TurvoKeyValue `json:"size,omitempty"`
	Weight         float64       `json:"weight,omitempty"`
	WeightUnits    TurvoKeyValue `json:"weightUnits,omitempty"`
	Temp           float64       `json:"temp,omitempty"`
	TempUnits      TurvoKeyValue `json:"tempUnits,omitempty"`
	ShipmentLength float64       `json:"shipmentLength,omitempty"`
}

type TurvoLane struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type TurvoGlobalRouteStop struct {
	GlobalShipLocationSourceId string                       `json:"globalShipLocationSourceId,omitempty"`
	Name                       string                       `json:"name,omitempty"`
	SchedulingType             TurvoKeyValue                `json:"schedulingType,omitempty"`
	StopType                   TurvoKeyValue                `json:"stopType,omitempty"`
	Timezone                   string                       `json:"timezone,omitempty"`
	Location                   TurvoLocationReference       `json:"location,omitempty"`
	SegmentSequence            int                          `json:"segmentSequence,omitempty"`
	LayoverTime                *TurvoTimeWithUnits          `json:"layoverTime,omitempty"`
	Sequence                   int                          `json:"sequence,omitempty"`
	State                      string                       `json:"state,omitempty"`
	Appointment                *TurvoAppointment            `json:"appointment,omitempty"`
	AppointmentConfirmation    bool                         `json:"appointmentConfirmation,omitempty"`
	PlannedAppointmentDate     *TurvoPlannedAppointmentDate `json:"plannedAppointmentDate,omitempty"`
	Services                   []TurvoKeyValue              `json:"services,omitempty"`
	PoNumbers                  []string                     `json:"poNumbers,omitempty"`
	Notes                      string                       `json:"notes,omitempty"`
	CustomerOrder              []TurvoRouteCustomerOrder    `json:"customerOrder,omitempty"`
	CarrierOrder               []TurvoRouteCarrierOrder     `json:"carrierOrder,omitempty"`
	Transportation             *TurvoTransportation         `json:"transportation,omitempty"`
	FragmentDistance           *TurvoDistanceWithUnits      `json:"fragmentDistance,omitempty"`
	Distance                   *TurvoDistanceWithUnits      `json:"distance,omitempty"`
	StopLevelFragmentDistance  float64                      `json:"stop_level_fragment_distance,omitempty"`
}

type TurvoLocationReference struct {
	ID int `json:"id"`
}

type TurvoTimeWithUnits struct {
	Value float64       `json:"value"`
	Units TurvoKeyValue `json:"units"`
}

type TurvoAppointment struct {
	Date     string `json:"date"`
	Timezone string `json:"timezone"`
	TimeZone string `json:"timeZone,omitempty"` // Response uses timeZone (camelCase)
	Flex     int    `json:"flex,omitempty"`
	HasTime  bool   `json:"hasTime,omitempty"`
}

type TurvoPlannedAppointmentDate struct {
	SchedulingType TurvoKeyValue     `json:"schedulingType,omitempty"`
	Appointment    *TurvoAppointment `json:"appointment,omitempty"`
}

type TurvoAppointmentRange struct {
	From *TurvoAppointment `json:"from,omitempty"`
	To   *TurvoAppointment `json:"to,omitempty"`
}

type TurvoRouteCustomerOrder struct {
	CustomerID            int `json:"customerId,omitempty"`
	CustomerOrderSourceID int `json:"customerOrderSourceId,omitempty"`
}

type TurvoRouteCarrierOrder struct {
	CarrierID            int `json:"carrierId,omitempty"`
	CarrierOrderSourceID int `json:"carrierOrderSourceId,omitempty"`
}

type TurvoTransportation struct {
	Mode        TurvoKeyValue `json:"mode,omitempty"`
	ServiceType TurvoKeyValue `json:"serviceType,omitempty"`
}

type TurvoDistanceWithUnits struct {
	Value float64       `json:"value"`
	Units TurvoKeyValue `json:"units"`
}

type TurvoModeInfo struct {
	Operation             int                `json:"_operation,omitempty"`
	SourceSegmentSequence string             `json:"sourceSegmentSequence,omitempty"`
	Mode                  TurvoKeyValue      `json:"mode,omitempty"`
	ServiceType           TurvoKeyValue      `json:"serviceType,omitempty"`
	TotalSegmentValue     *TurvoSegmentValue `json:"totalSegmentValue,omitempty"`
}

type TurvoSegmentValue struct {
	Sync     bool          `json:"sync,omitempty"`
	Value    float64       `json:"value"`
	Currency TurvoKeyValue `json:"currency,omitempty"`
}

type TurvoCreateCustomerOrder struct {
	CustomerOrderSourceID int               `json:"customerOrderSourceId"`
	Customer              TurvoAccount      `json:"customer,omitempty"`
	Items                 []TurvoOrderItem  `json:"items,omitempty"`
	Costs                 *TurvoOrderCosts  `json:"costs,omitempty"`
	ExternalIds           []TurvoExternalId `json:"externalIds,omitempty"`
}

type TurvoOrderItem struct {
	Dimensions           *TurvoDimensions             `json:"dimensions,omitempty"`
	ItemCategory         TurvoKeyValue                `json:"itemCategory,omitempty"`
	Qty                  float64                      `json:"qty,omitempty"`
	Unit                 TurvoKeyValue                `json:"unit,omitempty"`
	HandlingQty          float64                      `json:"handlingQty,omitempty"`
	HandlingUnit         TurvoKeyValue                `json:"handlingUnit,omitempty"`
	Name                 string                       `json:"name,omitempty"`
	Notes                string                       `json:"notes,omitempty"`
	PickupLocation       []TurvoItemLocationReference `json:"pickupLocation,omitempty"`
	DeliveryLocation     []TurvoItemLocationReference `json:"deliveryLocation,omitempty"`
	Operation            int                          `json:"_operation,omitempty"`
	ItemNumber           string                       `json:"itemNumber,omitempty"`
	Nmfc                 string                       `json:"nmfc,omitempty"`
	NmfcSub              string                       `json:"nmfcSub,omitempty"`
	IsHazmat             bool                         `json:"isHazmat,omitempty"`
	Stackable            bool                         `json:"stackable,omitempty"`
	FreightClass         TurvoKeyValue                `json:"freightClass,omitempty"`
	Value                float64                      `json:"value,omitempty"`
	TotalValue           float64                      `json:"totalValue,omitempty"`
	Currency             TurvoKeyValue                `json:"currency,omitempty"`
	MinTemp              *TurvoTemperature            `json:"minTemp,omitempty"`
	MaxTemp              *TurvoTemperature            `json:"maxTemp,omitempty"`
	StackDimensionsLimit *TurvoStackDimensions        `json:"stackDimensionsLimit,omitempty"`
	LoadBearingCapacity  *TurvoWeightWithUnit         `json:"loadBearingCapacity,omitempty"`
	MaxStackCount        int                          `json:"maxStackCount,omitempty"`
}

type TurvoDimensions struct {
	Length float64       `json:"length"`
	Width  float64       `json:"width"`
	Height float64       `json:"height"`
	Units  TurvoKeyValue `json:"units"`
}

type TurvoItemLocationReference struct {
	GlobalShipLocationSourceID string `json:"globalShipLocationSourceId,omitempty"`
	Name                       string `json:"name,omitempty"`
}

type TurvoTemperature struct {
	Temp     float64       `json:"temp"`
	TempUnit TurvoKeyValue `json:"tempUnit"`
}

type TurvoStackDimensions struct {
	Height float64       `json:"height"`
	Width  float64       `json:"width"`
	Unit   TurvoKeyValue `json:"unit"`
}

type TurvoWeightWithUnit struct {
	Value float64       `json:"value"`
	Unit  TurvoKeyValue `json:"unit"`
}

type TurvoOrderCosts struct {
	TotalAmount float64             `json:"totalAmount"`
	LineItem    []TurvoCostLineItem `json:"lineItem,omitempty"`
}

type TurvoCostLineItem struct {
	Code     TurvoKeyValue `json:"code,omitempty"`
	Qty      float64       `json:"qty,omitempty"`
	Price    float64       `json:"price,omitempty"`
	Amount   float64       `json:"amount,omitempty"`
	Billable bool          `json:"billable,omitempty"`
	Notes    string        `json:"notes,omitempty"`
}

type TurvoExternalId struct {
	Type               TurvoKeyValue `json:"type,omitempty"`
	Value              string        `json:"value"`
	CopyToCarrierOrder bool          `json:"copyToCarrierOrder,omitempty"`
}

type TurvoCreateCarrierOrder struct {
	CarrierOrderSourceID int           `json:"carrierOrderSourceId,omitempty"`
	Carrier              TurvoAccount  `json:"carrier,omitempty"`
	Drivers              []TurvoDriver `json:"drivers,omitempty"`
}

type TurvoDriver struct {
	DriverID        int `json:"driverId,omitempty"`
	Operation       int `json:"_operation,omitempty"`
	SegmentSequence int `json:"segmentSequence,omitempty"`
}

type TurvoCreateParty struct {
	Account TurvoAccount `json:"account,omitempty"`
}

// TurvoKeyValue represents a key-value pair used throughout Turvo API
type TurvoKeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	ID    int    `json:"id,omitempty"`
}

// TurvoShipmentCreateResponse represents the response from Turvo's POST shipments endpoint
type TurvoShipmentCreateResponse struct {
	Status  string                     `json:"Status"`
	Details TurvoShipmentCreateDetails `json:"details"`
}

// TurvoShipmentResponse represents the response from Turvo's GET shipment endpoint
type TurvoShipmentResponse struct {
	Status  string                     `json:"Status"`
	Details TurvoShipmentCreateDetails `json:"details"`
}

type TurvoShipmentCreateErrorResponse struct {
	Status  string                                  `json:"Status"`
	Details TurvoShipmentCreateErrorResponseDetails `json:"details"`
}

type TurvoShipmentCreateErrorResponseDetails struct {
	ErrorMessage string `json:"errorMessage"`
	ErrorCode    string `json:"errorCode"`
}

type TurvoShipmentCreateDetails struct {
	ID               int                            `json:"id"`
	CustomID         string                         `json:"customId,omitempty"`
	LtlShipment      bool                           `json:"ltlShipment"`
	Phase            TurvoKeyValue                  `json:"phase,omitempty"`
	Services         []TurvoKeyValue                `json:"services,omitempty"`
	CustomAttributes map[string]interface{}         `json:"custom_attributes,omitempty"`
	StartDate        TurvoDateWithTimezone          `json:"startDate,omitempty"`
	EndDate          TurvoDateWithTimezoneAndFlex   `json:"endDate,omitempty"`
	Transportation   *TurvoTransportation           `json:"transportation,omitempty"`
	Status           *TurvoCreateStatus             `json:"status,omitempty"`
	Tracking         *TurvoTracking                 `json:"tracking,omitempty"`
	Margin           *TurvoMargin                   `json:"margin,omitempty"`
	Equipment        []TurvoEquipmentResponse       `json:"equipment,omitempty"`
	Contributors     []TurvoContributorResponse     `json:"contributors,omitempty"`
	Lane             TurvoLane                      `json:"lane,omitempty"`
	GlobalRoute      []TurvoGlobalRouteStopResponse `json:"globalRoute,omitempty"`
	ModeInfo         []TurvoModeInfoResponse        `json:"modeInfo,omitempty"`
	CustomerOrder    []TurvoCustomerOrderResponse   `json:"customerOrder,omitempty"`
	CarrierOrder     []TurvoCarrierOrderResponse    `json:"carrierOrder,omitempty"`
	Party            []TurvoPartyResponse           `json:"party,omitempty"`
	UseRoutingGuide  bool                           `json:"use_routing_guide,omitempty"`
}

type TurvoDateWithTimezoneAndFlex struct {
	Date     string `json:"date"`
	TimeZone string `json:"timeZone"`
	Flex     int    `json:"flex,omitempty"`
}

type TurvoTracking struct {
	IsTrackable bool             `json:"isTrackable,omitempty"`
	Deleted     bool             `json:"deleted,omitempty"`
	IsTracking  bool             `json:"isTracking,omitempty"`
	Description string           `json:"description,omitempty"`
	Source      string           `json:"source,omitempty"`
	Frequency   int              `json:"frequency,omitempty"`
	RouteSteps  *TurvoRouteSteps `json:"routeSteps,omitempty"`
}

type TurvoRouteSteps struct {
	VisitedGeoWayPoints string `json:"visitedGeoWayPoints,omitempty"`
	CountGeoWayPoints   int    `json:"countGeoWayPoints,omitempty"`
	StepsPolyline       string `json:"stepsPolyline,omitempty"`
}

type TurvoMargin struct {
	MinPay float64 `json:"minPay,omitempty"`
	MaxPay float64 `json:"maxPay,omitempty"`
}

type TurvoEquipmentResponse struct {
	Deleted        bool          `json:"deleted,omitempty"`
	ID             int           `json:"id,omitempty"`
	Type           TurvoKeyValue `json:"type,omitempty"`
	Size           TurvoKeyValue `json:"size,omitempty"`
	Weight         float64       `json:"weight,omitempty"`
	Temp           float64       `json:"temp,omitempty"`
	ShipmentLength float64       `json:"shipmentLength,omitempty"`
}

type TurvoContributorResponse struct {
	Deleted         bool          `json:"deleted,omitempty"`
	ID              int           `json:"id,omitempty"`
	ContributorUser TurvoUser     `json:"contributorUser,omitempty"`
	Type            TurvoKeyValue `json:"type,omitempty"`
}

type TurvoUser struct {
	Name string `json:"name,omitempty"`
	ID   int    `json:"id,omitempty"`
}

type TurvoGlobalRouteStopResponse struct {
	ID                         int                          `json:"id,omitempty"`
	Name                       string                       `json:"name,omitempty"`
	GlobalShipLocationSourceID string                       `json:"globalShipLocationSourceId,omitempty"`
	SchedulingType             TurvoKeyValue                `json:"schedulingType,omitempty"`
	StopType                   TurvoKeyValue                `json:"stopType,omitempty"`
	Timezone                   string                       `json:"timezone,omitempty"`
	Location                   TurvoLocationReference       `json:"location,omitempty"`
	Address                    *TurvoAddress                `json:"address,omitempty"`
	SegmentID                  string                       `json:"segmentId,omitempty"`
	SegmentSequence            int                          `json:"segmentSequence,omitempty"`
	Sequence                   int                          `json:"sequence,omitempty"`
	State                      string                       `json:"state,omitempty"`
	Appointment                *TurvoAppointment            `json:"appointment,omitempty"`
	Services                   []TurvoKeyValue              `json:"services,omitempty"`
	PoNumbers                  []string                     `json:"poNumbers,omitempty"`
	Notes                      string                       `json:"notes,omitempty"`
	Contact                    *TurvoContact                `json:"contact,omitempty"`
	CustomerOrder              []TurvoRouteOrderReference   `json:"customerOrder,omitempty"`
	CarrierOrder               []TurvoRouteOrderReference   `json:"carrierOrder,omitempty"`
	Deleted                    bool                         `json:"deleted,omitempty"`
	FragmentDistance           *TurvoDistanceWithUnits      `json:"fragmentDistance,omitempty"`
	StopLevelFragmentDistance  float64                      `json:"stop_level_fragment_distance,omitempty"`
	IsShipmentDwellTimeEdited  bool                         `json:"isShipmentDwellTimeEdited,omitempty"`
	ExpectedDwellTime          *TurvoTimeWithUnits          `json:"expectedDwellTime,omitempty"`
	LayoverTime                *TurvoTimeWithUnits          `json:"layoverTime,omitempty"`
	OriginalAppointmentDate    *TurvoPlannedAppointmentDate `json:"originalAppointmentDate,omitempty"`
	PlannedDate                *TurvoPlannedAppointmentDate `json:"plannedDate,omitempty"`
	Transportation             *TurvoTransportation         `json:"transportation,omitempty"`
}

type TurvoContact struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type TurvoRouteOrderReference struct {
	CustomerID int  `json:"customerId,omitempty"`
	CarrierID  int  `json:"carrierId,omitempty"`
	ID         int  `json:"id,omitempty"`
	Deleted    bool `json:"deleted,omitempty"`
}

type TurvoModeInfoResponse struct {
	ID                int                `json:"id,omitempty"`
	SegmentID         string             `json:"segmentId,omitempty"`
	Mode              TurvoKeyValue      `json:"mode,omitempty"`
	ServiceType       TurvoKeyValue      `json:"serviceType,omitempty"`
	Deleted           bool               `json:"deleted,omitempty"`
	TotalSegmentValue *TurvoSegmentValue `json:"totalSegmentValue,omitempty"`
}

type TurvoCustomerOrderResponse struct {
	ID          int                       `json:"id,omitempty"`
	Deleted     bool                      `json:"deleted,omitempty"`
	Customer    TurvoAccountResponse      `json:"customer,omitempty"`
	TotalMiles  float64                   `json:"totalMiles,omitempty"`
	Items       []TurvoOrderItemResponse  `json:"items,omitempty"`
	Route       []TurvoRouteStop          `json:"route,omitempty"`
	Costs       *TurvoOrderCostsResponse  `json:"costs,omitempty"`
	ExternalIds []TurvoExternalIdResponse `json:"externalIds,omitempty"`
	Contacts    []TurvoShipContact        `json:"contacts,omitempty"`
}

type TurvoAccountResponse struct {
	ID         int               `json:"id,omitempty"`
	Name       string            `json:"name,omitempty"`
	Owner      *TurvoUser        `json:"owner,omitempty"`
	Comissions []TurvoCommission `json:"comissions,omitempty"`
}

type TurvoCommission struct {
	UserID int    `json:"userId,omitempty"`
	Value  string `json:"value,omitempty"`
}

type TurvoOrderItemResponse struct {
	Deleted          bool                        `json:"deleted,omitempty"`
	ID               int                         `json:"id,omitempty"`
	Dimensions       *TurvoDimensions            `json:"dimensions,omitempty"`
	ItemCategory     TurvoKeyValue               `json:"itemCategory,omitempty"`
	Qty              float64                     `json:"qty,omitempty"`
	Unit             TurvoKeyValue               `json:"unit,omitempty"`
	Name             string                      `json:"name,omitempty"`
	Notes            string                      `json:"notes,omitempty"`
	PickupLocation   []TurvoItemLocationResponse `json:"pickupLocation,omitempty"`
	DeliveryLocation []TurvoItemLocationResponse `json:"deliveryLocation,omitempty"`
	ItemNumber       string                      `json:"itemNumber,omitempty"`
	Nmfc             string                      `json:"nmfc,omitempty"`
	NmfcSub          string                      `json:"nmfcSub,omitempty"`
	IsHazmat         bool                        `json:"isHazmat,omitempty"`
	NetWeight        *TurvoNetWeight             `json:"netWeight,omitempty"`
	Weight           float64                     `json:"weight,omitempty"`
	FreightClass     TurvoKeyValue               `json:"freightClass,omitempty"`
	Value            float64                     `json:"value,omitempty"`
	TotalValue       float64                     `json:"totalValue,omitempty"`
	Currency         TurvoKeyValue               `json:"currency,omitempty"`
}

type TurvoItemLocationResponse struct {
	ID                         int    `json:"id,omitempty"`
	GlobalShipLocationID       int    `json:"globalShipLocationId,omitempty"`
	Name                       string `json:"name,omitempty"`
	GlobalShipLocationSourceID string `json:"globalShipLocationSourceId,omitempty"`
}

type TurvoNetWeight struct {
	Weight float64 `json:"weight"`
}

type TurvoRouteStop struct {
	Deleted     bool                   `json:"deleted,omitempty"`
	ID          int                    `json:"id,omitempty"`
	StopType    TurvoKeyValue          `json:"stopType,omitempty"`
	Location    TurvoLocationInfo      `json:"location,omitempty"`
	Address     *TurvoAddress          `json:"address,omitempty"`
	Email       string                 `json:"email,omitempty"`
	Extension   string                 `json:"extension,omitempty"`
	Phone       string                 `json:"phone,omitempty"`
	Sequence    int                    `json:"sequence,omitempty"`
	State       string                 `json:"state,omitempty"`
	Appointment *TurvoRouteAppointment `json:"appointment,omitempty"`
	Services    []TurvoKeyValue        `json:"services,omitempty"`
	PoNumbers   []string               `json:"poNumbers,omitempty"`
	Notes       string                 `json:"notes,omitempty"`
	Contact     *TurvoContact          `json:"contact,omitempty"`
}

type TurvoLocationInfo struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type TurvoAddress struct {
	ID    string  `json:"id,omitempty"`
	Line1 string  `json:"line1,omitempty"`
	City  string  `json:"city,omitempty"`
	State string  `json:"state,omitempty"`
	Zip   string  `json:"zip,omitempty"`
	Lon   float64 `json:"lon,omitempty"`
	Lat   float64 `json:"lat,omitempty"`
}

type TurvoRouteAppointment struct {
	Start      string        `json:"start,omitempty"`
	TimeZone   string        `json:"timeZone,omitempty"`
	Flex       int           `json:"flex,omitempty"`
	HasTime    bool          `json:"hasTime,omitempty"`
	Scheduling TurvoKeyValue `json:"scheduling,omitempty"`
}

type TurvoOrderCostsResponse struct {
	SubTotal    float64                     `json:"subTotal,omitempty"`
	TotalTax    float64                     `json:"totalTax,omitempty"`
	TotalAmount float64                     `json:"totalAmount,omitempty"`
	Deleted     bool                        `json:"deleted,omitempty"`
	Notes       string                      `json:"notes,omitempty"`
	LineItem    []TurvoCostLineItemResponse `json:"lineItem,omitempty"`
}

type TurvoCostLineItemResponse struct {
	Deleted  bool          `json:"deleted,omitempty"`
	ID       int           `json:"id,omitempty"`
	Code     TurvoKeyValue `json:"code,omitempty"`
	Qty      float64       `json:"qty,omitempty"`
	Price    float64       `json:"price,omitempty"`
	Amount   float64       `json:"amount,omitempty"`
	Billable bool          `json:"billable,omitempty"`
	Notes    string        `json:"notes,omitempty"`
}

type TurvoExternalIdResponse struct {
	Deleted            bool          `json:"deleted,omitempty"`
	ID                 int           `json:"id,omitempty"`
	Type               TurvoKeyValue `json:"type,omitempty"`
	Value              string        `json:"value"`
	CopyToCarrierOrder bool          `json:"copyToCarrierOrder,omitempty"`
}

type TurvoShipContact struct {
	Deleted       bool          `json:"deleted,omitempty"`
	ID            int           `json:"id,omitempty"`
	Role          TurvoKeyValue `json:"role,omitempty"`
	ShipContactID int           `json:"shipContactId,omitempty"`
	Title         string        `json:"title,omitempty"`
	Name          string        `json:"name,omitempty"`
	Phone         []TurvoPhone  `json:"phone,omitempty"`
	Email         []TurvoEmail  `json:"email,omitempty"`
}

type TurvoPhone struct {
	ID        string        `json:"id,omitempty"`
	Extension string        `json:"extension,omitempty"`
	Number    string        `json:"number,omitempty"`
	Type      TurvoKeyValue `json:"type,omitempty"`
	Deleted   bool          `json:"deleted,omitempty"`
	IsPrimary bool          `json:"isPrimary,omitempty"`
	Primary   bool          `json:"primary,omitempty"`
	Country   TurvoKeyValue `json:"country,omitempty"`
}

type TurvoEmail struct {
	ID        string        `json:"id,omitempty"`
	Email     string        `json:"email"`
	Type      TurvoKeyValue `json:"type,omitempty"`
	Deleted   bool          `json:"deleted,omitempty"`
	IsPrimary bool          `json:"isPrimary,omitempty"`
	Primary   bool          `json:"primary,omitempty"`
}

type TurvoCarrierOrderResponse struct {
	ID          int                       `json:"id,omitempty"`
	Deleted     bool                      `json:"deleted,omitempty"`
	Carrier     TurvoAccountResponse      `json:"carrier,omitempty"`
	Costs       *TurvoOrderCostsResponse  `json:"costs,omitempty"`
	ExternalIds []TurvoExternalIdResponse `json:"externalIds,omitempty"`
	Contacts    []TurvoShipContact        `json:"contacts,omitempty"`
	Drivers     []TurvoDriverResponse     `json:"drivers,omitempty"`
}

type TurvoDriverResponse struct {
	SegmentSequence        int                 `json:"segmentSequence,omitempty"`
	SegmentID              string              `json:"segmentId,omitempty"`
	ContextType            string              `json:"contextType,omitempty"`
	Context                *TurvoDriverContext `json:"context,omitempty"`
	Phone                  *TurvoPhone         `json:"phone,omitempty"`
	Email                  *TurvoEmail         `json:"email,omitempty"`
	DriverAssignmentStatus TurvoKeyValue       `json:"driverAssignmentStatus,omitempty"`
	DriverAssignmentID     int                 `json:"driverAssignmentId,omitempty"`
	Deleted                bool                `json:"deleted,omitempty"`
}

type TurvoDriverContext struct {
	ID               int    `json:"id,omitempty"`
	Name             string `json:"name,omitempty"`
	LinkedUserID     int    `json:"linkedUserId,omitempty"`
	MessagingEnabled bool   `json:"messagingEnabled,omitempty"`
}

type TurvoPartyResponse struct {
	ID        int                      `json:"id,omitempty"`
	Deleted   bool                     `json:"deleted,omitempty"`
	Account   TurvoAccountResponse     `json:"account,omitempty"`
	Costs     *TurvoOrderCostsResponse `json:"costs,omitempty"`
	Payment   []TurvoPayment           `json:"payment,omitempty"`
	Deduction []TurvoDeduction         `json:"deduction,omitempty"`
}

type TurvoPayment struct {
	ReferenceNo string        `json:"referenceNo,omitempty"`
	Deleted     bool          `json:"deleted,omitempty"`
	Amount      float64       `json:"amount,omitempty"`
	PaymentDate string        `json:"paymentDate,omitempty"`
	Method      TurvoKeyValue `json:"method,omitempty"`
	Notes       string        `json:"notes,omitempty"`
	Invoice     *TurvoInvoice `json:"invoice,omitempty"`
	ID          int           `json:"id,omitempty"`
}

type TurvoInvoice struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	InvoiceNo string `json:"invoiceNo,omitempty"`
	InvoiceID string `json:"invoiceId,omitempty"`
}

type TurvoDeduction struct {
	Deleted     bool    `json:"deleted,omitempty"`
	ID          int     `json:"id,omitempty"`
	Amount      float64 `json:"amount,omitempty"`
	Date        string  `json:"date,omitempty"`
	Description string  `json:"description,omitempty"`
	IssueDate   string  `json:"issueDate,omitempty"`
}
