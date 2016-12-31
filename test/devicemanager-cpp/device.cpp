
#include <string>
#include <random>

#include "device.hpp"
#include "devicemanager-cpp/device_manager.hpp"

jvstest::TestDevice::TestDevice(jvs::DeviceManager& dm) :
    jvs::Device(dm) {
}

jvstest::TestDevice::~TestDevice() {
}

void jvstest::TestDevice::setup(const std::string& id) {
    proto::Device d;
    d.set_id(id);
    d.set_is_active(true);
    d.set_model_id("TD01");
    d.set_model_name("Test Device");
    d.set_model_description("Test Device");
    d.set_manufacturer("Widgets Inc.");

    d.set_address("/devices/" + id);

    proto::DeviceConfig* dc = d.mutable_config();
    assert(dc != nullptr);
    dc->set_name("Test device " + id);
    dc->set_description("Test device description");

    proto::DeviceState *ds = d.mutable_state();
    assert(ds != nullptr);

    ds->set_is_reachable(true);

    proto::DeviceState_BinaryState *dsbs = ds->mutable_binary();
    assert(dsbs != nullptr);

    dsbs->set_is_on(true);

    proto::DeviceState_RangeState *dsrs = ds->mutable_range();
    assert(dsrs != nullptr);

    dsrs->set_value(37);

    updateSelf(d);
}

bool jvstest::TestDevice::setConfig(proto::DeviceConfig& config) {
    proto::Device d = getDevice();
    d.mutable_config()->CopyFrom(config);
    
    updateSelf(d);
    return true;
}

bool jvstest::TestDevice::setState(proto::DeviceState& state) {
    proto::Device d = getDevice();
    d.mutable_state()->CopyFrom(state);

    updateSelf(d);
    return true;
}

void jvstest::TestDevice::randomUpdate() {
    std::random_device rd;
    std::mt19937 gen(rd());
    std::uniform_int_distribution<> red(0, 255);
    std::uniform_int_distribution<> green(0, 255);
    std::uniform_int_distribution<> blue(0, 255);
    std::uniform_int_distribution<> level(0, 100);

    proto::DeviceState ds = getDevice().state();
    ds.mutable_range()->set_value(level(gen));
    ds.mutable_color_rgb()->set_red(red(gen));
    ds.mutable_color_rgb()->set_green(green(gen));
    ds.mutable_color_rgb()->set_blue(blue(gen));

    setState(ds);
}
