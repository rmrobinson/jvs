
#include "bridge.hpp"
#include "device_manager.hpp"

void jvs::Bridge::updateSelf(const proto::Bridge& bridge,
                             bool suppressNotification) {
  // TODO: checks and stuff.

  _bridge = bridge;

  if (suppressNotification) {
    return;
  }

  proto::WatchBridgesResponse update;
  update.set_action(proto::WatchBridgesResponse::CHANGED);
  update.mutable_bridge()->CopyFrom(_bridge);

  _deviceManager.broadcast(update);
}
