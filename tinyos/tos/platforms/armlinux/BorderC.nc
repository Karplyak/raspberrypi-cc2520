/* Wire together the components needed to create a border router style packet
 * exit and entry point.
 *
 * @author: Brad Campbell <bradjc@umich.edu>
 */

#include "border.h"

configuration BorderC {
}

implementation {
  components BorderP;
  components TunC;

  components MainC;
  MainC.SoftwareInit -> BorderP.SoftwareInit;

  components IPForwardingEngineP;
  IPForwardingEngineP.IPForward[ROUTE_IFACE_TUN] -> TunC.IPForward;

  components IPStackC;
  BorderP.ForwardingTable -> IPStackC.ForwardingTable;
  // BorderP.ForwardingNoLoop -> IPStackC.ForwardingEvents[ROUTE_IFACE_TUN];

  // components IPNeighborDiscoveryC;
  // BorderP.NeighborDiscovery -> IPNeighborDiscoveryC.NeighborDiscovery;

#if RPL_ROUTING
  components RplBorderRouterP, IPPacketC, IPNeighborDiscoveryC;
  RplBorderRouterP.NeighborDiscovery -> IPNeighborDiscoveryC.NeighborDiscovery;
  RplBorderRouterP.ForwardingEvents -> IPStackC.ForwardingEvents[ROUTE_IFACE_TUN];
  RplBorderRouterP.IPPacket -> IPPacketC.IPPacket;
#endif

  components DumpHopHopP;
  DumpHopHopP.NeighborDiscovery -> IPNeighborDiscoveryC.NeighborDiscovery;
  DumpHopHopP.ForwardingEvents -> IPStackC.ForwardingEvents[ROUTE_IFACE_154];
  DumpHopHopP.IPPacket -> IPPacketC.IPPacket;
}
