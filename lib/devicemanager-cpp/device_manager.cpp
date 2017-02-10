
#include <chrono>
#include <thread>

#include <grpc/grpc.h>
#include <grpc++/server.h>
#include <grpc++/server_builder.h>
#include <grpc++/server_context.h>
#include <grpc++/security/server_credentials.h>

#include "device_manager.hpp"

jvs::DeviceManager::DeviceManager() : _isRunning(false)
{}

grpc::Status jvs::DeviceManager::GetBridges(grpc::ServerContext* context,
                                            const proto::GetBridgesRequest* request,
                                            proto::GetBridgesResponse* response) {

    assert(context != nullptr);
    assert(request != nullptr);
    assert(response != nullptr);

    (void) context;
    (void) request;

    std::unique_lock<std::mutex> lock(_bridgesMutex);

    for (auto const& pair : _bridges) {
        proto::Bridge* b = response->add_bridges();
        assert (b != nullptr);
        b->CopyFrom(pair.second->getBridge());
    }

    return grpc::Status::OK;
}

grpc::Status jvs::DeviceManager::WatchBridges(grpc::ServerContext* context,
                                              const proto::WatchBridgesRequest* request,
                                              grpc::ServerWriter<proto::WatchBridgesResponse>* writer) {
    assert(context != nullptr);
    assert(request != nullptr);
    assert(writer != nullptr);
    
    (void) context;
    (void) request;

    ConcurrentNotifier<proto::WatchBridgesResponse>::Watcher w (*this);
    bool keepWriting = true;

    while (keepWriting) {
      proto::WatchBridgesResponse resp;
      w.wait_and_pop(resp);

      keepWriting = writer->Write(resp);
    }

    return grpc::Status::OK;
}

grpc::Status jvs::DeviceManager::GetDevices(grpc::ServerContext* context,
                                            const proto::GetDevicesRequest* request,
                                            proto::GetDevicesResponse* response) {

    assert(context != nullptr);
    assert(request != nullptr);
    assert(response != nullptr);

    (void) context;
    (void) request;

    std::unique_lock<std::mutex> lock(_devicesMutex);

    for (auto const& pair : _devices) {
        proto::Device* d = response->add_devices();
        assert (d != nullptr);
        d->CopyFrom(pair.second->getDevice());
    }

    return grpc::Status::OK;
}

grpc::Status jvs::DeviceManager::GetDevice(grpc::ServerContext* context,
                                           const proto::GetDeviceRequest* request,
                                           proto::GetDeviceResponse* response) {

    assert(context != nullptr);
    assert(request != nullptr);
    assert(response != nullptr);

    (void) context;
    (void) request;

    std::unique_lock<std::mutex> lock(_devicesMutex);

    auto pair = _devices.find(request->address());
  
    if (pair == _devices.end()) {
      return grpc::Status(grpc::StatusCode::NOT_FOUND, "Device does not exist");
    }

    assert(pair->second != nullptr);

    response->mutable_device()->CopyFrom(pair->second->getDevice());
    return grpc::Status::OK;
}

grpc::Status jvs::DeviceManager::WatchDevices(grpc::ServerContext* context,
                                              const proto::WatchDevicesRequest* request,
                                              grpc::ServerWriter<proto::WatchDevicesResponse>* writer) {
    assert(context != nullptr);
    assert(request != nullptr);
    assert(writer != nullptr);

    (void) context;
    (void) request;

    ConcurrentNotifier<proto::WatchDevicesResponse>::Watcher w (*this);
    bool keepWriting = true;

    while (keepWriting) {
      proto::WatchDevicesResponse resp;
      w.wait_and_pop(resp);

      keepWriting = writer->Write(resp);
    }

    return grpc::Status::OK;
}

grpc::Status jvs::DeviceManager::SetDeviceConfig(grpc::ServerContext* context,
                                                 const proto::SetDeviceConfigRequest* request,
                                                 proto::SetDeviceConfigResponse* response) {
  assert(context != nullptr);
  assert(request != nullptr);
  assert(response != nullptr);

  std::unique_lock<std::mutex> lock(_devicesMutex);

  auto pair = _devices.find(request->address());
  
  if (pair == _devices.end()) {
    return grpc::Status(grpc::StatusCode::NOT_FOUND, "Device does not exist");
  }

  assert(pair->second != nullptr);

  if (request->version() != pair->second->getVersion()) {
    return grpc::Status(grpc::StatusCode::ABORTED, "Device has been updated since this request was issued");
  }

  proto::DeviceConfig c;
  c.CopyFrom(request->config());

  bool ret = pair->second->setConfig(c);

  if (!ret) {
      return grpc::Status(grpc::StatusCode::UNKNOWN, "SetDeviceConfig call failed");
  }

  response->mutable_device()->mutable_config()->CopyFrom(c);

  return grpc::Status::OK;
}

grpc::Status jvs::DeviceManager::SetDeviceState(grpc::ServerContext* context,
                                                const proto::SetDeviceStateRequest* request,
                                                proto::SetDeviceStateResponse* response) {
  assert(context != nullptr);
  assert(request != nullptr);
  assert(response != nullptr);

  std::unique_lock<std::mutex> lock(_devicesMutex);

  auto pair = _devices.find(request->address());
  
  if (pair == _devices.end()) {
    return grpc::Status(grpc::StatusCode::NOT_FOUND, "Device does not exist");
  }

  assert(pair->second != nullptr);

  if (request->version() != pair->second->getVersion()) {
    return grpc::Status(grpc::StatusCode::ABORTED, "Device has been updated since this request was issued");
  }

  proto::DeviceState s;
  s.CopyFrom(request->state());

  bool ret = pair->second->setState(s);

  if (!ret) {
      return grpc::Status(grpc::StatusCode::UNKNOWN, "SetDeviceState call failed");
  }

  response->mutable_device()->mutable_state()->CopyFrom(s);
  return grpc::Status::OK;
}

void jvs::DeviceManager::start(uint16_t port) {
  std::string address("0.0.0.0:" + std::to_string(port));

  grpc::ServerBuilder::ServerBuilder builder;
  builder.AddListeningPort(address, grpc::InsecureServerCredentials());
  builder.RegisterService(static_cast<proto::BridgeManager::Service*> (this));
  builder.RegisterService(static_cast<proto::DeviceManager::Service*> (this));

  _grpcServer = builder.BuildAndStart();

  // Keep running forever; when we kill the main process this process will also die.
  _grpcThread = std::thread([this](){
    _isRunning = true;
    _grpcServer->Wait();
  });

  _grpcThread.detach();
}

void jvs::DeviceManager::addBridge(std::shared_ptr<Bridge> bridge) {
  std::unique_lock<std::mutex> lock(_bridgesMutex);
  _bridges.emplace(bridge->getId(), bridge);

  proto::WatchBridgesResponse notification;
  notification.set_action(proto::WatchBridgesResponse::ADDED);
  notification.mutable_bridge()->CopyFrom(bridge->getBridge());

  broadcast(notification);
}

void jvs::DeviceManager::removeBridge(const std::string& id) {
  std::unique_lock<std::mutex> lock(_bridgesMutex);

  auto pair = _bridges.find(id);
  
  if (pair == _bridges.end()) {
    return;
  }

  assert(pair->second != nullptr);

  proto::WatchBridgesResponse notification;
  notification.set_action(proto::WatchBridgesResponse::REMOVED);
  notification.mutable_bridge()->CopyFrom(pair->second->getBridge());

  _bridges.erase(id);

  broadcast(notification);
}

void jvs::DeviceManager::addDevice(std::shared_ptr<Device> device) {
  std::unique_lock<std::mutex> lock(_devicesMutex);
  _devices.emplace(device->getAddress(), device);

  proto::WatchDevicesResponse notification;
  notification.set_action(proto::WatchDevicesResponse::ADDED);
  notification.mutable_device()->CopyFrom(device->getDevice());

  broadcast(notification);
}

void jvs::DeviceManager::removeDevice(const std::string& address) {
  std::unique_lock<std::mutex> lock(_devicesMutex);

  auto pair = _devices.find(address);
  
  if (pair == _devices.end()) {
    return;
  }

  assert(pair->second != nullptr);

  proto::WatchDevicesResponse notification;
  notification.set_action(proto::WatchDevicesResponse::REMOVED);
  notification.mutable_device()->CopyFrom(pair->second->getDevice());

  _devices.erase(address);

  broadcast(notification);
}
