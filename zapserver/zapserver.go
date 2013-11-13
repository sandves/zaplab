package main

import (
	"flag"
	"fmt"
	"github.com/sandves/zaplab/chzap"
	"github.com/sandves/zaplab/ztorage"
	"net"
	"os"
	"os/signal"
	"runtime/pprof"
	"strings"
	"time"
)

var zaps ztorage.Zaps
var sliceZaps *ztorage.SliceZaps

var sock *net.UDPConn
var listener *net.TCPListener

var task = flag.String("task", "", "Specify which task you want to run (a - g)")
var memprofile = flag.String("memprofile", "", "Write memory profile to specified file")

func main() {

	flag.Parse()

	udpAddr, err := net.ResolveUDPAddr("udp", "224.0.1.130:10000")
	checkError(err)
	sock, err = net.ListenMulticastUDP("udp", nil, udpAddr)
	checkError(err)

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":12110")
	checkError(err)
	listener, err = net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	zaps = ztorage.NewZapStore()
	sliceZaps = ztorage.NewSliceZapStore()

	determineTaskToRun()

	writeMemProfifle()
}

func determineTaskToRun() {
	switch *task {
	case "a":
		go handleZaps(sock, nil)
	case "b":
		fmt.Println("Take a look at chzap/chzap.go")
	case "c":
		go handleZaps(sock, ztorage.Zapper(sliceZaps))
		go computeViewers("NRK1", ztorage.Zapper(sliceZaps))
		go computeViewers("TV2 Norge", ztorage.Zapper(sliceZaps))
		go computeZaps(ztorage.Zapper(sliceZaps))
	case "d":
		fmt.Println("Take a look at ztorage/slize.go")
	case "e":
		fmt.Println("Take a look at ztorage/ztorage.go")
	case "f", "g":
		go handleZaps(sock, zaps)
		go handleClient(listener, zaps)
	}
}

func handleZaps(conn *net.UDPConn, z ztorage.Zapper) {
	for {
		var buf [1024]byte
		n, _, err := conn.ReadFromUDP(buf[0:])
		checkError(err)
		str := string(buf[:n])
		strSlice := strings.Split(str, ", ")
		if len(strSlice) == 5 {
			if z != nil {
				var channelZap *chzap.ChZap = chzap.NewChZap(str)
				z.StoreZap(*channelZap)
			} else {
				fmt.Println(str)
			}
		}
	}
}

func handleClient(listener *net.TCPListener, z ztorage.Zaps) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go Subscribe(conn, z)
	}
}

func Stats(zs ztorage.Zaps) string {
	topTen := zs.TopTenChannels()
	var topTenStr string
	for i := range topTen {
		topTenStr += fmt.Sprintf("Channel %d: %s\n", (i + 1), topTen[i])
	}
	avgZapDur := zs.AverageZapDuration().String()
	topTenStr += fmt.Sprintf("\nAverage zap duration: %s\n", avgZapDur)
	return topTenStr
}

func Subscribe(conn net.Conn, z ztorage.Zaps) {
	var stats string
	for _ = range time.Tick(1 * time.Second) {
		stats = Stats(z)
		_, err := conn.Write([]byte(stats))
		if err != nil {
			conn.Close()
			break
		}
	}
}

func computeViewers(chName string, z ztorage.Zapper) {
	for _ = range time.Tick(1 * time.Second) {
		numberOfViewers := z.ComputeViewers(chName)
		fmt.Printf("%s: %d\n", chName, numberOfViewers)
	}
}

func computeZaps(z ztorage.Zapper) {
	for _ = range time.Tick(5 * time.Second) {
		fmt.Printf("Total number of zaps: %d\n", z.ComputeZaps())
	}
}

//if the memprofile flag was specified, write a heap profile to file
func writeMemProfifle() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Kill, os.Interrupt)
	<-signalChan
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		checkError(err)
		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
