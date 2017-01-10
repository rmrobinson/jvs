#pragma once

namespace jvs {
namespace openzwaved {

struct Message {
    enum Type {
        Command,
        Shutdown
    };

    Type type;
};

}
}
