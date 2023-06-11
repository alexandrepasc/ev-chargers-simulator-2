package general

type Model struct {
	HostIP      string `json:"hostIp"`      // The IP address of the machine that will run the application.
	ConnTimeout int64  `json:"connTimeout"` // The connection timeout value for all the simulators.
	CSAddr      string `json:"cSAddr"`      // Central system ip address that the simulators will communicate to.
	CSPort      string `json:"cSPort"`      // Central system port the that simulators will communicate to.
	Lang        string `json:"lang"`        // Language used in the application
	APIAddr     string `json:"apiAddr"`     // Api server address
	APIPort     string `json:"apiPort"`     // Api server port
}
