# runner

用于执行后台任务或Cron作业。

```go
// Runner runs a set of tasks within a given timeout and can be
// shut down on an operating system interrupt.
type Runner struct {
	// interrupt channel reports a signal from the
	// operating system.
	interrupt chan os.Signal

	// complete channel reports that processing is done.
	complete chan error

	// timeout reports that time has run out.
	timeout <-chan time.Time

	// tasks holds a set of functions that are executed
	// synchronously in index order.
	tasks []func(int)
}

// ErrTimeout is returned when a value is received on the timeout channel.
var ErrTimeout = errors.New("received timeout")

// ErrInterrupt is returned when an event from the OS is received.
var ErrInterrupt = errors.New("received interrupt")

// New returns a new ready-to-use Runner.
func New(d time.Duration) *Runner {
	return &Runner{
		interrupt: make(chan os.Signal, 1),
		complete:  make(chan error),
		timeout:   time.After(d),
	}
}

// Add attaches tasks to the Runner. A task is a function that
// takes an int ID.
func (r *Runner) Add(tasks ...func(int)) {
	r.tasks = append(r.tasks, tasks...)
}

// Start runs all tasks and monitors channel events.
func (r *Runner) Start() error {
	// We want to receive all interrupt based signals.
	signal.Notify(r.interrupt, os.Interrupt)

	// Run the different tasks on a different goroutine.
	go func() {
		r.complete <- r.run()
	}()

	select {
	// Signaled when processing is done.
	case err := <-r.complete:
		return err

	// Signaled when we run out of time.
	case <-r.timeout:
		return ErrTimeout
	}
}

// run executes each registered task.
func (r *Runner) run() error {
	for id, task := range r.tasks {
		// Check for an interrupt signal from the OS.
		if r.gotInterrupt() {
			return ErrInterrupt
		}

		// Execute the registered task.
		task(id)
	}

	return nil
}

// gotInterrupt verifies if the interrupt signal has been issued.
func (r *Runner) gotInterrupt() bool {
	select {
	// Signaled when an interrupt event is sent.
	case <-r.interrupt:
		// Stop receiving any further signals.
		signal.Stop(r.interrupt)
		return true

	// Continue running as normal.
	default:
		return false
	}
}
```

main程序

```go
// timeout is the number of second the program has to finish.
const timeout = 3 * time.Second

// main is the entry point for the program.
func main() {
	log.Println("Starting work.")

	// Create a new timer value for this run.
	r := runner.New(timeout)

	// Add the tasks to be run.
	r.Add(createTask(), createTask(), createTask())

	// Run the tasks and handle the result.
	if err := r.Start(); err != nil {
		switch err {
		case runner.ErrTimeout:
			log.Println("Terminating due to timeout.")
			os.Exit(1)
		case runner.ErrInterrupt:
			log.Println("Terminating due to interrupt.")
			os.Exit(2)
		}
	}

	log.Println("Process ended.")
}

// createTask returns an example task that sleeps for the specified
// number of seconds based on the id.
func createTask() func(int) {
	return func(id int) {
		log.Printf("Processor - Task #%d.", id)
		time.Sleep(time.Duration(id) * time.Second)
	}
}
```

