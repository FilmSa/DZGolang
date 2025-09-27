package main

import (
	"strings"
	"testing"
)

func runUniqWithInput(input string, args []string) (string, string, error) {
	in := strings.NewReader(input)
	var out, errOut strings.Builder

	err := runUniq(in, &out, &errOut, args)

	return out.String(), errOut.String(), err
}

func TestUniqBasic(t *testing.T) {
	out, errOut, err := runUniqWithInput("a\na\na\nb\nc\nc\n", []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v, stderr=%q", err, errOut)
	}
	want := "a\nb\nc\n"
	if out != want {
		t.Errorf("got=%q want=%q", out, want)
	}
}

func TestCountFlag(t *testing.T) {
	out, errOut, err := runUniqWithInput("a\na\nb\n", []string{"-c"})
	if err != nil {
		t.Fatalf("unexpected error: %v, stderr=%q", err, errOut)
	}
	want := "2 a\n1 b\n"
	if out != want {
		t.Errorf("got=%q want=%q", out, want)
	}
}

func TestDuplicatesFlag(t *testing.T) {
	out, errOut, err := runUniqWithInput("a\na\nb\nc\nc\n", []string{"-d"})
	if err != nil {
		t.Fatalf("unexpected error: %v, stderr=%q", err, errOut)
	}
	want := "a\nc\n"
	if out != want {
		t.Errorf("got=%q want=%q", out, want)
	}
}

func TestUniqueFlag(t *testing.T) {
	out, errOut, err := runUniqWithInput("a\na\nb\nc\n", []string{"-u"})
	if err != nil {
		t.Fatalf("unexpected error: %v, stderr=%q", err, errOut)
	}
	want := "b\nc\n"
	if out != want {
		t.Errorf("got=%q want=%q", out, want)
	}
}

func TestIgnoreCase(t *testing.T) {
	out, errOut, err := runUniqWithInput("Hello\nhello\nHELLO\nWorld\n", []string{"-i"})
	if err != nil {
		t.Fatalf("unexpected error: %v, stderr=%q", err, errOut)
	}
	want := "Hello\nWorld\n"
	if out != want {
		t.Errorf("got=%q want=%q", out, want)
	}
}

func TestSkipFields(t *testing.T) {
	out, errOut, err := runUniqWithInput("x a\nx a\nx b\n", []string{"-f", "1"})
	if err != nil {
		t.Fatalf("unexpected error: %v, stderr=%q", err, errOut)
	}
	want := "x a\nx b\n"
	if out != want {
		t.Errorf("got=%q want=%q", out, want)
	}
}

func TestSkipChars(t *testing.T) {
	out, errOut, err := runUniqWithInput("xxa\nxxa\nxxb\n", []string{"-s", "2"})
	if err != nil {
		t.Fatalf("unexpected error: %v, stderr=%q", err, errOut)
	}
	want := "xxa\nxxb\n"
	if out != want {
		t.Errorf("got=%q want=%q", out, want)
	}
}

func TestConflictFlags(t *testing.T) {
	_, errOut, err := runUniqWithInput("a\na\n", []string{"-c", "-d"})
	if err == nil {
		t.Fatalf("ожидалось сообщение об ошибке, got errOut=%q", errOut)
	}
	if !strings.Contains(err.Error(), "нельзя одновременно") {
		t.Errorf("unexpected error message: %v", err)
	}
}
