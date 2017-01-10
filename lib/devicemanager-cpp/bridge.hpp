#pragma once

#include "bridge.pb.h"

namespace jvs {

class DeviceManager;

class Bridge {
public:
    Bridge(DeviceManager& dm) : _deviceManager(dm) {};
    virtual ~Bridge() {};

    inline const std::string& getId() const {
        return _bridge.id();
    }

    // The following functions implement the behaviours which Manager requires vis a vis bridgemanager
    inline const proto::Bridge& getBridge() const {
        return _bridge;
    };

    virtual bool setConfig(proto::BridgeConfig& config) = 0;

protected:
    void updateSelf(const proto::Bridge& bridge, bool suppressNotification = false);

    /// @brief Handle to the managing device manager.
    /// Used to propogate updates back to the watchers.
    DeviceManager& _deviceManager;

private:
    /// @brief The protocol representation of this bridge.
    proto::Bridge _bridge;
};

}
