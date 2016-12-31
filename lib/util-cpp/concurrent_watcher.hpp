#pragma once

#include <cassert>
#include <mutex>

#include "concurrent_queue.hpp"

namespace jvs {

template<typename T>
class ConcurrentNotifier
{
public:
    class Watcher : public ConcurrentQueue<T> {
    public:
        Watcher(ConcurrentNotifier<T>& notifier) : _parent(notifier) {
            _parent.addWatcher(this);
        }
        virtual ~Watcher() {
            _parent.removeWatcher(this);
        };

    private:
        ConcurrentNotifier<T>& _parent;
    };

  inline void broadcast(const T& update) {
      std::unique_lock<std::mutex> lock(_mutex);

      for (auto& w : _watchers) {
          assert(w != nullptr);
          w->push(update);
      }
  };

protected:
    inline void addWatcher(Watcher* watcher) {
        std::unique_lock<std::mutex> lock(_mutex);
        _watchers.push_back(watcher);
    };

    inline void removeWatcher(Watcher* watcher) {
        std::unique_lock<std::mutex> lock(_mutex);
        _watchers.remove(watcher);
    };

private:
    mutable std::mutex _mutex;

    std::list<Watcher *> _watchers;
};

}
