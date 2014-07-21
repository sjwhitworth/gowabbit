package gowabbit

import (
	"fmt"
	"testing"
	"time"
)

var Vowpal = NewWabbit(26542, 13, "/usr/local/bin/vw", "/model/directory/here", true)

func TestPresence(t *testing.T) {
	err := Vowpal.checkPresence()
	if err != nil {
		t.Error(err)
	}
}

func TestBackgroundCheck(t *testing.T) {
	err := Vowpal.StartDaemonWabbit(true)
	if err != nil {
		t.Error(err)
	}
}

func TestRun(t *testing.T) {
	err := Vowpal.StartDaemonWabbit(false)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Second * 2)
	val, err := Vowpal.Predict("hello horse| test:junk")
	fmt.Printf("Got back [%v, %s] from Vowpal\n", val.val, val.class)
	if err != nil {
		t.Error(err)
	}

	if val.class != "horse" {
		t.Errorf(`Expected "horse", got %s`, val.class)
	}

	err = Vowpal.KillDaemonWabbit()
	if err != nil {
		t.Error(err)
	}
}
