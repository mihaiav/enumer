package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLowerCamel(t *testing.T) {
	tests := map[string]string{
		"":          "",
		"String":    "string",
		"AType":     "aType",
		"camelCase": "camelCase",
		"UPPERCASE": "UPPERCASE",
		"E1":        "e1",
		"Éclair":    "éclair",
	}
	for input, want := range tests {
		if got := lowerCamel(input); got != want {
			t.Errorf("lowerCamel(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestGeneratedXDoAtContracts(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "state.go")
	if err := os.WriteFile(source, []byte(`package test

type State int

const (
	StateReady State = iota
	StateHTTP
	StateAlias = StateReady
)
`), 0644); err != nil {
		t.Fatal(err)
	}

	var g Generator
	g.parsePackage([]string{source})
	g.generate("State", true, true, false, false, false, true, "lowercamel", "State", false)
	got := string(g.format())

	for _, want := range []string{
		`var _StateStringValues = []string{"ready", "HTTP"}`,
		`func (State) Enum() []State`,
		`func (i State) Zero() any`,
		`if !i.IsAState()`,
		`text, ok := value.(string)`,
		`if len(xa) != 1`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("generated output is missing %q:\n%s", want, got)
		}
	}
}

func TestParsePackageForTypesIncludesTests(t *testing.T) {
	var g Generator
	g.parsePackageForTypes([]string{"./internaltest/testenum"}, []string{"TestOnly"}, true)
	g.generate("TestOnly", false, false, false, false, false, false, "noop", "TestOnly", false)
	got := string(g.format())
	if !strings.Contains(got, `var _TestOnlyStringValues = []string{"First", "Second"}`) {
		t.Fatalf("test-only enum was not generated from its test package:\n%s", got)
	}
}
