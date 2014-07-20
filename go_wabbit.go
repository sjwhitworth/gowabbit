// Package go_wabbit is a utility wrapper for Vowpal Wabbit.

package go_wabbit

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	NoVwError     = errors.New("Could not find VW on this box at the specified path.")
	CheckDuration = time.Duration(time.Microsecond * 500)
)

type Wabbit struct {
	tcpPort, children  int
	modelPath, binpath string
	quiet              bool
}

// What is returned from VW after calling Predict. If no class is given, it will return an empty string.
type Prediction struct {
	val   float64
	class string
}

// If continous, runs a goroutine in the background that checks that VW is still up. If not, it panics.
func (w Wabbit) StartDaemonWabbit(continualCheck bool) error {
	err := w.checkPresence()
	if err != nil {
		return err
	}

	comm := fmt.Sprintf("vw --daemon -i %s -t --port %d", w.modelPath, w.tcpPort)
	if w.quiet {
		comm += " --quiet"
	}

	_, err = runCommand(comm, false)

	if err != nil {
		return err
	}

	if continualCheck {
		go checkRunning(w)
	}

	return nil
}

// Stops the daemon by killing it completely.
func (w Wabbit) KillDaemonWabbit() error {
	_, err := runCommand("killall vw", false)
	if err != nil {
		return err
	}
	return nil
}

// Takes a VW formatted string and returns the raw output as a string.
// @todo: Write parser to split apart tag and prediction.
// @todo: Deal with cluster.
// @todo: Use TCP socket, not command line cruft.
// @todo: Deal with multiple predictions.
func (w Wabbit) Predict(command string) (*Prediction, error) {
	comm := ` echo "` + command + fmt.Sprintf(`" | nc localhost %v`, w.tcpPort)
	val, err := runCommand(comm, true)
	if err != nil {
		return nil, err
	}

	splitString := strings.Split(string(val), " ")
	flt, _ := strconv.ParseFloat(string(splitString[0]), 64)

	var class string
	if len(splitString) == 1 {
		class = ""
	} else {
		class = strings.TrimSpace(splitString[1])
	}

	pred := &Prediction{
		val:   flt,
		class: class,
	}

	return pred, nil
}

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
func (w Wabbit) checkPresence() error {
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
func checkRunning(w Wabbit) (duration time.Duration) {
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
