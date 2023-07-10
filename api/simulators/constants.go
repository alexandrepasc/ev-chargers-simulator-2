package simulators

const (
	simulatorsEp   string = "/simulators"                // simulators endpoint
	simulatorsIDEp string = "/simulators/:id"            // simulators with the asset id endpoint
	simModelsEp    string = simulatorsEp + "/models"     // models endpoint
	simModelsIDEp  string = simulatorsEp + "/models/:id" // models with the model id endpoint
	runEp          string = simulatorsEp + "/run"        // start endpoint
	stopEp         string = simulatorsEp + "/stop"       // stop endpoint
	simStatusEp    string = simulatorsEp + "/status"     // Get running assets status
)
