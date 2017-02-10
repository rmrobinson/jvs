
#include <cassert>
#include <ctime>

#include <openzwave/Notification.h>

#include "command_class.hpp"
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

    d.set_address(std::to_string(_homeId) + "/" + std::to_string(_nodeId));

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

    const OpenZWave::ValueID& id = notification->GetValueID();

    switch (notification->GetType()) {
        case OpenZWave::Notification::Type_ValueAdded:
            fprintf(stdout, "Value added on node %d: ", _nodeId);
            _valueIds.push_back(id);

            processValueId(id, true);
            break;

        case OpenZWave::Notification::Type_ValueRemoved:
            removeValueId(id);
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
        {
            uint8_t groupIdx = notification->GetGroupIdx();
            fprintf(stdout, "Group %d changed\n", groupIdx);
            break;
        }
        case OpenZWave::Notification::Type_NodeProtocolInfo:
            fprintf(stdout, "Protocol info\n");
            break;
        case OpenZWave::Notification::Type_NodeNaming:
        {
            const std::string name = _zwaveManager.GetNodeName(_homeId, _nodeId);
            const std::string manufacturer = _zwaveManager.GetNodeManufacturerName(_homeId, _nodeId);
            const std::string modelId = _zwaveManager.GetNodeProductType(_homeId, _nodeId);
            const std::string modelName = _zwaveManager.GetNodeProductId(_homeId, _nodeId);
            const std::string modelDescription =_zwaveManager.GetNodeProductName(_homeId, _nodeId);

            proto::Device d = getDevice();

            proto::DeviceConfig* c = d.mutable_config();
            assert(c != nullptr);
            c->set_name(name);

            d.set_model_id(modelId);
            d.set_model_name(modelName);
            d.set_model_description(modelDescription);
            d.set_manufacturer(manufacturer);

            updateSelf(d);

            break;
        }
        case OpenZWave::Notification::Type_NodeEvent:
        {
            uint8_t event = notification->GetEvent();
            fprintf(stdout, "Node event: %d\n", event);
            break;
        }

        default:
            fprintf(stdout, "Node %d unhandled notification type %d (%s)\n", _nodeId, notification->GetType(), notification->GetAsString().c_str());
            assert(false);
    }    
}

void jvs::openzwaved::Node::processValueId(const OpenZWave::ValueID& vid, bool isAddition) {
    proto::Device d = getDevice();
    bool hasDChanged = false;

    proto::DeviceState* state = d.mutable_state();
    assert(state != nullptr);

    fprintf(stdout, "%s\n", formatValueId(vid).c_str());

    const std::string label = _zwaveManager.GetValueLabel(vid);
    const std::string units = _zwaveManager.GetValueUnits(vid);

    // We can also get at the OpenZWave-supplied help string if we decide to start exposing this.
    //std::string help = _zwaveManager.GetValueHelp(vid);

    if (vid.GetCommandClassId() == CommandClass::SwitchBinary) {
        if (vid.GetIndex() == 0) {
            proto::DeviceState_BinaryState* binary = state->mutable_binary();
            assert(binary != nullptr);

            if (vid.GetType() != OpenZWave::ValueID::ValueType_Bool) {
                fprintf(stderr, "Invalid type retrieved for switch binary: %d\n", vid.GetType());
                return;
            }

            bool val = false;
            _zwaveManager.GetValueAsBool(vid, &val);
            binary->set_is_on(val);

            hasDChanged = true;
        }
    } else if (vid.GetCommandClassId() == CommandClass::SwitchMultilevel) {
        if (label.compare("Level") == 0) {
            proto::DeviceState_BinaryState* binary = state->mutable_binary();
            assert(binary != nullptr);
            proto::DeviceState_RangeState* range = state->mutable_range();
            assert(range != nullptr);
            proto::Device_RangeDevice* rangeDevice = d.mutable_range();
            assert(rangeDevice != nullptr);
            rangeDevice->set_minimum(_zwaveManager.GetValueMin(vid));
            rangeDevice->set_maximum(_zwaveManager.GetValueMax(vid));
    
            if (vid.GetType() != OpenZWave::ValueID::ValueType_Byte) {
                fprintf(stderr, "Invalid type retrieved for switch level: %d\n", vid.GetType());
                return;
            }
    
            uint8_t val = 0;
            _zwaveManager.GetValueAsByte(vid, &val);

            if (val != range->value() || (val > 0) != binary->is_on() || isAddition) {
                range->set_value(val);
                binary->set_is_on(val > 0);

                hasDChanged = true;

                if (!isAddition) {
                    // If we see a change in the level, schedule a check for a short period of time later.
                    // We have seen lights which change on a range report values part of the way along.
                    // So this will prevent us from leaving the range value partway changed.
                    _zwaveManager.RefreshValue(vid);
                }
            }
        }
    } else if (vid.GetCommandClassId() == CommandClass::SensorBinary) {
        if (label.compare("Sensor") == 0) {
            proto::DeviceState_PresenceState* presence = state->mutable_presence();
            assert(presence != nullptr);

            if (vid.GetType() != OpenZWave::ValueID::ValueType_Bool) {
                fprintf(stderr, "Invalid type retrieved for sensor binary: %d\n", vid.GetType());
                return;
            }

            bool val;
            _zwaveManager.GetValueAsBool(vid, &val);
            presence->set_is_present(val);

            hasDChanged = true;
        }
    } else if (vid.GetCommandClassId() == CommandClass::SensorMultilevel) {
        if (label.compare("Temperature") == 0) {
            proto::DeviceState_TemperatureState* temperature = state->mutable_temperature();
            assert(temperature != nullptr);

            if (vid.GetType() != OpenZWave::ValueID::ValueType_Decimal) {
                fprintf(stderr, "Invalid type retrieved for sensor multilevel (temperature): %d\n", vid.GetType());
                return;
            }

            float val;
            _zwaveManager.GetValueAsFloat(vid, &val);

            if (units.compare("F") == 0) {
                val = (val - 32) * 5 / 9;
            }
            temperature->set_temperature_celsius(val);

            hasDChanged = true;
        } else if (label.compare("Relative Humidity")) {
            float val;
            _zwaveManager.GetValueAsFloat(vid, &val);

        } else if (label.compare("Luminance")) {
            float val;
            _zwaveManager.GetValueAsFloat(vid, &val);

        } else if (label.compare("Ultraviolet")) {
            float val;
            _zwaveManager.GetValueAsFloat(vid, &val);
        }
    } else if (vid.GetCommandClassId() == CommandClass::Battery) {
        if (label.compare("Battery Level") == 0) {
            uint8_t val;
            _zwaveManager.GetValueAsByte(vid, &val);
        }
    } else if (vid.GetCommandClassId() == CommandClass::Alarm) {
        if (label.compare("Burglar") == 0) {
            uint8_t val;
            _zwaveManager.GetValueAsByte(vid, &val);
        }
    }

    if (hasDChanged) {
        updateSelf(d);
    }
}

void jvs::openzwaved::Node::removeValueId(const OpenZWave::ValueID& vid) {
    _valueIds.erase(std::remove(_valueIds.begin(), _valueIds.end(), vid), _valueIds.end());
}

std::string jvs::openzwaved::Node::formatValueId(const OpenZWave::ValueID& vid) const {
    const std::string label = _zwaveManager.GetValueLabel(vid);
    const std::string help = _zwaveManager.GetValueHelp(vid);
    const std::string units = _zwaveManager.GetValueUnits(vid);

    std::string ret;

    ret += label + " (" + help + "), CC=" + std::to_string(vid.GetCommandClassId()) + ", idx=" + std::to_string(vid.GetIndex()) + " equals";

    switch (vid.GetType()) {
        case OpenZWave::ValueID::ValueType_Bool:
        {
            bool val;
            _zwaveManager.GetValueAsBool(vid, &val);
            ret += " bool " + std::to_string(val);
            break;
        }
        case OpenZWave::ValueID::ValueType_Byte:
        {
            uint8_t val;
            _zwaveManager.GetValueAsByte(vid, &val);
            ret += " byte " + std::to_string(val);
            break;
        }
        case OpenZWave::ValueID::ValueType_Decimal:
        {
            float val;
            _zwaveManager.GetValueAsFloat(vid, &val);
            ret += " float " + std::to_string(val);
            break;
        }
        case OpenZWave::ValueID::ValueType_Int:
        {
            int32_t val;
            _zwaveManager.GetValueAsInt(vid, &val);
            ret += " int " + std::to_string(val);
            break;
        }
        case OpenZWave::ValueID::ValueType_List:
        {
            ret += " str list [";
            std::vector<std::string> val;
            _zwaveManager.GetValueListItems(vid, &val);
            for (size_t i = 0; i < val.size(); i++) {
                if (i > 0) {
                    ret += '|';
                }

                ret += val[i];
            }
            ret += "]";
            break;
        }
        case OpenZWave::ValueID::ValueType_String:
        {
            std::string val;
            _zwaveManager.GetValueAsString(vid, &val);
            ret += " str " + val;
            break;
        }
        case OpenZWave::ValueID::ValueType_Short:
        {
            int16_t val;
            _zwaveManager.GetValueAsShort(vid, &val);
            ret += " short " + std::to_string(val);
            break;
        }
        case OpenZWave::ValueID::ValueType_Button:
            ret += " button";
            break;
        default:
            ret += " unhandled type " + std::to_string(vid.GetType());
    }
 
    ret += " " + units;

    if (_zwaveManager.IsValueReadOnly(vid)) {
      ret += " readonly";
    }

    return ret;
}

bool jvs::openzwaved::Node::setConfig(proto::DeviceConfig& config) {
    std::unique_lock<std::mutex> lock(_mutex);

    if (!_isActive) {
        return false;
    }

    _zwaveManager.SetNodeName(_homeId, _nodeId, config.name());

    proto::Device d = getDevice();
    d.mutable_config()->CopyFrom(config);

    updateSelf(d);
    return true;
}

bool jvs::openzwaved::Node::setState(proto::DeviceState& state) {
    std::unique_lock<std::mutex> lock(_mutex);

    if (!_isActive) {
        return false;
    }

    proto::Device d = getDevice();

    for (size_t i = 0; i < _valueIds.size(); i++) {
        const OpenZWave::ValueID& vid = _valueIds[i];
        std::string label = _zwaveManager.GetValueLabel(vid);

        switch(vid.GetCommandClassId()) {
            case CommandClass::SwitchMultilevel:
            {
                // A multi-level switch should have both binary and range properties.
                assert(d.state().has_binary());
                assert(d.state().has_range());

                // Here we set the dim level, based on both the binary and range states.
                if (label.compare("Level") == 0) {
                    uint8_t level = 0; // default to off

                    if (state.has_binary() && state.binary().is_on()) {
                        if (state.has_range()) {
                            // We shouldn't have gotten to this point, but confirm it anyways.
                            if (state.range().value() < 0 || state.range().value() > UINT8_MAX) {
                                return false;
                            }

                            level = static_cast<uint8_t>(state.range().value());
                        }

                        // this is a bit of a special case.
                        // if we are saying "on" but the supplied level is 0, we assume
                        // that this is because the level was previously set to 0 when turned off
                        // so 'on' is overriding this value.
                        if (level == 0) {
                            level = 99; // on, if we were previously off and level isn't set.
                        }
                    }

                    if (!_zwaveManager.SetValue(vid, level)) {
                        return false;
                    }

                    if (state.has_range()) {
                        d.mutable_state()->mutable_range()->CopyFrom(state.range());
                    }
                    d.mutable_state()->mutable_binary()->CopyFrom(state.binary());
                }
                break;
            }
            default:
                continue;
        }
    }

    // TODO: we should just be updating the 'processing' flag in the device
    //updateSelf(d);
    return true;
}
