package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

const FILENAME = "./what_was_I_doing.txt"
const DELIMITER = "║"
const DATE_FORMAT = "02/01/06 15:04:05"

const (
	BLACK   = "30"
	RED     = "31"
	GREEN   = "32"
	YELLOW  = "33"
	BLUE    = "34"
	MAGENTA = "35"
	CYAN    = "36"
	WHITE   = "37"
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

func recent() {
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

		date, err := time.Parse(DATE_FORMAT, strings.TrimSpace(split[0]))
		if err != nil {
			fmt.Println("Error: Could not parse date", err)
			continue
		}

		text := split[1]
		textLength := utf8.RuneCountInString(text)
		doneTime := ""
		if strings.Contains(text, "@done") {
			regexp := regexp.MustCompile(`@done\((.*)\)`)
			match := regexp.FindStringSubmatch(text)
			if len(match) > 1 {
				doneDate, err := time.Parse(DATE_FORMAT, match[1])
				if err != nil {
					fmt.Println("Error: Could not parse date", err)
					continue
				}
				doneTime = doneDate.Sub(date).String()
			}
			text = regexp.ReplaceAllString(text, cit("@done(", RED)+cit("$1", MAGENTA)+cit(")", RED))
		}

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
		tabs := strings.Repeat(" ", 88-textLength)

		fmt.Println(
			"  "+cit(cuteDate, CYAN), DELIMITER, text,
			tabs, cit("[", MAGENTA)+"Currently"+cit("]", MAGENTA),
			cit(doneTime, YELLOW),
		)
	}
}

func done() {
	var lines []string

	file, err := os.OpenFile(FILENAME, os.O_RDWR, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		if !strings.Contains(line, "@done") {
			lines[i] = fmt.Sprintf("%s @done(%s)", line, time.Now().Format(DATE_FORMAT))
			text := strings.ReplaceAll(strings.Split(lines[i], DELIMITER)[1], "@done", cit("@done", RED))

			fmt.Println("\t", cit("Tagged:", CYAN), "added tag", cit("@done", RED), "to", strings.TrimSpace(text))
			_, err = fmt.Fprintf(file, "%s\n", strings.Join(lines, "\n"))
			if err != nil {
				fmt.Println(err)
			}
			break
		}
	}
}

func last() {
	var line string

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
		tmp := scanner.Text()
		if !strings.Contains(tmp, "@done") {
			line = tmp
		}
	}
	split := strings.Split(line, DELIMITER)
	if len(split) > 1 {
		doing := strings.TrimSpace(split[1])
		date, err := time.Parse(DATE_FORMAT, strings.TrimSpace(split[0]))
		if err != nil {
			fmt.Println("Error: Could not parse date", err)
			return
		}

		timestamp := strings.Split(formatDate(date), " ")
		day := ""
		time := timestamp[0]
		if len(timestamp) > 1 {
			day = " on " + timestamp[0]
			time = timestamp[1]
		}

		fmt.Printf("%s (at %s%s)\n", doing, time, day)
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

			_, err = fmt.Fprintf(file, "%s %s %s\n", date.Format(DATE_FORMAT), DELIMITER, text)
			if err != nil {
				fmt.Println(err)
			}
		case "done", "did":
			done()
		case "recent":
			recent()
		case "last":
			last()
		default:
			fmt.Println("Error: Command not found")
		}
	} else {
		recent()
	}
}
