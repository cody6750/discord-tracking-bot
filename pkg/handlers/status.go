package handlers

var (
	currentStatus string = "not set"
)

func EnableStatus(status string) string {
	currentStatus = status
	return currentStatus
}
