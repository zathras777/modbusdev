# modbusdev [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT) [![GoDoc](https://godoc.org/github.com/goburrow/modbus?status.svg)](https://godoc.org/github.com/zathras777/modbusdev)

This package grew from a desire to have a simplified way of accessing the various registers on modebus devices I have presently. The aim was to build on top of the [modbus](https://github.com/goburrow/modbus) package rather than try to replace it.

The package only contains data for devices I need, but others are easily added. Additionally at present the access is read-only but this will change as I start to interact with the devices in different ways.

When returning values, where possible the raw values are returned, but there is an option to also get factored values based on known factors for the device registers. Units for the values stored in registers can also be accessed.

## Numbering

In order to allow for an index of registers I chose to go retro and use the Modicon convention address numbering. This allows for the appropiate register access method to be used and gives a unique index for each configured register. Register 0 has number 1 in this system.

## Usage

To use the package,

1. Create the modbus client in the usual way.
2. Create an instance of a Reader using the NewReader function
3. Read the registers you are interested in.

## Example

```go
package main

import (
    "log"

    "github.com/goburrow/modbus"
    "github.com/zathras777/modbusdev"
)

func main() {
    client := modbus.TCPClient("192.168.1.100:502")
    solax, err := modbusdev.NewReader(client, "solaxx1hybrid")
    if err != nil {
        log.Fatal(err)
    }
    solax.Dump(true)
}
```

This sample output was done after dark :-)

```
  30001: Grid Voltage                                     0.00 V
  30002: Grid Current                                     0.00 A
  30003: Inverter Power                                   0.00 W
  30004: PV1 Voltage                                      0.00 V
  30005: PV2 Voltage                                      0.00 V
  30006: PV1 Current                                      0.00 A
  30007: PV2 Current                                      0.00 A
  30008: Grid Frequency                                   0.00 Hz
  30009: Inner Temp                                       0.00 C
  30010: Run Mode                                         0.00 
  30011: PV1 Power                                        0.00 W
  30012: PV2 Power                                        0.00 W
  30021: Battery Voltage                                  0.00 V
  ...
```

## Device List

I don't have many devices :-)

- [Eastron SDM230-Modbus Power Meter](http://www.eastrongroup.com/productsview/72.html)
- [Solax X1 Hybrid Inverter](https://www.solaxpower.com/single-phase-hybrid/)

## Simple Database Access

I've been using a PostgreSQL database, so have added a simple interface to allow for easier recording of data from Map() results across my projects that are using modbusdev.

Given this example JSON configuration,
```
{
    "database": {
        "Host": "localhost",
        "Port": 5432,
        "User": "dbusername",
        "Password": "dbuserpassword",
        "Name": "modbusdatabase"
        "Query": "INSERT INTO table (time%s) VALUES (NOW()%s)"
        "fields": [
            {"name": "pv1", "code": 30011},
            {"name": "pv2", "code": 30012},
            {"name": "inverter", "code": 30003},
            ...
        ]
    },
```

The supplied query should have 2 string placeholders (%s) which will be replaced with the field names and suitable query markers for the database query.

The first step is to import the configuration and decode the JSON.

```
...
struct jsonConfig struct {
    Database modbusdev.DatabaseConfig
}

func main() {
    client := modbus.TCPClient("192.168.1.100:502")
    solax, err := modbusdev.NewReader(client, "solaxx1hybrid")
    if err != nil {
        log.Fatal(err)
    }

	jsonCfg := jsonConfig{}
	jsonFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(jsonFile)
	err = decoder.Decode(&jsonCfg)
	jsonFile.Close()
	if err != nil {
		log.Fatal(err)
	}
```

With the database configuration parsed, we set up a loop to constantly read the data and insert into the database as follows,

```
    err = jsonCfg.Database.OpenDatabase()
    if err != nil {
        log.Fatal(err)
    }
    defer jsonCfg.Database.Close()

    for {
        mapData := solax.Map(true)
        if err := cfg.Database.Execute(mapData); err != nil {
			log.Print(err)
			break
		}
        time.Sleep(5 * time.Second)
    }
}

```

This is primarily written to simplify my home workflow so will likely not be useful for many folks!

## Bugs & Improvements

Always happy to have bugs found. Even happier to have pull requests submitted :-)

If it's useful to you and you want additional devices added, submit the pull request and I'll merge them in.
