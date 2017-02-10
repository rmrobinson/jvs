#pragma once

#include <openzwave/Manager.h>

#include "devicemanager-cpp/bridge.hpp"

#include "node.hpp"

namespace jvs {
namespace openzwaved {

/// @brief An instance of an OpenZWave driver.
/// This class also implements the JVS Bridge interface.
class Driver : public Bridge {
public:
    enum ControllerMode {
        Primary = 1,
        Secondary,
        StaticUpdate
    };

    Driver(DeviceManager& deviceManager,
           OpenZWave::Manager& zwaveManager,
           uint32_t homeId,
           uint8_t nodeId);

    inline bool isActive() const {
        return _isActive;
    };

    inline void disable() {
        _isActive = false;
    };

    inline std::shared_ptr<Node>& getNode(uint8_t id) {
        return _nodes[id];
    }

    void processZwaveNotification(const OpenZWave::Notification* notification);

    // The following functions implement the behaviours of the parent class.
    virtual bool setConfig(proto::BridgeConfig& config) override;

private:
    void refreshDriverData();

    OpenZWave::Manager& _zwaveManager;

    bool _isActive;

    uint32_t _homeId;
    uint8_t _nodeId;
    ControllerMode _mode;

    /// @brief The collection of nodes exposed by this driver.
    /// Since there are a fixed number of possible nodes per driver,
    /// these will be allocated on driver allocation and marked as inactive.
    std::vector<std::shared_ptr<Node> > _nodes;
};

}
}

