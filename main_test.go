package main

import (
	"testing"
)

type LauncherStub struct {
	expectedCommand string
	t               *testing.T
}

func (l LauncherStub) Exec(cmd string) {
	if l.expectedCommand != cmd {
		l.t.Fail()
	}
}

func TestReadConfig(t *testing.T) {
	config := ReadConfig("test.toml")

	if config.Rule == nil {
		t.Error("ReadCondig not working")
	}
}

func TestCheckConfiguration(t *testing.T) {
	config := ReadConfig("test.toml")

	checked, err := CheckConfiguration(config)
	if !checked {
		t.Fail()
	}

	if err != nil {
		t.Fail()
	}
}

func TestCheckConfigurationWrongConfig(t *testing.T) {
	config := ReadConfig("err.toml")

	checked, err := CheckConfiguration(config)
	if checked {
		t.Fail()
	}

	if err == nil {
		t.Fail()
	}
}

func TestRunQuery(t *testing.T) {
	config := ReadConfig("test.toml")
	var laucher LauncherStub
	laucher.expectedCommand = "open https://COMPANY.atlassian.net/browse/SB-1234"
	laucher.t = t
	RunQuery("SB-1234", laucher, config)
}

func TestRunQueryMatches(t *testing.T) {
	config := ReadConfig("test.toml")
	var laucher LauncherStub
	laucher.expectedCommand = "open https://github.com/bastos/super"
	laucher.t = t
	RunQuery("gh:bastos/super", laucher, config)
}

func TestRunQueryEscape(t *testing.T) {
	config := ReadConfig("test.toml")
	var laucher LauncherStub
	laucher.expectedCommand = "open https://google.com/?q=Tiago+Bastos"
	laucher.t = t
	RunQuery("google:Tiago Bastos", laucher, config)
}
