package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Options struct {
	C, D, U, I bool
	F, S       int
	InputFile  string
	OutputFile string
}

func parseOptions(args []string, errOut io.Writer) (Options, error) {
	var opt Options
	fs := flag.NewFlagSet("uniq", flag.ContinueOnError)
	fs.SetOutput(errOut)

	fs.BoolVar(&opt.C, "c", false, "подсчитать количество повторений")
	fs.BoolVar(&opt.D, "d", false, "вывести только повторяющиеся")
	fs.BoolVar(&opt.U, "u", false, "вывести только уникальные")
	fs.BoolVar(&opt.I, "i", false, "игнорировать регистр")
	fs.IntVar(&opt.F, "f", 0, "пропустить num полей")
	fs.IntVar(&opt.S, "s", 0, "пропустить num символов")

	if err := fs.Parse(args); err != nil {
		return opt, err
	}

	rest := fs.Args()
	if len(rest) > 0 {
		opt.InputFile = rest[0]
	}
	if len(rest) > 1 {
		opt.OutputFile = rest[1]
	}
	if len(rest) > 2 {
		return opt, fmt.Errorf("слишком много аргументов")
	}

	cnt := 0
	if opt.C {
		cnt++
	}
	if opt.D {
		cnt++
	}
	if opt.U {
		cnt++
	}
	if cnt > 1 {
		return opt, fmt.Errorf("нельзя одновременно использовать -c, -d и -u")
	}

	return opt, nil
}

func prepare(line string, opt Options) string {
	if opt.I {
		line = strings.ToLower(line)
	}
	if opt.F > 0 {
		parts := strings.Fields(line)
		if len(parts) > opt.F {
			line = strings.Join(parts[opt.F:], " ")
		} else {
			line = ""
		}
	}
	if opt.S > 0 {
		if len(line) > opt.S {
			line = line[opt.S:]
		} else {
			line = ""
		}
	}
	return line
}

func runUniq(in io.Reader, out io.Writer, errOut io.Writer, args []string) error {
	opt, err := parseOptions(args, errOut)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(in)
	writer := bufio.NewWriter(out)
	defer writer.Flush()

	var prevLine, prevKey string
	var count int
	hasPrev := false

	flush := func() {
		if !hasPrev {
			return
		}
		if opt.C {
			fmt.Fprintf(writer, "%d %s\n", count, prevLine)
		} else if opt.D {
			if count > 1 {
				fmt.Fprintln(writer, prevLine)
			}
		} else if opt.U {
			if count == 1 {
				fmt.Fprintln(writer, prevLine)
			}
		} else {
			fmt.Fprintln(writer, prevLine)
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		key := prepare(line, opt)

		if hasPrev && key == prevKey {
			count++
		} else {
			flush()
			prevLine = line
			prevKey = key
			count = 1
			hasPrev = true
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Ошибка чтения: %v", err)
	}

	flush()
	return nil
}

func main() {
	opt, err := parseOptions(os.Args[1:], os.Stderr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка:", err)
		os.Exit(1)
	}

	var in io.Reader = os.Stdin
	if opt.InputFile != "" {
		f, err := os.Open(opt.InputFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Ошибка открытия входного файла:", err)
			os.Exit(1)
		}
		defer f.Close()
		in = f
	}

	var out io.Writer = os.Stdout
	if opt.OutputFile != "" {
		f, err := os.Create(opt.OutputFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Ошибка открытия выходного файла:", err)
			os.Exit(1)
		}
		defer f.Close()
		out = f
	}

	if err := runUniq(in, out, os.Stderr, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
