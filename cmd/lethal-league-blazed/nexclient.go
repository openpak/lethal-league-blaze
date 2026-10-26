package main

import (
	"fmt"
	"net"
	"time"
)

// client is the middleware client stub for this title's NEX stack: it speaks
// just enough UDP to tell whether the auth (RendezVous) and secure listeners
// served by switch-nexd are up. It performs no RendezVous login and no RMC —
// full protocol coverage lives in switch-nexd behind the registry entry
// proposed in docs/nex-registry.md.
type client struct {
	authAddr   string
	secureAddr string
	timeout    time.Duration
}

func newClient(authAddr, secureAddr string) *client {
	return &client{authAddr: authAddr, secureAddr: secureAddr, timeout: 2 * time.Second}
}

// ping reports nil only when both NEX listeners accept UDP traffic. A refused
// connection (ICMP port unreachable) means nothing is listening; a timeout
// means packets go unanswered, which for a UDP listener still counts as
// "something is there" only when the dial itself succeeded — here we treat a
// clean dial as reachable, since PRUDP servers never answer a bare datagram.
func (c *client) ping() error {
	if err := probeUDP(c.authAddr, c.timeout); err != nil {
		return fmt.Errorf("nex auth %s: %w", c.authAddr, err)
	}
	if err := probeUDP(c.secureAddr, c.timeout); err != nil {
		return fmt.Errorf("nex secure %s: %w", c.secureAddr, err)
	}
	return nil
}

func probeUDP(addr string, timeout time.Duration) error {
	conn, err := net.DialTimeout("udp", addr, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	// A single empty datagram: enough to surface ICMP port-unreachable on
	// loopback when nothing listens. The write itself rarely fails; the
	// following read distinguishes "refused" from "up but silent".
	_ = conn.SetDeadline(time.Now().Add(timeout))
	if _, err := conn.Write([]byte{}); err != nil {
		return err
	}
	var one [1]byte
	if _, err := conn.Read(one[:]); err != nil {
		if ne, ok := err.(net.Error); ok && ne.Timeout() {
			return nil // silent listener: up
		}
		return err // refusedReset/closed: nothing there
	}
	return nil
}
