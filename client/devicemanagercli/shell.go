package main

import (
	"fmt"
	"strconv"

	"faltung.ca/jvs/lib/proto-go"
	"github.com/abiosoft/ishell"
)

type shell struct {
	handle *ishell.Shell

	conns []conn
}

func newShell() shell {
	return shell{
		handle: ishell.New(),
		conns:  make([]conn, 8),
	}
}

func (s *shell) registerConnCommands() {
	s.handle.Register("conn", func(args ...string) (ret string, err error) {
		if len(args) < 1 {
			ret = "Invalid command or argument\nconn (add|remove|list) ..."
			return
		}

		switch args[0] {
		case "add":
			if len(args) < 2 {
				ret = "Invalid command or argument\nconn add {address}"
				return
			}

			c, err2 := connect(args[1])

			if err2 != nil {
				err = err2
				return
			} else {
				s.conns = append(s.conns, c)
			}

		case "remove":
			if len(args) < 2 {
				ret = "Invalid command or argument\nconn remove {address}"
				return
			}

			for idx, conn := range s.conns {
				if conn.addr != args[1] {
					continue
				}

				conn.close()

				// Remove from list
				s.conns = append(s.conns[:idx], s.conns[idx+1:]...)
			}

		case "list":
			ret = "Active connections\n"
			found := false
			for _, conn := range s.conns {
				if conn.conn != nil {
					if found {
						ret += "\n"
					}

					found = true
					ret += conn.addr
				}
			}

			if !found {
				ret += "(none)"
			}

		default:
			ret = "Unknown argument"
		}

		return
	})
}

func (s *shell) registerBridgeCommands() {
	s.handle.Register("bridge", func(args ...string) (ret string, err error) {
		if len(args) < 1 {
			ret = "Invalid command or argument\nbridge (list|watch) ..."
			return
		}

		switch args[0] {
		case "list":
			ret, err = listBridges(s.conns)

		case "watch":
			if len(args) < 2 {
				ret = "Invalid command or argument\nbridge watch (start|stop)"
				return
			}

			switch args[1] {
			case "start":
				watchBridges(s.conns)
				ret = "Ok"
			case "stop":
				stopWatchBridges(s.conns)
				ret = "Ok"
			default:
				ret = "Unknown argument"
			}

		default:
			ret = "Unknown argument"
		}

		return
	})
}

func (s *shell) registerDeviceCommands() {
	s.handle.Register("devices", func(args ...string) (ret string, err error) {
		if len(args) < 1 {
			ret = "Invalid command or argument\ndevices (list|watch) ..."
			return
		}

		switch args[0] {
		case "list":
			ret, err = listDevices(s.conns)

		case "watch":
			if len(args) < 2 {
				ret = "Invalid command or argument\ndevice watch (start|stop)"
				return
			}

			switch args[1] {
			case "start":
				watchDevices(s.conns)
				ret = "Ok"
			case "stop":
				stopWatchDevices(s.conns)
				ret = "Ok"
			default:
				ret = "Unknown argument"
			}

		default:
			ret = "Unknown argument"
		}

		return
	})

	s.handle.Register("device", func(args ...string) (ret string, err error) {
		if len(args) < 2 {
			ret = "Invalid command or argument\ndevice {id} (get|set) ..."
			return
		}

		id := args[0]
		switch args[1] {
		case "get":
			var d proto.Device
			d, err = getDevice(s.conns, id)

			if err != nil {
				return
			}

			ret += fmt.Sprintf("%v", d)

		case "set":
			if len(args) < 3 {
				ret = "Invalid command or argument\ndevice {id} set (config|state) ..."
				return
			}

			switch args[2] {
			case "state":
				if len(args) < 5 {
					ret = "Invalid command or argument\ndevice {id} set state (isOn|level|colour) {value}..."
					return
				}

				var device proto.Device
				device, err = getDevice(s.conns, id)

				if err != nil {
					return
				}

				var state proto.DeviceState
				state = *device.State

				switch args[3] {
				case "isOn":
					if state.Binary == nil {
						ret = "Invalid argument for device, does not support binary commands"
						return
					}

					if args[4] == "on" || args[4] == "ON" || args[4] == "On" || args[4] == "true" || args[4] == "yes" {
						state.Binary.IsOn = true
					} else {
						state.Binary.IsOn = false
					}

					break

				case "level":
					if state.Range == nil {
						ret = "Invalid argument for device, does not support range commands"
						return
					}

					arg, convErr := strconv.Atoi(args[4])

					if convErr != nil {
						ret = "Invalid argument for level, unable to parse as int: " + convErr.Error()
						return
					}

					level := int32(arg)

					if level < device.Range.Minimum || level > device.Range.Maximum {
						ret = "Invalid argument for level, must be between " + strconv.Itoa(int(device.Range.Minimum)) + " and " + strconv.Itoa(int(device.Range.Maximum))
						return
					}

					state.Range.Value = level

					break

				case "colour":
					if state.ColorRgb == nil {
						ret = "Invalid argument for device, does not support colour commands"
						return
					}

					if len(args) < 7 {
						ret = "Invalid command or argument\ndevice {id} set state color {red} {green} {blue}"
						return
					}

					var red, green, blue int
					var convErr error
					red, convErr = strconv.Atoi(args[4])
					green, convErr = strconv.Atoi(args[5])
					blue, convErr = strconv.Atoi(args[6])

					if convErr != nil {
						ret = "Invalid argument for colour, unable to parse as int: " + convErr.Error()
						return
					} else if red < 0 || red > 255 ||
						green < 0 || green > 255 ||
						blue < 0 || blue > 255 {
						ret = "Invalid argument for colour, values must be between 0 and 255"
						return
					}

					state.ColorRgb.Red = int32(red)
					state.ColorRgb.Green = int32(green)
					state.ColorRgb.Blue = int32(blue)

					break

				default:
					ret = "Unknown argument"
					return
				}

				_, err = setDeviceState(s.conns, id, state)

				ret = "Ok"

			default:
				ret = "Unknown argument"
			}

		default:
			ret = "Unknown argument"
		}

		return
	})
}

func (s *shell) run() {
	s.handle.SetHomeHistoryPath(".dm_history")

	s.handle.Println("-- JVS CLI --")

	s.registerConnCommands()
	s.registerBridgeCommands()
	s.registerDeviceCommands()

	s.handle.Start()
}
