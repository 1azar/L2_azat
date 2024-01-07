package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func errWrapper(cmdName string, err error, info string) string {
	return fmt.Sprintf("ошибка обработки команды \"%s\":%s\nИспользование команды: %s", cmdName, err.Error(), info)
}

// ProcessKill не предупреждая убивает процесс по PID.
// Args: kill ... <pid>
func processKill(args []string) error {
	pidStr := args[len(args)-1]

	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return err
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	err = proc.Kill()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("supported commands: ")
mainLoop:
	for {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print("[DEV08-SHELL-ᓚᘏᗢ]: ", wd, "> ")
		scanner.Scan()
		line := scanner.Text()

		commands := strings.Split(line, "|")

		var cmd *exec.Cmd
		var lastCmd *exec.Cmd

		for i, cmdStr := range commands {
			cmdStr = strings.TrimSpace(cmdStr)
			args := strings.Fields(cmdStr)

			switch args[0] {
			case "kill":
				err = processKill(args)
				if err != nil {
					fmt.Println(errWrapper("kill", err, "kill <pid>"))
					continue mainLoop
				}
				continue mainLoop
			case "pwd":
				dir, err := os.Getwd()
				if err != nil {
					fmt.Println(errWrapper("pwd", err, "pwd"))
					continue mainLoop
				}
				fmt.Println(dir)
				continue mainLoop
			case "echo":
				fmt.Println(strings.Join(args[1:], " "))
				continue mainLoop
			case "cd":
				if len(args) < 2 {
					continue mainLoop
				}
				err = os.Chdir(args[1])
				if err != nil {
					fmt.Println(errWrapper("cd", err, "cd <path>"))
				}
				continue mainLoop
			case "ls":
				files, err := os.ReadDir(wd)
				if err != nil {
					fmt.Println(errWrapper("ls", err, "ls"))
					return
				}
				for _, file := range files {
					fmt.Println(file.Name())
				}
				continue mainLoop
			case "/quit":
				return
			default:
				switch runtime.GOOS {
				case "windows":
					if args[0] == "ps" {
						args[0] = "tasklist"
					}
					// оборачиваем внешнюю команду для выполнения в стандартном шелле windows
					newArgs := make([]string, 0, len(args)+1)
					newArgs = append(newArgs, "/C")
					newArgs = append(newArgs, args...)
					cmd = exec.Command("cmd", newArgs...)
				case "linux", "mac":
					cmd = exec.Command(args[0], args[1:]...)
				}
			}

			if lastCmd != nil {
				// Связываем стандартный вывод предыдущей команды с входом текущей
				cmd.Stdin, _ = lastCmd.StdoutPipe()
			}

			if i == len(commands)-1 {
				// Последняя команда - связываем ее стандартный вывод с STDOUT
				cmd.Stdout = os.Stdout
			} else {
				// Создаем pipe для связи между командами
				pr, pw := io.Pipe()
				cmd.Stdout = pw
				lastCmd.Stdin = pr
			}
			err = cmd.Start()
			if err != nil {
				fmt.Println(errWrapper("", err, ""))
				break
			}

			lastCmd = cmd
		}

		// Ждем завершения выполнения всех процессов
		if lastCmd != nil {
			lastCmd.Wait()
		}
	}
}
