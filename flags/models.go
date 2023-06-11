package flags

/*
A structure with variables to map with the flags that the application supports.

ForceUpdate	-	Force the update of the configurations with the values in the flags.

HostAddr	-	The ip address of the machine that will run the application.

ConnTimeout	-	The connection timeout value for all the simulators.

CSAddr		-	Central system ip address that the simulators will communicate to.

CSPort		-	Central system port the that simulators will communicate to.

GCFolder	-	Folder where the application general configuration files are located.

SCFolder	-	Folder where the simulators configuration files are located.

Language	-	The language used by the application.
*/
type Flags struct {
	ForceUpdate bool   // Force the update of the configurations with the values in the flags
	HostAddr    string // The ip address of the machine that will run the application.
	ConnTimeout int64  // The connection timeout value for all the simulators.
	CSAddr      string // Central system ip address that the simulators will communicate to.
	CSPort      string // Central system port the that simulators will communicate to.
	GCFolder    string // Folder where the application general configuration files are located.
	SCFolder    string // Folder where the simulators configuration files are located.
	Language    string // The language used by the application.
	APIAddr     string // The ip address/hostname to serve the http server
	APIPort     string // The port to serve the http server
}
