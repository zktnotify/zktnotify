package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	units "github.com/docker/go-units"
	"github.com/google/go-github/github"
	"github.com/leaftree/ctnotify/pkg/config"
	"github.com/urfave/cli"

	version "github.com/zktnotify/zktnotify/pkg/version"
)

var Upgrade = cli.Command{
	Name:  "upgrade",
	Usage: "upgrade server to a new version",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "force, f",
			Usage: "-y/-yes",
		},
	},
	Action: actionUpgrade,
}

var (
	gcli *github.Client
)

func askForConfirmation(prompt string, _default bool) bool {
	var response string
	fmt.Print(prompt)
	_, err := fmt.Scanln(&response)
	if err != nil {
		return _default
	}
	response = strings.ToLower(response)

	switch response {
	case "y", "ye", "yes":
		return true
	case "n", "no", "not":
		return false
	default:
		return askForConfirmation(prompt, _default)
	}
}

func upgradeComfirmReplace(c *cli.Context) bool {
	serverHost := hostname()
	started, err := isServerStartup()
	if err != nil {
		fmt.Println("get server status failed:", err)
		return false
	}
	if started {
		if !c.Bool("force") {
			fmt.Printf("server(%s) is started up ...\n", serverHost)
			return askForConfirmation("Would you like to upgrade [Y/n]? ", false)
		}
	}
	return true
}

func upgradeComfirmVersion(release *github.RepositoryRelease) bool {
	if *release.TagName < version.Version() {
		fmt.Printf("current version: %s, latest version in github.com: %s\n", version.Version(), *release.TagName)
		return askForConfirmation("Current version looks even higher, confirm to upgrade [Y/n]? ", false)
	} else if *release.TagName == version.Version() {
		fmt.Println("this is the latest version")
		return false
	}
	return true
}

func githubLatestVersion(owner, repo string) (release *github.RepositoryRelease, err error) {
	grep, _, err := gcli.Repositories.GetLatestRelease(context.Background(), owner, repo)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return grep, nil
}

func downloadProgramPackage(release *github.RepositoryRelease) (string, error) {
	namePrefix := fmt.Sprintf("zktnotify-v[0-9]{1,2}.[0-9]{1,2}.[0-9]{1,3}-%s-%s", runtime.GOOS, runtime.GOARCH)
	reg := regexp.MustCompile(namePrefix)

	var asset *github.ReleaseAsset

	for _, val := range release.Assets {
		if reg.MatchString(*val.Name) {
			asset = &val
			break
		}
	}
	if asset == nil {
		return "", errors.New("No program package match you machine")
	}

	fmt.Println(*asset.Name, "will be downloaded,", units.HumanSize(float64(*asset.Size)), "...")

	reader, url, err := gcli.Repositories.DownloadReleaseAsset(context.Background(), owner, repo, *asset.ID)
	if err != nil {
		return "", err
	}
	if url != "" {
		resp, err := http.Get(url)
		if err != nil {
			return "", err
		}
		reader = resp.Body
	}
	defer reader.Close()

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	// create version save directory
	os.Mkdir(filepath.Join(config.WorkDir, ".version"), 0755)

	fileName := filepath.Join(config.WorkDir, ".version", *asset.Name)
	return fileName, ioutil.WriteFile(fileName, data, 0754)
}

func duplicateFile(new, old string) error {
	fstat, err := os.Stat(old)
	if os.IsNotExist(err) {
		return err
	}

	freader, rerr := os.OpenFile(old, os.O_RDWR, 0775)
	if rerr != nil {
		return rerr
	}
	defer freader.Close()

	fwriter, werr := os.OpenFile(new, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0775)
	if werr != nil {
		return werr
	}
	defer fwriter.Close()

	written, err := io.Copy(fwriter, freader)
	if err != nil {
		return err
	}
	if written != fstat.Size() {
		return errors.New("duplicate abort")
	}
	return nil
}

func lookupExec(pathname string) string {
	pattern := "zktnotify"
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		return ""
	}
	for _, dir := range rd {
		if dir.Mode().Perm()&os.ModeSymlink != 0 && dir.Mode().Perm()&0700 != 0 {
			if ok, _ := regexp.MatchString(pattern, dir.Name()); ok {
				return dir.Name()
			}
		}
	}
	return ""
}

func updateExecuteFiles(c *cli.Context, file string) error {
	var exec = lookupExec(config.WorkDir)

	if exec != "" {
		os.Remove(exec)
	}

	return os.Symlink(file, filepath.Join(config.WorkDir, "zktnotify"))
}

func actionUpgrade(c *cli.Context) error {
	config.NewConfig(c.String("conf"))
	gcli = github.NewClient(httpClient())

	grep, err := githubLatestVersion(owner, repo)
	if err != nil {
		return err
	}
	if !upgradeComfirmVersion(grep) {
		return nil
	}

	file, err := downloadProgramPackage(grep)
	if err != nil {
		return err
	}

	if !upgradeComfirmReplace(c) {
		return nil
	}

	if err := updateExecuteFiles(c, file); err != nil {
		fmt.Println(err)
		return err
	}

	// TODO: restart
	// TODO: process bar of downloading file
	// TODO: download in sections, don't write the whole file after downloading
	// all response
	// TODO: check package has been downloaded and verification downloaded package
	return nil
}
