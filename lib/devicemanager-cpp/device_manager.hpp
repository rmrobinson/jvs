#pragma once

#include <mutex>
#include <thread>
#include <unordered_map>

#include <grpc++/server.h>

#include "util-cpp/concurrent_watcher.hpp"

#include "bridge.grpc.pb.h"
#include "device.grpc.pb.h"

#include "bridge.hpp"
#include "device.hpp"

namespace jvs {

class DeviceManager :
    public proto::BridgeManager::Service,
    public proto::DeviceManager::Service,
    public ConcurrentNotifier<proto::WatchBridgesResponse>,
    public ConcurrentNotifier<proto::WatchDevicesResponse> {
public:
    DeviceManager();

    using ConcurrentNotifier<proto::WatchBridgesResponse>::broadcast;
    using ConcurrentNotifier<proto::WatchDevicesResponse>::broadcast;

    /// @brief Implements the GetBridges function call of proto::Bridgemanager
    /// This function blocks on the _bridgesMutex to ensure consistent access.
    grpc::Status GetBridges(grpc::ServerContext* context,
                            const proto::GetBridgesRequest* request,
                            proto::GetBridgesResponse* response) override;

    /// @brief Implements the WatchBridges function call of proto::BridgeManager
    /// This function does not block on the _bridgesMutex since updates are passed in
    /// by value, negating consistency issues (order issues may still be present).
    grpc::Status WatchBridges(grpc::ServerContext* context,
                              const proto::WatchBridgesRequest* request,
                              grpc::ServerWriter<proto::WatchBridgesResponse>* writer) override;


    grpc::Status GetDevices(grpc::ServerContext* context,
                            const proto::GetDevicesRequest* request,
                            proto::GetDevicesResponse* response) override;

    grpc::Status WatchDevices(grpc::ServerContext* context,
                              const proto::WatchDevicesRequest* request,
                              grpc::ServerWriter<proto::WatchDevicesResponse>* writer) override;

    grpc::Status GetDevice(grpc::ServerContext* context,
                           const proto::GetDeviceRequest* request,
                           proto::GetDeviceResponse* response) override;

    grpc::Status SetDeviceConfig(grpc::ServerContext* context,
                                 const proto::SetDeviceConfigRequest* request,
                                 proto::SetDeviceConfigResponse* response) override;

    grpc::Status SetDeviceState(grpc::ServerContext* context,
                                const proto::SetDeviceStateRequest* request,
                                proto::SetDeviceStateResponse* response) override;


    void addBridge(std::shared_ptr<Bridge> bridge);
    void removeBridge(const std::string& id);

    void addDevice(std::shared_ptr<Device> device);
    void removeDevice(const std::string& address);

    void start(uint16_t port);

    inline bool isRunning() const {
        return _isRunning;
    }

private:
    /// @brief The collection of registered bridges, keyed by their IDs.
    std::unordered_map<std::string, std::shared_ptr<Bridge> > _bridges;
    /// @brief Serializes accesses to the _bridges collection.
    mutable std::mutex _bridgesMutex;

    /// @brief The collection of registered devices, keyed by the device address.
    /// It is up to the implementation to prevent collections between device addresss
    /// with the same internal IDs, i.e. by prepending the bridge ID and a delimiter
    /// to any internal IDs.
    /// A simple exmaple would be <bridgeAddress>/<deviceId>
    std::unordered_map<std::string, std::shared_ptr<Device> > _devices;
    /// @brief Serializes access to the _devices collection.
    mutable std::mutex _devicesMutex;

    /// @brief The thread running the GRPC listener.
    std::thread _grpcThread;
    /// @brief The GRPC server instance.
    std::unique_ptr<grpc::Server> _grpcServer;

    /// @brief A flag keeping track of whether the GRPC thread is running.
    bool _isRunning;

};

}