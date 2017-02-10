
#include <cassert>

#include "server.hpp"

jvs::openzwaved::Server::Server() : _zwaveManager(nullptr), _options(nullptr), _done(false) {}

jvs::openzwaved::Server::~Server() {
}

void jvs::openzwaved::Server::run(const std::vector<std::string>& devicePaths, const uint16_t port) {
    if (devicePaths.size() < 1) {
        return;
    }

    _options = OpenZWave::Options::Create("./config/", "", "--SaveConfiguration=true --DumpTriggerLevel=0");
    assert (_options != nullptr);

    _options->AddOptionBool("ConsoleOutput", false);
    // TODO: It isn't clear these options are needed, but keep for now in case they come in handy.
    //_options->AddOptionInt("PollInterval", 500);
    //_options->AddOptionBool("IntervalBetweenPolls", true);
    //_options->AddOptionBool("ValidateValueChanges", true);
    _options->Lock();

    _zwaveManager = OpenZWave::Manager::Create();
    assert (_zwaveManager != nullptr);

    _zwaveManager->AddWatcher(jvs::openzwaved::Server::onZwaveNotification, this);

    for (size_t i = 0; i < devicePaths.size(); i++) {
        const std::string& path = devicePaths[i];

        if (path.length() < 1) {
            continue;
        //} else if (path.find("usb") != std::string::npos) {
        //    _zwaveManager->AddDriver("HID Controller", OpenZWave::Driver::ControllerInterface_Hid);
        } else {
            _zwaveManager->AddDriver(path);
        }
    }

    _deviceManager.start(port);


    while (!_done) {
        // TODO: we should really be looping somehwere else, this is gross.
    }

    assert(_zwaveManager != nullptr);
    _zwaveManager->RemoveWatcher(jvs::openzwaved::Server::onZwaveNotification, nullptr);

    _zwaveManager = nullptr;
    OpenZWave::Manager::Destroy();

    _options = nullptr;
    OpenZWave::Options::Destroy();

    return;
}

void jvs::openzwaved::Server::stop() {
    _done = true;
}

void jvs::openzwaved::Server::onZwaveNotification(OpenZWave::Notification const* notification, void* context) {
    assert(notification != nullptr);
    assert(context != nullptr);

    jvs::openzwaved::Server* server = static_cast<jvs::openzwaved::Server*>(context);
    assert(server != nullptr);

    if (!server->_done) {
        server->processZwaveNotification(notification);
    }
}

void jvs::openzwaved::Server::processZwaveNotification(const OpenZWave::Notification* notification) {
    assert(notification != nullptr);
    assert(_zwaveManager != nullptr);

    std::unique_lock<std::mutex> lock(_driversMutex);

    //const OpenZWave::ValueID& id = msg.event->GetValueID();

    switch (notification->GetType()) {
        case OpenZWave::Notification::Type_DriverReset:
            fprintf(stdout, "Driver %x reset\n", notification->GetHomeId());

            _drivers.erase(notification->GetHomeId());
            // We intentionally do not break here so as to go and follow the same flow for creating a bridge.

        case OpenZWave::Notification::Type_DriverReady:
        {
            uint32_t homeId = notification->GetHomeId();
            uint8_t nodeId = notification->GetNodeId();

            fprintf(stdout, "Driver %x created\n", homeId);

            std::shared_ptr<Driver> d = std::make_shared<Driver>(_deviceManager, *_zwaveManager, homeId, nodeId);
            assert(d != nullptr);

            _drivers.emplace(homeId, d);
            _deviceManager.addBridge(std::shared_ptr<jvs::Bridge>(d));

            break;
        }

        case OpenZWave::Notification::Type_ValueAdded:
        case OpenZWave::Notification::Type_ValueRemoved:
        case OpenZWave::Notification::Type_ValueChanged:
        case OpenZWave::Notification::Type_ValueRefreshed:
        case OpenZWave::Notification::Type_Group:
        case OpenZWave::Notification::Type_NodeNew:
        case OpenZWave::Notification::Type_NodeAdded:
        case OpenZWave::Notification::Type_NodeRemoved:
        case OpenZWave::Notification::Type_NodeProtocolInfo:
        case OpenZWave::Notification::Type_NodeNaming:
        case OpenZWave::Notification::Type_NodeEvent:
        case OpenZWave::Notification::Type_PollingDisabled:
        case OpenZWave::Notification::Type_PollingEnabled:
        case OpenZWave::Notification::Type_SceneEvent:
        case OpenZWave::Notification::Type_NodeReset:
        {
            auto dItr = _drivers.find(notification->GetHomeId());

            if (dItr == _drivers.end()) {
                fprintf(stderr, "Received driver message without driver existing\n");
                break;
            }

            assert(dItr->second != nullptr);
            dItr->second->processZwaveNotification(notification);
            break;
        }

        case OpenZWave::Notification::Type_CreateButton:
        case OpenZWave::Notification::Type_DeleteButton:
        case OpenZWave::Notification::Type_ButtonOn:
        case OpenZWave::Notification::Type_ButtonOff:
            fprintf(stdout, "Button command!!!\n");
            break;

        case OpenZWave::Notification::Type_DriverFailed:
            fprintf(stdout, "Unable to load driver for a device, ignoring\n");
            break;

        case OpenZWave::Notification::Type_EssentialNodeQueriesComplete:
        case OpenZWave::Notification::Type_NodeQueriesComplete:
        case OpenZWave::Notification::Type_AwakeNodesQueried:
        case OpenZWave::Notification::Type_AllNodesQueriedSomeDead:
        case OpenZWave::Notification::Type_AllNodesQueried:
            fprintf(stdout, "Node query complete with some state\n");
            break;

        case OpenZWave::Notification::Type_DriverRemoved:
            _drivers.erase(notification->GetHomeId());
            break;

        case OpenZWave::Notification::Type_Notification:
        case OpenZWave::Notification::Type_ControllerCommand:
            fprintf(stdout, "Notification or controller command result\n");
            break;

        default:
            fprintf(stdout, "Unhandled type %d\n", notification->GetType());
            assert(false);
    }

}
