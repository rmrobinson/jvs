#include <csignal>
#include <cstdio>
#include <string>
#include <thread>

#include <unistd.h>

#include "server.hpp"

static bool done = false;

static void signal_handler(int signal) {
  if (signal == SIGINT) {
    done = true;
  }
}

int main(int argc, char* argv[]) {
    GOOGLE_PROTOBUF_VERIFY_VERSION;

    int port;

    int opt;
    while ((opt = getopt(argc, argv, "p:")) != -1) {
      switch (opt) {
        case 'p':
          port = atoi(optarg);
          break;
        default:
          fprintf(stderr, "Usage: %s -p listenPort\n", argv[0]);
          return EXIT_FAILURE;
      }
    }

    if (port < 0 || port > UINT16_MAX) {
      fprintf(stderr, "Invalid port specified: %d\n", port);
      return EXIT_FAILURE;
    }

    std::signal(SIGINT, signal_handler);

    fprintf(stdout, "Starting server on port %d\n", port);

    jvstest::TestServer s;
    s.setup();
    s.run(port);

    while (!done) {
    }

    return EXIT_SUCCESS;
}
