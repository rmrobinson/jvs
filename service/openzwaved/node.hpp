#pragma once

#include <list>

#include <openzwave/Manager.h>

#include "device.pb.h"

#include "message.hpp"

namespace jvs {
namespace openzwaved {

class Node {
public:
    Node(OpenZWave::Manager& manager, uint32_t homeId, uint8_t nodeId);

    void activate();

    inline bool isActive() const {
        return _isActive;
    };

    void addValue(const OpenZWave::ValueID& id);

    void removeValue(const OpenZWave::ValueID& id);

    uint32_t lastChangedTime() const;

    inline uint8_t nodeId() const {
        return _nodeId;
    }

    inline const proto::Device& getDeviceData() const {
        return _deviceData;
    }

    void processEventMessage(const Message& msg);

private:
    void processValueId(const OpenZWave::ValueID& id);

    void removeValueId(const OpenZWave::ValueID& id);

    OpenZWave::Manager& _manager;

    uint32_t _homeId;
    uint8_t _nodeId;

    bool _isActive;

    time_t _lastModifiedTime;

    std::list<OpenZWave::ValueID> _valueIds;

    proto::Device _deviceData;

    // TODO: add in group information here.
};

}
}