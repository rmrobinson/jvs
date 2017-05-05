package modules

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"faltung.ca/jvs/lib/proto-go"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"faltung.ca/jvs/service/commander"
	"strconv"
)

type handle struct {
	conn   *grpc.ClientConn
	device proto.DeviceManagerClient
	bridge proto.BridgeManagerClient
}

type DeviceCommand struct {
	commander.Command
	handlesLock sync.RWMutex
	handles     map[string]handle
}

func (dc *DeviceCommand) Name() string {
	return "device"
}

func (dc *DeviceCommand) Help() string {
	return "Device module commands"
}

func (dc *DeviceCommand) Execute(ctx context.Context, args []string) (string, error) {
	if len(args) < 1 {
		return "", errors.New("Missing argument. bridge|get|set supported")
	}

	switch args[0] {
	case "conn":
		return dc.handleConnCmd(ctx, args[1:])
	case "bridge":
		return dc.handleBridgeCmd(ctx, args[1:])
	case "get":
		return dc.handleGetCmd(ctx, args[1:])
	case "set":
		return dc.handleSetCmd(ctx, args[1:])
	}

	return "", errors.New("Unsupported command")
}

// device conn add <URL>
// device conn remove <URL>
// device conn list
func (dc *DeviceCommand) handleConnCmd(ctx context.Context, args []string) (string, error) {
	if len(args) < 1 {
		return "", errors.New("Missing argument. add|remove|list supported")
	}
	switch args[0] {
	case "add":
		if len(args) < 2 {
			return "", errors.New("Missing argument. URL required")
		}

		var err error
		var h handle
		var opts []grpc.DialOption
		opts = append(opts, grpc.WithInsecure())
		h.conn, err = grpc.Dial(args[1], opts...)

		if err != nil {
			return "", err
		}

		h.bridge = proto.NewBridgeManagerClient(h.conn)
		h.device = proto.NewDeviceManagerClient(h.conn)

		dc.handlesLock.Lock()
		defer dc.handlesLock.Unlock()

		if dc.handles == nil {
			dc.handles = map[string]handle{}
		}

		dc.handles[args[1]] = h
		return "Connected to " + args[1], nil

	case "remove":
		if len(args) < 2 {
			return "", errors.New("Missing argument. URL required")
		}

		dc.handlesLock.Lock()
		defer dc.handlesLock.Unlock()

		if h, ok := dc.handles[args[1]]; ok {
			h.conn.Close()
			delete(dc.handles, args[1])
		}

		return "Disconnected from " + args[1], nil

	case "list":
		dc.handlesLock.Lock()
		defer dc.handlesLock.Unlock()

		if len(dc.handles) < 1 {
			return "No connections", nil
		}

		ret := "Connected to: "
		for addr := range dc.handles {
			ret += addr + ", "
		}

		return strings.TrimSuffix(ret, ", "), nil
	}

	return "", errors.New("Not implemented)")
}

func (dc *DeviceCommand) handleBridgeCmd(ctx context.Context, args []string) (string, error) {
	return "", errors.New("Not implemented)")
}

// device get <Address>
func (dc *DeviceCommand) handleGetCmd(ctx context.Context, args []string) (string, error) {
	if len(args) < 1 {
		return "", errors.New("Missing argument. Device address required")
	}
	dc.handlesLock.RLock()
	defer dc.handlesLock.RUnlock()

	var wg sync.WaitGroup
	results := map[string]*proto.Device{}

	for id, h := range dc.handles {
		wg.Add(1)

		go func(id string, h handle) {
			defer wg.Done()

			req := &proto.GetDeviceRequest{
				Address: args[0],
			}
			resp, err := h.device.GetDevice(ctx, req)

			if err != nil {
				return
			}

			results[id] = resp.Device
		}(id, h)
	}

	wg.Wait()

	for _, result := range results {
		if result != nil {
			return fmt.Sprintf("%v\n", result), nil
		}
	}

	return "", errors.New("Device address not found")
}

// device set <Address> <config|state>
// if config:
//
// if state:
// isOn <true|false>
// level <int in range of min&max defined for device>
// colour <red between 0 and 255> <green between 0 and 255> <blue between 0 and 255>
func (dc *DeviceCommand) handleSetCmd(ctx context.Context, args []string) (string, error) {
	if len(args) < 1 {
		return "", errors.New("Missing argument. Device address required")
	}
	dc.handlesLock.RLock()
	defer dc.handlesLock.RUnlock()

	var wg sync.WaitGroup
	results := map[string]string{}

	for id, h := range dc.handles {
		wg.Add(1)

		go func(id string, h handle) {
			defer wg.Done()

			req := &proto.GetDeviceRequest{
				Address: args[0],
			}
			resp, err := h.device.GetDevice(ctx, req)

			// If we don't find the device here, that's cool, we assume another connection has it.
			if err != nil {
				return
			}

			device := resp.Device

			if len(args) < 2 {
				results[id] = "Missing argument. Device field to be set required"
				return
			}

			switch args[1] {
			case "config":
				results[id] = "Invalid argument. Config not supported"
				return

			case "state":
				var state proto.DeviceState
				state = *device.State

				if len(args) < 3 {
					results[id] = "Missing argument. Field to set required"
					return
				}

				switch args[2] {
				case "isOn":
					if state.Binary == nil {
						results[id] = "Invalid argument. Device does not support binary commands"
						return
					}

					if len(args) < 4 {
						results[id] = "Missing argument. isOn requires a flag to be set"
						return
					}

					if args[3] == "on" || args[3] == "ON" || args[3] == "On" || args[3] == "true" || args[3] == "yes" {
						state.Binary.IsOn = true
					} else {
						state.Binary.IsOn = false
					}

					break

				case "level":
					if state.Range == nil {
						results[id] ="Invalid argument. Device does not support range commands"
						return
					}

					if len(args) < 4 {
						results[id] = "Missing argument. level requires a value to be set"
						return
					}

					arg, convErr := strconv.Atoi(args[3])

					if convErr != nil {
						results[id] = "Invalid argument. Value unable to parse as int: " + convErr.Error()
						return
					}

					level := int32(arg)

					if level < device.Range.Minimum || level > device.Range.Maximum {
						results[id] = "Invalid argument. Level must be between " + strconv.Itoa(int(device.Range.Minimum)) + " and " + strconv.Itoa(int(device.Range.Maximum))
						return
					}

					state.Range.Value = level

					break

				case "colour":
					if state.ColorRgb == nil {
						results[id] = "Invalid argument. Device does not support colour commands"
						return
					}

					if len(args) < 6 {
						results[id] = "Invalid argument. rgb requires {red} {green} {blue} values"
						return
					}

					var red, green, blue int
					var convErr error
					red, convErr = strconv.Atoi(args[4])
					green, convErr = strconv.Atoi(args[5])
					blue, convErr = strconv.Atoi(args[6])

					if convErr != nil {
						results[id] = "Invalid argument. Colour field unable to parse as int: " + convErr.Error()
						return
					} else if red < 0 || red > 255 ||
						green < 0 || green > 255 ||
						blue < 0 || blue > 255 {
						results[id] = "Invalid argument. Color field must be between 0 and 255"
						return
					}

					state.ColorRgb.Red = int32(red)
					state.ColorRgb.Green = int32(green)
					state.ColorRgb.Blue = int32(blue)

					break

				default:
					results[id] = "Unknown argument"
					return
				}

				stateReq := &proto.SetDeviceStateRequest{
					Address: args[0],
					State: &state,
				}

				_, err := h.device.SetDeviceState(ctx, stateReq)

				if err != nil {
					results[id] = "Error setting: " + err.Error()
				} else {
					results[id] = "Successfully set"
				}

				return
			default:
				results[id] = "Unknown argument"
				return
			}

		}(id, h)
	}

	wg.Wait()

	for _, result := range results {
		if len(result) > 0 {
			return fmt.Sprintf("%v\n", result), nil
		}
	}

	return "", errors.New("Device address not found")
}
