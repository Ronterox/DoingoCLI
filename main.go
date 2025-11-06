package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

const FILENAME = "./what_was_I_doing.txt"
const BLACK = "30"
const RED = "31"
const GREEN = "32"
const YELLOW = "33"
const BLUE = "34"
const PURPLE = "35"
const CYAN = "36"
const WHITE = "37"

func cit(text string, color string) string {
	return "\033[" + color + "m" + text + "\033[0m"
}

func showRecent() {
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

		tabs := strings.Repeat(" ", max(88-len(date)-len(text), 0))
		fmt.Println("\t"+cit(date, CYAN), "║", text, tabs, cit("[", PURPLE)+"Currently"+cit("]", PURPLE))
	}
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "now", "later":
			date := time.Now().Format("15:04am")
			text := strings.Join(os.Args[2:], " ")
			fmt.Println("\t", cit("New entry:", CYAN), "added", date+": \"", text, "\" to Currently")

			file, err := os.OpenFile(FILENAME, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				fmt.Println(err)
			}
			defer file.Close()
			fmt.Fprintf(file, "%s ║ %s [Currently]\n", date, text)
		case "done", "did":
			panic("TODO: Implement done")
		case "recent", "last":
			showRecent()
		default:
			fmt.Println("Error: Command not found")
		}
	} else {
		showRecent()
	}
}
