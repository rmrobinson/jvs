package deviceclient_go

import (
	"errors"
	"faltung.ca/jvs/lib/proto-go"
	"fmt"
	"golang.org/x/net/context"
	"io"
	"log"
)

func ListBridges(conns []Conn) (ret string, err error) {
	ret = ""

	for _, conn := range conns {
		if conn.Conn == nil {
			continue
		}

		resp, err2 := conn.BridgeClient.GetBridges(context.Background(), &proto.GetBridgesRequest{})
		if err2 != nil {
			err = err2
			return
		}

		for _, bridge := range resp.Bridges {
			ret += fmt.Sprintf("%v\n", bridge)
		}
	}

	return
}

func WatchBridges(conns []Conn) {
	for idx, conn := range conns {
		if conn.Conn == nil || conn.cancelBridgeWatcher != nil {
			continue
		}

		cancelConn, cancelConnFunc := context.WithCancel(context.Background())
		conns[idx].cancelBridgeWatcher = cancelConnFunc

		fmt.Printf("Watching bridges on %s\n", conn.Addr)
		stream, err := conn.BridgeClient.WatchBridges(cancelConn, &proto.WatchBridgesRequest{})

		if err != nil {
			return
		}

		go func() {
			for {
				msg, err := stream.Recv()
				if err == io.EOF {
					break
				} else if err != nil {
					log.Printf("Error while watching bridges: %v", err)
					break
				}

				log.Printf("Change: %v, Bridge: %+v\n", msg.Action, msg.Bridge)
			}
		}()
	}
}

func StopWatchBridges(conns []Conn) {
	for idx, conn := range conns {
		if conn.Conn == nil || conn.cancelBridgeWatcher == nil {
			continue
		}

		conn.cancelBridgeWatcher()
		conns[idx].cancelBridgeWatcher = nil
	}
}

func ListDevices(conns []Conn) (ret string, err error) {
	ret = ""

	for _, conn := range conns {
		if conn.Conn == nil {
			continue
		}

		resp, err2 := conn.DeviceClient.GetDevices(context.Background(), &proto.GetDevicesRequest{})
		if err2 != nil {
			err = err2
			return
		}

		for _, device := range resp.Devices {
			ret += fmt.Sprintf("%v\n", device)
		}
	}

	return
}

func WatchDevices(conns []Conn) {
	for idx, conn := range conns {
		if conn.Conn == nil || conn.cancelDeviceWatcher != nil {
			continue
		}

		cancelConn, cancelConnFunc := context.WithCancel(context.Background())
		conns[idx].cancelDeviceWatcher = cancelConnFunc

		fmt.Printf("Watching devices on %s\n", conn.Addr)
		stream, err := conn.DeviceClient.WatchDevices(cancelConn, &proto.WatchDevicesRequest{})

		if err != nil {
			return
		}

		go func() {
			for {
				msg, err := stream.Recv()
				if err == io.EOF {
					break
				} else if err != nil {
					log.Printf("Error while watching devices: %v", err)
					break
				}

				log.Printf("Change: %v, Device: %+v\n", msg.Action, msg.Device)
			}
		}()
	}
}

func StopWatchDevices(conns []Conn) {
	for idx, conn := range conns {
		if conn.Conn == nil || conn.cancelDeviceWatcher == nil {
			continue
		}

		conn.cancelDeviceWatcher()
		conns[idx].cancelDeviceWatcher = nil
	}
}

func GetDevice(conns []Conn, id string) (ret proto.Device, err error) {
	ret = proto.Device{}

	for _, conn := range conns {
		if conn.Conn == nil {
			continue
		}

		resp, err2 := conn.DeviceClient.GetDevice(context.Background(), &proto.GetDeviceRequest{
			Id: id,
		})
		if err2 != nil {
			log.Printf("Error: %s\n", err2)
			err = err2
			return
		}

		ret = *resp.Device
		return
	}

	return
}

func SetDeviceState(conns []Conn, id string, state proto.DeviceState) (ret proto.Device, err error) {
	ret = proto.Device{}

	for _, conn := range conns {
		if conn.conn == nil {
			continue
		}

		resp, err2 := conn.DeviceClient.SetDeviceState(context.Background(), &proto.SetDeviceStateRequest{
			Id:    id,
			State: &state,
		})

		if err2 != nil {
			err = err2
			return
		} else if len(resp.Error) > 0 {
			err = errors.New(resp.Error)
			return
		}

		ret = *resp.Device
	}

	return
}
