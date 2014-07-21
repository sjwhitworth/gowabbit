##GoWabbit

A wrapper library for interacting with Vowpal Wabbit from Go. Provides functionality to:

* Start VW as a daemon, from Go
* Continually check if VW is in a healthy state
* Send predictions to VW and return the value and the class
* Stop the daemon

##Getting Started

```
// Instantiate a new Wabbit with configuration.
Wabbit := NewWabbit(26542, 10, "/usr/local/bin/vw", "model.model", true)

// Start VW as a daemon, with continuous checking to see if it is still running.
err := Vowpal.StartDaemonWabbit(true)

if err != nil {
	fmt.Println("omg the wabbits are taking over")
}

// Send a prediction to your model, and print the class, and the value.
pred, err := Vowpal.Predict("wabbits| height:32.0 weight:200.0")
fmt.Printf("Vowpal predicted a value of %v and a class of %v", pred.Val(), pred.Class())

// Kill the daemon
err = Vowpal.KillDaemonWabbit()
if err != nil {
	fmt.Println("This wabbit is immortal")
}