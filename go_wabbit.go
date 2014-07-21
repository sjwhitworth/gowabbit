// Package gowabbit is a utility wrapper for Vowpal Wabbit.

package gowabbit

import (
	"errors"
	"fmt"
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

// Initialise a new Wabbit object. This controls VWP.
func NewWabbit(tcpPort, children int, binPath, modelPath string, quiet bool) *Wabbit {
	return &Wabbit{
		tcpPort:   tcpPort,
		children:  children,
		binpath:   binPath,
		modelPath: modelPath,
		quiet:     quiet,
	}
}

// If continous, runs a goroutine in the background that checks that VW is still up. If not, it panics.
func (w *Wabbit) StartDaemonWabbit(continualCheck bool) error {
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
func (w *Wabbit) KillDaemonWabbit() error {
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
func (w *Wabbit) Predict(command string) (*Prediction, error) {
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

func (p *Prediction) Class() string {
	return p.class
}

func (p *Prediction) Val() string {
	return p.val
}
