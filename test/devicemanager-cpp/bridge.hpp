#pragma once

#include "devicemanager-cpp/bridge.hpp"

#include "device.hpp"

namespace jvstest {

class TestBridge : public jvs::Bridge {
public:
    TestBridge(jvs::DeviceManager& manager);
    virtual ~TestBridge();

    void setup(const std::string& id);
    void run();

    // The following functions implement the behaviours of the parent class.
    virtual bool setConfig(proto::BridgeConfig& config) override;

private:
    /// @brief The collection of devices currenty part of this bridge.
    std::vector<std::shared_ptr<TestDevice> > _devices;

};

}
