#include <cstdio>

// Link against these to confirm we can find the required include files.
#include <grpc++/grpc++.h>
#include <openzwave/Manager.h>

#include "version.pb.h"

int main(int argc, char* argv[]) {
    (void) argc;
    (void) argv;

    printf("OpenZWave version: %s\n", OpenZWave::Manager::getVersionAsString().c_str());

    return 0;
}
