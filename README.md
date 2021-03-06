# esb-bridge


## WORK IN PROGRESS
This is work in progress and not considered functional, as long as this notice exists

## Packages in this module

### internal/usbprotocol
Handles USB serial connection to the esb-bridge-fw (device running the firmware). Should not be used directly

### pkg/esbbridge
Communicates with the firmware to transfer ESB messages to targets, and to receive incoming ESB messages

### cmd/server
CLI tool that Provides an interface to the esb bridge over the network (TCP socket). This is necessary because only one process can access the USB serial port. Also, there is only one physical instance of this device connected to the PC running the server, but there may be several nodes distributed across the network which want to access the ESB devices.

### pkg/client
Talks to the server over TCP socket in order to send and receive ESB messages. This component can be used by end-point implementations, meaning packages that provide access to a class of ESB device (e.g. binary sensor, switch, light etc) or more general packages like a MQTT-to-esb-bridge

## Get it running

The server is the only component of this repository which is intended to run directly
Examples below: ESB bridge device on /dev/ttyACM0, server on port 9815

### Run Go executable of the server
```
$ go run cmd/server/main.go -d /dev/ttyACM0 -p 9815
```

### Docker
Build the Docker image
```
$ docker build . -t esb-bridge-server
```

Run the server
```
$ docker run -d --device /dev/ttyACM0 -p 9815:9815 --hostname esbbridgeserver esb-bridge-server
```
Note: The parameter `--hostname esbbridgeserver` is important, without it clients cannot connect to the server. This will be be resolved in a future version (hopefully)


## Limitations
* ESB connection parameters are fixed in esb-bridge firmware, cannot be changed

## License
If not mentioned otherwise in individual source files, the MIT license (see LINCESE file) is applicable