#include "linux_atomic.h"

module McuSleepP @safe() {
  provides {
    interface McuSleep;
    interface McuPowerState;
  }
  uses {
    interface McuPowerOverride;
    interface BlockingIO;
  }
}

implementation {

  mcu_power_t getPowerState() {
    return 0;
  }

  async command void McuSleep.sleep() {
    __nesc_enable_interrupt();
    call BlockingIO.waitForIO();
    __nesc_disable_interrupt();
  }

  async command void McuPowerState.update() {
  }

  default async command mcu_power_t McuPowerOverride.lowestState() {
    return 1;
  }

}
