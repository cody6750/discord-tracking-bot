package handlers

var (
	currentStatus string = "not set"
)

func SetStatus(status string) string {
	currentStatus = status
	return currentStatus
}
