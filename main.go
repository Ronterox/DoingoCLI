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

const FILENAME = "what_was_I_doing.txt"
const DELIMITER = "║"
const DATE_FORMAT = "02/01/06 15:04:05"
const LONGEST = len("Yesterday")
const PADDING = 88

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
		return "Yesterday " + date.Format("3:04pm")
	} else if hours < 24*7 {
		return date.Format("Mon 3:04pm")
	} else if hours < 24*365 {
		return date.Format("02/01 3:04pm")
	}

	return date.Format("02/01/06 3:04pm")
}

func padding(length int) string {
	return strings.Repeat(" ", max(PADDING-length, 0))
}

func printTask(line string) {
	split := strings.Split(line, DELIMITER)

	date, err := time.Parse(DATE_FORMAT, strings.TrimSpace(split[0]))
	if err != nil {
		fmt.Println("Error: Could not parse date", err)
		return
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
				return
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
	cuteDate = strings.Repeat(" ", LONGEST-len(first)) + "  " + first + " " + rest

	end := min(PADDING, textLength)
	line = text[0:end]
	fmt.Println(
		cit(cuteDate, CYAN), DELIMITER, "> "+line,
		padding(textLength), cit("[", MAGENTA)+"Currently"+cit("]", MAGENTA),
		cit(doneTime, YELLOW),
	)

	for end < textLength {
		start := end
		end = min(end+PADDING, textLength)

		line = text[start:end]
		padding := strings.Repeat(" ", utf8.RuneCountInString(cuteDate))

		fmt.Println(padding, DELIMITER, ">", line)
	}
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

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		printTask(line)
	}
}

func done(pattern string) {
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

	// Reverse search for @done
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		if !strings.Contains(line, "@done") && strings.Contains(line, pattern) {
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

func show(day string) bool {
	file, err := os.OpenFile(FILENAME, os.O_RDWR, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println(err)
	}
	defer file.Close()

	var targetDate time.Time
	now := time.Now()

	weekday := func(index time.Weekday) time.Time {
		// Calculate days to subtract: 0 for Tue, 1 for Wed, ..., 6 for Mon.
		// We use (weekday + 5) % 7 because Tuesday (time.Tuesday) is 2, so (2+5)%7 = 0
		daysToSubtract := -int((now.Weekday() + index) % 7)
		return now.AddDate(0, 0, daysToSubtract)
	}

	switch day {
	case "today":
		targetDate = now
	case "yesterday":
		targetDate = now.AddDate(0, 0, -1) // Subtract 1 day
	case "monday", "mon":
		targetDate = weekday(6)
	case "tuesday", "tue":
		targetDate = weekday(5)
	case "wednesday", "wed":
		targetDate = weekday(4)
	case "thursday", "thu":
		targetDate = weekday(3)
	case "friday", "fri":
		targetDate = weekday(2)
	case "saturday", "sat":
		targetDate = weekday(1)
	case "sunday", "sun":
		targetDate = weekday(0)
	default:
		return false
	}

	match := targetDate.Format("02/01/06")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		date := strings.Split(line, " ")[0]

		if date == match {
			printTask(line)
		}
	}

	return true
}

// FIX: Urgent archive destroys old ones
func archive() {
	var lines []string

	file, err := os.OpenFile(FILENAME, os.O_RDWR, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		fmt.Println(err)
	}
	defer file.Close()

	date := time.Now().Format("20060102")
	fileName := fmt.Sprintf("%s_%s.txt", FILENAME, date)

	outFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer outFile.Close()

	amount := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "@done") {
			outFile.WriteString(line + "\n")
			amount++
		} else {
			lines = append(lines, line)
		}
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = file.Truncate(0)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Fprintln(file, strings.Join(lines, "\n"))
	fmt.Println("Archived", amount, "tasks to", fileName)
}

func main() {
	if len(os.Args) > 1 {
		command := strings.ToLower(os.Args[1])
		switch command {
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
			done(strings.Join(os.Args[2:], " "))
		case "edit":
			panic("TODO: Implement edit")
		case "recent":
			recent()
		case "today", "yesterday":
			show(command)
		case "archive", "move":
			archive()
		case "begin", "reset":
			panic("TODO: Implement begin and resume, restarting an specific task")
		case "undo":
			panic("TODO: Implement undo. a binary tree that undoes the last command")
		case "check":
			panic("TODO: Implement check. this will anotate the moment that I got distracted")
		case "last":
			last()
		default:
			if !show(command) {
				fmt.Println("Error: Command not found")
			}
		}
	} else {
		recent()
	}
}
