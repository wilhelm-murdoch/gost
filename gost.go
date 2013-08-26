package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"code.google.com/p/goauth2/oauth"
	"github.com/droundy/goopt"
	"github.com/google/go-github/github"
)

var (
	VERSION     = "1.0.1"
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
		log.Fatalln("Please specify a valid file with -f or --file")
	}

	bytes, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatalln("Invalid file specified;", err)
	}
	content := string(bytes)

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

	input := &github.Gist{
		Description: description,
		Public:      public,
		Files: map[github.GistFilename]github.GistFile{
			github.GistFilename(*name): github.GistFile{Content: &content},
		},
	}

	fmt.Print("Gosting Gist ... ")

	gist, _, err := client.Gists.Create(input)
	if err != nil {
		log.Fatalln("Unable to create gist:", err)
	}

	fmt.Println("Done!")
	fmt.Println("Gist URL:", string(*gist.HTMLURL))
}
