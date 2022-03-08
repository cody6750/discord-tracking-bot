package handlers

var (
	mediaPath string
)

//SetHandlerMediaPath used to pass in the media path from the trackingbot package to avoid circular dependency
func SetHandlerMediaPath(path string) {
	mediaPath = path
}
