#pragma once

#include "devicemanager-cpp/device_manager.hpp"

#include "bridge.hpp"

namespace jvstest {

class TestServer {
public:
    void setup();

    std::string printable() const;

    void run(uint16_t port);

private:
    jvs::DeviceManager _mgr;

    std::vector<std::shared_ptr<TestBridge> > _bridges;
};

}
