# XRayInstance

### A simple XRay instance code to use xray-core

## Requirements

- GoLang v1.20

## How This Code Works

This code loads a JSON configuration file and passes it to InboundHandlerConfig and OutboundHandlerConfig structs using the LoadConfig function.

The ConfigLoader function takes a filename as input and returns a core.Config struct or an error if there was an issue opening or reading the file.

To utilize additional features, you will need to uncomment the imported packages in order to include them in the build.

---

The main function first calls ShowVersion which prints the version statement of the Xray instance.

Next, it calls ConfigLoader to load the JSON configuration file and create a new Xray instance based on the given configuration. If there is an error, it will log the error and exit.

Then it calls the Start function to start the Xray instance, which includes all registered features. If there is an error, it will log the error and exit.

The Stop function is called to close the Xray instance once the program exits.

Finally, the program will keep running until it receives an os Interrupt or syscall sigterm signal.

During execution, the program explicitly triggers GC to remove garbage from the configuration loading, and also frees up OS memory using debug.FreeOSMemory().

---

ShowVersion
The ShowVersion function prints the version statement of the Xray instance.

Stop
The Stop function is used to close the Xray instance.

Start
The Start function is used to start the Xray instance. It returns an error if there is an issue starting the instance.

## How to run

Simply build the code and run it with your own configuration :

```
CGO_ENABLED=0 go build -v -ldflags="-s -w" .
XRayInstance config.json
```
