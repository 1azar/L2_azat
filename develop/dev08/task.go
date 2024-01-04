package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

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

			switch runtime.GOOS {
			case "windows":
				switch args[0] {
				case "pwd":
					args[0] = "cd"
				case "cd":
					err := os.Chdir(args[1])
					if err != nil {
						fmt.Println("Ошибка смены директории: ", err)
					}
					continue
				case "kill":
					args = make([]string, 2)
					args[0], args[1] = "echo", "( ´･･)ﾉ(._.`) not implemented"
				case "fork":
					args = make([]string, 2)
					args[0], args[1] = "echo", "( ´･･)ﾉ(._.`) not implemented"
				}
				newArgs := make([]string, 0, len(args)+1)
				newArgs = append(newArgs, "/C")
				newArgs = append(newArgs, args...)
				cmd = exec.Command("cmd", newArgs...)
			default: // linux mac
				cmd = exec.Command(args[0], args[1:]...)
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
			err := cmd.Start()
			if err != nil {
				fmt.Println("Ошибка выполнения команды:", err)
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

//package main
//
//import (
//	"bufio"
//	"fmt"
//	"log"
//	"os"
//	"os/exec"
//	"runtime"
//	"strings"
//)
//
//func main() {
//	scanner := bufio.NewScanner(os.Stdin)
//
//	for {
//		wd, err := os.Getwd()
//		if err != nil {
//			log.Fatal(err)
//		}
//		fmt.Print("dev08-shell: ", wd, "> ")
//		scanner.Scan()
//		line := scanner.Text()
//
//		// Обработка команды
//		args := strings.Fields(line)
//		if len(args) == 0 {
//			continue
//		}
//
//		switch args[0] {
//		case "cd":
//			if len(args) > 1 {
//				err := os.Chdir(args[1])
//				if err != nil {
//					fmt.Println("Error changing directory:", err)
//				}
//			} else {
//				fmt.Println("Usage: cd <directory>")
//			}
//		case "pwd":
//			cwd, err := os.Getwd()
//			if err != nil {
//				fmt.Println("Error getting current directory:", err)
//			} else {
//				fmt.Println(cwd)
//			}
//		case "echo":
//			fmt.Println(strings.Join(args[1:], " "))
//		case "kill":
//			fmt.Println("( ´･･)ﾉ(._.`) not implemented")
//		case "fork":
//			fmt.Println("( ´･･)ﾉ(._.`) not implemented")
//		case "ps":
//			cmd := exec.сommand("ps", "aux")
//			cmd.Stdout = os.Stdout
//			cmd.Stderr = os.Stderr
//			err := cmd.Run()
//			if err != nil {
//				fmt.Println("Error running ps command:", err)
//			}
//		default:
//			runCommand(args)
//		}
//	}
//}
//
//func runCommand(args []string) {
//	var cmd *exec.Cmd
//	switch runtime.GOOS {
//	case "windows":
//		newArgs := make([]string, 0, len(args)+1)
//		newArgs = append(newArgs, "/C")
//		newArgs = append(newArgs, args...)
//		cmd = exec.сommand("cmd", newArgs...)
//	default:
//		cmd = exec.сommand(args[0], args[1:]...)
//	}
//	cmd.Stdout = os.Stdout
//	cmd.Stderr = os.Stderr
//
//	err := cmd.Run()
//	cmd.Start()
//	if err != nil {
//		fmt.Println("Error running command:", err)
//	}
//}
