package main

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path"
	"regexp"
	"strconv"
	"strings"
)

const Version = "0.1.0"

type rule struct {
	Name    string
	Regex   string
	Prefix  string
	Command string
	Escape  bool
}

type rules struct {
	Rule []rule
}

type CommanderLauncher interface {
	Exec(cmd string)
}

type Launcher struct{}

func (l Launcher) Exec(cmd string) {
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head, parts...).Output()

	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s", out)
}

func CheckConfiguration(config rules) (bool, error) {
	for _, s := range config.Rule {
		_, err := regexp.Compile(s.Regex)

		if err != nil {
			return false, errors.New("Rule regex does not compile:" + s.Name)
		}
	}

	return true, nil
}

func RunQuery(query string, launcher CommanderLauncher, config rules) {
	var found = false
	for _, s := range config.Rule {
		var newQuery = ""

		if s.Prefix != "" {
			if strings.HasPrefix(query, s.Prefix) {
				newQuery = strings.TrimPrefix(query, s.Prefix)
			} else {
				continue
			}
		} else {
			newQuery = query
		}

		regex, err := regexp.Compile(s.Regex)

		if err != nil {
			fmt.Printf("Rule regex does not compile: %s\n", s.Name)
			os.Exit(1)
		}

		if regex.MatchString(newQuery) {
			found = true

			matches := regex.FindStringSubmatch(newQuery)
			command := s.Command

			if s.Escape {
				newQuery = url.QueryEscape(newQuery)
			}

			for index, match := range matches {
				if s.Escape {
					match = url.QueryEscape(match)
				}

				command = strings.Replace(command, "$"+strconv.Itoa(index), match, 1)
			}

			launcher.Exec(command)
		}
	}

	if !found {
		fmt.Printf("Rule not found for %s\n", query)
	}
}

func ReadConfig(file string) rules {
	var config rules

	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Error reading the config file")
		os.Exit(1)
	}

	if _, err := toml.Decode(string(data), &config); err != nil {
		fmt.Println("Error decoding the config file")
		os.Exit(1)
	}

	return config
}

func ReadFromPipe() string {
	pipedQuery := ""

	fi, err := os.Stdin.Stat()
	if err != nil {
		fmt.Println("Error reading pipe data")
		os.Exit(1)
	}

	if fi.Mode()&os.ModeNamedPipe != 0 {
		buffer, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println("Error reading pipe data")
			os.Exit(1)
		}

		pipedQuery = strings.TrimSpace(string(buffer))
	}

	return pipedQuery
}

func main() {
	argsWithoutProg := os.Args[1:]
	query := strings.Join(argsWithoutProg, " ")
	pipedQuery := ReadFromPipe()
	usr, err := user.Current()

	if err != nil {
		fmt.Println("Could not get the user's home directory")
		os.Exit(1)
	}

	configFile := path.Join(usr.HomeDir, ".super.toml")

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Printf("No such file or directory: %s\n", configFile)
		os.Exit(1)
	}

	if query == "" && pipedQuery == "" {
		fmt.Println("You should provide a query")
		os.Exit(1)
	}

	config := ReadConfig(configFile)

	switch query {
	case "-v":
		fmt.Printf("Version %s\n", Version)
		os.Exit(0)
	case "--check":
		_, err := CheckConfiguration(config)

		if err != nil {
			fmt.Print(err)
		} else {
			fmt.Printf("Everything ok\n")
		}
	default:
		var launcher Launcher

		if pipedQuery != "" {
			RunQuery(pipedQuery, launcher, config)
		} else {
			RunQuery(query, launcher, config)
		}
	}
}
