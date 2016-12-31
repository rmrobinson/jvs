#pragma once

#include <condition_variable>
#include <queue>
#include <mutex>

namespace jvs {

// See https://www.justsoftwaresolutions.co.uk/threading/implementing-a-thread-safe-queue-using-condition-variables.html
template<typename t>
class ConcurrentQueue
{
public:
    inline void push(t const& data)
    {
        std::unique_lock<std::mutex> lock(_mutex);
        _q.push(data);
        lock.unlock();
        _condvar.notify_one();
    }

    inline bool empty() const
    {
        std::unique_lock<std::mutex> lock(_mutex);
        return _q.empty();
    }

    inline bool try_pop(t& popped_value)
    {
        std::unique_lock<std::mutex> lock(_mutex);
        if (_q.empty())
        {
            return false;
        }

        popped_value = _q.front();
        _q.pop();
        return true;
    }

    inline void wait_and_pop(t& popped_value)
    {
        std::unique_lock<std::mutex> lock(_mutex);
        while (_q.empty())
        {
            _condvar.wait(lock);
        }

        popped_value = _q.front();
        _q.pop();
    }

private:
    std::queue<t> _q;
    mutable std::mutex _mutex;
    std::condition_variable _condvar;
};

}

