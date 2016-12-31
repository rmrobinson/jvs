#include <random>

#include "server.hpp"

void jvstest::TestServer::setup() {
    std::random_device rd;
    std::mt19937 mt(rd());
    std::uniform_real_distribution<double> dist(1.0, 10.0);

    //const int count = dist(mt);
    const int count = 1;

    for (int i = 0; i < count; i++) {
        std::shared_ptr<TestBridge> b = std::make_shared<TestBridge>(_mgr);
        b->setup(std::to_string(i));

        _bridges.push_back(b);
        _mgr.addBridge(std::shared_ptr<jvs::Bridge>(b));
    }
}

void jvstest::TestServer::run(uint16_t port) {
    for (size_t i = 0; i < _bridges.size(); i++) {
        if (_bridges[i] == nullptr) {
          continue;
        }

        _bridges[i]->run();
    }

    _mgr.start(port);
}
