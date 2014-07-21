package gowabbit

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Wrapper function to preserve my sanity from having to write "sh -c" repeatedly.
// If sync is true, it runs the command synchronously.
// If not, it runs the command, and ignores STDIN.
func runCommand(command string, sync bool) ([]byte, error) {
	if sync {
		val, err := exec.Command("sh", "-c", command).Output()
		if err != nil {
			return []byte{}, err
		}
		return val, nil
	}

	err := exec.Command("sh", "-c", command).Start()
	if err != nil {
		return []byte{}, err
	}
	return []byte{}, nil
}

// Checks for the presence of Vowpal Wabbit on the box to the according binpath.
// It does not check that it is running correctly, or any other errors.
func (w *Wabbit) checkPresence() error {
	command := fmt.Sprintf("which %s", w.binpath)
	val, err := runCommand(command, true)
	if err != nil {
		return err
	}

	if len(val) < 1 {
		return NoVwError
	}

	return nil
}

// Checks if the VW daemon is actively running on the box
func checkRunning(w *Wabbit) (duration time.Duration) {
	pathToCheck := "pgrep vw | wc -l"
	// Sleep for a small period of time to let VW start
	time.Sleep(time.Millisecond * 500)
	for {
		select {
		case <-time.After(CheckDuration):
			val, _ := runCommand(pathToCheck, true)
			flt, err := strconv.ParseInt(strings.TrimSpace(string(val)), 0, 64)
			if err != nil {
				fmt.Println("Error parsing int")
			}
			if flt == 0 {
				panic("Vowpal is not running")
			}
		}
	}
}
