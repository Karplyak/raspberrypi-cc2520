package main

import (
	"encoding/json"
	"fmt"
	"net"
	"log"
	"strings"
	"os/exec"
	"sync"
	"github.com/lab11/go-tuntap/tuntap"
	"code.google.com/p/gcfg"
)

var prefixes *PrefixManager
var tunids *TunManager

type mutexWrap struct {
	lock *sync.Mutex
}

// Map of locks that prevent the same client from having two sessions
// at the same time.
var client_locks map[string]*mutexWrap
// Lock protects the non-threadsafe lock map
var client_locks_lock sync.Mutex


func lockClient (id string) {
	client_locks_lock.Lock()

	mw := client_locks[id]
	// On first sighting of this client create a new lock
	if mw == nil {
		mw = &mutexWrap{lock: new(sync.Mutex)}
	}
	client_locks[id] = mw
	mw.lock.Lock()

	client_locks_lock.Unlock()
}

func unlockClient (id string) {
	client_locks_lock.Lock()
	mw := client_locks[id]
	mw.lock.Unlock()
	client_locks_lock.Unlock()
}

// Block on the TCP socket waiting for the client to tunnel us IPv6 packets.
func clientTCP (tcpc net.Conn, tcp_ch chan []byte, quit_ch chan int) {
	for {
		buf := make([]byte, 4096)
		rlen, err := tcpc.Read(buf)
		if err != nil {
			// Disconnect
			quit_ch <- 1
			break
		}
		tcp_ch <- buf[0:rlen]
	}
}

// Block on reading from the TUN device
// After receiving data from the TUN device it checks to see if the client
// has disconnected and if so quits.
func clientTUN (tun *tuntap.Interface, tun_ch chan []byte, quit_ch chan int,
	quit_ch2 chan int) {
	for {
		pkt, err := tun.ReadPacket()
		if err != nil {
			log.Fatal(err)
		}

		// Check if there is data in the quit channel that tells us to stop
		select {
		case quit_tun := <- quit_ch:
			if quit_tun == 1 {
				tun.Close()
				<- quit_ch2
				return
			}
		default:
		}

		tun_ch <- pkt.Packet
	}
}

// Takes care of interacting with a client
func handleClient (tcpc net.Conn) {
	buf := make([]byte, 4096)
	var err error

	var newclient ClientIdentifier

	for {
		// Read in a message
		// This should be a JSON blob
		rlen, err := tcpc.Read(buf)
		if err != nil {
			log.Fatal(err)
		}

		// Parse the JSON blob into a ClientIdentifer
		err = json.Unmarshal(buf[0:rlen], &newclient)
		if err != nil {
			continue
		}

		break
	}

	fmt.Println("Client connected", newclient.Id)

	lockClient(newclient.Id)

	// Get the unique prefix for this client
	var prefix ClientPrefix
	prefix.Prefix, err = prefixes.getPrefix(newclient.Id)
	if err != nil { log.Fatal(err) }

	// Send the client the prefix
	pbuf, err := json.Marshal(prefix)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(prefix, pbuf)
	tcpc.Write(pbuf)

	// Setup a tun interface
	tunname := tunids.getNewTunName()
	tun, err := tuntap.Open(tunname, tuntap.DevTun)
	if err != nil { log.Fatal(err) }

	// Remove the /64 portion
	prefixbase := strings.Split(prefix.Prefix, "/")

	// Enable the tun interface
	exec.Command("ifconfig", tunname, "up").Run()

	// Route all packets for that prefix to the tun interface
	exec.Command("ip", "-6", "route", "add", prefix.Prefix, "dev",
		tunname).Run()

	// Setup and run the goroutines that will handle interaction with the client
	tcp_ch := make(chan []byte)
	quit_tcp_ch := make(chan int)
	go clientTCP(tcpc, tcp_ch, quit_tcp_ch)

	tun_ch := make(chan []byte)
	tun_quit_ch := make(chan int, 10)
	tun_quit_ch2 := make(chan int)
	go clientTUN(tun, tun_ch, tun_quit_ch, tun_quit_ch2)

	// Loop while shuffling packets around
	for {
		select {
		case newpkt := <- tcp_ch:
			// Received a packet from the client via the TCP connection
			// Send it to the TUN device
			var tuntappkt tuntap.Packet
			tuntappkt.Packet = newpkt
			tun.WritePacket(&tuntappkt)

		case tunpkt := <- tun_ch:
			// Got a packet from the TUN device, destined for the client
			// Write it into the TCP connection
			tcpc.Write(tunpkt)

		case quit_tcp := <- quit_tcp_ch:
			// Got a quit signal from the TCP listener, this means that the
			// client disconnected.
			// Shut everything down
			if quit_tcp == 1 {
				fmt.Println("Client disconnected", newclient.Id)
				// Send a quit message to the TUN listener. This will block
				// until the TUN listener gets it, which will only happen
				// once a packet is routed to the now disconnected client.
				tun_quit_ch <- 1
				// Now send a UDP packet to get the TUN listener to wake it up
				serverAddr, _ := net.ResolveUDPAddr("udp6", "[" + prefixbase[0] + "3]:8765")
				con, _ := net.DialUDP("udp6", nil, serverAddr)
				con.Write([]byte("1"))
				con.Close()
				// Block until the tun thread has received the UDP packet
				tun_quit_ch2 <- 1
				// Finish up
				tunids.unsetTunName(tunname)
				unlockClient(newclient.Id)
				return
			}
		}
	}
}

// Sit in a loop accepting TCP connections from tunnel clients
func acceptTcp (tcpl net.Listener) {
	for {
		c, err := tcpl.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleClient(c)
	}
}

func main () {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	// Dummy channel that keeps this application alive
	main_quit := make(chan int)

	client_locks = make(map[string]*mutexWrap)

	// Parse the config file
	var cfg ConfigIni
	err := gcfg.ReadFileInto(&cfg, "config.ini")
	if err != nil {
		log.Fatal(err)
	}

	prefixes = Create(cfg.Server.Assignments, cfg.Server.Prefixrange)
	tunids = CreateTunIds()

	// Start the TCP listener
	l, err := net.Listen("tcp", cfg.Server.Localhost + ":" + cfg.Server.Listenport)
	if err != nil {
		log.Fatal(err)
	}

	go acceptTcp(l)

	// Wait on the accept tcp goroutine
	// This keeps the application from exiting
	<- main_quit
}
