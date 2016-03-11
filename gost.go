package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/docopt/docopt.go"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
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
          -h --help                      Will display this help screen.
          -v --version                   Displays the current version of Gost.`
)

const DEFAULT_GIST_NAME = "gostfile"

func contentFromFile(file interface{}) (string, string, error) {
	bytes, err := ioutil.ReadFile(file.(string))
	if err != nil {
		return "", "", errors.New("Invalid file specified")
	}

	name := path.Base(file.(string))

	return string(bytes), name, nil
}

func contentFromStdin() (string, string, error) {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return "", "", errors.New("Cannot read from Stdin")
	}

	if fi.Mode()&os.ModeNamedPipe != 0 {
		stdin, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return "", "", errors.New("Cannot read from Stdin")
		}

		return string(stdin), DEFAULT_GIST_NAME, nil
	}

	return "", "", nil
}

func contentFromClip() (string, string, error) {
	content, err := clipboard.ReadAll()
	return content, DEFAULT_GIST_NAME, err
}

func main() {
	arguments, err := docopt.Parse(usage, nil, true, version, false)

	if err != nil {
		fmt.Println("Could not properly execute command; exiting ...")
		os.Exit(1)
	}

	var file string
	var name string
	var content string

	switch {
	case arguments["--file"] != nil:
		file = arguments["--file"].(string)
		content, name, err = contentFromFile(file)
	case arguments["--clip"]:
		content, name, err = contentFromClip()
	default:
		content, name, err = contentFromStdin()
	}

	if err != nil {
		fmt.Println(err, "; exiting ...")
		os.Exit(1)
	}

	if len(strings.TrimSpace(content)) == 0 {
		fmt.Println("No content to gost; exiting ...")
		os.Exit(1)
	}

	if arguments["--name"] != nil {
		name = arguments["--name"].(string)
	}

	token := arguments["--token"]
	if token == nil {
		token = os.Getenv("GOST")
	}

	ghc := github.NewClient(nil)
	if len(strings.TrimSpace(token.(string))) > 0 {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token.(string)},
		)

		ghc = github.NewClient(
			oauth2.NewClient(oauth2.NoContext, ts),
		)
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
			github.GistFilename(name): github.GistFile{Content: &content},
		},
	}

	fmt.Println("Gosting Gist ... ")

	gist, _, err := ghc.Gists.Create(input)
	if err != nil {
		fmt.Println("Unable to create gist:", err)
		os.Exit(1)
	}

	fmt.Println("Done!")
	fmt.Println("Gist URL:", string(*gist.HTMLURL))
	os.Exit(0)
}
