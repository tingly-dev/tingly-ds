package main

import (
	"net"
	"strconv"
	"testing"
	"time"
)

func TestInstanceAddressStableAndHighPort(t *testing.T) {
	a := instanceAddress("Tingly DS")
	b := instanceAddress("Tingly DS")

	if a != b {
		t.Fatalf("instanceAddress not deterministic: %q vs %q", a, b)
	}
	host, portStr, err := net.SplitHostPort(a)
	if err != nil {
		t.Fatalf("SplitHostPort: %v", err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("parse port: %v", err)
	}
	if host != "127.0.0.1" {
		t.Fatalf("expected loopback host, got %q", host)
	}
	if port < 49152 || port > 65535 {
		t.Fatalf("port %d outside dynamic range 49152-65535", port)
	}

	// Different app names should not collide trivially.
	if other := instanceAddress("Another App"); other == a {
		t.Fatalf("distinct app names collided on %q", a)
	}
}

func TestDetectPrimaryThenSecond(t *testing.T) {
	const app = "tingly-ds-test-detect"

	// Reset shared state captured by DetectInstance so prior tests don't leak.
	primaryListener = nil

	// First caller wins the lock and is detected as primary.
	if mode := DetectInstance(app); mode != InstancePrimary {
		t.Fatalf("first DetectInstance = %d, want InstancePrimary", mode)
	}
	if primaryListener == nil {
		t.Fatal("DetectInstance(Primary) left nil listener")
	}
	listener := primaryListener

	// Wire up the primary's show handler.
	show := make(chan struct{}, 1)
	ServeInstance(func() { show <- struct{}{} })

	// A second instance cannot rebind the port; DetectInstance must observe that
	// and return Second after signalling the primary.
	if mode := DetectInstance(app); mode != InstanceSecond {
		t.Fatalf("second DetectInstance = %d, want InstanceSecond", mode)
	}

	select {
	case <-show:
		// primary was asked to show its window — expected.
	case <-time.After(time.Second):
		t.Fatal("primary did not receive show signal within 1s")
	}

	// A second call to ServeInstance is a no-op because the listener was already
	// taken: it must not spawn a second accept loop.
	ServeInstance(func() { t.Error("second ServeInstance invoked handler") })

	if err := listener.Close(); err != nil {
		t.Fatalf("close listener: %v", err)
	}
	// Give the accept loop a moment to observe the closed listener.
	time.Sleep(50 * time.Millisecond)
}

func TestSignalPrimaryNobodyListening(t *testing.T) {
	// A free, very-high port almost certainly has nobody listening.
	addr := instanceAddress("tingly-ds-test-nobody-" + time.Now().Format("150405"))
	if signalPrimary(addr) {
		t.Fatal("signalPrimary returned true with no listener")
	}
}
