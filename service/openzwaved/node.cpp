
#include <cassert>
#include <ctime>

#include <openzwave/Notification.h>

#include "node.hpp"

//jvs::openzwaved::Node::Node() : _manager(nullptr), _homeId(0), _nodeId(0), _isActive(false), _lastModifiedTime(time(0)), _deviceData() {}

jvs::openzwaved::Node::Node(OpenZWave::Manager& manager, uint32_t homeId, uint8_t nodeId) :
 _manager(manager), _homeId(homeId), _nodeId(nodeId), _isActive(true), _lastModifiedTime(time(0)), _deviceData() {
  _deviceData.set_id(std::to_string(nodeId));
  _deviceData.set_address(std::to_string(nodeId - 1));
}

void jvs::openzwaved::Node::activate() {
    _isActive = true;

    std::string desc;
    if (_manager.IsNodeZWavePlus(_homeId, _nodeId)) {
        desc = _manager.GetNodeDeviceTypeString(_homeId, _nodeId);
      } else {
        desc = _manager.GetNodeType(_homeId, _nodeId);
      }

      proto::DeviceConfig* conf(_deviceData.mutable_config());
      assert(conf != nullptr);
      conf->set_description(desc);
}

void jvs::openzwaved::Node::processEventMessage(const Message& msg) {
    assert(_isActive);
    assert(msg.event != nullptr);
    assert(_nodeId == msg.event->GetNodeId());

    fprintf(stdout, "Node %d handling event: %s\n", _nodeId, msg.event->GetAsString().c_str());

    const OpenZWave::ValueID& id = msg.event->GetValueID();

    switch (msg.event->GetType()) {
        case OpenZWave::Notification::Type_ValueAdded:
            fprintf(stdout, "Value added on node %d: ", msg.event->GetNodeId());
            _valueIds.push_back(id);

            processValueId(id);
            break;

        case OpenZWave::Notification::Type_ValueRemoved:
            fprintf(stdout, "Value removed\n");
            break;

        case OpenZWave::Notification::Type_ValueChanged:
            fprintf(stdout, "Value changed on node %d: ", msg.event->GetNodeId());
            processValueId(id);
            break;

        case OpenZWave::Notification::Type_ValueRefreshed:
            fprintf(stdout, "Value refreshed on node %d: ", msg.event->GetNodeId());
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

            const std::string name = _manager.GetNodeName(_homeId, _nodeId);
            const std::string modelName = _manager.GetNodeProductName(_homeId, _nodeId);
            const std::string modelId = _manager.GetNodeProductType(_homeId, _nodeId);
            const std::string manufacturer = _manager.GetNodeManufacturerName(_homeId, _nodeId);

            proto::DeviceConfig* conf(_deviceData.mutable_config());
            assert(conf != nullptr);
            conf->set_name(name);

            _deviceData.set_model_id(modelId);
            _deviceData.set_model_name(modelName);
            _deviceData.set_manufacturer(manufacturer);

            break;
        }

        default:
            fprintf(stdout, "Unhandled type %d\n", msg.event->GetType());
            assert(false);
    }    
}

void jvs::openzwaved::Node::processValueId(const OpenZWave::ValueID& id) {
    proto::DeviceState* state = _deviceData.mutable_state();
    assert(state != nullptr);

    std::string label = _manager.GetValueLabel(id);

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
            _manager.GetValueAsByte(id, &val);
            range->set_value(val);
            binary->set_is_on(val > 0);
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
            _manager.GetValueAsBool(id, &val);
            presence->set_is_present(val);
        }
        // TODO: add the other sensor values in here.
    }

    // TODO: everything below is simply for improved debugging.
    fprintf(stdout, "%s (CC %d, index %d) =", label.c_str(), id.GetCommandClassId(), id.GetIndex());

    switch (id.GetType()) {
        case OpenZWave::ValueID::ValueType_Bool:
        {
            bool val;
            _manager.GetValueAsBool(id, &val);
            fprintf(stdout, " bool %d", val);
            break;
        }
        case OpenZWave::ValueID::ValueType_Byte:
        {
            uint8_t val;
            _manager.GetValueAsByte(id, &val);
            fprintf(stdout, " byte %d", val);
            break;
        }
        case OpenZWave::ValueID::ValueType_Decimal:
        {
            float val;
            _manager.GetValueAsFloat(id, &val);
            fprintf(stdout, " float %f", val);
            break;
        }
        case OpenZWave::ValueID::ValueType_Int:
        {
            int val;
            _manager.GetValueAsInt(id, &val);
            fprintf(stdout, " int %d", val);
            break;
        }
        case OpenZWave::ValueID::ValueType_List:
        {
            std::vector<std::string> val;
            _manager.GetValueListItems(id, &val);
            for (size_t i = 0; i < val.size(); i++) {
                fprintf(stdout, " %s", val[i].c_str());
            }
            break;
        }
        case OpenZWave::ValueID::ValueType_String:
        {
            std::string val;
            _manager.GetValueAsString(id, &val);
            fprintf(stdout, " %s", val.c_str());
            break;
        }
        case OpenZWave::ValueID::ValueType_Short:
        {
            short val;
            _manager.GetValueAsShort(id, &val);
            fprintf(stdout, " %d", val);
            break;
        }
        case OpenZWave::ValueID::ValueType_Button:
            fprintf(stdout, " button");
            break;
        default:
            fprintf(stdout, "Unhandled type %d", id.GetType());
    }

    std::string units = _manager.GetValueUnits(id);

    fprintf(stdout, " %s", units.c_str());

    if (_manager.IsValueReadOnly(id)) {
      fprintf(stdout, " readonly");
    }

    fprintf(stdout, "\n");
}
