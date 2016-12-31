
#include "bridge.hpp"
#include "device_manager.hpp"

void jvs::Bridge::updateSelf(const proto::Bridge& bridge) {
  // TODO: checks and stuff.

  _bridge = bridge;

  proto::WatchBridgesResponse update;
  update.set_action(proto::WatchBridgesResponse::CHANGED);
  update.mutable_bridge()->CopyFrom(_bridge);

  _dm.broadcast(update);
}
