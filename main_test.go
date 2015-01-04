package main

import (
	"testing"
)

type launcherStub struct {
	expectedCommand string
	t               *testing.T
}

func (l launcherStub) Exec(cmd string) {
	if l.expectedCommand != cmd {
		l.t.Fail()
	}
}

func TestReadConfig(t *testing.T) {
	config := readConfig("test.toml")

	if config.Rule == nil {
		t.Error("ReadCondig not working")
	}
}

func TestCheckConfiguration(t *testing.T) {
	config := readConfig("test.toml")

	checked, err := checkConfiguration(config)
	if !checked {
		t.Fail()
	}

	if err != nil {
		t.Fail()
	}
}

func TestCheckConfiguration_Error(t *testing.T) {
	config := readConfig("err.toml")

	checked, err := checkConfiguration(config)
	if checked {
		t.Fail()
	}

	if err == nil {
		t.Fail()
	}
}

func TestRunQuery(t *testing.T) {
	config := readConfig("test.toml")
	var laucher launcherStub
	laucher.expectedCommand = "open https://COMPANY.atlassian.net/browse/SB-1234"
	laucher.t = t
	runQuery("SB-1234", laucher, config)
}

func TestRunQuery_Matches(t *testing.T) {
	config := readConfig("test.toml")
	var laucher launcherStub
	laucher.expectedCommand = "open https://github.com/bastos/super"
	laucher.t = t
	runQuery("gh:bastos/super", laucher, config)
}

func TestRunQuery_Escape(t *testing.T) {
	config := readConfig("test.toml")
	var laucher launcherStub
	laucher.expectedCommand = "open https://google.com/?q=Tiago+Bastos"
	laucher.t = t
	runQuery("google:Tiago Bastos", laucher, config)
}
