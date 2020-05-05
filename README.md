# modbusdev

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

This data was recorded after dark :-)

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

## Bugs & Improvements

Always happy to have bugs found. Even happier to have pull requests submitted :-) 

If it's useful to you and you want additional devices added, submit the pull request and I'll merge them in.