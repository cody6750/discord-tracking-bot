package handlers

var (
	currentStatus string = "not set"
)

//SetStatus used to pass in the status from the trackingbot package to avoid circular dependency
func SetStatus(status string) string {
	currentStatus = status
	return currentStatus
}
