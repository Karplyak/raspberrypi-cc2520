
/* Wrapper config for interfaces that operate on the metadata level of a packet.
 * This is designed for BLIP, and any radio that wishes to support BLIP needs
 * to create this configuration and provide these interfaces.
 *
 * @author Brad Campbell <bradjc@umich.edu>
 */

configuration RadioPacketMetadataC {
  provides {
    interface LowPowerListening;
    interface PacketLink;
    interface PacketAcknowledgements;
  }
}

implementation {
  components CC2520RpiRadioC;
  components RadioPacketMetadataP;

  RadioPacketMetadataP.PacketMetadata -> CC2520RpiRadioC.PacketMetadata;
  RadioPacketMetadataP.RadioAddress -> CC2520RpiRadioC.RadioAddress;

  LowPowerListening = CC2520RpiRadioC.LowPowerListening;
  PacketLink = RadioPacketMetadataP.PacketLink;
  PacketAcknowledgements = RadioPacketMetadataP.PacketAcknowledgements;
}
