package lms

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ziutek/telnet"
)

// StateType represents the state of a given audio instance. One of:
// PLAY, STOP, PAUSE, INVALID
type StateType uint8

const (
	// PLAY stream is Playing
	PLAY StateType = iota + 1
	// STOP stream is Stopped
	STOP
	// PAUSE stream is paused
	PAUSE
	// INVALID stream is in an unknown state, or the device may not even exist.
	INVALID
)

// Our telnet connection
var tel *telnet.Conn

// Mutex to sync access to the telnet connection
var mux sync.Mutex

// Default telnet timeout
const timeout = 2 * time.Second

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Clamp an integer between min/max
func clamp(v int, min int, max int) int {
	if v < min {
		v = min
	} else if v > max {
		v = max
	}
	return v
}

// Send text to the given telnet connection
func sendln(t *telnet.Conn, s string) {
	check(t.SetWriteDeadline(time.Now().Add(timeout)))
	buf := make([]byte, len(s)+1)
	copy(buf, s)
	buf[len(s)] = '\n'
	_, err := t.Write(buf)
	check(err)
}

func expect(d ...string) ([]byte, error) {
	check(tel.SetReadDeadline(time.Now().Add(timeout)))
	check(tel.SkipUntil(d...))
	data, err := tel.ReadBytes('\n')
	return data, err
}

// PauseStream Pause/Unpause the LMS stream for the given MAC
func PauseStream(s string, p bool) {
	pause := "1"
	if !p {
		pause = "0"
	}
	mux.Lock()
	sendln(tel, s+" pause "+pause)
	mux.Unlock()
}

// GetStreamState Get the stream state for the given MAC
func GetStreamState(s string) StateType {
	mux.Lock()
	sendln(tel, s+" mode ?")
	data, err := expect("mode ")
	mux.Unlock()
	check(err)

	var state StateType
	st := strings.TrimSpace(string(data[:]))
	if st == "pause" {
		state = PAUSE
	} else if st == "play" {
		state = PLAY
	} else if st == "stop" {
		state = STOP
	} else if st == "%3F" {
		state = INVALID
	}

	return state
}

// GetVolume returns the volume for a given MAC
func GetVolume(s string) int {
	mux.Lock()
	sendln(tel, s+" mixer volume ?")
	data, err := expect("mixer volume ")
	mux.Unlock()
	check(err)

	st := strings.TrimSpace(string(data[:]))
	if st != "%3F" {
		vol, err := strconv.Atoi(st)
		check(err)
		return vol
	}

	return 0
}

// SetVolume sets the volume of a given MAC
func SetVolume(s string, v int) {
	mux.Lock()
	sendln(tel, s+" mixer volume "+strconv.Itoa(clamp(v, 0, 100)))
	mux.Unlock()
}

// Connect Connects to the given Logitech Media Server on "IP:PORT"
func Connect(ip string) {
	var err error
	mux.Lock()
	tel, err = telnet.Dial("tcp", ip)
	check(err)
	tel.SetUnixWriteMode(true)
	mux.Unlock()
}
