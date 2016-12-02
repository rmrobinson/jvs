#pragma once

#include <openzwave/Manager.h>

namespace jvs {
namespace openzwaved {

class Notification {
public:
    enum Type {
        Driver = 1,
        Button,
        Group,
        Node,
        Scene
    };

    const std::string& printable() const;

    Type getType() const;
    uint32_t getReceivedTime() const;

private:
    Type _type;

    uint32_t _receivedTime;

};

class ButtonNotification : public Notification {
public:
    uint8_t getButtonId() const;

private:
    uint8_t _buttonId;
};

class SceneNotification : public Notification {
    uint8_t getSceneId() const;

private:
    uint8_t _sceneId;
};

class DriverNotification : public Notification {
public:
    enum ControllerMode {
        Primary = 1,
        Secondary,
        StaticUpdate
    };

    DriverNotification(uint32_t homeId, uint8_t nodeId, ControllerMode mode);

    uint32_t getHomeId() const;

    uint8_t getNodeId() const;

    ControllerMode getMode() const;

private:
    uint32_t _homeId;

    uint8_t _nodeId;
    ControllerMode _mode;

};

class NodeNotification : public Notification {
public:
    enum Type {
        ValueAdded = 1,
        ValueRemoved,
        ValueChanged,
        ValueRefreshed,
        // GroupChanged, -- contained in the GroupNotification class
        Added,
        Removed,
        ProtocolInfo,
        Naming,
        Event,
        PollingEnabled,
        PollingDisabled,
    };

    Type getType() const;

    uint32_t getHomeId() const;

    uint8_t getNodeId() const;

    const OpenZWave::ValueID& getValueId() const;

private:
    Type _type;
    uint32_t _homeId;
    uint8_t _nodeId;

    OpenZWave::ValueID _valueId;
};

class GroupNotification : public NodeNotification {
public:
    uint8_t getGroupIndex() const;

private:
    uint8_t _groupIdx;
};


}
}

