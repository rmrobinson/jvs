
#include <cassert>

#include "server.hpp"

jvs::openzwaved::Server::Server() : _manager(nullptr), _options(nullptr), _done(false) {}

jvs::openzwaved::Server::~Server() {
}

void jvs::openzwaved::Server::run(const std::vector<std::string>& devicePaths, const uint16_t port) {
    // TODO: remove once we add in the gRPC endpoint code.
    (void)port;

    if (devicePaths.size() < 1) {
        return;
    }

    _options = OpenZWave::Options::Create("./config/", "", "--SaveConfiguration=true --DumpTriggerLevel=0");
    assert (_options != nullptr);

    _options->AddOptionBool("ConsoleOutput", false);
    _options->AddOptionInt("PollInterval", 500);
    _options->AddOptionBool("IntervalBetweenPolls", true);
    _options->AddOptionBool("ValidateValueChanges", true);
    _options->Lock();

    _manager = OpenZWave::Manager::Create();
    assert (_manager != nullptr);

    _manager->AddWatcher(jvs::openzwaved::Server::onZwaveNotification, this);

    for (size_t i = 0; i < devicePaths.size(); i++) {
        const std::string& path = devicePaths[i];

        if (path.length() < 1) {
            continue;
        //} else if (path.find("usb") != std::string::npos) {
        //    _manager->AddDriver("HID Controller", OpenZWave::Driver::ControllerInterface_Hid);
        } else {
            _manager->AddDriver(path);
        }
    }

    // TODO: initialize the gRPC endpoint contained in the devicemanager-cpp library.

    while (!_done) {
        Message m;
        _messages.wait_and_pop(m);

        if (m.type == Message::Event) {
            processEventMessage(m);
            // TODO: we need to delete the notification, but we'll leak it now since the destructor is private.
            // delete m.event;
        } else if (m.type == Message::Shutdown) {
            fprintf(stdout, "Shutting down\n");
        } else {
            fprintf(stdout, "Unknown message\n");
        }
    }

    assert(_manager != nullptr);
    _manager->RemoveWatcher(jvs::openzwaved::Server::onZwaveNotification, nullptr);

    _manager = nullptr;
    OpenZWave::Manager::Destroy();

    _options = nullptr;
    OpenZWave::Options::Destroy();

    return;
}

void jvs::openzwaved::Server::stop() {
    _done = true;

    Message m;
    m.type = Message::Shutdown;
    _messages.push(m);
}

void jvs::openzwaved::Server::sendMessage(const Message& msg) {
    if (!_done) {
        _messages.push(msg);
    }
}

void jvs::openzwaved::Server::onZwaveNotification(OpenZWave::Notification const* ozwn, void* context) {
    assert(ozwn != nullptr);
    assert(context != nullptr);

    jvs::openzwaved::Server* server = static_cast<jvs::openzwaved::Server*>(context);

    OpenZWave::Notification* n = new OpenZWave::Notification(*ozwn);
    assert(n != nullptr);

    Message m;
    m.type = Message::Event;
    m.event = n;

    server->sendMessage(m);
}

void jvs::openzwaved::Server::processEventMessage(const Message& msg) {
    assert(msg.event != nullptr);
    assert(_manager != nullptr);

    //const OpenZWave::ValueID& id = msg.event->GetValueID();

    switch (msg.event->GetType()) {
        case OpenZWave::Notification::Type_DriverReset:
            fprintf(stdout, "Bridge %x reset\n", msg.event->GetHomeId());

            _bridges.erase(msg.event->GetHomeId());
            // We intentionally do not break here so as to go and follow the same flow for creating a bridge.

        case OpenZWave::Notification::Type_DriverReady:
        {
            uint32_t homeId = msg.event->GetHomeId();
            uint8_t nodeId = msg.event->GetNodeId();

            fprintf(stdout, "Bridge %x created\n", homeId);

            _bridges.emplace(homeId, Bridge(*_manager, homeId, nodeId));
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
            auto bItr = _bridges.find(msg.event->GetHomeId());

            if (bItr == _bridges.end()) {
                fprintf(stderr, "Received bridge message without bridge existing\n");
                break;
            }

            bItr->second.processEventMessage(msg);
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
            _bridges.erase(msg.event->GetHomeId());
            break;

        case OpenZWave::Notification::Type_Notification:
        case OpenZWave::Notification::Type_ControllerCommand:
            fprintf(stdout, "Notification or controller command result\n");
            break;

        default:
            fprintf(stdout, "Unhandled type %d\n", msg.event->GetType());
            assert(false);
    }

}
