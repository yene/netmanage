package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	wol "github.com/sabhiram/go-wol"
	"github.com/tatsushid/go-fastping"
)

func main() {
	http.HandleFunc("/ping/", ping)
	http.HandleFunc("/wake/", wake)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// /wake/?mac=c8:2a:14:2c:e6:73
func wake(w http.ResponseWriter, r *http.Request) {
	// bcastInterface can be "eth0", "eth1", etc.. An empty string implies
	// that we use the default interface when sending the UDP packet (nil)
	bcastInterface := ""
	macAddr := r.URL.Query()["mac"][0]
	UDPPort := "9"
	broadcastIP := "255.255.255.255"

	err := wol.SendMagicPacket(macAddr, broadcastIP+":"+UDPPort, bcastInterface)
	if err != nil {
		fmt.Fprintf(w, "ERROR: %s\n", err.Error())
	}
	fmt.Fprintf(w, "Magic packet sent successfully to %s\n", macAddr)

}

// /ping/?ip=192.168.1.1
func ping(w http.ResponseWriter, r *http.Request) {
	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", r.URL.Query()["ip"][0])
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		fmt.Fprintf(w, "IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
	}
	p.OnIdle = func() {
		fmt.Fprintf(w, "finish")
	}
	err = p.Run()
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
}
