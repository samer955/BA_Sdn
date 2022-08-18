package metrics

import (
	"bufio"
	"os/exec"
	"runtime"
	"strings"
)

func GetNumberOfOnlineUser() int {
	if runtime.GOOS == "linux" {
		output, _ := exec.Command("who").Output()
		return outputToIntLinux(string(output))
	}
	if runtime.GOOS == "windows" {
		output, _ := exec.Command("query", "user").Output()
		return outputToIntWindows(string(output))
	}
	return 0
}

func outputToIntWindows(output string) int {
	var users = 0

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		words := strings.Fields(line)
		if strings.HasPrefix(words[3], "Active") {
			users++
		}
	}
	err := scanner.Err()
	if err != nil {
		return 0
	}
	return users
}

func outputToIntLinux(output string) int {
	var users = 0

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		users++
	}
	err := scanner.Err()
	if err != nil {
		return 0
	}
	return users
}
