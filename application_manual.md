# Electric Vehicle Charging Simulator Manual

- [Electric Vehicle Charging Simulator Manual](#electric-vehicle-charging-simulator-manual)
  - [Introduction](#introduction)
  - [Installation](#installation)
    - [First Run Default Configurations](#first-run-default-configurations)
    - [First Run Define Configurations](#first-run-define-configurations)
    - [First Run](#first-run)
  - [Use the Application](#use-the-application)
    - [Change Configurations](#change-configurations)
      - [By Flags](#by-flags)
      - [By Restful API](#by-restful-api)
        - [/configs/general](#configsgeneral)
        - [/configs/simulators](#configssimulators)
    - [Setup Simulators](#setup-simulators)
      - [Manage Models](#manage-models)
        - [/simulators/models](#simulatorsmodels)
        - [/simulators/models/:id](#simulatorsmodelsid)
      - [Manage Simulator](#manage-simulator)
        - [/simulators](#simulators)
        - [/simulators/:id](#simulatorsid)
      - [Run The Simulators](#run-the-simulators)
  - [Appendix](#appendix)
    - [Flags](#flags)
    - [Asset Type (type)](#asset-type-type)
    - [Protocols (protocol)](#protocols-protocol)
    - [Phase Rotation (phaseRotation)](#phase-rotation-phaserotation)
    - [Current Type (curentType)](#current-type-curenttype)
    - [Charging States (chargingState)](#charging-states-chargingstate)
    - [Error Code (errorCode)](#error-code-errorcode)

|         |             |
| ------- | ----------- |
| Product | **v0.14.7** |
| Manual  | **Rev1**    |

## Introduction

## Installation
There are no installation file the application will do the necessery acions in the first execution to setup. There are a some files and folders that need to be in place to execute the application, for them to be created the user that executes the application need to have write permissions to the folder where the application file is located, and to the folder(s) where we want to locate the other files.

The app doesn't have any GUI but in the other hand has a **Restful API** that enables the user to control the application, and have some arguments that can be set during the execution of the app.

### First Run Default Configurations
The first time the application is executed, and without any arguments, it will generate the `configs.json` file, the `settings` folder and in it the `general.json` file, the `simConf` folder, and the `/simConf/models` folder. The previous files and folders are essential for the application to run and can't be moved or deleted. Some of this configurations are set using the default values defined in the application and can be changed on the first execution or during the usage of the application.

### First Run Define Configurations
To know the configurations that can be changes on the execution we can run the application with the flag *-h*, `./ev-chargers-simulator-linux-amd64 -h`, that will print the list of all the flags the are supported as all the default values.

If you whant to change any of the listed configurations ([Appendix](#appendix)) during the 1st execution, you can simply set the *flag* and the value after.

Change the **General configurations folder**:
- `./ev-chargers-simulator-linux-amd64 -gs /full/path/to/the/location`

Change the **Rest API http server port**:
- `./ev-chargers-simulator-linux-amd64 -ap 3333`

Change the **General and Simulators folder**:
- `./ev-chargers-simulator-linux-amd64 -gs /full/path/to/the/location -ss /full/path`

If only part of the *flags* are set in the 1st run the application will set the other configurations with the default values, as in the previous section, that are shown when the user execute the application using the *flag -h*.

### First Run
After the user start the first time the application it will have the same behaviour as the rest of the executions. It will start the *Restful API Server*, load the configurations, look for the the *models* and *simulators* configuration files and load them, and will wait for the user commands using the *API*.

## Use the Application
There are a couple of informations that we need to know to be able to be able to use this application to it's full potential, as how to change the configurations, create/edit/remove models and simulators, what is a model, etc... In this section we will dive into explaining how to use it, and the interfaces that are available to the user.

### Change Configurations
There are multiple configurations some of them are not changeable but others are, and the user has two ways to change them.

#### By Flags
As described in the [Installation](#installation) there are some *command line* arguments that can be used to define the configurations in the first execution, the full list of arguments can be found in the [Appendix Flags](#flags) section.

After the *installation* process we can still use the *flags* to change the configurations, but to do it we need to use an extra one `-i`. This is the argument that triggers the application to save the new information.

With this we can run the application and at the same time change the configurations and save them for the next execution. The next examples are the same the were written previously but with the `-i` *flag*, so the configurations are changed and saved.

Change the **General configurations folder**:
- `./ev-chargers-simulator-linux-amd64 -i -gs /full/path/to/the/location`

Change the **Rest API http server port**:
- `./ev-chargers-simulator-linux-amd64 -i -ap 3333`

Change the **General and Simulators folder**:
- `./ev-chargers-simulator-linux-amd64 -i -gs /full/path/to/the/location -ss /full/path`

#### By Restful API
Executing the application it will start an *HTTP* service that will enable changing the configurations. There are two endpoints that can be used to change the application configurations, the `/configs/general` and `/configs/simulators`.

##### /configs/general
With this endpoint we can get the current configuration as change them. The list of all the configurations that can be changed are listed down.
- `hostIp`: The IP address of the conputer that is running the application
- `connTimeout`: Connection timeout for all the simulators
- `cSAddr`: OCPP central system IP address
- `cSPort`: OCPP central system port that is leastning for the charge points
- `lang`: Language
- `apiAddr`: IP address that the *HTTP* service will run on
- `apiPort`: IP port that the *HTTP* service will be listening on

To list the current configurations values execute the request:
```
GET http://{apiAddr}:{apiPort}/configs/general
```

The response will return all the configurations listed previously in a *JSON* format.

The same format returned by the list request should be used to execute the update configurations request.
```
PUT http://{apiAddr}:{apiPort}/configs/general
{
    "hostIp": "0.0.0.0",
    "connTimeout": 70,
    "cSAddr": "localhost",
    "cSPort": "1234",
    "lang": "pt-PT",
    "apiAddr": "localhost",
    "apiPort": "8001"
}
```
**WARNING:** The update of any of the configurations and the update of the other endpoints needs to have the full body, could not be a partial udate. If only one field needs to be changed the full body needs to be sent.

##### /configs/simulators
In this endpoint we can change the path were the configuration files of the simulators will be located.

To get the current path we can use the:
```
GET http://{apiAddr}:{apiPort}/configs/simulators
```

To update as in the previous endpoint we can use the *JSON* response structure.
```
PUT http://{apiAddr}:{apiPort}/configs/simulators
{
    "simFolder": "/full/path/to/folder/"
}
```

### Setup Simulators
In this section it is described how to setup and use the simulators. In oposition to the previous section the only way to do this is interfacing with the *Restful API*.

#### Manage Models
The **Models** are used to add configurations to the simulators, for example the *OCPP Chargers* have information that are retreived by the central system that is related with the equipement. For the *OCPP Chargers* in the boot request the charger sends the brand, model, etc... and this type of information is configured in the **Model**. The properties are listed bellow:

- `id`: Model identifier
- `name`: Model name
- `type`: Type of asset that this config can be used (ocpp, modbus)
- `ocpp`: Ocpp structure
  - `serialNumb`: Equipment serial number
  - `model`: Equipment model name
  - `vendor`: Equipment vendor
  - `vendorId`: Vendor identifier
  - `fwVersion`: Firmware version
  - `meterSerialNumb`: Power meter serial number
  - `modem`: Modem information
    - `iccid`: SIM card identifier
    - `imsi`: International Mobile Subscriber Identity
- `modbus`: Modbus structure

To use a *model* it needs to be configured in the *simulator* configuration. One simulator can be associated with one *model* but one *model* can have multiple *simulators*.

##### /simulators/models
This endpoint gives the ability to list and create *models*.

To list the existing models use the following method:
```
GET http://{apiAddr}:{apiPort}/simulators/models
```

It will return the total of existing *models* and the list of them.
```
{
    "total": 1,
    "models": [
        {
            "id": "3dfb64f3-d26b-44f5-8b1a-c56b56fa1b15",
            "name": "model1",
            "type": "ocpp",
            "ocpp": {
                "serialNumb": "",
                "model": "",
                "vendor": "",
                "vendorId": "",
                "fwVersion": "",
                "meterSerialNumb": "",
                "modem": {
                    "iccid": "",
                    "imsi": ""
                }
            },
            "modbus": {}
        }
    ]
}
```

Create a new *model* with the same endpoint but now using the `POST` method:
```
POST http://{apiAddr}:{apiPort}/simulators/models
{
    "name": "model1",
    "type": "ocpp",
    "ocpp": {
        "serialNumb": "qweqwe",
        "model": "EV Charger 1",
        "vendor": "Simulator",
        "fwVersion": "1.0.0",
        "meterSerialNumb": "789asd",
        "modem": {
            "iccid": "111111111111",
            "imsi": "2222222222222"
        }
    }
}
```

The service will return the same body structure sent in the request, but with the `ID` generated.

The properties `name` and `type` are required, the `name` is a free field but the `type` is not. It can only have the values `ocpp` or `modbus`, and it will be used to handle the type of protocol the the associated simulators can have.

##### /simulators/models/:id
Besides listing and creating the *models*, the user has the ability to edit and delete a model. This endpoint has the same path as the previous one but with the `ID` (`:id`) at the end that will identify the *model* that will be edited or deleted.

To edit a model get the `ID`, add it to the end of the endpoint path, and the full body of the *model* with the edited values.
```
PUT http://{apiAddr}:{apiPort}/simulators/models/{model-identifier-here}
{
    "name": "model1",
    "type": "ocpp",
    "ocpp": {
        "serialNumb": "qweqwe",
        "model": "EV Charger 1",
        "vendor": "Simulator",
        "fwVersion": "1.0.0",
        "meterSerialNumb": "changed-value",
        "modem": {
            "iccid": "111111111111",
            "imsi": "2222222222222"
        }
    }
}
```

Now deleting a *model* use the following:
```
DELETE http://{apiAddr}:{apiPort}/simulators/models/{model-identifier-here}
```

#### Manage Simulator
Now looking into the core of this application, the **Simulators**. There is on endpoint that is used to manage the **Simulator** assets and their configurations. In here will explaine how to add, edit, delete, and configure the assets. All the configurable properties are listed bellow:

- `name`: Simulator name
- `type`: Type of the asset that the configuration will be used to (evc, pm)
- `protocol`: Protocol used by the asset
- `model`: Id of the configuration file set in the model's folder
- `port`: Communication ip port
- `cPId`: Charge point id to identify the unit (used in the ocpp protocol)
- `startCharging`: Set the asset to start charging behaviour by itself
- `phases`: Phases number
- `phaseRotation`: The asset phase rotation, if the asset is DC the value should be NotApplicable
- `curentType`: Type of current of the asset (AC or DC)
- `authorizeRemote`: Configurataion AuthorizeRemoteTxRequests
- `authList`: Enable or disable authorization local list
- `evses`: List of evses that the asset has
  - `id`: Evse identifier number
  - `connectors`: The list of connectors of the evse
    - `id`: Connector identifier number
	- `data`: The loop of data
    	- `duration`: The duration in seconds that the current data will be in place
        - `chargingState`: The state of charging for the current data
        - `errorCode`: Error code
        - `powerFactor`: The power factor
        - `power`: Power in W
        - `powerExport`: Power exported by the ev to the cp in W, vehicle to grid
        - `voltage`: Array of voltages each entry for each phase in V
        - `startSoC`: Start ev charge state in percentage
        - `endSoC`: End ev charge state in percentage

##### /simulators
This endpoint can be used to create and list all the *simulators*.

To list the existing *simulators* use the following method:
```
GET http://{apiAddr}:{apiPort}/simulators
```

It will return the total number of number of existing *simulators* and their list.
```
{
    "total": 1,
    "assets": [
        {
            "simId": "585e8e82-b296-49f7-aea7-f3478a0c3a51",
            "name": "sim1",
            "type": "evc",
            "protocol": "ocpp16",
            "model": "00000000-0000-0000-0000-000000000000",
            "cPId": "123789",
            "startCharging": false,
            "phases": 1,
            "phaseRotation": "NotApplicable",
            "curentType": "dc",
            "authorizeRemote": true,
            "authList": true,
            "evses": [
                {
                    "id": 1,
                    "connectors": [
                        {
                            "id": 1,
                            "omitempty": false,
                            "data": [
                                {
                                    "duration": 10,
                                    "chargingState": 1,
                                    "errorCode": 0,
                                    "powerFactor": 900,
                                    "power": 0,
                                    "powerExport": 0,
                                    "voltage": [
                                        230,
                                        0,
                                        0
                                    ],
                                    "startSoC": 0,
                                    "endSoC": 0
                                }
                            ]
                        }
                    ]
                }
            ]
        }
    ]
}
```

To create a new *simulator* use the same endpoint but with a new method and sending the body described down. There are a couple of rules that are required to be followed so the application can work normally, the `data` array field in the configuration needs to have the 1st position with the `Available` (1) state and the last position with `Finishing` (6). The `chargingState` available values are exposed in the [Charging States appendix](#charging-states). 
```
POST http://{apiAddr}:{apiPort}/simulators
{
    "name": "sim1",
    "type": "evc",
    "protocol": "ocpp16",
    "model": "00000000-0000-0000-0000-000000000000",
    "cPId": "123789",
    "startCharging": false,
    "phases": 1,
    "phaseRotation": "NotApplicable",
    "curentType": "dc",
    "authorizeRemote": true,
    "authList": true,
    "evses": [
        {
            "id": 1,
            "connectors": [
                {
                    "id": 1,
                    "omitempty": false,
                    "data": [
                        {
                            "duration": 10,
                            "chargingState": 1,
                            "errorCode": 0,
                            "powerFactor": 900,
                            "power": 0,
                            "powerExport": 0,
                            "voltage": [
                                230,
                                0,
                                0
                            ],
                            "startSoC": 0,
                            "endSoC": 0
                        }
                    ]
                }
            ]
        }
    ]
}
```

It will return the created simulator information as response.

##### /simulators/:id
With this endpoint we are able to edit the configuration of one simulator as able to delete one. To use it replace the `:id` by the identifier of the *simulator* that is present in the property `simId`.

To update a simulator use the following example, to update any of the properties the full *simulator* body.
```
PUT http://{apiAddr}:{apiPort}/simulators/{simulator-identifier-here}
{
    "name": "sim1",
    "type": "evc",
    "protocol": "ocpp16",
    "model": "00000000-0000-0000-0000-000000000000",
    "cPId": "123789",
    "startCharging": false,
    "phases": 1,
    "phaseRotation": "NotApplicable",
    "curentType": "dc",
    "authorizeRemote": true,
    "authList": true,
    "evses": [
        {
            "id": 1,
            "connectors": [
                {
                    "id": 1,
                    "omitempty": false,
                    "data": [
                        {
                            "duration": 10,
                            "chargingState": 1,
                            "errorCode": 0,
                            "powerFactor": 900,
                            "power": 0,
                            "powerExport": 0,
                            "voltage": [
                                230,
                                0,
                                0
                            ],
                            "startSoC": 0,
                            "endSoC": 0
                        }
                    ]
                }
            ]
        }
    ]
}
```

To delete a simulator use the *delete* method:
```
DELETE http://{apiAddr}:{apiPort}/simulators/{simulator-identifier-here}
```

NOTE: the simulator data needs to have the Available state in the 1st item and Finishing in the last item
/simulators
/simulators/:id
#### Run The Simulators
/run
/stop
/simulators/status

using the endpoint /configs the configurations will be saved but the changes will only be set in place after the application is restarted

to update any configuation using the API the full body needs to be sent in the request

## Appendix

### Flags
```
  -ai string
    	Address used to serve the api http server (default "localhost")
  -ap string
    	Port used to serve the api http server (default "8000")
  -csi string
    	Central system ip address (default "iot-gate-imx8.lan")
  -csp string
    	Central system port (default "49443")
  -gs string
    	General configuration folder, store the application general configurations (default "/home/alex/Downloads/ev-chargers-simulator-linux-amd64/build/settings")
  -i	Force the update of the configurations with the values in the flags
  -l string
    	Language used by the application [en-GB, pt-PT] (default "en-GB")
  -s string
    	Simulator host ip address (default "84.5.7.64")
  -ss string
    	Simulators configuration folder, store the simulators configurations (default "/home/alex/Downloads/ev-chargers-simulator-linux-amd64/build/simConf")
  -t int
    	Connection timeout (seconds) (default 70)
  -v	Return the current application version
```

### Asset Type (type)
- `evc`: Type electric vehicle charger
- `pm`: Type power meter (not supported)

### Protocols (protocol)
- `ocpp16`: Ocpp version 1.6
- `ocpp201`: Ocpp version 2.0.1 (not supported)
- `modbus`: Modbus protocol (not supported)

### Phase Rotation (phaseRotation)
- `NotApplicable`: Not applicable for dc chargers
- `Unknown`: Not able to retrieve the rotation
- `RST`: L1 L2 L3
- `RTS`: L1 L3 L2
- `SRT`: L2 L1 L3
- `STR`: L2 L3 L1
- `TRS`: L3 L1 L2
- `TSR`: L3 L2 L1

### Current Type (curentType)
- `ac`: Alternating current
- `dc`: Direct current

### Charging States (chargingState)
- `1`: Available
- `2`: Preparing
- `3`: Charging
- `4`: SuspendedEV
- `5`: SuspendedEVSE
- `6`: Finishing
- `7`: Reserved
- `8`: Unavailable
- `9`: Faulted

### Error Code (errorCode)
- `0`: NoError
- `1`: ConnectorLockFailure
- `2`: EVCommunicationError
- `3`: GroundFailure
- `4`: HighTemperature
- `5`: InternalError
- `6`: LocalListConflict
- `7`: OtherError
- `8`: OverCurrentFailure
- `9`: OverVoltage
- `10`: PowerMeterFailure
- `11`: PowerSwitchFailure
- `12`: ReaderFailure
- `13`: ResetFailure
- `14`: UnderVoltage
- `15`: WeakSignal
