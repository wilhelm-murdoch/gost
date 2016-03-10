package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"code.google.com/p/goauth2/oauth"
	"github.com/atotto/clipboard"
	"github.com/docopt/docopt.go"
	"github.com/google/go-github/github"
)

var (
	version = "Gost 1.2.0"
	usage   = `Gost - A simple command line utility for easily creating Gists for Github

        Usage:
         gost [--file=<file>] [--clip] [--name=<name>] [--description=<description>] [--token=<token>] [--public] [--paste]
         gost (--help | --version)

        Options:
          -t --token=<token>             Optional Github API authentication token. If excluded, your Gist will be created anonymously.
          -f --file=<file>               Create a Gist from file.
          -n --name=<name>               Optional name for your new Gist.
          -d --description=<description> Optional description for your new Gist.
          -c --clip                      Create a Gist from the contents of your clipboard.
          -p --public                    Make this Gist public [default: false].
	  -P --paste                     Will paste your latest gist to stdout and local clipboard.
          -h --help                      Will display this help screen.
          -v --version                   Displays the current version of Gost.`
)

func contentFromFile(file interface{}) (string, error) {
	fmt.Print(file)
	fmt.Print("...")
	bytes, err := ioutil.ReadFile(file.(string))
	if err != nil {
		return "", errors.New("Invalid file specified")
	}

	return string(bytes), nil
}

func contentFromStdin() (string, error) {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return "", errors.New("Cannot read from Stdin")
	}

	if fi.Mode()&os.ModeNamedPipe != 0 {
		stdin, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return "", errors.New("Cannot read from Stdin")
		}

		return string(stdin), nil
	}

	return "", nil
}

func contentFromClip() (string, error) {
	return clipboard.ReadAll()
}

func main() {
	arguments, err := docopt.Parse(usage, nil, true, version, false)

	if err != nil {
		fmt.Println("Could not properly execute command; exiting ...")
		os.Exit(1)
	}

	file := arguments["--file"]

	var content string
	switch {
	case len(file.(string)) > 0:
		content, err = contentFromFile(file)
	case arguments["--clip"]:
		content, err = contentFromClip()
	default:
		content, err = contentFromStdin()
	}

	if err != nil {
		fmt.Println(err, "; exiting ...")
		os.Exit(1)
	}

	if len(strings.TrimSpace(content)) == 0 {
		fmt.Println("No content to gost; exiting ...")
		os.Exit(1)
	}

	name := arguments["--name"]
	if name == nil {
		if arguments["--file"] != nil {
			name = path.Base(file.(string))
		} else {
			name = "gistfile"
		}
	}

	token := arguments["--token"]
	if token == nil {
		token = os.Getenv("GOST")
	}

	client := github.NewClient(nil)
	if len(strings.TrimSpace(token.(string))) > 0 {
		t := &oauth.Transport{
			Token: &oauth.Token{AccessToken: token.(string)},
		}

		client = github.NewClient(t.Client())
	}

	description := arguments["--description"]
	if description == nil {
		description = ""
	}

	public := arguments["--public"].(bool)
	desc := description.(string)

	input := &github.Gist{
		Description: &desc,
		Public:      &public,
		Files: map[github.GistFilename]github.GistFile{
			github.GistFilename(name.(string)): github.GistFile{Content: &content},
		},
	}

	fmt.Println("Gosting Gist ... ")
	fmt.Print(content)
	os.Exit(0)
	gist, _, err := client.Gists.Create(input)
	if err != nil {
		fmt.Println("Unable to create gist:", err)
		os.Exit(1)
	}

	fmt.Println("Done!")
	fmt.Println("Gist URL:", string(*gist.HTMLURL))
	os.Exit(0)
}
