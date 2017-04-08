package lms

import (
	"strconv"
	"sync"
	"testing"
)

var zones = []string{
	"00:00:00:00:00:01",
	"00:00:00:00:00:02",
	"00:00:00:00:00:03",
	"00:00:00:00:00:04",
	"00:00:00:00:00:05",
	"00:00:00:00:00:06",
	"00:00:00:00:00:07",
	"00:00:00:00:00:08",
}

func TestAsync1(t *testing.T) {
	t.Log("Test Async 1")
	var wg sync.WaitGroup

	Connect("10.10.10.10:9090")
	for i := 1; i < 2; i++ {
		for zone := range zones {
			var z = zones[zone]
			wg.Add(1)
			go func() {
				GetStreamState(z)
				GetVolume(z)
				wg.Done()
			}()
		}
	}
	wg.Wait()
}

func TestAsyncVolume(t *testing.T) {
	t.Log("Test Async Volume")
	var wg sync.WaitGroup

	Connect("10.10.10.10:9090")
	for i := 1; i < 2; i++ {
		for zone := range zones {
			var z = zones[zone]
			wg.Add(1)
			go func() {
				SetVolume(z, 73)
				vol := GetVolume(z)
				if vol != 73 {
					t.Error("Zone: " + z + " volume test failed expected 73 got " + strconv.Itoa(vol))
					wg.Done()
					return
				}
				wg.Done()
			}()
		}
	}
	wg.Wait()
}

func TestClamp(t *testing.T) {
	t.Log("Testing Clamp")
	c := clamp(-10, 0, 5)
	if c != 0 {
		t.Errorf("Expected %d got %d", 0, c)
	}

	c = clamp(10, 0, 5)
	if c != 5 {
		t.Errorf("Expected %d got %d", 5, c)
	}

	c = clamp(3, 0, 5)
	if c != 3 {
		t.Errorf("Expected %d got %d", 3, c)
	}
}
