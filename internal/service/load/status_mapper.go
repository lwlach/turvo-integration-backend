package load

import "strings"

// TurvoStatusToAPI maps Turvo status code/value to API status
func TurvoStatusToAPI(statusKey, statusValue string) string {
	// Map Turvo status codes to standardized API statuses
	switch statusKey {
	case "2100":
		return "quote_active"
	case "2101":
		return "tendered"
	case "2102":
		return "covered"
	case "2103":
		return "dispatched"
	case "2104":
		return "at_pickup"
	case "2105":
		return "en_route"
	case "2106":
		return "at_delivery"
	case "2107":
		return "delivered"
	case "2108":
		return "ready_for_billing"
	case "2109":
		return "processing"
	case "2110":
		return "carrier_paid"
	case "2111":
		return "customer_paid"
	case "2112":
		return "completed"
	case "2113":
		return "canceled"
	case "2114":
		return "quote_inactive"
	case "2115":
		return "picked_up"
	case "2116":
		return "route_complete"
	case "2117":
		return "tender_offered"
	case "2118":
		return "tender_accepted"
	case "2119":
		return "tender_rejected"
	case "2120":
		return "draft"
	case "2121":
		return "shipment_ready"
	case "2123":
		return "acquiring_location"
	case "2124":
		return "customs_hold"
	case "2125":
		return "arrived"
	case "2126":
		return "available"
	case "2127":
		return "out_gated"
	case "2129":
		return "in_gated"
	case "2131":
		return "arriving_to_port"
	case "2132":
		return "berthing"
	case "2133":
		return "unloading"
	case "2134":
		return "ramped"
	case "2135":
		return "deramped"
	case "2136":
		return "departed"
	case "2137":
		return "held"
	case "2138":
		return "out_for_delivery"
	case "2139":
		return "in_transshipment"
	case "2140":
		return "on_hold"
	case "2141":
		return "interline"
	default:
		// If no mapping found, return lowercase version of the value
		return normalizeStatus(statusValue)
	}
}

// APIToTurvoStatus maps API status to Turvo status code
func APIToTurvoStatus(apiStatus string) (string, string) {
	// Map API status to Turvo status code and value
	switch apiStatus {
	case "quote_active":
		return "2100", "Quote active"
	case "tendered":
		return "2101", "Tendered"
	case "covered":
		return "2102", "Covered"
	case "dispatched":
		return "2103", "Dispatched"
	case "at_pickup":
		return "2104", "At pickup"
	case "en_route":
		return "2105", "En route"
	case "at_delivery":
		return "2106", "At delivery"
	case "delivered":
		return "2107", "Delivered"
	case "ready_for_billing":
		return "2108", "Ready for billing"
	case "processing":
		return "2109", "Processing"
	case "carrier_paid":
		return "2110", "Carrier paid"
	case "customer_paid":
		return "2111", "Customer paid"
	case "completed":
		return "2112", "Completed"
	case "canceled":
		return "2113", "Canceled"
	case "quote_inactive":
		return "2114", "Quote inactive"
	case "picked_up":
		return "2115", "Picked up"
	case "route_complete":
		return "2116", "Route Complete"
	case "tender_offered":
		return "2117", "Tender - offered"
	case "tender_accepted":
		return "2118", "Tender - accepted"
	case "tender_rejected":
		return "2119", "Tender - rejected"
	case "draft":
		return "2120", "Draft"
	case "shipment_ready":
		return "2121", "Shipment Ready"
	case "acquiring_location":
		return "2123", "Acquiring Location"
	case "customs_hold":
		return "2124", "Customs Hold"
	case "arrived":
		return "2125", "Arrived"
	case "available":
		return "2126", "Available"
	case "out_gated":
		return "2127", "Out Gated"
	case "in_gated":
		return "2129", "In Gated"
	case "arriving_to_port":
		return "2131", "Arriving to Port"
	case "berthing":
		return "2132", "Berthing"
	case "unloading":
		return "2133", "Unloading"
	case "ramped":
		return "2134", "Ramped"
	case "deramped":
		return "2135", "Deramped"
	case "departed":
		return "2136", "Departed"
	case "held":
		return "2137", "Held"
	case "out_for_delivery":
		return "2138", "Out for Delivery"
	case "in_transshipment":
		return "2139", "In TransShipment"
	case "on_hold":
		return "2140", "On Hold"
	case "interline":
		return "2141", "Interline"
	case "pending":
		// Default for new loads
		return "2120", "Draft"
	default:
		// Try to find by matching value (case-insensitive)
		return findTurvoStatusByValue(apiStatus)
	}
}

// normalizeStatus converts a status string to a normalized format
func normalizeStatus(status string) string {
	// Convert to lowercase and replace spaces with underscores
	normalized := strings.ToLower(status)
	normalized = strings.ReplaceAll(normalized, " ", "_")
	normalized = strings.ReplaceAll(normalized, "-", "_")
	return normalized
}

// findTurvoStatusByValue tries to find Turvo status by matching value
func findTurvoStatusByValue(value string) (string, string) {
	// Common mappings for values that might come in different formats
	valueLower := strings.ToLower(value)

	statusMap := map[string][2]string{
		"quote active":      {"2100", "Quote active"},
		"tendered":          {"2101", "Tendered"},
		"covered":           {"2102", "Covered"},
		"dispatched":        {"2103", "Dispatched"},
		"at pickup":         {"2104", "At pickup"},
		"en route":          {"2105", "En route"},
		"at delivery":       {"2106", "At delivery"},
		"delivered":         {"2107", "Delivered"},
		"ready for billing": {"2108", "Ready for billing"},
		"processing":        {"2109", "Processing"},
		"carrier paid":      {"2110", "Carrier paid"},
		"customer paid":     {"2111", "Customer paid"},
		"completed":         {"2112", "Completed"},
		"canceled":          {"2113", "Canceled"},
		"quote inactive":    {"2114", "Quote inactive"},
		"picked up":         {"2115", "Picked up"},
		"route complete":    {"2116", "Route Complete"},
		"tender - offered":  {"2117", "Tender - offered"},
		"tender - accepted": {"2118", "Tender - accepted"},
		"tender - rejected": {"2119", "Tender - rejected"},
		"draft":             {"2120", "Draft"},
		"shipment ready":    {"2121", "Shipment Ready"},
	}

	if mapping, found := statusMap[valueLower]; found {
		return mapping[0], mapping[1]
	}

	// Default to Draft if not found
	return "2120", "Draft"
}
