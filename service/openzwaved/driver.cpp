
#include <cassert>

#include <openzwave/Notification.h>

#include "devicemanager-cpp/device_manager.hpp"

#include "driver.hpp"

jvs::openzwaved::Driver::Driver(jvs::DeviceManager& dm,
                                OpenZWave::Manager& zwm,
                                uint32_t homeId,
                                uint8_t nodeId) :
    jvs::Bridge(dm),
    _zwaveManager(zwm),
    _isActive(true),
    _homeId(homeId),
    _nodeId(nodeId),
    _mode(Secondary) {

    for (size_t i = 0; i < UINT8_MAX; i++) {
      std::shared_ptr<Node> n = std::make_shared<Node>(dm, zwm, _homeId, i + 1);
      _nodes.push_back(n);
    }

    if (_zwaveManager.IsStaticUpdateController(_homeId)) {
        _mode = Driver::StaticUpdate;
    } else if (_zwaveManager.IsPrimaryController(_homeId)) {
        _mode = Driver::Primary;
    }

    std::string id;
    std::stringstream ss;
    ss << std::hex << _homeId;
    ss >> id;

    proto::Bridge b;
    b.set_id(id);
    b.set_is_active(true);
    b.set_model_id("");
    b.set_model_name("");
    b.set_model_description("");
    b.set_manufacturer("");

    proto::BridgeConfig* bc = b.mutable_config();
    assert(bc != nullptr);
    bc->set_name("");
    bc->set_timezone("UTC");

    proto::Address* addr = bc->mutable_address();
    assert(addr != nullptr);

    proto::Address_Usb* usb = addr->mutable_usb();
    assert(usb != nullptr);

    usb->set_path("/dev/testDevice" + id);

    proto::BridgeState *bs = b.mutable_state();
    assert(bs != nullptr);

    bs->set_is_paired(true);
    
    proto::BridgeState_Version* version = bs->mutable_version();
    assert(version != nullptr);

    version->set_api("0.01");
    version->set_sw("alpha");

    // We don't want a CHANGED notification going out after creation.
    // When we add the bridge to the device manager the necessary notification
    // will be sent.
    updateSelf(b, true);

    // Driver data is mostly exposed by the controller node, so,
    // once we have it in processZwaveNotification we'll fill in those fields.
}

void jvs::openzwaved::Driver::processZwaveNotification(const OpenZWave::Notification* notification) {
    assert(notification != nullptr);
    assert(_homeId == notification->GetHomeId());

    //const OpenZWave::ValueID& id = notification->GetValueID();

    switch (notification->GetType()) {
        case OpenZWave::Notification::Type_ValueAdded:
        case OpenZWave::Notification::Type_ValueRemoved:
        case OpenZWave::Notification::Type_ValueChanged:
        case OpenZWave::Notification::Type_ValueRefreshed:
        case OpenZWave::Notification::Type_Group:
        case OpenZWave::Notification::Type_NodeProtocolInfo:
        case OpenZWave::Notification::Type_NodeNaming:
        case OpenZWave::Notification::Type_NodeEvent:
        {
            const uint8_t nodeId = notification->GetNodeId();
            std::shared_ptr<Node> n = _nodes[nodeId - 1];
            assert(n != nullptr);

            if (n->isActive()) {
              n->processZwaveNotification(notification);
            } else {
              fprintf(stdout, "Received node notification without node existing\n");
            }

            // If we've gotten updated info on the controller node, we should refresh ourselves.
            if (n->nodeId() == _nodeId) {
              refreshDriverData();
            }

            break;
        }

        case OpenZWave::Notification::Type_NodeNew:
            fprintf(stdout, "New node available\n");
            break;

        case OpenZWave::Notification::Type_NodeAdded:
        {
            const uint8_t nodeId = notification->GetNodeId();
            std::shared_ptr<Node> n = _nodes[nodeId - 1];

            if (!n->isActive()) {
                fprintf(stdout, "Driver %x, node %d active\n", _homeId, nodeId);
                n->activate();

                _deviceManager.addDevice(std::shared_ptr<Device>(n));
            } else {
                fprintf(stdout, "Node already active: %d\n", nodeId);
            }

            break;
        }

        case OpenZWave::Notification::Type_NodeRemoved:
        {
            const uint8_t nodeId = notification->GetNodeId();
            std::shared_ptr<Node> n = _nodes[nodeId - 1];

            if (n->isActive()) {
                fprintf(stdout, "Node %d removed\n", nodeId);
                n->deactivate();
            } else {
                fprintf(stdout, "Node is already deactivated: %d\n", nodeId);
            }
            break;
        }

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
            fprintf(stdout, "Node %d reset\n", notification->GetNodeId());
            break;

        default:
            fprintf(stdout, "Unhandled notification type %d\n", notification->GetType());
            assert(false);
    }    
}

bool jvs::openzwaved::Driver::setConfig(proto::BridgeConfig& config) {
    // TODO: actually enact changes on this driver.
    proto::Bridge b = getBridge();
    b.mutable_config()->CopyFrom(config);
    
    updateSelf(b);
    return true;
}

void jvs::openzwaved::Driver::refreshDriverData() {
    proto::Bridge b = getBridge();

    std::shared_ptr<Node>& n = getNode(_nodeId);
    assert(n != nullptr);

    b.set_manufacturer(n->getDevice().manufacturer());
    b.set_model_id(n->getDevice().model_id());
    b.set_model_name(n->getDevice().model_name());

    if (b.model_description().length() < 1) {
        b.set_model_description("ZWave Bridge");
    }

    proto::BridgeConfig* c = b.mutable_config();
    assert(c != nullptr);
    c->set_name(n->getDevice().config().name());

    updateSelf(b);
}
