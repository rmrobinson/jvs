#pragma once

namespace jvs {
namespace openzwaved {

// Sourced from multiple places.
// http://wiki.micasaverde.com/index.php/ZWave_Command_Classes
// https://github.com/OpenZWave/open-zwave-control-panel/blob/master/zwavelib.cpp

enum CommandClass {
  NoOperation = 0,
  Basic = 32,
  SwitchBinary = 37,
  SwitchMultilevel = 38,
  SwitchMultilevelV2 = 38,
  SwitchAll = 39,
  SwitchToggleBinary = 40,
  SwitchToggleMultilevel = 41,
  SensorBinary = 48,
  SensorMultilevel = 49,
  SensorMultilevelV2 = 49,
  CentralScene = 88,
  Alarm = 113,
  Battery = 128
};

}
}