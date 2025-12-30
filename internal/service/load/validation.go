package load

import (
	"fmt"
	"strconv"

	"github.com/lwlach/turvo-integration-backend/internal/models"
)

// ValidateTurvoShipment validates required fields for Turvo shipment creation
func ValidateTurvoShipment(shipment *models.TurvoShipmentCreate) error {
	var errors []string

	// Validate ltlShipment (must be explicitly set, but can be false)
	// This is a boolean field, so it's always set (defaults to false)

	// Validate startDate
	if shipment.StartDate.Date == "" {
		errors = append(errors, "startDate.date is required")
	}

	// Validate endDate
	if shipment.EndDate.Date == "" {
		errors = append(errors, "endDate.date is required")
	}

	// Validate lane
	if shipment.Lane.Start == "" {
		errors = append(errors, "lane.start is required")
	}
	if shipment.Lane.End == "" {
		errors = append(errors, "lane.end is required")
	}

	// Validate customerOrder
	if len(shipment.CustomerOrder) == 0 {
		errors = append(errors, "customerOrder is required")
	} else {
		for i, order := range shipment.CustomerOrder {
			if order.Customer.ID == 0 {
				errors = append(errors, fmt.Sprintf("customerOrder[%d].customer.id is required", i))
			}
			if order.CustomerOrderSourceID == 0 {
				errors = append(errors, fmt.Sprintf("customerOrder[%d].customerOrderSourceId is required", i))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation failed: %v", errors)
	}

	return nil
}

func ValidateLoad(load *models.Load) error {
	var errors []string

	// Validate customer (required for customerOrder)
	if load.Customer == nil {
		errors = append(errors, "customer is required")
	} else {
		// Customer.ExternalTMSId is required for Turvo API (maps to customerOrder.customer.id)
		if load.Customer.ExternalTMSId == "" {
			errors = append(errors, "customer.externalTMSId is required")
		}
	}

	// Validate pickup (required for lane.start and startDate)
	if load.Pickup == nil {
		errors = append(errors, "pickup is required")
	} else {
		// Pickup must have ReadyTime or ApptTime for startDate mapping
		if load.Pickup.ReadyTime == nil && load.Pickup.ApptTime == nil {
			errors = append(errors, "pickup.readyTime or pickup.apptTime is required (for startDate)")
		}

		// Pickup must have City/State or Name for lane.start mapping
		hasCityState := load.Pickup.City != "" && load.Pickup.State != ""
		hasName := load.Pickup.Name != ""
		if !hasCityState && !hasName {
			errors = append(errors, "pickup must have either (city and state) or name (for lane.start)")
		}
	}

	// Validate consignee (required for lane.end and endDate)
	if load.Consignee == nil {
		errors = append(errors, "consignee is required")
	} else {
		// Consignee must have ApptTime for endDate mapping
		if load.Consignee.ApptTime == nil {
			errors = append(errors, "consignee.apptTime is required (for endDate)")
		}

		// Consignee must have City/State or Name for lane.end mapping
		hasCityState := load.Consignee.City != "" && load.Consignee.State != ""
		hasName := load.Consignee.Name != ""
		if !hasCityState && !hasName {
			errors = append(errors, "consignee must have either (city and state) or name (for lane.end)")
		}
	}

	// Validate carrier (optional, but if provided, externalTMSId is required)
	if load.Carrier != nil {
		// Carrier.ExternalTMSId is required if carrier is provided (maps to carrierOrder.carrier.id)
		if load.Carrier.ExternalTMSId == "" {
			errors = append(errors, "carrier.externalTMSId is required when carrier is provided")
		} else {
			// Validate that externalTMSId can be parsed as an integer (required for Turvo API)
			if _, err := strconv.Atoi(load.Carrier.ExternalTMSId); err != nil {
				errors = append(errors, "carrier.externalTMSId must be a valid integer")
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation failed: %v", errors)
	}

	return nil
}
