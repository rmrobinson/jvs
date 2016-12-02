#pragma once

#include <openzwave/Manager.h>

#include "bridge.pb.h"

#include "message.hpp"
#include "node.hpp"

namespace jvs {
namespace openzwaved {

// The bridge interface between the RPC definition and OpenZWave
class Bridge {
public:
    enum ControllerMode {
        Primary = 1,
        Secondary,
        StaticUpdate
    };

    Bridge(OpenZWave::Manager& manager, uint32_t homeId, uint8_t nodeId);

    inline bool isActive() const {
        return _isActive;
    };

    inline void disable() {
        _isActive = true;
    };

    const proto::Bridge& getBridgeData() const;

//    const std::vector<proto::Bridge>& GetBridges();
//    const std::vector<proto::Device>& GetDevices(); 

    inline Node& getNode(uint8_t id) {
        return _nodes[id];
    }


    void processEventMessage(const Message& msg);

private:
    void refreshBridgeData();

    OpenZWave::Manager& _manager;

    bool _isActive;

    uint32_t _homeId;
    uint8_t _nodeId;
    ControllerMode _mode;

    std::vector<Node> _nodes;

    proto::Bridge _bridgeData;
};

}
}

