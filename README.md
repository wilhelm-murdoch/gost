### Gost
***
Gost is a small command line utility written in Go. It is the result of me being frustrated (and lazy) with having to leave my terminal to create gists on Github.

It does two things:

1. Uploads a specified file, or the contents of your clipboard, as a new Gist.
2. Returns the resulting URL.


### Installation
You'll have to compile this on your own, so make sure you have the Go compiler installed on your machine. This utility was written with version `1.1.2` and it might work with earlier versions, though I have not tested it yet.

1. Clone into your `$GOPATH/src` directory.
2. Fetch all external dependancies with `go get -v`
3. Navigate into the `$GOPATH/src/github.com/wilhelm-murdoch/gost` directory and run `go install`

If all went well, the executable should now reside within `$GOPATH/bin`. If you want it available throughout your system, just add `$GOPATH/bin` to your systems' `$PATH`.

### Setup
Gost will create gists for you anonymously out of the box. However, if you want to pair your Github account with your gists, you'll first have to grab a personal API token from Github. You can get one of those [from here](https://github.com/settings/applications).

Then, you will have to do one of the following:

1. Create an enviromental variable entitled `GOST` and assign your token to it. Gost will automatically find this variable and use it for your gists.
2. Use the `--token` flag every time you invoke gost from the command line. Otherwise, your gists will be anonymous and private by default.

### Usage

You can find usage documentation with the following command:

```
Gost - A simple command line utility for easily creating Gists for Github

Usage:
 gost (--file=<file> | --clip) [--name=<name>] [--description=<description>] [--token=<token>] [--public]
 gost (--help | --version)

Options:
  -t --token=<token>             Optional Github API authentication token. If excluded, your Gist will be created anonymously.
  -f --file=<file>               Create a Gist from file.
  -n --name=<name>               Optional name for your new Gist.
  -d --description=<description> Optional description for your new Gist.
  -c --clip                      Create a Gist from the contents of your clipboard.
  -p --public                    Make this Gist public [default: false].
  -h --help                      Will display this help screen.
  -v --version                   Displays the current version of Gost.
```

### Examples

Create a private gist:

```
$: gost --file /path/to/by/file.txt
Gosting Gist ... Done!
Gist URL: https://gist.github.com/234234232
```

Create a public gist:

```
$: gost --file /path/to/by/file.txt --public
Gosting Gist ... Done!
Gist URL: https://gist.github.com/234234232
```

Create a public gist with a custom name and description:

```
$: gost --file /path/to/by/file.txt --public --name 'My Gosted Gist' --description 'This is quite handy!'
Gosting Gist ... Done!
Gist URL: https://gist.github.com/234234232
```

Create a gist from the contents of your clipboard:

```
$: gost --clip
Gosting Gist ... Done!
Gist URL: https://gist.github.com/234234232

```
