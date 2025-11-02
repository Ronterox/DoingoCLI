package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	const FILENAME = "./what_was_I_doing.txt"

	if len(os.Args) > 1 {
		fmt.Println(os.Args[1])
		switch os.Args[1] {
		case "now":
			date := time.Now().Format("15:04am")
			text := strings.Join(os.Args[2:], " ")
			fmt.Println("New entry: added", date+": \"", text, "\" to Currently")

			file, err := os.OpenFile(FILENAME, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				fmt.Println(err)
			}
			defer file.Close()
			fmt.Fprintf(file, "%s ║ %s [Currently]\n", date, text)
		default:
			fmt.Println("Error: Command not found")
		}
	} else {
		file, err := os.Open(FILENAME)
		if err != nil {
			if os.IsNotExist(err) {
				return
			}
			fmt.Println(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			split := strings.Split(line, "║")

			date := split[0]
			text := strings.Replace(split[1], "[Currently]", "", 1)

			tabs := strings.Repeat("\t", 22-len(date)-len(text)/3)
			fmt.Println("\t", date, "║", text, tabs, "[Currently]")
		}
	}
}
