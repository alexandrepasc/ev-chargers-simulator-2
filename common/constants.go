package common

const (
	Version           string = "0.3.2"             // Application version
	DefSimIP          string = "84.5.7.64"         // Default simulator ip address from the machine that is running the application
	DefTimeout        int64  = 70                  // Default simulators timeout
	DefCSIP           string = "iot-gate-imx8.lan" // Default central system ip address
	DefCSPort         string = "49443"             // Default central system port address
	defGSFolder       string = "/settings"         // Default general settings folder
	defSCFolder       string = "/simConf"          // Default simulators configurations folder
	FolderPermissions int    = 0o777               // Folder permissions used to the settings folders in octal
	FilePermissions   int    = 0o600               // File permissions used to the configuration files in octal
)

var (
	DefGSPath = GetThePath(defGSFolder) // Default general settings path
	DefSCPath = GetThePath(defSCFolder) // Default simulators configurations path
)
