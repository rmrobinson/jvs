
#include <cassert>

#include <openzwave/Notification.h>

#include "bridge.hpp"

jvs::openzwaved::Bridge::Bridge(OpenZWave::Manager& manager, uint32_t homeId, uint8_t nodeId) :
  _manager(manager), _isActive(true), _homeId(homeId), _nodeId(nodeId), _mode(Secondary) {
    for (size_t i = 0; i < UINT8_MAX; i++) {
      Node n (_manager, _homeId, i + 1);
      _nodes.emplace_back(n);
    }

    if (_manager.IsStaticUpdateController(_homeId)) {
        _mode = Bridge::StaticUpdate;
    } else if (_manager.IsPrimaryController(_homeId)) {
        _mode = Bridge::Primary;
    }

    std::string id;
    std::stringstream ss;
    ss << std::hex << _homeId;
    ss >> id;

    _bridgeData.set_id(id);
    _bridgeData.set_is_active(true);

    // Bridge data is mostly returned in the controller node, so,
    // once we have it in processEventMessage we'll fill in those fields.
}

void jvs::openzwaved::Bridge::processEventMessage(const Message& msg) {
    assert(msg.event != nullptr);
    assert(_homeId == msg.event->GetHomeId());

    //const OpenZWave::ValueID& id = msg.event->GetValueID();

    switch (msg.event->GetType()) {
        case OpenZWave::Notification::Type_ValueAdded:
        case OpenZWave::Notification::Type_ValueRemoved:
        case OpenZWave::Notification::Type_ValueChanged:
        case OpenZWave::Notification::Type_ValueRefreshed:
        case OpenZWave::Notification::Type_Group:
        case OpenZWave::Notification::Type_NodeProtocolInfo:
        case OpenZWave::Notification::Type_NodeNaming:
        case OpenZWave::Notification::Type_NodeEvent:
        {
            Node& n = _nodes[msg.event->GetNodeId() - 1];

            if (n.isActive()) {
              n.processEventMessage(msg);
            } else {
              fprintf(stdout, "Received node message without node existing\n");
            }

            // If we've gotten updated info on the controller node, we should refresh ourselves.
            if (n.nodeId() == _nodeId) {
              refreshBridgeData();
            }

            break;
        }

        case OpenZWave::Notification::Type_NodeNew:
            fprintf(stdout, "New node available\n");
            break;

        case OpenZWave::Notification::Type_NodeAdded:
        {
            Node& n = _nodes[msg.event->GetNodeId() - 1];

            if (!n.isActive()) {
                fprintf(stdout, "Bridge %x, node %d active\n", _homeId, msg.event->GetNodeId());
                n.activate();
            } else {
                fprintf(stdout, "Node already active: %d\n", msg.event->GetNodeId());
            }

            break;
        }

        case OpenZWave::Notification::Type_NodeRemoved:
            fprintf(stdout, "Node %d removed\n", msg.event->GetNodeId());
            break;

        case OpenZWave::Notification::Type_PollingDisabled:
            fprintf(stdout, "Polling disabled\n");
            break;

        case OpenZWave::Notification::Type_PollingEnabled:
            fprintf(stdout, "Polling enabled\n");
            break;

        case OpenZWave::Notification::Type_SceneEvent:
            fprintf(stdout, "Scene event\n");
            break;

        case OpenZWave::Notification::Type_NodeReset:
            fprintf(stdout, "Node %d reset\n", msg.event->GetNodeId());
            break;

        default:
            fprintf(stdout, "Unhandled type %d\n", msg.event->GetType());
            assert(false);
    }    
}

void jvs::openzwaved::Bridge::refreshBridgeData() {
    _bridgeData.set_manufacturer(getNode(_nodeId).getDeviceData().manufacturer());
    _bridgeData.set_model_id(getNode(_nodeId).getDeviceData().model_id());
    _bridgeData.set_model_name(getNode(_nodeId).getDeviceData().model_name());
}
