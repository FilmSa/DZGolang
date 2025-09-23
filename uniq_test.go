package main

import (
	"strings"
	"testing"
)

func runUniqWithInput(input string, args []string) (string, string) {
	in := strings.NewReader(input)
	var out, errOut strings.Builder

	runUniq(in, &out, &errOut, args)

	return out.String(), errOut.String()
}

func TestUniqBasic(t *testing.T) {
	out, _ := runUniqWithInput("a\na\na\nb\nc\nc\n", []string{})
	want := "a\nb\nc\n"
	if out != want {
		t.Errorf("got=%q want=%q", out, want)
	}
}

func TestCountFlag(t *testing.T) {
	out, _ := runUniqWithInput("a\na\nb\n", []string{"-c"})
	want := "2 a\n1 b\n"
	if out != want {
		t.Errorf("got=%q want=%q", out, want)
	}
}

func TestDuplicatesFlag(t *testing.T) {
	out, _ := runUniqWithInput("a\na\nb\nc\nc\n", []string{"-d"})
	want := "a\nc\n"
	if out != want {
		t.Errorf("got=%q want=%q", out, want)
	}
}

func TestUniqueFlag(t *testing.T) {
	out, _ := runUniqWithInput("a\na\nb\nc\n", []string{"-u"})
	want := "b\nc\n"
	if out != want {
		t.Errorf("got=%q want=%q", out, want)
	}
}

func TestIgnoreCase(t *testing.T) {
	out, _ := runUniqWithInput("Hello\nhello\nHELLO\nWorld\n", []string{"-i"})
	want := "Hello\nWorld\n"
	if out != want {
		t.Errorf("got=%q want=%q", out, want)
	}
}

func TestSkipFields(t *testing.T) {
	out, _ := runUniqWithInput("x a\nx a\nx b\n", []string{"-f", "1"})
	want := "x a\nx b\n"
	if out != want {
		t.Errorf("got=%q want=%q", out, want)
	}
}

func TestSkipChars(t *testing.T) {
	out, _ := runUniqWithInput("xxa\nxxa\nxxb\n", []string{"-s", "2"})
	want := "xxa\nxxb\n"
	if out != want {
		t.Errorf("got=%q want=%q", out, want)
	}
}

func TestConflictFlags(t *testing.T) {
	_, errOut := runUniqWithInput("a\na\n", []string{"-c", "-d"})
	if !strings.Contains(errOut, "Ошибка") {
		t.Errorf("ожидалось сообщение об ошибке, got=%q", errOut)
	}
}
