#pragma once

#include "devicemanager-cpp/device.hpp"

namespace jvstest {

class TestDevice : public jvs::Device {
public:
    TestDevice(jvs::DeviceManager& dm);
    virtual ~TestDevice();

    void setup(const std::string& id);

    void randomUpdate();

    // The following functions implement the behaviours of the parent class.
    virtual bool setConfig(proto::DeviceConfig& config) override;

    virtual bool setState(proto::DeviceState& state) override;

private:

};

}
