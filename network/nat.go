package network

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	dnetwork "github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	"github.com/multiformats/go-multiaddr"
)

type NAT struct {
	ProtocolID string
	Port       int
	Host       host.Host

	IsServerMode bool
}

func (n *NAT) setupHost() error {
	log.Println("setting up host, and listening network interface")

	sourceMultiAddr, err := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/0")
	if err != nil {
		return fmt.Errorf("failed to generate multi address: %+v", err)
	}

	h, err := libp2p.New(libp2p.ListenAddrs(sourceMultiAddr))
	if err != nil {
		return fmt.Errorf("failed to create p2p client: %+v", err)
	}

	log.Printf("Node ID: %v", h.ID().String())
	log.Printf("Listening Network Interface: %s", h.Addrs())

	n.Host = h

	return nil
}

func (n *NAT) initDHT(ctx context.Context) (*dht.IpfsDHT, error) {
	if len(n.Host.Addrs()) == 0 {
		return nil, errors.New("trying to initialize distributed hash table, but Host is not setted up")
	}

	log.Println("Initialize distributed hash table")

	var (
		options        = make([]dht.Option, 0)
		bootstrapPeers = make([]multiaddr.Multiaddr, 0)
	)

	if n.IsServerMode {
		log.Println("Node running under server mode")
		options = append(options, dht.Mode(dht.ModeServer))
	} else {
		defaultBootstrap, err := multiaddr.NewMultiaddr("/ip4/128.199.135.207/tcp/37545/p2p/12D3KooWBxeptFzwATdqNkPi3VZqRF5fDSNmcSAsBKuAAJxdNnSB")
		if err != nil {
			return nil, fmt.Errorf("failed to define default bootstrap peer: %+v", err)
		}

		bootstrapPeers = append(bootstrapPeers, defaultBootstrap)
	}

	kademliaDHT, err := dht.New(ctx, n.Host, options...)
	if err != nil {
		return nil, err
	}
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	for _, peerAddr := range bootstrapPeers {
		log.Printf("Connecting to bootstrap peer")

		addr, err := peer.AddrInfoFromP2pAddr(peerAddr)
		if err != nil {
			log.Fatalf("failed to get bootstrap address info from P2P: %+v", err)
		}

		wg.Add(1)
		go func(addr *peer.AddrInfo) {
			defer wg.Done()
			if err := n.Host.Connect(context.Background(), *addr); err != nil {
				log.Fatalf("failed to connect to boostrap peer: %+v", err)
			} else {
				log.Println("Successfully connect to bootstrap peer")
			}
		}(addr)
	}
	wg.Wait()

	log.Println("Finished to initialize distributed hash table")

	return kademliaDHT, nil
}

func (n *NAT) discoverNetwork(ctx context.Context) error {
	if len(n.Host.Addrs()) == 0 {
		return errors.New("trying to discovering network, but Host is not setted up")
	}

	log.Println("Discovering peer on network")

	dht, err := n.initDHT(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize DHT: %+v", err)
	}

	routeDiscover := routing.NewRoutingDiscovery(dht)
	dutil.Advertise(ctx, routeDiscover, n.ProtocolID)

	isConnected := false
	for !isConnected && !n.IsServerMode {
		log.Println("Searching for peers...")
		peerChan, err := routeDiscover.FindPeers(ctx, n.ProtocolID)
		if err != nil {
			log.Fatalf("failed to search peer on network: %+v", err)
		}

		for peer := range peerChan {
			if peer.ID == n.Host.ID() {
				continue
			}

			if n.Host.Network().Connectedness(peer.ID) == dnetwork.Connected {
				isConnected = true
				continue
			}

			if err := n.Host.Connect(ctx, peer); err == nil {
				log.Printf("successfullt connect to peer %s", peer.ID)
				isConnected = true
			}
		}

		time.Sleep(time.Second)
	}

	log.Println("Peer discovery complete")

	return nil
}

func initializeNAT(ctx context.Context) (*NAT, error) {
	nat := &NAT{
		ProtocolID: "/discovery/1.0.0",
	}

	flag.IntVar(&nat.Port, "port", 0, "Used to defining port")
	flag.BoolVar(&nat.IsServerMode, "server-mode", false, "Used to run in server mode")
	flag.Parse()

	if err := nat.setupHost(); err != nil {
		log.Fatalf("failed to setup host: %+v", err)
	}

	if err := nat.discoverNetwork(ctx); err != nil {
		return nil, err
	}

	return nat, nil
}
