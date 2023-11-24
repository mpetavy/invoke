package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	amount  = flag.Int("n", 1, "amount of parallel invocations")
	timeout = flag.Int("t", 0, "timeout before terminate command (ms)")
)

func Error(err error) bool {
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	return err != nil
}

func killAll(processes []*os.Process) {
	for _, process := range processes {
		fmt.Printf("Kill process with PID %d\n", process.Pid)

		Error(process.Signal(os.Kill))
	}
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		Error(fmt.Errorf("Please provide a custom command to invoke. Sample: invoke -n 3 cmd /c echo Hello world"))
	}

	args := os.Args[len(os.Args)-flag.NArg():]

	fmt.Print("Invoke command: ")
	for _, c := range args {
		fmt.Printf("\"%s\" ", c)
	}
	fmt.Println()

	fmt.Println()

	var processes []*os.Process

	start := time.Now()

	for i := 0; i < *amount; i++ {
		var cmd *exec.Cmd

		if len(args) > 0 {
			cmd = exec.Command(args[0], args[1:]...)
		} else {
			cmd = exec.Command(args[0])
		}

		if *amount == 1 {
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}

		err := cmd.Start()
		if err != nil {
			Error(err)
		}

		fmt.Printf("Started process with PID %d\n", cmd.Process.Pid)

		if err != nil {
			if _, ok := err.(*exec.ExitError); !ok {
				log.Fatalf("cannot run command %s: %v", args[0], err)
			}
		}

		processes = append(processes, cmd.Process)
	}

	if *timeout > 0 {
		d := time.Duration(*timeout) * time.Millisecond

		fmt.Printf("Wait max %v before terminating the processes...\n", d)
		fmt.Printf("%s\n", strings.Repeat("-", 80))

		time.Sleep(d)

		fmt.Printf("%s\n", strings.Repeat("-", 80))
		fmt.Printf("Terminate all processes\n")

		killAll(processes)
	} else {
		code := -1
		for _, process := range processes {
			state, err := process.Wait()
			if err != nil {
				Error(err)
			}

			if code == -1 {
				code = state.ExitCode()
			}
		}

		fmt.Printf("%s\n", strings.Repeat("-", 80))
		fmt.Printf("Invoke measured time: %s\n", time.Since(start))

		os.Exit(code)
	}
}
