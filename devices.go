package modbusdev

import (
	"fmt"
	"strings"
)

var sdm230 = map[int]Register{
	30001: {"Line to neutral volts", "V", 0x0000, "ieee32", 1},
	30007: {"Current", "A", 0x0006, "ieee32", 1},
	30013: {"Active Power", "W", 0x000C, "ieee32", 1},
	30019: {"Apparent Power", "VA", 0x0012, "ieee32", 1},
	30025: {"Reactive Power", "VAr", 0x0018, "ieee32", 1},
	30031: {"Power Factor", "", 0x001E, "ieee32", 1},
	30037: {"Phase Angle", "Degrees", 0x0024, "ieee32", 1},
	30071: {"Frequency", "Hz", 0x0046, "ieee32", 1},
	30073: {"Import Active Energy", "kWh", 0x0048, "ieee32", 1},
	30075: {"Export Active Energy", "kWh", 0x004A, "ieee32", 1},
	30077: {"Import Reactive Energy", "kVArh", 0x004C, "ieee32", 1},
	30079: {"Export Reactive Energy", "kVArh", 0x004E, "ieee32", 1},
	30085: {"Total system power demand", "W", 0x0054, "ieee32", 1},
	30087: {"Maximum total system power demand", "W", 0x0056, "ieee32", 1},
	30089: {"Current system positive power demand", "W", 0x0058, "ieee32", 1},
	30091: {"Maximum system positive power demand", "W", 0x005A, "ieee32", 1},
	30093: {"Current system reverse power demand", "W", 0x005C, "ieee32", 1},
	30095: {"Maximum system reverse power demand", "W", 0x005E, "ieee32", 1},
	30259: {"Current demand", "Amps", 0x0102, "ieee32", 1},
	30265: {"Maximum current Demand", "A", 0x0108, "ieee32", 1},
	30343: {"Total Active Energy", "kWh", 0x0156, "ieee32", 1},
	30345: {"Total Reactive Energy", "kVArh", 0x0158, "ieee32", 1},
}

// Additional registers that may be of interest to some.
var sdm230Ex = map[int]Register{
	// Included as example in protocol document?
	//	40001:  {"Demand Time", "ms", 0x0000, "ieee32", 1},
	40013:  {"Relay Pulse Width", "ms", 0x000C, "ieee32", 1},
	40019:  {"Network Parity Stop", "", 0x0012, "ieee32", 1},
	40021:  {"Network Node", "", 0x0014, "ieee32", 1},
	40029:  {"Network Baud Rate", "", 0x001c, "ieee32", 1},
	462721: {"Screen Settings", "", 0xf500, "u32", 1},
	463761: {"System Power", "", 0xf910, "u32", 1},
	463776: {"Measurement Mode", "", 0xf91f, "u32", 1},
	463792: {"Pulse Indicators", "", 0xf92f, "u32", 1},
}

/*
 * Solax register information from
 * https://github.com/wills106/homeassistant-config/blob/43365e6eed685e82763f786e7a46c387083a93b5/packages/solax.yaml
 */
var solaxX1Hybrid = map[int]Register{
	30001: {"Grid Voltage", "V", 0, "u16", 0.1},
	30002: {"Grid Current", "A", 0x01, "s16", 0.1},
	30003: {"Inverter Power", "W", 0x02, "s16", 1},
	30004: {"PV1 Voltage", "V", 0x03, "u16", 0.1},
	30005: {"PV2 Voltage", "V", 0x04, "u16", 0.1},
	30006: {"PV1 Current", "A", 0x05, "u16", 0.1},
	30007: {"PV2 Current", "A", 0x06, "u16", 0.1},
	30008: {"Grid Frequency", "Hz", 0x07, "u16", .01},
	30009: {"Inner Temp", "C", 0x08, "s16", 1},
	// 0 - waiting, 1 - checking, 2 - normal, 3 - off, 7 - eps, 9 - idle
	30010: {"Run Mode", "", 0x09, "u16", 1},
	30011: {"PV1 Power", "W", 0x0a, "u16", 1},
	30012: {"PV2 Power", "W", 0x0b, "u16", 1},
	30021: {"Battery Voltage", "V", 0x14, "s16", .1},
	30022: {"Battery Current", "A", 0x15, "s16", .1},
	30023: {"Battery Power", "W", 0x16, "s16", 1},
	30024: {"Charger Board Temperature", "C", 0x17, "s16", 1},
	30025: {"Battery Temperature", "C", 0x18, "s16", 1},
	30026: {"Charger Boost Temperature", "C", 0x19, "s16", 1},
	30029: {"Battery Capacity", "%", 0x1C, "u16", 1},
	30030: {"Battery Energy Charged", "W", 0x1D, "u32", 1},
	30032: {"BMS Warning", "", 0x1F, "u16", 1},
	30033: {"Battery Energy Discharged", "W", 0x20, "u32", 1},
	// ???
	30036: {"Battery State of Health", "", 0x23, "u16", 1},
	30065: {"Inverter Fault", "", 0x40, "u32", 1},
	30067: {"Charger Fault", "", 0x42, "u16", 1},
	// 512 when meter fault present
	30068: {"Manager Fault", "", 0x43, "u16", 1},
	30071: {"Measured Power", "W", 0x46, "s32", .001},
	30073: {"Feed In Energy", "kWh", 0x48, "u32", .01},
	30075: {"Consumed Energy", "kWh", 0x4A, "u32", .01},
	30077: {"EPS Voltage", "V", 0x4C, "u16", .1},
	30078: {"EPS Current", "A", 0x4D, "u16", .1},
	30079: {"EPS VA", "VA", 0x4E, "u16", .1},
	30080: {"EPS Frequency", "Hz", 0x4F, "u16", 1},
	30081: {"Energy Today", "kW", 0x50, "u16", .1},
	30082: {"Energy Total", "kW", 0x51, "u32", .001},
}

// Additional registers that may be of interest to some.
var solaxX1HybridEx = map[int]Register{
	// The following registers can be read to give the described values,
	// but writing to the holding registers requires different information?
	// Advanced Grid Settings
	40026: {"Vac Lower", "V", 0x19, "u16", .1},
	40027: {"Vac Upper", "V", 0x1a, "u16", .1},
	40028: {"FEC Lower", "Hz", 0x1b, "u16", .01},
	40029: {"FEC Upper", "Hz", 0x1c, "u16", .01},
	40032: {"Vac 10M Avg", "V", 0x1f, "u16", .1},
	40033: {"Vac Lower Slow", "V", 0x20, "u16", .1},
	40034: {"Vac Upper Slow", "V", 0x21, "u16", .1},
	40035: {"FEC Lower Slow", "Hz", 0x22, "u16", .01},
	40036: {"FEC Upper Slow", "Hz", 0x23, "u16", .01},

	// Current Date & Time
	40135: {"Minutes", "", 0x86, "u16", 1},
	40136: {"Hours", "", 0x87, "u16", 1},
	40137: {"Day", "", 0x88, "u16", 1},
	40138: {"Month", "", 0x89, "u16", 1},
	40139: {"Year", "", 0x8A, "u16", 1},

	// Write as register 34
	40140: {"Min Charger Capacity", "%", 0x8C, "u16", 1},
	// Write as register 36
	40145: {"Charge Max Current", "A", 0x90, "u16", .1},
	// Write as register 37
	40146: {"Discharge Max Current", "A", 0x91, "u16", .1},

	// Times for Force Time Use
	40147: {"Charge Period 1 Start Hour", "", 0x92, "u16", 1},
	40148: {"Charge Period 1 Start Minutes", "", 0x93, "u16", 1},
	40149: {"Charge Period 1 Finish Hour", "", 0x94, "u16", 1},
	40150: {"Charge Period 1 Finish Minutes", "", 0x95, "u16", 1},
	40155: {"Charge Period 2 Start Hour", "", 0x9A, "u16", 1},
	40156: {"Charge Period 2 Start Minutes", "", 0x9B, "u16", 1},
	40157: {"Charge Period 2 Finish Hour", "", 0x9C, "u16", 1},
	40158: {"Charge Period 2 Finish Minutes", "", 0x9D, "u16", 1},

	// MAC Address is stored in 3 registers
	40163: {"MAC Address #1", "", 0xA2, "u16", 1},
	40164: {"MAC Address #2", "", 0xA3, "u16", 1},
	40165: {"MAC Address #3", "", 0xA4, "u16", 1},

	40183: {"Max Export Power", "W", 0xB6, "u16", 1},
	40187: {"Rated Power", "kW", 0xBA, "u16", .001},
	40223: {"Battery version number", "", 0xDE, "u16", .01},
	40225: {"Admin Password", "", 0xE0, "u16", 1},

	// Times when Work Mode set to Backup
	40255: {"Backup Start Hour", "", 0xFE, "u16", 1},
	40256: {"Backup Start Minute", "", 0xFF, "u16", 1},
	40257: {"Backup Finish Hour", "", 0x100, "u16", 1},
	40258: {"Backup finish Minute", "", 0x101, "u16", 1},

	// Modbus Information
	40265: {"Use Meter", "", 0x108, "u16", 1},
	40266: {"Meter 1 ID", "", 0x109, "u16", 1},
	40267: {"Meter 2 ID", "", 0x10A, "u16", 1},
}

// Not sure if there is a better way to do this, but it works for now.
func joinMaps(aaa, bbb map[int]Register) map[int]Register {
	regMap := make(map[int]Register)
	for k, v := range aaa {
		regMap[k] = v
	}
	for k, v := range bbb {
		regMap[k] = v
	}
	return regMap
}

// RegistersByName Given a device string, return the approrpriate map of registers.
func RegistersByName(device string) (registers map[int]Register, err error) {
	switch strings.ToLower(device) {
	case "sdm230":
		registers = sdm230
	case "sdm230ex":
		registers = joinMaps(sdm230, sdm230Ex)
	case "solaxx1hybrid":
		registers = solaxX1Hybrid
	case "solaxx1hybridex":
		registers = joinMaps(solaxX1Hybrid, solaxX1HybridEx)
	default:
		err = fmt.Errorf("Device '%s' is not known. Add the details and then update reader.go to include it", device)
	}
	return
}
