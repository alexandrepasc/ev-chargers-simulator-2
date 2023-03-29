package flags

/*
A structure with variables to map with the flags that the application supports.

HostAddr	-	The ip address of the machine that will run the application.

ConnTimeout	-	The connection timeout value for all the simulators.

CSAddr		-	Central system ip address that the simulators will communicate to.

CSPort		-	Central system port the that simulators will communicate to.

ConfFolder	-	Folder where the simulators configuration files are located.
*/
type Flags struct {
	HostAddr    string // The ip address of the machine that will run the application.
	ConnTimeout int64  // The connection timeout value for all the simulators.
	CSAddr      string // Central system ip address that the simulators will communicate to.
	CSPort      string // Central system port the that simulators will communicate to.
	ConfFolder  string // Folder where the simulators configuration files are located.
}
