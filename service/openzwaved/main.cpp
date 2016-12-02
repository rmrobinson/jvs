#include <csignal>
#include <cstdio>
#include <string>

#include <unistd.h>

#include "server.hpp"

static jvs::openzwaved::Server s;

static void signal_handler(int signal) {
  if (signal == SIGINT) {
    s.stop();
  }
}

int main(int argc, char* argv[]) {
    GOOGLE_PROTOBUF_VERIFY_VERSION;

    std::string usbPath;
    int port;

    int opt;
    while ((opt = getopt(argc, argv, "d:p:")) != -1) {
      switch (opt) {
        case 'd':
          usbPath = std::string(optarg);
          break;
        case 'p':
          port = atoi(optarg);
          break;
        default:
          fprintf(stderr, "Usage: %s -d devicePath -p listenPort\n", argv[0]);
          return EXIT_FAILURE;
      }
    }

    if (port < 0 || port > UINT16_MAX) {
      fprintf(stderr, "Invalid port specified: %d\n", port);
      return EXIT_FAILURE;
    }

    std::signal(SIGINT, signal_handler);

    fprintf(stdout, "Starting server on port %d\n", port);

    std::vector<std::string> usbPaths;
    usbPaths.push_back(usbPath);

    s.run(usbPaths, static_cast<uint16_t>(port));

    return EXIT_SUCCESS;
}
