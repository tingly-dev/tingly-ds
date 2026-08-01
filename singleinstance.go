package main

import (
	"hash/fnv"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

// Single-instance guard.
//
// Wails v3 has no built-in single-instance lock (unlike Electron's
// requestSingleInstanceLock), so a process-level guard is implemented here on
// top of the standard net package. The mechanism is the same on every platform:
//
//  1. The first instance binds a local TCP listener on 127.0.0.1:<port>. The
//     bind itself is the lock; whoever owns the port is the primary instance.
//  2. A second instance fails to bind the same port, dials the primary instead,
//     writes a "show" instruction over the connection, and exits immediately.
//  3. The primary's accept loop receives the instruction and reveals its window.
//
// TCP is chosen over Unix domain sockets because the latter need platform
// specific path/permission handling (and a TCP fallback on Windows anyway),
// whereas plain loopback TCP behaves identically across macOS, Windows and
// Linux with no extra dependencies.

// InstanceMode reports whether the current process is the primary instance.
type InstanceMode int

const (
	// InstancePrimary means no other instance was detected; the caller owns the
	// lock and should keep running.
	InstancePrimary InstanceMode = iota
	// InstanceSecond means another instance already owns the lock. The caller
	// has already signalled it and should exit.
	InstanceSecond
)

const (
	// instanceShowCmd is the single IPC instruction a second instance sends to
	// ask the primary to surface its window.
	instanceShowCmd = "show\n"
	// instanceDialTimeout bounds how long a second instance waits while trying
	// to reach the primary before giving up.
	instanceDialTimeout = 2 * time.Second
)

// instanceAddress returns the loopback host:port shared by every instance of
// this build. The port is derived from appName so different apps do not collide
// and is mapped into the high dynamic range (49152–65535), which is least
// likely to overlap with registered services.
func instanceAddress(appName string) string {
	h := fnv.New32a()
	_, _ = h.Write([]byte(appName))
	const (
		dynamicStart = 49152
		dynamicSpan  = 65535 - dynamicStart + 1
	)
	port := dynamicStart + int(h.Sum32()%dynamicSpan)
	return net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
}

// DetectInstance decides whether the current process may run.
//
// It must be called as early as possible in main, before the window or app are
// constructed. When it returns InstanceSecond the caller should exit at once —
// the primary has already been asked to surface its window. When it returns
// InstancePrimary the caller owns the lock and should later call ServeInstance
// to begin accepting show requests once its window is ready.
//
// Any failure to detect (for example the loopback interface being unavailable)
// is logged and treated as InstancePrimary so a detection glitch can never
// block the application from starting; the worst case degrades to the previous
// multi-instance behaviour.
func DetectInstance(appName string) InstanceMode {
	addr := instanceAddress(appName)

	listener, err := net.Listen("tcp", addr)
	if err == nil {
		// Nobody is listening yet — we are the primary. Hold the listener for
		// ServeInstance to take over.
		primaryListener = listener
		return InstancePrimary
	}

	// The port is taken. Either a previous instance owns it or an unrelated
	// process happens to occupy the same derived port. Try to contact it: if it
	// speaks our protocol, treat this as a second instance; if not, fall back to
	// running as primary (multi-instance) rather than refusing to start.
	if signalPrimary(addr) {
		return InstanceSecond
	}
	log.Printf("single-instance: port %s busy by unknown process; continuing as primary", addr)
	primaryListener = nil
	return InstancePrimary
}

// primaryListener holds the listener captured by DetectInstance so ServeInstance
// can reuse it instead of re-binding. It is nil when DetectInstance did not win
// the lock.
var primaryListener net.Listener

// signalPrimary connects to the primary and asks it to show its window. It
// reports whether a primary instance acknowledged the request.
func signalPrimary(addr string) bool {
	conn, err := net.DialTimeout("tcp", addr, instanceDialTimeout)
	if err != nil {
		return false
	}
	defer conn.Close()

	if _, err := conn.Write([]byte(instanceShowCmd)); err != nil {
		return false
	}
	// Wait briefly for the primary to close the connection, which signals it
	// received and handled the command.
	_ = conn.SetReadDeadline(time.Now().Add(instanceDialTimeout))
	_, _ = io.Copy(io.Discard, conn)
	return true
}

// ServeInstance begins accepting single-instance show requests on the lock
// captured by DetectInstance. onShow is invoked (from a goroutine) each time a
// second instance signals; call it after the application window exists so onShow
// can safely reveal it. It is a no-op if DetectInstance did not win the lock.
func ServeInstance(onShow func()) {
	listener := primaryListener
	if listener == nil {
		return
	}
	go serveInstanceLoop(listener, onShow)
}

func serveInstanceLoop(listener net.Listener, onShow func()) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			// Listener closed (process is shutting down); stop the loop.
			return
		}
		go handleInstanceConn(conn, onShow)
	}
}

func handleInstanceConn(conn net.Conn, onShow func()) {
	defer conn.Close()

	// Read the command. A small bounded buffer is enough: we only ever expect
	// the single "show\n" line, and we cap the read so a misbehaving peer cannot
	// stream unbounded data.
	buf := make([]byte, 64)
	_ = conn.SetReadDeadline(time.Now().Add(instanceDialTimeout))
	n, err := conn.Read(buf)
	if err != nil && n == 0 {
		return
	}
	cmd := strings.TrimSpace(string(buf[:n]))

	// Close promptly so the signalling peer's drain returns without delay; we do
	// not need to consume the rest of its output.

	if cmd == "show" {
		if onShow != nil {
			onShow()
		}
	}
}
