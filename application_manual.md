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
		- [Setup Simulators](#setup-simulators)
			- [Manage Models](#manage-models)
			- [Manage Simulator](#manage-simulator)
			- [Run The Simulators](#run-the-simulators)
	- [Appendix](#appendix)
		- [Flags](#flags)

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
/configs/general
/configs/simulators

### Setup Simulators
#### Manage Models
/simulators/models
/simulators/models/:id
#### Manage Simulator
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
