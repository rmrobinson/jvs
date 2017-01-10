
#include <cassert>
#include <ctime>

#include <openzwave/Notification.h>

#include "node.hpp"

jvs::openzwaved::Node::Node(jvs::DeviceManager& deviceManager,
                            OpenZWave::Manager& zwaveManager,
                            uint32_t homeId,
                            uint8_t nodeId) :
    jvs::Device(deviceManager),
    _zwaveManager(zwaveManager),
    _homeId(homeId),
    _nodeId(nodeId),
    _isActive(false),
    _lastModifiedTime(time(0)) {

    proto::Device d;
    d.set_id(std::to_string(_nodeId));
    d.set_is_active(_isActive);
    d.set_model_id("");
    d.set_model_name("");
    d.set_model_description("");
    d.set_manufacturer("");

    d.set_address(std::to_string(_homeId) + "/" + std::to_string(_nodeId - 1));

    proto::DeviceConfig* dc = d.mutable_config();
    assert(dc != nullptr);
    dc->set_name("");
    dc->set_description("");

    proto::DeviceState *ds = d.mutable_state();
    assert(ds != nullptr);

    ds->set_is_reachable(false);

    // We wish to suppress notifications when we create the node.
    // Being added to the bridge will trigger the necessary notifications.
    updateSelf(d, true);
}

void jvs::openzwaved::Node::activate() {
    _isActive = true;

    proto::Device d = getDevice();

    std::string desc;
    if (_zwaveManager.IsNodeZWavePlus(_homeId, _nodeId)) {
        desc = _zwaveManager.GetNodeDeviceTypeString(_homeId, _nodeId);
    } else {
        desc = _zwaveManager.GetNodeType(_homeId, _nodeId);
     }

    proto::DeviceConfig* c = d.mutable_config();
    assert(c != nullptr);
    c->set_description(desc);

    // We wish to suppress notifications when we activate the node,
    /// since the create message will not be sent until after activation succeeds.
    updateSelf(d, true);
}

void jvs::openzwaved::Node::deactivate() {
    _isActive = false;
}

void jvs::openzwaved::Node::processZwaveNotification(const OpenZWave::Notification* notification) {
    assert(notification != nullptr);
    assert(_isActive);
    assert(_nodeId == notification->GetNodeId());

    fprintf(stdout, "Node %d handling event: %s\n", _nodeId, notification->GetAsString().c_str());

    const OpenZWave::ValueID& id = notification->GetValueID();

    switch (notification->GetType()) {
        case OpenZWave::Notification::Type_ValueAdded:
            fprintf(stdout, "Value added on node %d: ", _nodeId);
            _valueIds.push_back(id);

            processValueId(id);
            break;

        case OpenZWave::Notification::Type_ValueRemoved:
            fprintf(stdout, "Value removed\n");
            break;

        case OpenZWave::Notification::Type_ValueChanged:
            fprintf(stdout, "Value changed on node %d: ", _nodeId);
            processValueId(id);
            break;

        case OpenZWave::Notification::Type_ValueRefreshed:
            fprintf(stdout, "Value refreshed on node %d: ", _nodeId);
            processValueId(id);
            break;

        case OpenZWave::Notification::Type_Group:
            fprintf(stdout, "Group\n");
            break;
        case OpenZWave::Notification::Type_NodeProtocolInfo:
            fprintf(stdout, "Protocol info\n");
            break;
        case OpenZWave::Notification::Type_NodeNaming:
        {
            fprintf(stdout, "Naming\n");

            const std::string name = _zwaveManager.GetNodeName(_homeId, _nodeId);
            const std::string manufacturer = _zwaveManager.GetNodeManufacturerName(_homeId, _nodeId);
            const std::string modelId = _zwaveManager.GetNodeProductType(_homeId, _nodeId);
            const std::string modelName = _zwaveManager.GetNodeProductName(_homeId, _nodeId);

            proto::Device d = getDevice();

            proto::DeviceConfig* c = d.mutable_config();
            assert(c != nullptr);
            c->set_name(name);

            d.set_model_id(modelId);
            d.set_model_name(modelName);
            d.set_manufacturer(manufacturer);

            updateSelf(d);

            break;
        }

        default:
            fprintf(stdout, "Unhandled notification type %d\n", notification->GetType());
            assert(false);
    }    
}

void jvs::openzwaved::Node::processValueId(const OpenZWave::ValueID& id) {
    proto::Device d = getDevice();
    bool hasDChanged = false;

    proto::DeviceState* state = d.mutable_state();
    assert(state != nullptr);

    std::string label = _zwaveManager.GetValueLabel(id);

    // Light switch with range.
    if (id.GetCommandClassId() == 38) {
        if (id.GetIndex() == 0) {
            proto::DeviceState_BinaryState* binary = state->mutable_binary();
            assert(binary != nullptr);
            proto::DeviceState_RangeState* range = state->mutable_range();
            assert(range != nullptr);
    
            if (id.GetType() != OpenZWave::ValueID::ValueType_Byte) {
                fprintf(stderr, "Invalid type retrieved for light level: %d\n", id.GetType());
                return;
            }
    
            uint8_t val = 0;
            _zwaveManager.GetValueAsByte(id, &val);
            range->set_value(val);
            binary->set_is_on(val > 0);

            hasDChanged = true;
        }
    } else if (id.GetCommandClassId() == 49) {
        if (id.GetIndex() == 0) {
            proto::DeviceState_PresenceState* presence = state->mutable_presence();
            assert(presence != nullptr);

            if (id.GetType() != OpenZWave::ValueID::ValueType_Bool) {
                fprintf(stderr, "Invalid type retrieved for sensor presence: %d\n", id.GetType());
                return;
            }

            bool val;
            _zwaveManager.GetValueAsBool(id, &val);
            presence->set_is_present(val);

            hasDChanged = true;
        }
        // TODO: add the other sensor values in here.
    }

    if (hasDChanged) {
        updateSelf(d);
    }

    // TODO: everything below is simply for improved debugging.
    fprintf(stdout, "%s (CC %d, index %d) =", label.c_str(), id.GetCommandClassId(), id.GetIndex());

    switch (id.GetType()) {
        case OpenZWave::ValueID::ValueType_Bool:
        {
            bool val;
            _zwaveManager.GetValueAsBool(id, &val);
            fprintf(stdout, " bool %d", val);
            break;
        }
        case OpenZWave::ValueID::ValueType_Byte:
        {
            uint8_t val;
            _zwaveManager.GetValueAsByte(id, &val);
            fprintf(stdout, " byte %d", val);
            break;
        }
        case OpenZWave::ValueID::ValueType_Decimal:
        {
            float val;
            _zwaveManager.GetValueAsFloat(id, &val);
            fprintf(stdout, " float %f", val);
            break;
        }
        case OpenZWave::ValueID::ValueType_Int:
        {
            int val;
            _zwaveManager.GetValueAsInt(id, &val);
            fprintf(stdout, " int %d", val);
            break;
        }
        case OpenZWave::ValueID::ValueType_List:
        {
            std::vector<std::string> val;
            _zwaveManager.GetValueListItems(id, &val);
            for (size_t i = 0; i < val.size(); i++) {
                fprintf(stdout, " %s", val[i].c_str());
            }
            break;
        }
        case OpenZWave::ValueID::ValueType_String:
        {
            std::string val;
            _zwaveManager.GetValueAsString(id, &val);
            fprintf(stdout, " %s", val.c_str());
            break;
        }
        case OpenZWave::ValueID::ValueType_Short:
        {
            short val;
            _zwaveManager.GetValueAsShort(id, &val);
            fprintf(stdout, " %d", val);
            break;
        }
        case OpenZWave::ValueID::ValueType_Button:
            fprintf(stdout, " button");
            break;
        default:
            fprintf(stdout, "Unhandled type %d", id.GetType());
    }

    std::string units = _zwaveManager.GetValueUnits(id);

    fprintf(stdout, " %s", units.c_str());

    if (_zwaveManager.IsValueReadOnly(id)) {
      fprintf(stdout, " readonly");
    }

    fprintf(stdout, "\n");
}

bool jvs::openzwaved::Node::setConfig(proto::DeviceConfig& config) {
    // TODO: Actually enact changes on the node.
    proto::Device d = getDevice();
    d.mutable_config()->CopyFrom(config);

    updateSelf(d);
    return true;
}

bool jvs::openzwaved::Node::setState(proto::DeviceState& state) {
    // TODO: actually enact changes on the node.
    proto::Device d = getDevice();
    d.mutable_state()->CopyFrom(state);

    updateSelf(d);
    return true;
}
