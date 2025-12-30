package load

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lwlach/turvo-integration-backend/internal/models"
	"github.com/lwlach/turvo-integration-backend/internal/turvo"
)

type Service struct {
	turvoClient *turvo.Client
}

func NewService(turvoClient *turvo.Client) *Service {
	return &Service{
		turvoClient: turvoClient,
	}
}

// GetAllLoads fetches all loads from Turvo and converts them to Drumkit format
// Deprecated: Use GetLoads with filters instead
func (s *Service) GetAllLoads() ([]models.Load, error) {
	filters := models.LoadFilters{
		Page:  1,
		Limit: 1000, // Large limit to get all
	}
	response, err := s.GetLoads(filters)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// GetLoads fetches loads from Turvo with filtering and pagination
func (s *Service) GetLoads(filters models.LoadFilters) (*models.LoadListResponse, error) {
	// Map our filters to Turvo filters
	turvoFilters := s.mapToTurvoFilters(filters)

	// Fetch shipments from Turvo with filters
	turvoShipments, turvoPagination, err := s.turvoClient.ListShipmentsWithFiltersAndPagination(turvoFilters)
	if err != nil {
		return nil, err
	}

	// Convert to loads
	loads := make([]models.Load, len(turvoShipments))

	if filters.IncludeDetails {
		// Fetch detailed shipment information concurrently using goroutines
		var wg sync.WaitGroup
		var mu sync.Mutex

		for i, shipment := range turvoShipments {
			wg.Add(1)
			go func(index int, shipmentID int, baseShipment models.TurvoShipment) {
				defer wg.Done()

				// Fetch detailed shipment information
				detailedShipment, err := s.turvoClient.GetShipment(shipmentID)
				var load models.Load
				if err != nil {
					// If fetching details fails, fall back to basic conversion
					load = s.turvoToDrumkit(&baseShipment)
				} else {
					// Convert detailed shipment to load
					load = s.turvoDetailsToDrumkit(detailedShipment)
				}

				// Store result in the correct position
				mu.Lock()
				loads[index] = load
				mu.Unlock()
			}(i, shipment.ID, shipment)
		}

		// Wait for all goroutines to complete
		wg.Wait()
	} else {
		// Use basic list conversion sequentially
		for i, shipment := range turvoShipments {
			loads[i] = s.turvoToDrumkit(&shipment)
		}
	}

	// Calculate pagination from Turvo's response
	// Note: Turvo doesn't provide total count, only current page info
	// totalRecordsInPage represents the number of records in the current page response
	// When moreAvailable is false, this is the exact total count
	// When moreAvailable is true, we don't know the total count
	total := turvoPagination.TotalRecordsInPage

	// Calculate pages
	// When moreAvailable is false, we can calculate pages accurately
	// When moreAvailable is true, we set pages to current page + 1 to indicate more pages exist
	var pages int
	if turvoPagination.MoreAvailable {
		// More pages available, but we don't know total, so indicate at least one more page
		pages = filters.Page + 1
	} else {
		// No more pages, calculate pages from total
		pages = (total + filters.Limit - 1) / filters.Limit
		if pages == 0 {
			pages = 1
		}
	}

	return &models.LoadListResponse{
		Data: loads,
		Pagination: models.Pagination{
			Total: total,
			Pages: pages,
			Page:  filters.Page,
			Limit: filters.Limit,
		},
	}, nil
}

// mapToTurvoFilters maps our API filters to Turvo's filter format
func (s *Service) mapToTurvoFilters(filters models.LoadFilters) models.TurvoShipmentFilters {
	turvoFilters := models.TurvoShipmentFilters{
		PageSize: filters.Limit,
		Start:    (filters.Page - 1) * filters.Limit,
	}

	// Map status: convert API status to Turvo status code
	if filters.Status != "" {
		statusKey, _ := APIToTurvoStatus(filters.Status)
		turvoFilters.Status = statusKey
	}

	// Map customerId
	if filters.CustomerID != "" {
		turvoFilters.CustomerID = filters.CustomerID
	}

	// Map pickup date filters: convert to RFC3339 format in UTC
	if filters.PickupDateFrom != nil {
		// Convert to UTC and format
		utcDate := filters.PickupDateFrom.UTC()
		turvoFilters.PickupDateGte = utcDate.Format(time.RFC3339)
	}
	if filters.PickupDateTo != nil {
		// Convert to UTC and format
		utcDate := filters.PickupDateTo.UTC()
		turvoFilters.PickupDateLte = utcDate.Format(time.RFC3339)
	}

	return turvoFilters
}

// CreateLoad creates a new load in Turvo from a Drumkit load format
func (s *Service) CreateLoad(load *models.Load) (*models.LoadCreateResponse, error) {
	// Validate the load
	if err := ValidateLoad(load); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	turvoShipment := s.drumkitToTurvo(load)

	// Validate required fields before creating shipment
	if err := ValidateTurvoShipment(turvoShipment); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	response, err := s.turvoClient.CreateShipment(turvoShipment)
	if err != nil {
		return nil, err
	}

	// Return minimal response with only id and createdAt
	createdAt := time.Now()
	if response.Details.ID > 0 {
		return &models.LoadCreateResponse{
			ID:        fmt.Sprintf("%d", response.Details.ID),
			CreatedAt: createdAt,
		}, nil
	}

	return nil, fmt.Errorf("invalid response: shipment ID is missing")
}

// turvoCreateResponseToDrumkit converts a TurvoShipmentCreateResponse to Drumkit load format
// Uses the original load as a base to preserve fields that might not be in the response
func (s *Service) turvoCreateResponseToDrumkit(response *models.TurvoShipmentCreateResponse, originalLoad *models.Load) models.Load {
	load := *originalLoad // Start with original load to preserve all fields

	if response.Details.ID == 0 {
		return load // Invalid response
	}

	shipment := &response.Details

	// Map ID and customId
	if shipment.ID > 0 {
		load.ExternalTMSLoadID = fmt.Sprintf("%d", shipment.ID)
	}
	if shipment.CustomID != "" {
		load.FreightLoadID = shipment.CustomID
	}

	// Map status from new structure using status mapper
	if shipment.Status != nil && shipment.Status.Code.Value != "" {
		apiStatus := TurvoStatusToAPI(shipment.Status.Code.Key, shipment.Status.Code.Value)
		if apiStatus != "" {
			load.Status = apiStatus
		}
	}

	// Map dates
	if shipment.StartDate.Date != "" {
		date, err := time.Parse(time.RFC3339, shipment.StartDate.Date)
		if err == nil && load.Pickup != nil {
			load.Pickup.ReadyTime = &date
		}
	}
	if shipment.EndDate.Date != "" {
		date, err := time.Parse(time.RFC3339, shipment.EndDate.Date)
		if err == nil && load.Consignee != nil {
			load.Consignee.ApptTime = &date
		}
	}

	// Map customer order
	if len(shipment.CustomerOrder) > 0 {
		custOrder := shipment.CustomerOrder[0]
		if load.Customer == nil {
			load.Customer = &models.Customer{}
		}
		if custOrder.Customer.Name != "" {
			load.Customer.Name = custOrder.Customer.Name
		}
		if custOrder.Customer.ID > 0 {
			load.Customer.ExternalTMSId = fmt.Sprintf("%d", custOrder.Customer.ID)
		}
	}

	// Map carrier order
	if len(shipment.CarrierOrder) > 0 {
		carrierOrder := shipment.CarrierOrder[0]
		if load.Carrier == nil {
			load.Carrier = &models.Carrier{}
		}
		if carrierOrder.Carrier.Name != "" {
			load.Carrier.Name = carrierOrder.Carrier.Name
		}
		if carrierOrder.Carrier.ID > 0 {
			load.Carrier.ExternalTMSId = fmt.Sprintf("%d", carrierOrder.Carrier.ID)
		}
	}

	return load
}

// drumkitToTurvo converts a Drumkit load to Turvo shipment format for creation
func (s *Service) drumkitToTurvo(load *models.Load) *models.TurvoShipmentCreate {
	shipment := &models.TurvoShipmentCreate{
		LtlShipment: false,
	}

	// Map status using status mapper
	if load.Status != "" {
		statusKey, statusValue := APIToTurvoStatus(load.Status)
		shipment.Status = &models.TurvoCreateStatus{
			Code: models.TurvoStatusCode{
				Key:   statusKey,
				Value: statusValue,
			},
		}
	}

	// Map dates from pickup and consignee
	// Pickup maps to startDate, Consignee maps to endDate
	timezone := "America/New_York" // Default timezone
	if load.Pickup.Timezone != "" {
		timezone = load.Pickup.Timezone
	}

	// Map pickup to startDate (prefer ReadyTime, fallback to ApptTime)
	// Validation ensures at least one of ReadyTime or ApptTime exists
	var startDate models.TurvoDateWithTimezone
	if load.Pickup.ReadyTime != nil {
		startDate = models.TurvoDateWithTimezone{
			Date:     load.Pickup.ReadyTime.Format(time.RFC3339),
			TimeZone: timezone,
		}
	} else {
		// ApptTime is guaranteed to exist if ReadyTime doesn't (validated)
		startDate = models.TurvoDateWithTimezone{
			Date:     load.Pickup.ApptTime.Format(time.RFC3339),
			TimeZone: timezone,
		}
	}

	// Map consignee to endDate
	// Validation ensures ApptTime exists
	consigneeTimezone := timezone
	if load.Consignee.Timezone != "" {
		consigneeTimezone = load.Consignee.Timezone
	}
	endDate := models.TurvoDateWithTimezone{
		Date:     load.Consignee.ApptTime.Format(time.RFC3339),
		TimeZone: consigneeTimezone,
	}

	// Set dates (required fields)
	shipment.StartDate = startDate
	shipment.EndDate = endDate

	// Map customId (freightLoadID)
	if load.FreightLoadID != "" {
		// Note: TurvoShipmentCreate doesn't have CustomID field, it's set in response
		// We'll need to check if there's a way to set it during creation
	}

	// Map lane (required field)
	// Validation ensures at least one of (City and State) or Name exists for each
	var laneStart, laneEnd string
	if load.Pickup.City != "" && load.Pickup.State != "" {
		laneStart = fmt.Sprintf("%s, %s", load.Pickup.City, load.Pickup.State)
	} else {
		laneStart = load.Pickup.Name
	}
	if load.Consignee.City != "" && load.Consignee.State != "" {
		laneEnd = fmt.Sprintf("%s, %s", load.Consignee.City, load.Consignee.State)
	} else {
		laneEnd = load.Consignee.Name
	}
	shipment.Lane = models.TurvoLane{
		Start: laneStart,
		End:   laneEnd,
	}

	// Map customer order (required field) with full details
	customerOrder := models.TurvoCreateCustomerOrder{
		Customer: models.TurvoAccount{
			Name: load.Customer.Name,
		},
	}
	if id, err := strconv.Atoi(load.Customer.ExternalTMSId); err == nil {
		customerOrder.Customer.ID = id
	}
	// Generate a random customerOrderSourceId
	customerOrder.CustomerOrderSourceID = int(uuid.New().ID())

	// Map external IDs (PO numbers, ref numbers, etc.)
	if load.PoNums != "" || load.Customer.RefNumber != "" {
		externalIds := []models.TurvoExternalId{}
		if load.PoNums != "" {
			poNums := strings.Split(load.PoNums, ",")
			for _, poNum := range poNums {
				poNum = strings.TrimSpace(poNum)
				if poNum != "" {
					externalIds = append(externalIds, models.TurvoExternalId{
						Type:  models.TurvoKeyValue{Key: "1400", Value: "Purchase shipment #"},
						Value: poNum,
					})
				}
			}
		}
		if load.Customer.RefNumber != "" {
			externalIds = append(externalIds, models.TurvoExternalId{
				Type:  models.TurvoKeyValue{Key: "1401", Value: "Reference Number"},
				Value: load.Customer.RefNumber,
			})
		}
		customerOrder.ExternalIds = externalIds
	}
	shipment.CustomerOrder = []models.TurvoCreateCustomerOrder{customerOrder}

	// Map billTo as Party (if provided)
	if load.BillTo != nil {
		party := models.TurvoCreateParty{
			Account: models.TurvoAccount{
				Name: load.BillTo.Name,
			},
		}
		if load.BillTo.ExternalTMSId != "" {
			if id, err := strconv.Atoi(load.BillTo.ExternalTMSId); err == nil {
				party.Account.ID = id
			}
		}
		shipment.Party = []models.TurvoCreateParty{party}
	}

	// Map carrier order (optional, but if provided, ID is required)
	if load.Carrier != nil {
		carrierOrder := models.TurvoCreateCarrierOrder{
			Carrier: models.TurvoAccount{
				Name: load.Carrier.Name,
			},
		}
		// ExternalTMSId is validated to exist and be a valid integer
		if id, err := strconv.Atoi(load.Carrier.ExternalTMSId); err == nil {
			carrierOrder.Carrier.ID = id
		}
		shipment.CarrierOrder = []models.TurvoCreateCarrierOrder{carrierOrder}
	}

	// Map equipment (weight, temperature, etc.)
	if load.TotalWeight != nil || load.Specifications != nil {
		equipment := models.TurvoEquipment{}
		if load.TotalWeight != nil {
			equipment.Weight = *load.TotalWeight
			equipment.WeightUnits = models.TurvoKeyValue{Key: "1520", Value: "lb"}
		}
		if load.Specifications != nil {
			if load.Specifications.MinTempFahrenheit != nil && load.Specifications.MaxTempFahrenheit != nil {
				// Use average temperature
				avgTemp := (*load.Specifications.MinTempFahrenheit + *load.Specifications.MaxTempFahrenheit) / 2
				equipment.Temp = avgTemp
				equipment.TempUnits = models.TurvoKeyValue{Key: "1510", Value: "°F"}
			} else if load.Specifications.MinTempFahrenheit != nil {
				equipment.Temp = *load.Specifications.MinTempFahrenheit
				equipment.TempUnits = models.TurvoKeyValue{Key: "1510", Value: "°F"}
			} else if load.Specifications.MaxTempFahrenheit != nil {
				equipment.Temp = *load.Specifications.MaxTempFahrenheit
				equipment.TempUnits = models.TurvoKeyValue{Key: "1510", Value: "°F"}
			}
		}
		equipment.Type = models.TurvoKeyValue{Key: "1200", Value: "Van"} // Default equipment type
		shipment.Equipment = []models.TurvoEquipment{equipment}
	}

	// Map distance
	if load.RouteMiles != nil {
		shipment.SkipDistanceCalculation = false
		// Distance can be set per route stop if needed
	} else {
		shipment.SkipDistanceCalculation = true
	}

	return shipment
}

// turvoToDrumkit converts a Turvo shipment (from list) to Drumkit load format
func (s *Service) turvoToDrumkit(shipment *models.TurvoShipment) models.Load {
	// Convert ID to string
	idStr := ""
	if shipment.ID > 0 {
		idStr = fmt.Sprintf("%d", shipment.ID)
	}

	load := models.Load{
		ExternalTMSLoadID: idStr,
		FreightLoadID:     shipment.CustomID,
		Status:            TurvoStatusToAPI(shipment.Status.Code.Key, shipment.Status.Code.Value),
	}

	// Map customer from customerOrder array
	if len(shipment.CustomerOrder) > 0 && !shipment.CustomerOrder[0].Deleted {
		customerOrder := shipment.CustomerOrder[0]
		if customerOrder.Customer.Name != "" {
			load.Customer = &models.Customer{
				Name: customerOrder.Customer.Name,
			}
			// Convert customer ID to string for externalTMSId
			if customerOrder.Customer.ID > 0 {
				load.Customer.ExternalTMSId = fmt.Sprintf("%d", customerOrder.Customer.ID)
			}
		}
	}

	// Map carrier from carrierOrder array
	if len(shipment.CarrierOrder) > 0 && !shipment.CarrierOrder[0].Deleted {
		carrierOrder := shipment.CarrierOrder[0]
		if carrierOrder.Carrier.Name != "" {
			load.Carrier = &models.Carrier{
				Name: carrierOrder.Carrier.Name,
			}
			// Convert carrier ID to string for externalTMSId
			if carrierOrder.Carrier.ID > 0 {
				load.Carrier.MCNumber = fmt.Sprintf("%d", carrierOrder.Carrier.ID)
			}
		}
	}

	// Map party information (could be shipper, consignee, etc.)
	// Note: The party array structure may need adjustment based on actual Turvo API response
	for _, party := range shipment.Party {
		if !party.Deleted && party.Account.Name != "" {
			// If we don't have a consignee yet, use party as consignee
			if load.Consignee == nil {
				load.Consignee = &models.Consignee{
					Name: party.Account.Name,
				}
				if party.Account.ID > 0 {
					load.Consignee.ExternalTMSId = fmt.Sprintf("%d", party.Account.ID)
				}
			}
		}
	}

	return load
}

// turvoDetailsToDrumkit converts a detailed Turvo shipment (from GET shipment endpoint) to Drumkit load format
func (s *Service) turvoDetailsToDrumkit(shipment *models.TurvoShipmentCreateDetails) models.Load {
	load := models.Load{}

	// Map ID and customId
	if shipment.ID > 0 {
		load.ExternalTMSLoadID = fmt.Sprintf("%d", shipment.ID)
	}
	if shipment.CustomID != "" {
		load.FreightLoadID = shipment.CustomID
	}

	// Map status
	if shipment.Status != nil && shipment.Status.Code.Value != "" {
		apiStatus := TurvoStatusToAPI(shipment.Status.Code.Key, shipment.Status.Code.Value)
		if apiStatus != "" {
			load.Status = apiStatus
		}
	}

	// Map dates from startDate and endDate
	if shipment.StartDate.Date != "" {
		if date, err := time.Parse(time.RFC3339, shipment.StartDate.Date); err == nil {
			load.Pickup = &models.Pickup{
				ReadyTime: &date,
				Timezone:  shipment.StartDate.TimeZone,
			}
		}
	}
	if shipment.EndDate.Date != "" {
		if date, err := time.Parse(time.RFC3339, shipment.EndDate.Date); err == nil {
			if load.Consignee == nil {
				load.Consignee = &models.Consignee{}
			}
			load.Consignee.ApptTime = &date
			load.Consignee.Timezone = shipment.EndDate.TimeZone
		}
	}

	// Map lane start and end to pickup and consignee
	// Lane format is typically "City, State" or location name
	if shipment.Lane.Start != "" {
		if load.Pickup == nil {
			load.Pickup = &models.Pickup{}
		}
		// If name is not set, use lane start as name
		if load.Pickup.Name == "" {
			load.Pickup.Name = shipment.Lane.Start
		}
		// Try to parse "City, State" format
		if parts := strings.Split(shipment.Lane.Start, ","); len(parts) == 2 {
			city := strings.TrimSpace(parts[0])
			state := strings.TrimSpace(parts[1])
			// Only set if not already set from globalRoute
			if load.Pickup.City == "" {
				load.Pickup.City = city
			}
			if load.Pickup.State == "" {
				load.Pickup.State = state
			}
		}
	}
	if shipment.Lane.End != "" {
		if load.Consignee == nil {
			load.Consignee = &models.Consignee{}
		}
		// If name is not set, use lane end as name
		if load.Consignee.Name == "" {
			load.Consignee.Name = shipment.Lane.End
		}
		// Try to parse "City, State" format
		if parts := strings.Split(shipment.Lane.End, ","); len(parts) == 2 {
			city := strings.TrimSpace(parts[0])
			state := strings.TrimSpace(parts[1])
			// Only set if not already set from globalRoute
			if load.Consignee.City == "" {
				load.Consignee.City = city
			}
			if load.Consignee.State == "" {
				load.Consignee.State = state
			}
		}
	}

	// Map globalRoute to pickup and consignee
	for _, stop := range shipment.GlobalRoute {
		if stop.Deleted {
			continue
		}

		// Determine if this is pickup or delivery based on stopType
		isPickup := stop.StopType.Value == "Pickup" || stop.StopType.Key == "1500"

		if isPickup {
			if load.Pickup == nil {
				load.Pickup = &models.Pickup{}
			}
			load.Pickup.Name = stop.Name
			if stop.Timezone != "" {
				load.Pickup.Timezone = stop.Timezone
			}
			if stop.Notes != "" {
				load.Pickup.ApptNote = stop.Notes
			}
			if stop.Appointment != nil && stop.Appointment.Date != "" {
				if apptTime, err := time.Parse(time.RFC3339, stop.Appointment.Date); err == nil {
					load.Pickup.ApptTime = &apptTime
					// If we don't have ReadyTime yet, use ApptTime
					if load.Pickup.ReadyTime == nil {
						load.Pickup.ReadyTime = &apptTime
					}
				}
			}
			if len(stop.PoNumbers) > 0 {
				load.PoNums = strings.Join(stop.PoNumbers, ", ")
			}
			// Map address if available
			if stop.Address != nil {
				load.Pickup.AddressLine1 = stop.Address.Line1
				load.Pickup.City = stop.Address.City
				load.Pickup.State = stop.Address.State
				load.Pickup.Zipcode = stop.Address.Zip
			}
			// Map contact information
			if stop.Contact != nil {
				load.Pickup.Contact = stop.Contact.Name
			}
		} else {
			// Delivery/Consignee
			if load.Consignee == nil {
				load.Consignee = &models.Consignee{}
			}
			load.Consignee.Name = stop.Name
			if stop.Timezone != "" {
				load.Consignee.Timezone = stop.Timezone
			}
			if stop.Notes != "" {
				load.Consignee.ApptNote = stop.Notes
			}
			if stop.Appointment != nil && stop.Appointment.Date != "" {
				if apptTime, err := time.Parse(time.RFC3339, stop.Appointment.Date); err == nil {
					load.Consignee.ApptTime = &apptTime
				}
			}
			// Map address if available
			if stop.Address != nil {
				load.Consignee.AddressLine1 = stop.Address.Line1
				load.Consignee.City = stop.Address.City
				load.Consignee.State = stop.Address.State
				load.Consignee.Zipcode = stop.Address.Zip
			}
			// Map contact information
			if stop.Contact != nil {
				load.Consignee.Contact = stop.Contact.Name
			}
		}
	}

	// Map customer order with detailed information
	if len(shipment.CustomerOrder) > 0 && !shipment.CustomerOrder[0].Deleted {
		custOrder := shipment.CustomerOrder[0]
		load.Customer = &models.Customer{
			Name: custOrder.Customer.Name,
		}
		if custOrder.Customer.ID > 0 {
			load.Customer.ExternalTMSId = fmt.Sprintf("%d", custOrder.Customer.ID)
		}
		// Map route stops for address and contact information
		if len(custOrder.Route) > 0 {
			for _, routeStop := range custOrder.Route {
				if routeStop.Deleted {
					continue
				}
				if routeStop.StopType.Value == "Pickup" || routeStop.StopType.Key == "1500" {
					if load.Pickup != nil {
						if routeStop.Address != nil {
							load.Pickup.AddressLine1 = routeStop.Address.Line1
							load.Pickup.City = routeStop.Address.City
							load.Pickup.State = routeStop.Address.State
							load.Pickup.Zipcode = routeStop.Address.Zip
						}
						if routeStop.Phone != "" {
							load.Pickup.Phone = routeStop.Phone
						}
						if routeStop.Email != "" {
							load.Pickup.Email = routeStop.Email
						}
						if routeStop.Contact != nil {
							load.Pickup.Contact = routeStop.Contact.Name
						}
					}
				} else {
					if load.Consignee != nil {
						if routeStop.Address != nil {
							load.Consignee.AddressLine1 = routeStop.Address.Line1
							load.Consignee.City = routeStop.Address.City
							load.Consignee.State = routeStop.Address.State
							load.Consignee.Zipcode = routeStop.Address.Zip
						}
						if routeStop.Phone != "" {
							load.Consignee.Phone = routeStop.Phone
						}
						if routeStop.Email != "" {
							load.Consignee.Email = routeStop.Email
						}
						if routeStop.Contact != nil {
							load.Consignee.Contact = routeStop.Contact.Name
						}
					}
				}
			}
		}
		// Map external IDs (PO numbers, ref numbers)
		if len(custOrder.ExternalIds) > 0 {
			poNums := []string{}
			for _, extId := range custOrder.ExternalIds {
				if !extId.Deleted {
					if extId.Type.Key == "1400" || extId.Type.Value == "Purchase shipment #" {
						poNums = append(poNums, extId.Value)
					} else if extId.Type.Key == "1401" || extId.Type.Value == "Reference Number" {
						if load.Customer != nil {
							load.Customer.RefNumber = extId.Value
						}
					}
				}
			}
			if len(poNums) > 0 {
				load.PoNums = strings.Join(poNums, ", ")
			}
		}
	}

	// Map billTo from Party
	if len(shipment.Party) > 0 {
		for _, party := range shipment.Party {
			if !party.Deleted && party.Account.Name != "" {
				if load.BillTo == nil {
					load.BillTo = &models.BillTo{}
				}
				load.BillTo.Name = party.Account.Name
				if party.Account.ID > 0 {
					load.BillTo.ExternalTMSId = fmt.Sprintf("%d", party.Account.ID)
				}
				break // Use first non-deleted party as billTo
			}
		}
	}

	// Map carrier order with detailed information
	if len(shipment.CarrierOrder) > 0 {
		for _, carrierOrder := range shipment.CarrierOrder {
			if !carrierOrder.Deleted {
				if load.Carrier == nil {
					load.Carrier = &models.Carrier{}
				}
				load.Carrier.Name = carrierOrder.Carrier.Name
				if carrierOrder.Carrier.ID > 0 {
					load.Carrier.ExternalTMSId = fmt.Sprintf("%d", carrierOrder.Carrier.ID)
				}
				// Map drivers if available
				if len(carrierOrder.Drivers) > 0 {
					driver := carrierOrder.Drivers[0]
					if driver.Context != nil {
						load.Carrier.FirstDriverName = driver.Context.Name
						if driver.Phone != nil {
							load.Carrier.FirstDriverPhone = driver.Phone.Number
						}
						if driver.Email != nil {
							load.Carrier.Email = driver.Email.Email
						}
					}
					if len(carrierOrder.Drivers) > 1 {
						driver2 := carrierOrder.Drivers[1]
						if driver2.Context != nil {
							load.Carrier.SecondDriverName = driver2.Context.Name
							if driver2.Phone != nil {
								load.Carrier.SecondDriverPhone = driver2.Phone.Number
							}
						}
					}
				}
				break // Use first non-deleted carrier order
			}
		}
	}

	// Map equipment to specifications and weight
	if len(shipment.Equipment) > 0 {
		equip := shipment.Equipment[0]
		if equip.Weight > 0 {
			load.TotalWeight = &equip.Weight
		}
		if equip.Temp != 0 {
			if load.Specifications == nil {
				load.Specifications = &models.Specifications{}
			}
			// Use temp as both min and max if only one value provided
			load.Specifications.MinTempFahrenheit = &equip.Temp
			load.Specifications.MaxTempFahrenheit = &equip.Temp
		}
	}

	// Map specifications from services in globalRoute
	if load.Specifications == nil {
		load.Specifications = &models.Specifications{}
	}
	for _, stop := range shipment.GlobalRoute {
		if stop.Deleted {
			continue
		}
		isPickup := stop.StopType.Value == "Pickup" || stop.StopType.Key == "1500"
		for _, service := range stop.Services {
			serviceValue := strings.ToLower(service.Value)
			if strings.Contains(serviceValue, "liftgate") {
				if isPickup {
					val := true
					load.Specifications.LiftgatePickup = &val
				} else {
					val := true
					load.Specifications.LiftgateDelivery = &val
				}
			}
			if strings.Contains(serviceValue, "inside") {
				if isPickup {
					val := true
					load.Specifications.InsidePickup = &val
				} else {
					val := true
					load.Specifications.InsideDelivery = &val
				}
			}
			if strings.Contains(serviceValue, "tarp") {
				val := true
				load.Specifications.Tarps = &val
			}
			if strings.Contains(serviceValue, "hazmat") {
				val := true
				load.Specifications.Hazmat = &val
			}
			if strings.Contains(serviceValue, "strap") {
				val := true
				load.Specifications.Straps = &val
			}
			if strings.Contains(serviceValue, "seal") {
				val := true
				load.Specifications.Seal = &val
			}
		}
	}

	// Map carrier external IDs (MC number, DOT number, etc.)
	if len(shipment.CarrierOrder) > 0 {
		for _, carrierOrder := range shipment.CarrierOrder {
			if !carrierOrder.Deleted && len(carrierOrder.ExternalIds) > 0 {
				if load.Carrier == nil {
					load.Carrier = &models.Carrier{}
				}
				for _, extId := range carrierOrder.ExternalIds {
					if !extId.Deleted {
						// Map based on external ID type
						extType := strings.ToLower(extId.Type.Value)
						if strings.Contains(extType, "truck") || extId.Type.Key == "7605" {
							load.Carrier.ExternalTMSTruckId = extId.Value
						} else if strings.Contains(extType, "trailer") || extId.Type.Key == "7606" {
							load.Carrier.ExternalTMSTrailerId = extId.Value
						} else if strings.Contains(extType, "bol") || extId.Type.Key == "7602" {
							// BOL number - could map to sealNumber or keep separate
							load.Carrier.SealNumber = extId.Value
						} else if strings.Contains(extType, "pro") || extId.Type.Key == "7603" {
							// PRO number - could map to a field if we have one
						}
					}
				}
			}
		}
	}

	// Map rate data from customerOrder costs
	if len(shipment.CustomerOrder) > 0 && !shipment.CustomerOrder[0].Deleted {
		custOrder := shipment.CustomerOrder[0]
		if custOrder.Costs != nil && !custOrder.Costs.Deleted {
			if load.RateData == nil {
				load.RateData = &models.RateData{}
			}
			if custOrder.Costs.TotalAmount > 0 {
				// Map total amount - might need to determine if it's customer or carrier rate
			}
		}
	}

	// Map rate data from carrierOrder costs
	if len(shipment.CarrierOrder) > 0 {
		for _, carrierOrder := range shipment.CarrierOrder {
			if !carrierOrder.Deleted && carrierOrder.Costs != nil && !carrierOrder.Costs.Deleted {
				if load.RateData == nil {
					load.RateData = &models.RateData{}
				}
				if carrierOrder.Costs.TotalAmount > 0 {
					// Map carrier costs
				}
			}
		}
	}

	// Map distance from customerOrder totalMiles
	if len(shipment.CustomerOrder) > 0 && !shipment.CustomerOrder[0].Deleted {
		if shipment.CustomerOrder[0].TotalMiles > 0 {
			load.RouteMiles = &shipment.CustomerOrder[0].TotalMiles
		}
	}

	return load
}
