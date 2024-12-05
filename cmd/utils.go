package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

var pidPath = "app.pid"
var logPath = "output.log"
var cofigPath = "/usr/share/command.json"

func writePidFile(pid int, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%d\n", pid))
	if err != nil {
		return err
	}

	return nil
}

func getCallerFuncName() string {
	counter, _, _, success := runtime.Caller(1)

	if !success {
		println("functionName: runtime.Caller: failed")
		os.Exit(1)
	}

	return runtime.FuncForPC(counter).Name()
}

func commandErrorMessage(stderr bytes.Buffer, program string) string {
	message := string(stderr.Bytes())

	if len(message) == 0 {
		message = "the command doesn't exist: " + program + "\n"
	}

	return message
}

func printCommandError(stderr bytes.Buffer, callerFunc string, program string, args ...string) {
	printCommandErrorUbication(callerFunc, program, args...)
	fmt.Fprintf(os.Stderr, "%s", commandErrorMessage(stderr, program))
}

func printCommandErrorUbication(callerFunc string, program string, args ...string) {
	format := "error at: %s: %s %s\n"
	argsJoined := strings.Join(args, " ")
	fmt.Fprintf(os.Stderr, format, callerFunc, program, argsJoined)
}

func runCommand(callerFunc string, program string) {
	command := exec.Command("bash", "-c", program)

	stdout, err := command.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating StdoutPipe:", err)
		os.Exit(1)
	}

	stderr, err := command.StderrPipe()
	if err != nil {
		fmt.Println("Error creating StderrPipe:", err)
		os.Exit(1)
	}

	if err := command.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		os.Exit(1)
	}

	// Save the PID of the running application
	pid := command.Process.Pid
	if err := writePidFile(pid, pidPath); err != nil {
		fmt.Println("Error writing PID file:", err)
		os.Exit(1)
	}

	// // Ensure the PID file is deleted on exit
	// defer func() {
	// 	if err := os.Remove(pidPath); err != nil {
	// 		fmt.Println("Error removing PID file:", err)
	// 	}
	// }()

	fmt.Printf("Running application PID: %d\n", pid)
	// Open log file
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		os.Exit(1)
	}
	defer logFile.Close()

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	go func() {
		io.Copy(multiWriter, stdout)
	}()

	var stderrBuffer bytes.Buffer
	go func() {
		io.Copy(io.MultiWriter(&stderrBuffer, multiWriter), stderr)
	}()

	// Release the process
	if err := command.Process.Release(); err != nil {
		fmt.Println("Error releasing process:", err)
		os.Exit(1)
	}

	// if err := command.Wait(); err != nil {
	// 	printCommandError(stderrBuffer, callerFunc, program, args...)
	// 	os.Exit(1)
	// }
}

func readPidFile(pidSearchCommand string) (int, error) {
	// pid, err := os.ReadFile(pidPath)
	// if err != nil {
	// 	return 0, err
	// }
	out := Cmd(pidSearchCommand, true)

	out = bytes.TrimSpace(out)
	pid, _ := strconv.Atoi(string(out))

	return pid, nil
}

func killProcess(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return process.Kill()
}

func Cmd(cmd string, shell bool) []byte {

	if shell {
		out, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			panic("some error found")
		}
		return out
	} else {
		out, err := exec.Command(cmd).Output()
		if err != nil {
			panic("some error found")
		}
		return out
	}
}
