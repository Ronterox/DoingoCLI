package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

const FILENAME = "./what_was_I_doing.txt"
const DELIMITER = "║"
const DATE_FORMAT = "02/01/06 15:04:05"

const (
	BLACK  = "30"
	RED    = "31"
	GREEN  = "32"
	YELLOW = "33"
	BLUE   = "34"
	PURPLE = "35"
	CYAN   = "36"
	WHITE  = "37"
)

func cit(text string, color string) string {
	return "\033[" + color + "m" + text + "\033[0m"
}

func formatDate(date time.Time) string {
	duration := time.Since(date)
	hours := duration.Hours()

	if hours < 24 {
		return date.Format("3:04pm")
	} else if hours < 48 {
		return "yesterday " + date.Format("3:04pm")
	} else if hours < 24*7 {
		return date.Format("Mon 3:04pm")
	} else if hours < 24*365 {
		return date.Format("02/01 3:04pm")
	}

	return date.Format("02/01/06 3:04pm")
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

	const LONGEST = len("yesterday")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Split(line, DELIMITER)

		date, err := time.Parse(DATE_FORMAT, strings.Trim(split[0], " "))
		if err != nil {
			fmt.Println("Error: Could not parse date", err)
			continue
		}

		text := strings.Replace(split[1], "[Currently]", "", 1)
		cuteDate := formatDate(date)

		dateParts := strings.Split(cuteDate, " ")
		first := " "
		rest := dateParts[0]
		if len(dateParts) > 1 {
			first = dateParts[0]
			rest = dateParts[1]
		}
		rest = strings.Repeat(" ", len("03:04pm")-len(rest)) + rest

		cuteDate = strings.Repeat(" ", LONGEST-len(first)) + first + " " + rest

		tabs := strings.Repeat(" ", max(88-len(cuteDate)-len(text), 0))
		fmt.Println("  "+cit(cuteDate, CYAN), DELIMITER, text, tabs, cit("[", PURPLE)+"Currently"+cit("]", PURPLE))
	}
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "now", "later":
			date := time.Now()
			text := strings.Join(os.Args[2:], " ")
			fmt.Println("\t", cit("New entry:", CYAN), "added", formatDate(date)+": \"", text, "\" to Currently")

			file, err := os.OpenFile(FILENAME, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				fmt.Println(err)
			}
			defer file.Close()
			fmt.Fprintf(file, "%s %s %s [Currently]\n", date.Format(DATE_FORMAT), DELIMITER, text)
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
