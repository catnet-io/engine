package fingerprint

// DeviceType defines the category of a network device.
type DeviceType string

const (
	DeviceWorkstation DeviceType = "workstation"
	DeviceServer      DeviceType = "server"
	DeviceRouter      DeviceType = "router"
	DeviceIoT         DeviceType = "iot"
	DeviceMobile      DeviceType = "mobile"
	DeviceUnknown     DeviceType = "unknown"
)

// FingerprintResult holds the detected properties of a host.
type FingerprintResult struct {
	OS         string
	OSFamily   string // "windows", "linux", "macos", "unix", "unknown"
	DeviceType DeviceType
	Vendor     string
	Confidence int // 0-100
}
