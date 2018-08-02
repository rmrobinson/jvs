package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	"github.com/rmrobinson/jvs/service/building"
	"github.com/rmrobinson/jvs/service/building/pb"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
)

var (
	// ErrUnableToSetupBottlerocket is returned if the supplied bridge configuration fails to properly initialize bottlerocketImpl.
	ErrUnableToSetupBottlerocket = errors.New("unable to set up bottlerocketImpl")
	// ErrUnableToSetupMonopAmp is returned if the supplied bridge configuration fails to properly initialize monop.
	ErrUnableToSetupMonopAmp     = errors.New("unable to set up monoprice amp")
	// ErrBridgeConfigInvalid is returned if the supplied bridge configuration is invalid.
	ErrBridgeConfigInvalid = errors.New("bridge config invalid")

)

type rootConfig struct {
	Bridges []bridgeConfig `yaml:"bridges"`
}

type bridgeConfig struct {
	Address addrConfig `yaml:"address"`
	CachePath string `yaml:"cachePath"`
	Type string `yaml:"type"`
}

type addrConfig struct {
	USBPath string `yaml:"usbPath"`
	IPAddress string `yaml:"ipAddress"`
	Port int32 `yaml:"port"`
	Proto string `yaml:"proto"`
}

func (rc *rootConfig) toProto() ([]*pb.BridgeConfig, error) {
	var ret []*pb.BridgeConfig

	for _, bridge := range rc.Bridges {
		config := &pb.BridgeConfig{
			Name:      bridge.Type,
			CachePath: bridge.CachePath,
			Address:   &pb.Address{},
		}

		if len(bridge.Address.USBPath) > 0 {
			config.Address.Usb = &pb.Address_Usb{
				Path: bridge.Address.USBPath,
			}
		} else if len(bridge.Address.IPAddress) > 0 {
			config.Address.Ip = &pb.Address_Ip{
				Host: bridge.Address.IPAddress,
				Port: bridge.Address.Port,
			}
		}
		ret = append(ret, config)
	}

	return ret, nil
}

func main() {
	var (
		port      = flag.Int("port", 1337, "The port for the deviced process to listen on")
		configPath = flag.String("config", "", "The path to the config file")
	)

	flag.Parse()

	yamlFile, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Printf("Error opening config file '%s': %s\n", *configPath, err.Error())
		os.Exit(1)
	}

	config := rootConfig{}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Printf("Unable to parse config file '%s': %s\n", *configPath, err.Error())
		os.Exit(1)
	}

	hub := building.NewHub()

	bridgeConfigs, err := config.toProto()
	if err != nil {
		log.Printf("Unable to convert config file '%s': %s\n", *configPath, err.Error())
		os.Exit(1)
	}

	var toClose []io.Closer

	for _, bridgeConfig := range bridgeConfigs {
		log.Printf("Initializing module %s\n", bridgeConfig.Name)
		switch bridgeConfig.Name {
		case "br":
			br := &bottlerocketImpl{}
			if err := br.setup(bridgeConfig); err != nil {
				log.Printf("Unable to setup bottlerocketImpl: %s\n", err.Error())
				os.Exit(1)
			}
			toClose = append(toClose, br)
			hub.AddBridge(br.bridge, time.Second)
		case "monopamp":
			monop := &monopampImpl{}
			if err := monop.setup(bridgeConfig); err != nil {
				log.Printf("Unable to setup monoprice amp: %s\n", err.Error())
				os.Exit(1)
			}
			toClose = append(toClose, monop)
			hub.AddBridge(monop.bridge, time.Second)
		case "hue":
			hue := &hueImpl{}
			if err := hue.setup(bridgeConfig, hub); err != nil {
				log.Printf("Unable to setup hue: %s\n", err.Error())
				os.Exit(1)
			}
			toClose = append(toClose, hue)
			go hue.Run() // the bridges are added via Run(), not here.
		case "proxy":
			proxy := &proxyImpl{}
			if err := proxy.setup(bridgeConfig, hub); err != nil {
				log.Printf("Unable to setup proxyImpl: %s\n", err.Error())
				os.Exit(1)
			}
			toClose = append(toClose, proxy)
			// the proxyImpl setup adds itself to the hub
		default:
			log.Printf("Unsupported type (%s) specified, skipping\n", bridgeConfig.Name)
		}
	}

	connStr := fmt.Sprintf("%s:%d", "", *port)
	lis, err := net.Listen("tcp", connStr)
	if err != nil {
		log.Printf("Error initializing listener: %s\n", err.Error())
		os.Exit(1)
	}
	defer lis.Close()
	log.Printf("Listening on %s\n", connStr)

	api := building.NewAPI(hub)

	grpcServer := grpc.NewServer()
	pb.RegisterBridgeManagerServer(grpcServer, api)
	pb.RegisterDeviceManagerServer(grpcServer, api)
	grpcServer.Serve(lis)

	for _, impl := range toClose {
		impl.Close()
	}
}
