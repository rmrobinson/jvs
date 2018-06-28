// Code generated by protoc-gen-go. DO NOT EDIT.
// source: bridge.proto

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	bridge.proto
	device.proto

It has these top-level messages:
	Address
	BridgeState
	BridgeConfig
	BridgeSwUpdate
	Bridge
	AddBridgeRequest
	AddBridgeResponse
	GetBridgesRequest
	GetBridgesResponse
	WatchBridgesRequest
	BridgeUpdate
	DeviceConfig
	DeviceState
	Device
	GetDevicesRequest
	GetDevicesResponse
	GetDeviceRequest
	GetDeviceResponse
	WatchDevicesRequest
	WatchDevicesResponse
	SetDeviceConfigRequest
	SetDeviceConfigResponse
	SetDeviceStateRequest
	SetDeviceStateResponse
*/
package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type BridgeType int32

const (
	BridgeType_Generic      BridgeType = 0
	BridgeType_Loopback     BridgeType = 1
	BridgeType_Proxy        BridgeType = 2
	BridgeType_Hue          BridgeType = 3
	BridgeType_Bottlerocket BridgeType = 4
	BridgeType_MonopriceAmp BridgeType = 5
)

var BridgeType_name = map[int32]string{
	0: "Generic",
	1: "Loopback",
	2: "Proxy",
	3: "Hue",
	4: "Bottlerocket",
	5: "MonopriceAmp",
}
var BridgeType_value = map[string]int32{
	"Generic":      0,
	"Loopback":     1,
	"Proxy":        2,
	"Hue":          3,
	"Bottlerocket": 4,
	"MonopriceAmp": 5,
}

func (x BridgeType) String() string {
	return proto.EnumName(BridgeType_name, int32(x))
}
func (BridgeType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type BridgeUpdate_Action int32

const (
	BridgeUpdate_ADDED   BridgeUpdate_Action = 0
	BridgeUpdate_CHANGED BridgeUpdate_Action = 1
	BridgeUpdate_REMOVED BridgeUpdate_Action = 2
)

var BridgeUpdate_Action_name = map[int32]string{
	0: "ADDED",
	1: "CHANGED",
	2: "REMOVED",
}
var BridgeUpdate_Action_value = map[string]int32{
	"ADDED":   0,
	"CHANGED": 1,
	"REMOVED": 2,
}

func (x BridgeUpdate_Action) String() string {
	return proto.EnumName(BridgeUpdate_Action_name, int32(x))
}
func (BridgeUpdate_Action) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{10, 0} }

type Address struct {
	Ip  *Address_Ip  `protobuf:"bytes,1,opt,name=ip" json:"ip,omitempty"`
	Usb *Address_Usb `protobuf:"bytes,2,opt,name=usb" json:"usb,omitempty"`
}

func (m *Address) Reset()                    { *m = Address{} }
func (m *Address) String() string            { return proto.CompactTextString(m) }
func (*Address) ProtoMessage()               {}
func (*Address) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Address) GetIp() *Address_Ip {
	if m != nil {
		return m.Ip
	}
	return nil
}

func (m *Address) GetUsb() *Address_Usb {
	if m != nil {
		return m.Usb
	}
	return nil
}

type Address_Ip struct {
	Host    string `protobuf:"bytes,1,opt,name=host" json:"host,omitempty"`
	Netmask string `protobuf:"bytes,2,opt,name=netmask" json:"netmask,omitempty"`
	Gateway string `protobuf:"bytes,3,opt,name=gateway" json:"gateway,omitempty"`
	Port    int32  `protobuf:"varint,4,opt,name=port" json:"port,omitempty"`
	ViaDhcp bool   `protobuf:"varint,10,opt,name=via_dhcp,json=viaDhcp" json:"via_dhcp,omitempty"`
}

func (m *Address_Ip) Reset()                    { *m = Address_Ip{} }
func (m *Address_Ip) String() string            { return proto.CompactTextString(m) }
func (*Address_Ip) ProtoMessage()               {}
func (*Address_Ip) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

func (m *Address_Ip) GetHost() string {
	if m != nil {
		return m.Host
	}
	return ""
}

func (m *Address_Ip) GetNetmask() string {
	if m != nil {
		return m.Netmask
	}
	return ""
}

func (m *Address_Ip) GetGateway() string {
	if m != nil {
		return m.Gateway
	}
	return ""
}

func (m *Address_Ip) GetPort() int32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *Address_Ip) GetViaDhcp() bool {
	if m != nil {
		return m.ViaDhcp
	}
	return false
}

type Address_Usb struct {
	Path string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
}

func (m *Address_Usb) Reset()                    { *m = Address_Usb{} }
func (m *Address_Usb) String() string            { return proto.CompactTextString(m) }
func (*Address_Usb) ProtoMessage()               {}
func (*Address_Usb) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 1} }

func (m *Address_Usb) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

type BridgeState struct {
	IsPaired bool                 `protobuf:"varint,1,opt,name=is_paired,json=isPaired" json:"is_paired,omitempty"`
	Version  *BridgeState_Version `protobuf:"bytes,100,opt,name=version" json:"version,omitempty"`
	Zigbee   *BridgeState_Zigbee  `protobuf:"bytes,110,opt,name=zigbee" json:"zigbee,omitempty"`
	Zwave    *BridgeState_Zwave   `protobuf:"bytes,111,opt,name=zwave" json:"zwave,omitempty"`
}

func (m *BridgeState) Reset()                    { *m = BridgeState{} }
func (m *BridgeState) String() string            { return proto.CompactTextString(m) }
func (*BridgeState) ProtoMessage()               {}
func (*BridgeState) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *BridgeState) GetIsPaired() bool {
	if m != nil {
		return m.IsPaired
	}
	return false
}

func (m *BridgeState) GetVersion() *BridgeState_Version {
	if m != nil {
		return m.Version
	}
	return nil
}

func (m *BridgeState) GetZigbee() *BridgeState_Zigbee {
	if m != nil {
		return m.Zigbee
	}
	return nil
}

func (m *BridgeState) GetZwave() *BridgeState_Zwave {
	if m != nil {
		return m.Zwave
	}
	return nil
}

type BridgeState_Version struct {
	Api string `protobuf:"bytes,1,opt,name=api" json:"api,omitempty"`
	Sw  string `protobuf:"bytes,2,opt,name=sw" json:"sw,omitempty"`
}

func (m *BridgeState_Version) Reset()                    { *m = BridgeState_Version{} }
func (m *BridgeState_Version) String() string            { return proto.CompactTextString(m) }
func (*BridgeState_Version) ProtoMessage()               {}
func (*BridgeState_Version) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 0} }

func (m *BridgeState_Version) GetApi() string {
	if m != nil {
		return m.Api
	}
	return ""
}

func (m *BridgeState_Version) GetSw() string {
	if m != nil {
		return m.Sw
	}
	return ""
}

type BridgeState_Zigbee struct {
	Channel int32 `protobuf:"varint,1,opt,name=channel" json:"channel,omitempty"`
}

func (m *BridgeState_Zigbee) Reset()                    { *m = BridgeState_Zigbee{} }
func (m *BridgeState_Zigbee) String() string            { return proto.CompactTextString(m) }
func (*BridgeState_Zigbee) ProtoMessage()               {}
func (*BridgeState_Zigbee) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 1} }

func (m *BridgeState_Zigbee) GetChannel() int32 {
	if m != nil {
		return m.Channel
	}
	return 0
}

type BridgeState_Zwave struct {
	HomeId string `protobuf:"bytes,1,opt,name=homeId" json:"homeId,omitempty"`
	Mode   string `protobuf:"bytes,2,opt,name=mode" json:"mode,omitempty"`
}

func (m *BridgeState_Zwave) Reset()                    { *m = BridgeState_Zwave{} }
func (m *BridgeState_Zwave) String() string            { return proto.CompactTextString(m) }
func (*BridgeState_Zwave) ProtoMessage()               {}
func (*BridgeState_Zwave) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 2} }

func (m *BridgeState_Zwave) GetHomeId() string {
	if m != nil {
		return m.HomeId
	}
	return ""
}

func (m *BridgeState_Zwave) GetMode() string {
	if m != nil {
		return m.Mode
	}
	return ""
}

type BridgeConfig struct {
	Name      string   `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Id        string   `protobuf:"bytes,2,opt,name=id" json:"id,omitempty"`
	Address   *Address `protobuf:"bytes,10,opt,name=address" json:"address,omitempty"`
	CachePath string   `protobuf:"bytes,11,opt,name=cache_path,json=cachePath" json:"cache_path,omitempty"`
	Timezone  string   `protobuf:"bytes,50,opt,name=timezone" json:"timezone,omitempty"`
}

func (m *BridgeConfig) Reset()                    { *m = BridgeConfig{} }
func (m *BridgeConfig) String() string            { return proto.CompactTextString(m) }
func (*BridgeConfig) ProtoMessage()               {}
func (*BridgeConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *BridgeConfig) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *BridgeConfig) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *BridgeConfig) GetAddress() *Address {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *BridgeConfig) GetCachePath() string {
	if m != nil {
		return m.CachePath
	}
	return ""
}

func (m *BridgeConfig) GetTimezone() string {
	if m != nil {
		return m.Timezone
	}
	return ""
}

type BridgeSwUpdate struct {
	IsAvailable bool   `protobuf:"varint,1,opt,name=is_available,json=isAvailable" json:"is_available,omitempty"`
	NotifyUser  bool   `protobuf:"varint,10,opt,name=notify_user,json=notifyUser" json:"notify_user,omitempty"`
	NotifyText  string `protobuf:"bytes,11,opt,name=notify_text,json=notifyText" json:"notify_text,omitempty"`
	NotifyUrl   string `protobuf:"bytes,12,opt,name=notify_url,json=notifyUrl" json:"notify_url,omitempty"`
}

func (m *BridgeSwUpdate) Reset()                    { *m = BridgeSwUpdate{} }
func (m *BridgeSwUpdate) String() string            { return proto.CompactTextString(m) }
func (*BridgeSwUpdate) ProtoMessage()               {}
func (*BridgeSwUpdate) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *BridgeSwUpdate) GetIsAvailable() bool {
	if m != nil {
		return m.IsAvailable
	}
	return false
}

func (m *BridgeSwUpdate) GetNotifyUser() bool {
	if m != nil {
		return m.NotifyUser
	}
	return false
}

func (m *BridgeSwUpdate) GetNotifyText() string {
	if m != nil {
		return m.NotifyText
	}
	return ""
}

func (m *BridgeSwUpdate) GetNotifyUrl() string {
	if m != nil {
		return m.NotifyUrl
	}
	return ""
}

type Bridge struct {
	Id               string        `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	IsActive         bool          `protobuf:"varint,2,opt,name=is_active,json=isActive" json:"is_active,omitempty"`
	Type             BridgeType    `protobuf:"varint,3,opt,name=type,enum=pb.BridgeType" json:"type,omitempty"`
	ModelId          string        `protobuf:"bytes,10,opt,name=model_id,json=modelId" json:"model_id,omitempty"`
	ModelName        string        `protobuf:"bytes,11,opt,name=model_name,json=modelName" json:"model_name,omitempty"`
	ModelDescription string        `protobuf:"bytes,12,opt,name=model_description,json=modelDescription" json:"model_description,omitempty"`
	Manufacturer     string        `protobuf:"bytes,13,opt,name=manufacturer" json:"manufacturer,omitempty"`
	IconUrl          []string      `protobuf:"bytes,20,rep,name=icon_url,json=iconUrl" json:"icon_url,omitempty"`
	Config           *BridgeConfig `protobuf:"bytes,100,opt,name=config" json:"config,omitempty"`
	State            *BridgeState  `protobuf:"bytes,101,opt,name=state" json:"state,omitempty"`
}

func (m *Bridge) Reset()                    { *m = Bridge{} }
func (m *Bridge) String() string            { return proto.CompactTextString(m) }
func (*Bridge) ProtoMessage()               {}
func (*Bridge) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *Bridge) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Bridge) GetIsActive() bool {
	if m != nil {
		return m.IsActive
	}
	return false
}

func (m *Bridge) GetType() BridgeType {
	if m != nil {
		return m.Type
	}
	return BridgeType_Generic
}

func (m *Bridge) GetModelId() string {
	if m != nil {
		return m.ModelId
	}
	return ""
}

func (m *Bridge) GetModelName() string {
	if m != nil {
		return m.ModelName
	}
	return ""
}

func (m *Bridge) GetModelDescription() string {
	if m != nil {
		return m.ModelDescription
	}
	return ""
}

func (m *Bridge) GetManufacturer() string {
	if m != nil {
		return m.Manufacturer
	}
	return ""
}

func (m *Bridge) GetIconUrl() []string {
	if m != nil {
		return m.IconUrl
	}
	return nil
}

func (m *Bridge) GetConfig() *BridgeConfig {
	if m != nil {
		return m.Config
	}
	return nil
}

func (m *Bridge) GetState() *BridgeState {
	if m != nil {
		return m.State
	}
	return nil
}

type AddBridgeRequest struct {
	Type   BridgeType    `protobuf:"varint,1,opt,name=type,enum=pb.BridgeType" json:"type,omitempty"`
	Config *BridgeConfig `protobuf:"bytes,2,opt,name=config" json:"config,omitempty"`
}

func (m *AddBridgeRequest) Reset()                    { *m = AddBridgeRequest{} }
func (m *AddBridgeRequest) String() string            { return proto.CompactTextString(m) }
func (*AddBridgeRequest) ProtoMessage()               {}
func (*AddBridgeRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *AddBridgeRequest) GetType() BridgeType {
	if m != nil {
		return m.Type
	}
	return BridgeType_Generic
}

func (m *AddBridgeRequest) GetConfig() *BridgeConfig {
	if m != nil {
		return m.Config
	}
	return nil
}

type AddBridgeResponse struct {
	Bridge *Bridge `protobuf:"bytes,1,opt,name=bridge" json:"bridge,omitempty"`
}

func (m *AddBridgeResponse) Reset()                    { *m = AddBridgeResponse{} }
func (m *AddBridgeResponse) String() string            { return proto.CompactTextString(m) }
func (*AddBridgeResponse) ProtoMessage()               {}
func (*AddBridgeResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *AddBridgeResponse) GetBridge() *Bridge {
	if m != nil {
		return m.Bridge
	}
	return nil
}

type GetBridgesRequest struct {
}

func (m *GetBridgesRequest) Reset()                    { *m = GetBridgesRequest{} }
func (m *GetBridgesRequest) String() string            { return proto.CompactTextString(m) }
func (*GetBridgesRequest) ProtoMessage()               {}
func (*GetBridgesRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

type GetBridgesResponse struct {
	Bridges []*Bridge `protobuf:"bytes,1,rep,name=bridges" json:"bridges,omitempty"`
}

func (m *GetBridgesResponse) Reset()                    { *m = GetBridgesResponse{} }
func (m *GetBridgesResponse) String() string            { return proto.CompactTextString(m) }
func (*GetBridgesResponse) ProtoMessage()               {}
func (*GetBridgesResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *GetBridgesResponse) GetBridges() []*Bridge {
	if m != nil {
		return m.Bridges
	}
	return nil
}

type WatchBridgesRequest struct {
}

func (m *WatchBridgesRequest) Reset()                    { *m = WatchBridgesRequest{} }
func (m *WatchBridgesRequest) String() string            { return proto.CompactTextString(m) }
func (*WatchBridgesRequest) ProtoMessage()               {}
func (*WatchBridgesRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

type BridgeUpdate struct {
	Action BridgeUpdate_Action `protobuf:"varint,1,opt,name=action,enum=pb.BridgeUpdate_Action" json:"action,omitempty"`
	Bridge *Bridge             `protobuf:"bytes,2,opt,name=bridge" json:"bridge,omitempty"`
}

func (m *BridgeUpdate) Reset()                    { *m = BridgeUpdate{} }
func (m *BridgeUpdate) String() string            { return proto.CompactTextString(m) }
func (*BridgeUpdate) ProtoMessage()               {}
func (*BridgeUpdate) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *BridgeUpdate) GetAction() BridgeUpdate_Action {
	if m != nil {
		return m.Action
	}
	return BridgeUpdate_ADDED
}

func (m *BridgeUpdate) GetBridge() *Bridge {
	if m != nil {
		return m.Bridge
	}
	return nil
}

func init() {
	proto.RegisterType((*Address)(nil), "pb.Address")
	proto.RegisterType((*Address_Ip)(nil), "pb.Address.Ip")
	proto.RegisterType((*Address_Usb)(nil), "pb.Address.Usb")
	proto.RegisterType((*BridgeState)(nil), "pb.BridgeState")
	proto.RegisterType((*BridgeState_Version)(nil), "pb.BridgeState.Version")
	proto.RegisterType((*BridgeState_Zigbee)(nil), "pb.BridgeState.Zigbee")
	proto.RegisterType((*BridgeState_Zwave)(nil), "pb.BridgeState.Zwave")
	proto.RegisterType((*BridgeConfig)(nil), "pb.BridgeConfig")
	proto.RegisterType((*BridgeSwUpdate)(nil), "pb.BridgeSwUpdate")
	proto.RegisterType((*Bridge)(nil), "pb.Bridge")
	proto.RegisterType((*AddBridgeRequest)(nil), "pb.AddBridgeRequest")
	proto.RegisterType((*AddBridgeResponse)(nil), "pb.AddBridgeResponse")
	proto.RegisterType((*GetBridgesRequest)(nil), "pb.GetBridgesRequest")
	proto.RegisterType((*GetBridgesResponse)(nil), "pb.GetBridgesResponse")
	proto.RegisterType((*WatchBridgesRequest)(nil), "pb.WatchBridgesRequest")
	proto.RegisterType((*BridgeUpdate)(nil), "pb.BridgeUpdate")
	proto.RegisterEnum("pb.BridgeType", BridgeType_name, BridgeType_value)
	proto.RegisterEnum("pb.BridgeUpdate_Action", BridgeUpdate_Action_name, BridgeUpdate_Action_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for BridgeManager service

type BridgeManagerClient interface {
	AddBridge(ctx context.Context, in *AddBridgeRequest, opts ...grpc.CallOption) (*AddBridgeResponse, error)
	GetBridges(ctx context.Context, in *GetBridgesRequest, opts ...grpc.CallOption) (*GetBridgesResponse, error)
	WatchBridges(ctx context.Context, in *WatchBridgesRequest, opts ...grpc.CallOption) (BridgeManager_WatchBridgesClient, error)
}

type bridgeManagerClient struct {
	cc *grpc.ClientConn
}

func NewBridgeManagerClient(cc *grpc.ClientConn) BridgeManagerClient {
	return &bridgeManagerClient{cc}
}

func (c *bridgeManagerClient) AddBridge(ctx context.Context, in *AddBridgeRequest, opts ...grpc.CallOption) (*AddBridgeResponse, error) {
	out := new(AddBridgeResponse)
	err := grpc.Invoke(ctx, "/pb.BridgeManager/AddBridge", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bridgeManagerClient) GetBridges(ctx context.Context, in *GetBridgesRequest, opts ...grpc.CallOption) (*GetBridgesResponse, error) {
	out := new(GetBridgesResponse)
	err := grpc.Invoke(ctx, "/pb.BridgeManager/GetBridges", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bridgeManagerClient) WatchBridges(ctx context.Context, in *WatchBridgesRequest, opts ...grpc.CallOption) (BridgeManager_WatchBridgesClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_BridgeManager_serviceDesc.Streams[0], c.cc, "/pb.BridgeManager/WatchBridges", opts...)
	if err != nil {
		return nil, err
	}
	x := &bridgeManagerWatchBridgesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type BridgeManager_WatchBridgesClient interface {
	Recv() (*BridgeUpdate, error)
	grpc.ClientStream
}

type bridgeManagerWatchBridgesClient struct {
	grpc.ClientStream
}

func (x *bridgeManagerWatchBridgesClient) Recv() (*BridgeUpdate, error) {
	m := new(BridgeUpdate)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for BridgeManager service

type BridgeManagerServer interface {
	AddBridge(context.Context, *AddBridgeRequest) (*AddBridgeResponse, error)
	GetBridges(context.Context, *GetBridgesRequest) (*GetBridgesResponse, error)
	WatchBridges(*WatchBridgesRequest, BridgeManager_WatchBridgesServer) error
}

func RegisterBridgeManagerServer(s *grpc.Server, srv BridgeManagerServer) {
	s.RegisterService(&_BridgeManager_serviceDesc, srv)
}

func _BridgeManager_AddBridge_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddBridgeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BridgeManagerServer).AddBridge(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.BridgeManager/AddBridge",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BridgeManagerServer).AddBridge(ctx, req.(*AddBridgeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BridgeManager_GetBridges_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBridgesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BridgeManagerServer).GetBridges(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.BridgeManager/GetBridges",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BridgeManagerServer).GetBridges(ctx, req.(*GetBridgesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BridgeManager_WatchBridges_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(WatchBridgesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BridgeManagerServer).WatchBridges(m, &bridgeManagerWatchBridgesServer{stream})
}

type BridgeManager_WatchBridgesServer interface {
	Send(*BridgeUpdate) error
	grpc.ServerStream
}

type bridgeManagerWatchBridgesServer struct {
	grpc.ServerStream
}

func (x *bridgeManagerWatchBridgesServer) Send(m *BridgeUpdate) error {
	return x.ServerStream.SendMsg(m)
}

var _BridgeManager_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.BridgeManager",
	HandlerType: (*BridgeManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddBridge",
			Handler:    _BridgeManager_AddBridge_Handler,
		},
		{
			MethodName: "GetBridges",
			Handler:    _BridgeManager_GetBridges_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "WatchBridges",
			Handler:       _BridgeManager_WatchBridges_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "bridge.proto",
}

func init() { proto.RegisterFile("bridge.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 960 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x55, 0xdf, 0x6e, 0xdb, 0xb6,
	0x17, 0x8e, 0xe4, 0xd8, 0xb2, 0x8f, 0xdd, 0xfc, 0x14, 0xb6, 0xc9, 0x4f, 0xf5, 0xb0, 0xcd, 0x15,
	0x56, 0x20, 0x58, 0x30, 0x6f, 0x4b, 0x2f, 0x06, 0x14, 0xd8, 0x85, 0x5b, 0x07, 0x69, 0x80, 0xa5,
	0x0b, 0xb8, 0xba, 0x03, 0x76, 0xe3, 0x51, 0x12, 0x63, 0x13, 0xb1, 0x45, 0x4e, 0xa4, 0xed, 0x38,
	0x7b, 0x8b, 0x61, 0x0f, 0xb4, 0xdb, 0xdd, 0xed, 0x21, 0xf6, 0x20, 0x03, 0xff, 0xc8, 0x56, 0xdd,
	0xad, 0x77, 0x3c, 0xe7, 0xfb, 0x78, 0xf8, 0xf1, 0x3b, 0x87, 0x12, 0x74, 0x92, 0x82, 0x65, 0x13,
	0xda, 0x17, 0x05, 0x57, 0x1c, 0xf9, 0x22, 0x89, 0xff, 0xf6, 0x20, 0x18, 0x64, 0x59, 0x41, 0xa5,
	0x44, 0x9f, 0x80, 0xcf, 0x44, 0xe4, 0xf5, 0xbc, 0x93, 0xf6, 0xd9, 0x41, 0x5f, 0x24, 0x7d, 0x07,
	0xf4, 0x2f, 0x05, 0xf6, 0x99, 0x40, 0x4f, 0xa0, 0xb6, 0x90, 0x49, 0xe4, 0x1b, 0xc2, 0xff, 0xaa,
	0x84, 0x91, 0x4c, 0xb0, 0xc6, 0xba, 0xbf, 0x82, 0x7f, 0x29, 0x10, 0x82, 0xfd, 0x29, 0x97, 0xca,
	0x94, 0x6a, 0x61, 0xb3, 0x46, 0x11, 0x04, 0x39, 0x55, 0x73, 0x22, 0x6f, 0x4d, 0x81, 0x16, 0x2e,
	0x43, 0x8d, 0x4c, 0x88, 0xa2, 0x2b, 0xb2, 0x8e, 0x6a, 0x16, 0x71, 0xa1, 0xae, 0x23, 0x78, 0xa1,
	0xa2, 0xfd, 0x9e, 0x77, 0x52, 0xc7, 0x66, 0x8d, 0x1e, 0x43, 0x73, 0xc9, 0xc8, 0x38, 0x9b, 0xa6,
	0x22, 0x82, 0x9e, 0x77, 0xd2, 0xc4, 0xc1, 0x92, 0x91, 0xe1, 0x34, 0x15, 0xdd, 0xc7, 0x50, 0x1b,
	0xc9, 0xc4, 0xec, 0x22, 0x6a, 0x5a, 0x9e, 0xae, 0xd7, 0xf1, 0x1f, 0x3e, 0xb4, 0x5f, 0x98, 0xbb,
	0xff, 0xa0, 0x88, 0xa2, 0xe8, 0x23, 0x68, 0x31, 0x39, 0x16, 0x84, 0x15, 0x34, 0x33, 0xc4, 0x26,
	0x6e, 0x32, 0x79, 0x6d, 0x62, 0xf4, 0x35, 0x04, 0x4b, 0x5a, 0x48, 0xc6, 0xf3, 0x28, 0x33, 0x77,
	0xfd, 0xbf, 0xbe, 0x6b, 0x65, 0x7b, 0xff, 0xad, 0x85, 0x71, 0xc9, 0x43, 0x7d, 0x68, 0xdc, 0xb3,
	0x49, 0x42, 0x69, 0x94, 0x9b, 0x1d, 0xc7, 0xbb, 0x3b, 0x7e, 0x32, 0x28, 0x76, 0x2c, 0x74, 0x0a,
	0xf5, 0xfb, 0x15, 0x59, 0xd2, 0x88, 0x1b, 0xfa, 0xd1, 0x7b, 0x74, 0x0d, 0x62, 0xcb, 0xe9, 0x9e,
	0x42, 0xe0, 0x0e, 0x44, 0x21, 0xd4, 0x88, 0x60, 0xee, 0x6a, 0x7a, 0x89, 0x0e, 0xc0, 0x97, 0x2b,
	0x67, 0xa9, 0x2f, 0x57, 0xdd, 0x18, 0x1a, 0xf6, 0x2c, 0xed, 0x6b, 0x3a, 0x25, 0x79, 0x4e, 0x67,
	0x86, 0x5f, 0xc7, 0x65, 0xd8, 0x7d, 0x06, 0x75, 0x73, 0x00, 0x3a, 0x86, 0xc6, 0x94, 0xcf, 0xe9,
	0x65, 0xe6, 0x2a, 0xba, 0x48, 0x5b, 0x38, 0xe7, 0x19, 0x75, 0x65, 0xcd, 0x3a, 0xfe, 0xdd, 0x83,
	0x8e, 0x95, 0xf8, 0x92, 0xe7, 0x37, 0x6c, 0xa2, 0x49, 0x39, 0x99, 0xd3, 0xd2, 0x67, 0xbd, 0xd6,
	0x6a, 0x58, 0x56, 0xaa, 0x61, 0x19, 0x7a, 0x0a, 0x01, 0xb1, 0x33, 0x62, 0x9a, 0xd5, 0x3e, 0x6b,
	0x57, 0xc6, 0x06, 0x97, 0x18, 0xfa, 0x18, 0x20, 0x25, 0xe9, 0x94, 0x8e, 0x4d, 0xe3, 0xda, 0x66,
	0x7b, 0xcb, 0x64, 0xae, 0x89, 0x9a, 0xa2, 0x2e, 0x34, 0x15, 0x9b, 0xd3, 0x7b, 0x9e, 0xd3, 0xe8,
	0xcc, 0x80, 0x9b, 0x58, 0xcb, 0x3a, 0x70, 0xce, 0xad, 0x46, 0x22, 0xd3, 0xcd, 0x7d, 0x02, 0x1d,
	0x26, 0xc7, 0x64, 0x49, 0xd8, 0x8c, 0x24, 0x33, 0xea, 0xfa, 0xdb, 0x66, 0x72, 0x50, 0xa6, 0xd0,
	0xa7, 0xd0, 0xce, 0xb9, 0x62, 0x37, 0xeb, 0xf1, 0x42, 0xd2, 0xc2, 0x0d, 0x12, 0xd8, 0xd4, 0x48,
	0xd2, 0xa2, 0x42, 0x50, 0xf4, 0x4e, 0x39, 0x49, 0x8e, 0xf0, 0x86, 0xde, 0x29, 0x2d, 0xb9, 0xac,
	0x50, 0xcc, 0xa2, 0x8e, 0x95, 0xec, 0x0a, 0x14, 0xb3, 0xf8, 0x2f, 0x1f, 0x1a, 0x56, 0x96, 0xf3,
	0xc4, 0xdb, 0x78, 0x62, 0x67, 0x8f, 0xa4, 0x8a, 0x2d, 0xad, 0xc3, 0x66, 0xf6, 0x06, 0x26, 0x46,
	0x31, 0xec, 0xab, 0xb5, 0xa0, 0xe6, 0x25, 0x1c, 0xd8, 0x57, 0x68, 0xcb, 0xbc, 0x59, 0x0b, 0x8a,
	0x0d, 0xa6, 0x9f, 0x80, 0xee, 0xc8, 0x6c, 0xcc, 0x32, 0xa3, 0xbc, 0x85, 0x03, 0x13, 0x5f, 0x66,
	0x5a, 0x95, 0x85, 0x4c, 0x67, 0x9c, 0x91, 0x26, 0xf3, 0x5a, 0xb7, 0xe7, 0x14, 0x0e, 0x2d, 0x9c,
	0x51, 0x99, 0x16, 0x4c, 0x28, 0x3d, 0xe3, 0x56, 0x7b, 0x68, 0x80, 0xe1, 0x36, 0x8f, 0x62, 0xe8,
	0xcc, 0x49, 0xbe, 0xb8, 0x21, 0xa9, 0x5a, 0x14, 0xb4, 0x88, 0x1e, 0x18, 0xde, 0x3b, 0x39, 0x2d,
	0x85, 0xa5, 0x3c, 0x37, 0x1e, 0x3c, 0xea, 0xd5, 0xb4, 0x14, 0x1d, 0x8f, 0x8a, 0x19, 0x3a, 0x81,
	0x46, 0x6a, 0x06, 0xc5, 0x3d, 0xa2, 0x70, 0x7b, 0x17, 0x3b, 0x40, 0xd8, 0xe1, 0xe8, 0x29, 0xd4,
	0xa5, 0x9e, 0xfa, 0x88, 0x6e, 0xbf, 0x2c, 0x95, 0xc7, 0x80, 0x2d, 0x1a, 0xff, 0x0c, 0xe1, 0x20,
	0xcb, 0x2c, 0x80, 0xe9, 0x2f, 0x0b, 0x2a, 0xd5, 0xc6, 0x2e, 0xef, 0x03, 0x76, 0x6d, 0x85, 0xf8,
	0x1f, 0x16, 0x12, 0x7f, 0x03, 0x87, 0x95, 0x13, 0xa4, 0xe0, 0xb9, 0xd4, 0x1d, 0x69, 0xd8, 0xaf,
	0xa6, 0xfb, 0x32, 0xc2, 0x76, 0x3b, 0x76, 0x48, 0xfc, 0x10, 0x0e, 0x2f, 0xa8, 0xb2, 0x49, 0xe9,
	0xb4, 0xc5, 0xcf, 0x01, 0x55, 0x93, 0xae, 0xdc, 0x67, 0x10, 0xd8, 0x4d, 0x32, 0xf2, 0x7a, 0xb5,
	0x9d, 0x7a, 0x25, 0x14, 0x1f, 0xc1, 0xc3, 0x1f, 0x89, 0x4a, 0xa7, 0x3b, 0x25, 0x7f, 0xdb, 0xbc,
	0x41, 0x37, 0xea, 0x5f, 0x42, 0x43, 0x0f, 0x12, 0xcf, 0x9d, 0x03, 0x95, 0x2f, 0x95, 0x65, 0xf4,
	0x07, 0x06, 0xc6, 0x8e, 0x56, 0xb9, 0x8d, 0xff, 0x9f, 0xb7, 0xf9, 0x02, 0x1a, 0x76, 0x17, 0x6a,
	0x41, 0x7d, 0x30, 0x1c, 0x9e, 0x0f, 0xc3, 0x3d, 0xd4, 0x86, 0xe0, 0xe5, 0xab, 0xc1, 0xeb, 0x8b,
	0xf3, 0x61, 0xe8, 0xe9, 0x00, 0x9f, 0x5f, 0x7d, 0xff, 0xf6, 0x7c, 0x18, 0xfa, 0x9f, 0x8f, 0x01,
	0xb6, 0x9e, 0x6b, 0xe8, 0x82, 0xe6, 0xb4, 0x60, 0x69, 0xb8, 0x87, 0x3a, 0xd0, 0xfc, 0x8e, 0x73,
	0x91, 0x90, 0xf4, 0x36, 0xf4, 0x74, 0xb5, 0xeb, 0x82, 0xdf, 0xad, 0x43, 0x1f, 0x05, 0x50, 0x7b,
	0xb5, 0xa0, 0x61, 0x0d, 0x85, 0xd0, 0x79, 0xc1, 0x95, 0x9a, 0xd1, 0x82, 0xa7, 0xb7, 0x54, 0x85,
	0xfb, 0x3a, 0x73, 0xc5, 0x73, 0x2e, 0x0a, 0x96, 0xd2, 0xc1, 0x5c, 0x84, 0xf5, 0xb3, 0x3f, 0x3d,
	0x78, 0x60, 0x4f, 0xb8, 0x22, 0x39, 0x99, 0xd0, 0x02, 0x3d, 0x87, 0xd6, 0xa6, 0x51, 0xe8, 0x91,
	0xfb, 0xa4, 0xbc, 0x33, 0x19, 0xdd, 0xa3, 0x9d, 0xac, 0xb5, 0x3f, 0xde, 0x43, 0xdf, 0x02, 0x6c,
	0xdb, 0x82, 0x0c, 0xed, 0xbd, 0xde, 0x75, 0x8f, 0x77, 0xd3, 0x95, 0xed, 0x9d, 0x6a, 0x67, 0x90,
	0x71, 0xfc, 0x5f, 0x7a, 0xd5, 0x0d, 0x77, 0x5b, 0x11, 0xef, 0x7d, 0xe5, 0x25, 0x0d, 0xf3, 0xeb,
	0x7d, 0xf6, 0x4f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xd0, 0xbc, 0xb0, 0xd2, 0x8a, 0x07, 0x00, 0x00,
}
