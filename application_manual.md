# Electric Vehicle Charging Simulator Manual

[[_TOC_]]

## Introduction

## Installation
There are no installation file the application will do the necessery acions in the first execution to setup. There are a some files and folders that need to be in place to execute the application, for them to be created the user that executes the application need to have write permissions to the folder where the application file is located, and to the folder(s) where we want to locate the other files.

The app doesn't have any GUI but in the other hand has a **Restful API** that enables the user to control the application, and have some arguments that can be set during the execution of the app.

The first time the application is executed, and without any arguments, it will generate the `configs.json` file

using the endpoint /configs the configurations will be saved but the changes will only be set in place after the application is restarted

to update any configuation using the API the full body needs to be sent in the request