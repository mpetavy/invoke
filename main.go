package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/mpetavy/common"
)

var (
	amount  *int
	timeout *int
)

func init() {
	common.Init(false, "1.0.3", "", "2017", "invokes a command to measure times", "mpetavy", fmt.Sprintf("https://github.com/mpetavy/%s", common.Title()), common.APACHE, nil, nil, run, 0)

	amount = flag.Int("n", 1, "amount of parallel invocations")
	timeout = flag.Int("t", 0, "timeout before terminate command (ms)")
}

func run() error {

	if flag.NArg() == 0 {
		_, err := fmt.Fprintf(os.Stdout, "Please provide a custom command to invoke. Sample: invoke -n 3 cmd /c echo Hello world")
		common.Error(err)
		common.Exit(0)
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
			panic(err)
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

		fmt.Println()
		fmt.Printf("Wait %v to terminate processes...\n", d)
		fmt.Println()

		time.Sleep(d)

		killAll(processes)
	} else {
		fmt.Printf("Wait for processes to terminate...\n")

		for _, process := range processes {
			_, err := process.Wait()
			if common.Error(err) {
				return err
			}
		}

		fmt.Printf("\nInvoke measured time: %s\n", time.Since(start))
	}

	return nil
}

func killAll(processes []*os.Process) {
	for _, process := range processes {
		fmt.Printf("Kill process with PID %d\n", process.Pid)

		common.Error(process.Signal(os.Kill))
	}
}

func main() {
	defer common.Done()

	common.Run(nil)
}
