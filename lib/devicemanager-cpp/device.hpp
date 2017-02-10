#pragma once

#include "device.pb.h"

namespace jvs {

class DeviceManager;

class Device {
public:
    Device(DeviceManager& dm) : _dm(dm) {};
    virtual ~Device() {};

    inline const std::string& getId() const {
      return _device.id();
    }

    inline const std::string& getAddress() const {
      return _device.address();
    }

    inline const std::string& getVersion() const {
      return _device.state().version();
    }

    // The following functions implement the behaviours which Manager requires vis a vis devicemanager
    inline const proto::Device& getDevice() const {
      return _device;
    };

    virtual bool setConfig(proto::DeviceConfig& config) = 0;
    virtual bool setState(proto::DeviceState& state) = 0;

protected:
    void updateSelf(const proto::Device& device, bool suppressNotification = false);

private:
    /// @brief Handle to the managing device manager.
    /// Used to propogate updates back to the watchers.
    DeviceManager& _dm;

    /// @brief The protocol representation of this device.
    proto::Device _device;
};

}
