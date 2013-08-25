package main

import (
	"code.google.com/p/goauth2/oauth"
	"fmt"
	goopt "github.com/droundy/goopt"
	"github.com/google/go-github/github"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var (
	VERSION     = "1.0.0"
	token       = goopt.String([]string{"-t", "--token"}, "", "Optional Github API authentication token. If excluded, your Gist will be created anonymously.")
	file        = goopt.String([]string{"-f", "--file"}, "", "Create a Gist from this file.")
	name        = goopt.String([]string{"-n", "--name"}, "", "Optional name of your new Gist.")
	description = goopt.String([]string{"-d", "--description"}, "", "Optional description of your new Gist.")
	public      = goopt.Flag([]string{"-P", "--public"}, []string{"-p", "--private"}, "Make this Gist public.", "Make this Gist private. (default)")
)

func main() {
	goopt.Description = func() string {
		return "A simple command line utility for easily creating Gists for Github."
	}
	goopt.Version = VERSION
	goopt.Summary = "Creates Gists from the command line."
	goopt.Parse(nil)

	if len(strings.TrimSpace(*file)) == 0 {
		fmt.Print("Please specify a valid file with -f or --file")
		return
	}

	if _, err := os.Stat(*file); os.IsNotExist(err) {
		fmt.Printf("No such file: %s", *file)
		return
	}

	bytes, err := ioutil.ReadFile(*file)

	content := string(bytes)

	if err != nil {
		fmt.Printf("Invalid file specified: %s", *file)
		return
	}

	if len(strings.TrimSpace(*name)) == 0 {
		*name = path.Base(*file)
	}

	if len(strings.TrimSpace(*token)) == 0 {
		*token = os.Getenv("GOST")
	}

	client := github.NewClient(nil)
	if len(strings.TrimSpace(*token)) > 0 {
		t := &oauth.Transport{
			Token: &oauth.Token{AccessToken: *token},
		}

		client = github.NewClient(t.Client())
	}

	input := new(github.Gist)
	input.Description = description
	input.Public = public
	input.Files = map[github.GistFilename]github.GistFile{
		github.GistFilename(*name): github.GistFile{Content: &content},
	}

	fmt.Print("Gosting Gist ... ")

	gist, _, err := client.Gists.Create(input)

	if err != nil {
		panic(err)
	}

	fmt.Print("Done!")
	fmt.Println("")
	fmt.Printf("Gist URL: %s", string(*gist.HTMLURL))
}
