#pragma once

#include <list>

#include <openzwave/Manager.h>

#include "devicemanager-cpp/device.hpp"

namespace jvs {
namespace openzwaved {

class Node : public Device {
public:
    Node(jvs::DeviceManager& deviceManager,
         OpenZWave::Manager& zwaveManager,
         uint32_t homeId,
         uint8_t nodeId);

    void activate();

    void deactivate();

    inline bool isActive() const {
        return _isActive;
    };

    void addValue(const OpenZWave::ValueID& id);

    void removeValue(const OpenZWave::ValueID& id);

    inline uint32_t lastChangedTime() const {
        return _lastModifiedTime;
    }

    inline uint8_t nodeId() const {
        return _nodeId;
    }

    void processZwaveNotification(const OpenZWave::Notification* notification);

    // The following functions implement the behaviours of the parent class.
    virtual bool setConfig(proto::DeviceConfig& config) override;

    virtual bool setState(proto::DeviceState& state) override;

private:
    void processValueId(const OpenZWave::ValueID& id);

    void removeValueId(const OpenZWave::ValueID& id);

    OpenZWave::Manager& _zwaveManager;

    uint32_t _homeId;
    uint8_t _nodeId;

    bool _isActive;

    time_t _lastModifiedTime;

    std::list<OpenZWave::ValueID> _valueIds;

    // TODO: add in group information here.
};

}
}