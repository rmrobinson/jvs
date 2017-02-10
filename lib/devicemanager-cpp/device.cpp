
#include "device.hpp"
#include "device_manager.hpp"

void jvs::Device::updateSelf(const proto::Device& device, bool suppressNotification) {
  // TODO: checks and stuff

  _device = device;

  if (suppressNotification) {
    return;
  }

  proto::WatchDevicesResponse update;
  update.set_action(proto::WatchDevicesResponse::CHANGED);
  update.mutable_device()->CopyFrom(_device);

  _dm.broadcast(update);
}
