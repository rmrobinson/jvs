
#include <string>
#include <random>

#include "bridge.hpp"
#include "devicemanager-cpp/device_manager.hpp"

jvstest::TestBridge::TestBridge(jvs::DeviceManager& dm) :
    jvs::Bridge(dm) {
}

jvstest::TestBridge::~TestBridge() {
}

void jvstest::TestBridge::setup(const std::string& id) {
    proto::Bridge b;
    b.set_id(id);
    b.set_is_active(true);
    b.set_model_id("T01");
    b.set_model_name("Testing");
    b.set_model_description("Test Bridge");
    b.set_manufacturer("Widgets Inc.");

    proto::BridgeConfig* bc = b.mutable_config();
    assert(bc != nullptr);
    bc->set_name("Test bridge " + id);
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

    updateSelf(b);

    std::random_device rd;
    std::mt19937 gen(rd());
    std::uniform_int_distribution<> dist(2, 5);

    const int count = dist(gen);

    for (int i = 0; i < count; i++) {
        std::shared_ptr<TestDevice> d = std::make_shared<TestDevice>(_dm);
        d->setup(std::to_string(i));

        _devices.push_back(d);

        _dm.addDevice(std::shared_ptr<jvs::Device>(d));
    }
}

void jvstest::TestBridge::run() {
    std::thread t = std::thread([this]() {
        size_t i = 0;
        while (true) {
            std::this_thread::sleep_for(std::chrono::seconds(1));

            std::shared_ptr<TestDevice> d = _devices[i % _devices.size()];
            d->randomUpdate();

            i++;

            if (i > 10) {
                _dm.removeDevice(d->getId());

                _dm.addDevice(std::shared_ptr<jvs::Device>(d));
                i = 0;
            }
        }
    });

    t.detach();   
}

bool jvstest::TestBridge::setConfig(proto::BridgeConfig& config) {
    proto::Bridge b = getBridge();
    b.mutable_config()->CopyFrom(config);
    
    updateSelf(b);
    return true;
}
