#pragma once

#include <openzwave/Manager.h>

namespace jvs {
namespace openzwaved {

struct Message {
    enum Type {
        Event,
        Command,
        Shutdown
    };

    Type type;

    OpenZWave::Notification* event;
};

}
}
