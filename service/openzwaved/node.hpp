#pragma once

#include <vector>

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

    /// @brief Send a request to the OpenZWave library to change the node state.
    ///
    /// Inspection of the OpenZWave code shows that their internal implementation
    /// is thread-safe, so we will use that to our advantage.
    /// This function will not persist its changes; the async callback from
    /// OpenZWave will trigger the updateSelf method.
    virtual bool setState(proto::DeviceState& state) override;

private:
    // TODO: this might be abstracted into an OpenZWave helper of some sort.
    std::string formatValueId(const OpenZWave::ValueID& vid) const;

    void processValueId(const OpenZWave::ValueID& vid, bool isAddition = false);

    void removeValueId(const OpenZWave::ValueID& vid);

    OpenZWave::Manager& _zwaveManager;

    const uint32_t _homeId;
    const uint8_t _nodeId;

    // This mutex serializes updates to the properties below.
    mutable std::mutex _mutex;

    bool _isActive;

    time_t _lastModifiedTime;

    std::vector<OpenZWave::ValueID> _valueIds;

    // TODO: add in group information here.
};

}
}