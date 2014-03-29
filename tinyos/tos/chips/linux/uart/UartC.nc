
/* Provides an easy way to read from the serial port /dev/ttyAMA0 on the RPi.
 * Be sure to check that the serial port isn't being used by the kernel, etc.
 * if you get errors.
 */

configuration UartC {
  provides {
    interface UartBuffer;
    interface UartConfig;
  }
}

implementation {
  components UartP;
  components MainC;
  components new IOFileC();
  components UnixTimeC;
  components HplBcm2835GeneralIOC;

  MainC.SoftwareInit -> UartP.SoftwareInit;

  UartP.IO -> IOFileC.IO;
  UartP.UnixTime -> UnixTimeC.UnixTime;
  UartP.RXPin -> HplBcm2835GeneralIOC.Port1_10;

  UartBuffer = UartP.UartBuffer;
  UartConfig = UartP.UartConfig;
}
