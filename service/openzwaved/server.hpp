#pragma once

#include <unordered_map>

#include <openzwave/Notification.h>
#include <openzwave/Manager.h>
#include <openzwave/Options.h>

#include "message.hpp"
#include "bridge.hpp"
#include "concurrent_queue.hpp"

namespace jvs {
namespace openzwaved {

class Server {
public:
    Server();
    ~Server();

    void run(const std::vector<std::string>& devicePaths, const uint16_t port);

    void stop();

    /// @brief Interface for other threads to access the server state.
    /// To prevent threading issues, all requests must be made in the form of a Message,
    /// which will be interpreted and a threadsafe copy of the data will be returned.
    /// @param message The requested data. See struct Message above.
    void sendMessage(const Message& message);

private:
    /// @brief Handler for events coming from the OpenZWave API.
    /// This will be called on one of several OpenZWave managed threads.
    /// Do not access any server internals from this function.
    static void onZwaveNotification(OpenZWave::Notification const* n, void* c);

    void processEventMessage(const Message& message);
    void processCommandMessage(const Message& message);

    /// @brief We will have a thread-safe queue here which the main thread pops events off to process.
    /// The OnZwaveNotification handler will slightly deserialize the notification then push it onto the queue.
    /// The gRPC handlers will format the commands as appropriate then push it onto the queue.
    /// The main thread will take all received messages and pass them to the appropriate methods.
    ConcurrentQueue<Message> _messages;

    /// @brief The server will hold a pointer to the manager and multiplex all of the inbound and outbound communication.
    /// We will use this handle, instead of the static Manager::Get(), to possibly support multiple handles if ozw changes in the future.
    OpenZWave::Manager* _manager;

    /// @brief The server will hold a pointer to the options used to initialize the manager.
    /// We will use this handle, instead of the static Options::Get(), to possibley support multiple handles if ozw changes in the future.
    OpenZWave::Options* _options;

    /// @brief Used to exit the message processing queue.
    bool _done;

    /// @brief Upon detecting a new ZWave home ID, we will create a bridge with that home ID and insert it into this map.
    /// All future updates will use this map to address the appropriate bridge.
    std::unordered_map<uint32_t, Bridge> _bridges;
    
};

}
}

