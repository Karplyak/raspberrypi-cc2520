raspberrypi-cc2520
==================

Code, hardware, and instructions to use the TI CC2520 with the Raspberry Pi.

[
![rpi-cc2520](https://raw.github.com/bradjc/raspberrypi-cc2520/master/media/rpi-cc2520_2ndgen_500px.jpg)
](https://raw.github.com/bradjc/raspberrypi-cc2520/master/media/rpi-cc2520_2ndgen.jpg)

The CC2520 is a 802.15.4 (Zigbee) radio commonly used in low power wireless
sensor networks. Typically 802.15.4 radios are used with WSN motes and embedded
microcontrollers. A side-effect of being used in WSNs, however, is the
microcontrollers are very memory constrained. This is often a problem when
certain nodes need to keep data about a large network of nodes, for instance
routing information for all of the nodes. To remedy this we propose using a
Raspberry Pi as a mote in a WSN. This provides both ample storage (memory) and
convenient control as all the utilities in Linux are available.


Hardware
--------

I have a small PCB that fits directly on to the main header on the
RPi. It contains the CC2520, an SMA connector for an antenna, the DS2411
id chip, and three LEDs. The eagle files, gerbers and BOM can be found in the
`hardware/eagle/rpi-cc2520` folder.


Software
--------

### Kernel Module

In order to support the CC2520
you need the kernel module from https://github.com/ab500/linux-cc2520-driver.
This kernel module is designed to run on top of the Raspbian Linux distribution.

### Custom Version of Raspbian

Instead of compiling and installing the CC2520 kernel module yourself, you can
download a pre-configured
[image](http://nuclear.eecs.umich.edu/public/raspbian_cc2520_2013-09-18.zip)
of raspbian and put that on the sdcard instead.
Download the img and then install it to an SD
card ([here](http://elinux.org/RPi_Easy_SD_Card_Setup) is a helpful guide). For
reference, I use:

    dcfldd bs=4M if=raspbian_cc2520_2013-06-27.img of=/dev/sdd statusinterval=4

### TinyOS

All of the code I have for the RPi/CC2520 is based on TinyOS. You can use the
CC2520 driver without TinyOS (see the tests in the linux-cc2520-driver repo),
but for my purposes TinyOS was the best option.

To setup my workflow,
on your non-RPI computer you need a copy of the main TinyOS repository and the
tinyos folder from this repo. You also need to install the dependencies for
TinyOS. There are instructions here:
[http://docs.tinyos.net/tinywiki/index.php/Installing_TinyOS]. If you want
the simple Linux install I use, look here:
[http://energy.eecs.umich.edu/wiki/doku.php?id=tinyos_install]. The TinyOS
applications are designed to be cross compiled for the RPI.

In order for the TinyOS build system to figure out all the correct paths you
need to help it along a bit. Add the following to your `.bashrc` file:

    export TINYOS_ROOT_DIR=<path to git repo>/tinyos-main
    export TINYOS_ROOT_DIR_ADDITIONAL=<path to git repo>/raspberrypi-cc2520/tinyos:$TINYOS_ROOT_DIR_ADDITIONAL


You will also need the correct cross compiler for the RPi:
`arm-linux-gnueabi-gcc`. On Ubuntu:

    sudo apt-get install gcc-arm-linux-gnueabi

#### NesC

NesC is the first pass compiler for TinyOS. This compiler converts .nc files
into c code that gcc can handle. This code requires nesc version 1.3.5+.

#### Supported TinyOS Features

  - Gpio
  - Interrupts (high latency, can't use for timing critical operations)
  - 802.15.4 Packets
  - Active Message
  - Timers
  - DS2411
  - Busy wait
  - Command line arguments
  - Printf
  - Random numbers
  - Unix timestamps
  - Uart receive
  - TUN interface


Setup
-----

Once you have an RPi setup with Raspbian, there are various changes you may
need to make to the Raspbian install depending on what you want to do.

### Use IPv6
To use IPv6 you need to enable IPV6 on the RPI:

    $ sudo vim /etc/modules
    add ipv6 on a newline

### Enable Interface Forwarding

By default, Linux will not forward packets between interfaces. This
functionality is critical, however, if you want the RPi to act as a border
router for a  wireless network. Once interface forwarding is enabled, Linux
considers the machine to be a
router. This causes it to no longer receive IPv6 router advertisements, because
routers are typically statically configured. In most cases we'd rather not deal
with that, so we would like Linux to both accept router advertisements and to
forward packets.
To enable these you need to do
the following:

    sudo vim /etc/sysctl.conf
    uncomment the line: net.ipv6.conf.all.forwarding=1
    add the line: net.ipv6.conf.eth0.accept_ra=2



Usage
-----

Now to test the TinyOS code and the kernel module.

### Blink

Assuming you have the correct compilers, you should be able to run the following
on your desktop and have it compile successfully:

    cd tinyos-main/apps/Blink
    make rpi

To install it to the RPi you can simply do:

    make rpi install scp.<ipaddress of the rpi>

Then on the rpi:

    sudo ./BlinkAppC

Pins 7, 12, and 13 should be toggling and the LEDs on the interface board should
be blinking.


### RadioCountToLeds

To test the radio with TinyOS, run the RadioCountToLeds app on the RPi and
another mote. The basic process is the same as with the blink app above.


### Debug

If you need to, you can run the TinyOS application with GDB. To add the debug
symbols:

    make rpi debug

Then on the RPi:

    sudo gdb ./AppC


Setting Up an Actual Border Router
----------------------------------

These are the instructions for how I'm testing the RPI as a border router.
Because IPv6 support is anything but universal, this requires more manual setup
than desired.

My setup consists of using Hurricane Electric to tunnel IPv6 to the RPi and all
computers I wish to send packets to.

### Setting up the Tunnel

Go to http://www.tunnelbroker.net/ to setup a tunnel to the IP address of the
RPI. The Hurricane Electric instructions will setup an IPv6 in IPv4 tunnel
device on the RPi.

Hurricane Electric conveniently provides every tunnel it creates a /64 prefix
that they route to the end of the tunnel. This is perfect for the BorderRouter
application as the WSN nodes can use IP addresses in this range.

### Setting Up the Border Router

The RPi will be running the BorderRouter app.
The border router is required to determine the 64 bit IPv6 prefix that the
network uses. How this prefix is assigned can be handled in several ways:
IPv6 Neighbor Discovery, DHCP, or static configuration.

#### Addressing: IPv6 Neighbor Discovery

To use router solicitations and router advertisements to assign a prefix
for the network, set `BLIP_ADDR_AUTOCONF=1` in the Makefile, as well
as `BLIP_SEND_ROUTER_SOLICITATIONS=1`. This will cause the border router
to send a router solicitation in order to receive a router solicitation
and then use the prefix that it receives from the router advertisement as the
network prefix.

The next step is something must respond to the router solicitations. To do
this, setup [radvd](https://github.com/reubenhwk/radvd) to respond to router
solictitations on the interface that the border router creates. To configure
the name of the TUN device the border router creates, run the border router
app like this:

    sudo ./BorderRouterC -i tunname

#### Addressing: DHCP

Blip supports both static and dynamic IP addresses. If you wish to reduce your
burden when flashing nodes and use dynamic addresses, you need to be
running a DHCP server. Ideally you could use any router's DHCP server, but in
the likely case that isn't available, you can run a DHCP server on the RPI.

I'm using [Dibbler](http://klub.com.pl/dhcpv6/). I couldn't figure out how to
cross compile it so I downloaded it to the RPi and built it on there (yeah it
took a little while). Alternatively, you can try getting
[these directions](http://klub.com.pl/dhcpv6/doxygen/dc/dec/compilation.html#compilationCross)
for cross-compiling Dibbler to work.

    tar xf dibbler-x.x.x.tar.gz
    ./configure
    make
    sudo make install

Running Dibbler is pretty straightforward. The last key is the configuration file
located at `/etc/dibbler/server.conf`. Here is mine:

    log-colors true

    iface relay1 {
     relay tun0

     class {
      pool <ipv6/64 from Hurricane Electric>/64
     }
    }

    iface "tun0" {
     class {
      pool <ipv6/64 from Hurricane Electric>/64
     }

     client link-local fe80::212:6d52:5000:1 {
      address <ipv6/64 from Hurricane Electric>:1
      prefix <ipv6/64 from Hurricane Electric>/64
     }
    }

#### Addressing: Static

To assign a static prefix, uncomment `components StaticIPAddressC;` from `BorderRouterC.nc`
and update the prefix in `brconfig.ini`.

### Compile and Install

In `/tinyos/apps/BorderRouter`:
    
    make rpi install scp.<ip address of rpi>

### Operation

Once running the BorderRouter application will forward all IPv6 packets that are not destined
for a node in the network to the TUN interface, and then Linux can route those packets
to the Internet.

