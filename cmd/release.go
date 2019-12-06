package cmd

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	units "github.com/docker/go-units"
	"github.com/google/go-github/github"
	gitconfig "github.com/tcnksm/go-gitconfig"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/oauth2"

	version "github.com/zktnotify/zktnotify/pkg/version"
)

var (
	owner = "zktnotify"
	repo  = "zktnotify"

	gitDraft           bool
	gitMessage         string
	gitTagName         string
	gitPrerelease      bool
	gitReleaseName     string
	gitTargetCommitish string
)

var Release = cli.Command{
	Name:   "release",
	Usage:  `release new version, upload to github.com`,
	Action: actionRelease,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "auth, a",
			Usage: "using name and password auth github.com",
		},
	},
	Description: `
Release a new version to github.com, require to login authentication.
There are 3 ways to authenticate:

  1. env variable: GITHUB_AUTH_TOKEN:
     GITHUB_AUTH_TOKEN=${token} zktnotify release

  2. global git config:
     git config --global github.token ${token}

  3. use github.com user name and password while releasing
`,
}

func lookupBranchName() string {
	err := error(nil)
	xpath, cpath := "", "."

	for {
		if cpath, err = filepath.Abs(cpath); err != nil {
			return xpath
		}
		if cpath == "/" {
			return xpath
		}

		files, err := ioutil.ReadDir(cpath)
		if err != nil {
			return xpath
		}
		for _, file := range files {
			if file.Name() == ".git" {
				data, _ := ioutil.ReadFile(cpath + "/.git/HEAD")
				if idx := bytes.LastIndex(data, []byte{'/'}); idx >= 0 {
					xpath = string(bytes.TrimSpace(data[idx+1:]))
				}
				return xpath
			}
		}
		cpath += "/.."
	}
}

func makeTag() error {
	var err error
	var data []byte

	gitTagName = version.Version()
	gitReleaseName = gitTagName

	data, err = input("Comment message for release and \x1b[90mCTRL-D\x1b[0m to exit:\n", true)
	if err != nil {
		return err
	}
	gitMessage = strings.TrimSpace(string(data))

	if gitTargetCommitish = lookupBranchName(); gitTargetCommitish == "" {
		return errors.New("not found git branch")
	}

	return nil
}

func input(title string, isText bool) ([]byte, error) {
	fmt.Printf("%s", title)
	var err error
	var body []byte

	if isText {
		body, err = ioutil.ReadAll(os.Stdin)
	} else {
		var data string
		fmt.Scanf("%s", &data)
		body = []byte(data)
	}
	if err != nil {
		return nil, err
	}

	return body, nil
}

func httpClient() *http.Client {
	token := ""
	ctx := context.Background()

	if token = os.Getenv("GITHUB_AUTH_TOKEN"); token == "" {
		if t, err := gitconfig.GithubToken(); err == nil && t != "" {
			token = t
		}
	}
	if token != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		return oauth2.NewClient(ctx, ts)
	}

	r := bufio.NewReader(os.Stdin)
	fmt.Print("Username for 'https://github.com': ")
	username, _ := r.ReadString('\n')

	fmt.Print("Password for 'https://github.com': ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)
	fmt.Println()

	tp := github.BasicAuthTransport{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
	}
	return tp.Client()
}

func createRelease(cli *github.Client) (*github.RepositoryRelease, error) {
	ctx := context.Background()

	release := &github.RepositoryRelease{
		TagName:         &gitTagName,
		TargetCommitish: &gitTargetCommitish,
		Name:            &gitReleaseName,
		Body:            &gitMessage,
		Draft:           &gitDraft,
		Prerelease:      &gitPrerelease,
	}

	grep, _, err := cli.Repositories.CreateRelease(ctx, owner, repo, release)
	return grep, err
}

func assets(pathname string) []string {
	// name-version-OS-ARCH[.suffix]
	// zktnotify-v1.0.0-linux-amd64
	// zktnotify-v1.0.0-windows-amd64.exe
	files := []string{}
	pattern := "zktnotify-v[0-9].[0-9].[0-9]-*"

	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		return files
	}
	for _, dir := range rd {
		if dir.Mode().IsRegular() {
			if ok, _ := regexp.MatchString(pattern, dir.Name()); ok {
				files = append(files, dir.Name())
			}
		}
	}
	return files
}

func uploadAssetes(cli *github.Client, rel *github.RepositoryRelease) error {
	ctx := context.Background()

	for _, name := range assets(".") {
		opt := &github.UploadOptions{
			Name: name,
		}
		var size float64
		if f, err := os.Stat(name); err == nil {
			size = float64(f.Size())
		}
		file, err := os.Open(name)
		if err != nil {
			return nil
		}

		fmt.Println("Uploading", name, units.HumanSize(size), "...")
		_, _, err = cli.Repositories.UploadReleaseAsset(ctx, owner, repo, *rel.ID, opt, file)
		if err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func actionRelease(c *cli.Context) error {
	if err := makeTag(); err != nil {
		log.Println(err)
		return err
	}

	gcli := github.NewClient(httpClient())

	rel, err := createRelease(gcli)
	if err != nil {
		log.Println(err)
		return err
	}

	if err := uploadAssetes(gcli, rel); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
