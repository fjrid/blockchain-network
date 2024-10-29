package network

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	p2pNetwork "github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/multiformats/go-multiaddr"
)

type NAT struct {
	InternalIP   net.IP
	InternalPort int
	ProtocolID   string

	Host host.Host
}

type discoveryNotifee struct {
	PeerChan chan peer.AddrInfo
}

func (dn *discoveryNotifee) HandlePeerFound(p peer.AddrInfo) {
	dn.PeerChan <- p
}

func (n *NAT) handleNewStream(stream p2pNetwork.Stream) {
	log.Println("New incoming stream")
}

func (n *NAT) registerMDNS(h *host.Host) (*discoveryNotifee, error) {
	notifee := &discoveryNotifee{PeerChan: make(chan peer.AddrInfo)}
	mdnsService := mdns.NewMdnsService(*h, "libp2p-discover", notifee)

	if err := mdnsService.Start(); err != nil {
		return nil, err
	}

	return notifee, nil
}

func (n *NAT) discoverNetwork() error {
	sourceMultiAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", n.InternalIP, n.InternalPort))
	if err != nil {
		return fmt.Errorf("failed to generate multi address: %+v", err)
	}

	h, err := libp2p.New(libp2p.ListenAddrs(sourceMultiAddr))
	if err != nil {
		return fmt.Errorf("failed to create p2p client: %+v", err)
	}

	h.SetStreamHandler(protocol.ID(n.ProtocolID), n.handleNewStream)

	log.Printf("Node ID: %s", h.ID().ShortString())
	log.Printf("Listening On: %s", h.Addrs())

	mdnsService, err := n.registerMDNS(&h)
	if err != nil {
		return fmt.Errorf("failed to start mdns: %+v", err)
	}

	// Discover Peer
	go func(dn *discoveryNotifee, h host.Host) {
		for peerInfo := range dn.PeerChan {
			err := h.Connect(context.Background(), peerInfo)
			if err != nil {
				log.Printf("failed to connect peer %s: %+v", peerInfo.Addrs, err)
			} else {
				log.Println("Successfully to connect peer")
			}
		}
	}(mdnsService, h)

	n.Host = h
	return nil
}

func initializeNAT() (*NAT, error) {
	internalPort, err := strconv.Atoi(os.Args[1])
	if err != nil {
		return nil, err
	}

	internalIP := os.Args[2]

	nat := &NAT{
		InternalPort: internalPort,
		InternalIP:   net.ParseIP(internalIP),
		ProtocolID:   "/discovery/1.0.0",
	}

	if err := nat.discoverNetwork(); err != nil {
		return nil, err
	}

	return nat, nil
}
