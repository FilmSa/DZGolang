package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var cFlag bool
var dFlag bool
var uFlag bool
var iFlag bool
var fFlag int
var sFlag int

func setupFlags() {
	flag.BoolVar(&cFlag, "c", false, "подсчитать количество повторений")
	flag.BoolVar(&dFlag, "d", false, "вывести только повторяющиеся")
	flag.BoolVar(&uFlag, "u", false, "вывести только уникальные")
	flag.BoolVar(&iFlag, "i", false, "игнорировать регистр")
	flag.IntVar(&fFlag, "f", 0, "пропустить num полей")
	flag.IntVar(&sFlag, "s", 0, "пропустить num символов")
}

func prepare(line string) string {
	if iFlag {
		line = strings.ToLower(line)
	}
	if fFlag > 0 {
		parts := strings.Fields(line)
		if len(parts) > fFlag {
			line = strings.Join(parts[fFlag:], " ")
		} else {
			line = ""
		}
	}
	if sFlag > 0 {
		if len(line) > sFlag {
			line = line[sFlag:]
		} else {
			line = ""
		}
	}
	return line
}

func runUniq(in io.Reader, out io.Writer, errOut io.Writer, args []string) error {
	flag.CommandLine = flag.NewFlagSet("uniq", flag.ContinueOnError)
	flag.CommandLine.SetOutput(errOut)
	setupFlags()
	if err := flag.CommandLine.Parse(args); err != nil {
		return err
	}

	cnt := 0
	if cFlag {
		cnt++
	}
	if dFlag {
		cnt++
	}
	if uFlag {
		cnt++
	}
	if cnt > 1 {
		fmt.Fprintln(errOut, "Ошибка: нельзя одновременно использовать -c, -d и -u")
		return fmt.Errorf("конфликт флагов")
	}

	scanner := bufio.NewScanner(in)
	writer := bufio.NewWriter(out)
	defer writer.Flush()

	var prevLine, prevKey string
	count := 0

	flush := func() {
		if prevLine == "" {
			return
		}
		if cFlag {
			fmt.Fprintf(writer, "%d %s\n", count, prevLine)
		} else if dFlag {
			if count > 1 {
				fmt.Fprintln(writer, prevLine)
			}
		} else if uFlag {
			if count == 1 {
				fmt.Fprintln(writer, prevLine)
			}
		} else {
			fmt.Fprintln(writer, prevLine)
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		key := prepare(line)

		if key == prevKey {
			count++
		} else {
			flush()
			prevLine = line
			prevKey = key
			count = 1
		}
	}
	flush()

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(errOut, "Ошибка чтения:", err)
		return err
	}

	return nil
}

func main() {
	if err := runUniq(os.Stdin, os.Stdout, os.Stderr, os.Args[1:]); err != nil {
		os.Exit(1)
	}
}
