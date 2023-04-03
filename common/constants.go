package common

const (
	Version     string = "0.2.0"             // Application version
	DefSimIP    string = "84.5.7.64"         // Default simulator ip address from the machine that is running the application
	DefTimeout  int64  = 70                  // Default simulators timeout
	DefCSIP     string = "iot-gate-imx8.lan" // Default central system ip address
	DefCSPort   string = "49443"             // Default central system port address
	defGSFolder string = "/settigs"          // Default general settings folder
	defSCFolder string = "/simConf"          // Default simulators configurations folder
)

var (
	DefGSPath = getThePath(defGSFolder) // Default general settings path
	DefSCPath = getThePath(defSCFolder) // Default simulators configurations path
)
