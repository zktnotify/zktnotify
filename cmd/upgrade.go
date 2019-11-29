package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/urfave/cli"
	"github.com/zktnotify/zktnotify/pkg/config"
	"github.com/zktnotify/zktnotify/pkg/xhttp"
)

var Upgrade = cli.Command{
	Name:  "upgrade",
	Usage: "upgrade update server binary",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "force, f",
			Usage: "-y/-yes",
		},
	},
	Action: actionUpdate,
}

func actionUpdate(c *cli.Context) error {
	config.NewConfig(c.String("conf"))

	return upgradeVersion(c)
}

type Release struct {
	Version   string `json:"tag_name"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
}

func githubLatestVersion(repo, name string) (release *Release, err error) {
	githubURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repo, name)

	token := githubAccessToken()
	if token == "" {
		return nil, fmt.Errorf("fetch latest version failed: github access token not configurated")
	}

	data, err := xhttp.Get(githubURL, map[string]interface{}{
		"Authorization": "token " + token,
	})
	if err != nil {
		return nil, err
	}

	fmt.Println(string(data))

	if err := json.Unmarshal(data, &release); err != nil {
		return nil, err
	}

	return release, nil
}

func githubAccessToken() string {
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		return token
	}

	return config.Config.XClient.Github.Token
}

func upgradeVersion(c *cli.Context) error {
	serverHost := hostname()
	started, err := isServerStartup()
	if err != nil {
		log.Println("get server status failed:", err)
		return err
	}
	if started {
		if !c.Bool("force") {
			log.Printf("server(%s) is started up ...\n", serverHost)
			if !askForConfirmation("Would you like to upgrade [Y/n]? ", false) {
				return nil
			}
		}
	}

	return githubUpdate(c)
}

func askForConfirmation(prompt string, _default bool) bool {
	var response string
	fmt.Print(prompt)
	_, err := fmt.Scanln(&response)
	if err != nil {
		return _default
	}
	response = strings.ToLower(response)

	if strings.HasPrefix("yes", response) {
		return true
	} else if strings.HasPrefix("not", response) {
		return false
	} else {
		return askForConfirmation(prompt, _default)
	}

	return false
}

func githubUpdate(c *cli.Context) error {
	repo, name := "zktnotify", "ctnotify"
	tag, err := githubLatestVersion(repo, name)
	if err != nil {
		fmt.Println("Update failed:", err)
		return err
	}
	/*
		if tag.Version == version {
			fmt.Println("No update available, already at the latest version!")
			return nil
		}
	*/

	fmt.Println("New version available -- ", tag.Version)
	fmt.Print(tag.Body)

	if !c.Bool("force") {
		if !askForConfirmation("Would you like to update [Y/n]? ", true) {
			return nil
		}
	}
	fmt.Printf("New version available: %s downloading ... \n", tag.Version)

	cleanVersion := tag.Version
	if strings.HasPrefix(cleanVersion, "v") {
		cleanVersion = cleanVersion[1:]
	}
	osArch := runtime.GOOS + "_" + runtime.GOARCH

	downloadURL := fmt.Sprintf("https://github.com/{repo}/{name}/releases/download/{tag}/{name}_{version}_{os_arch}.tar.gz", map[string]interface{}{
		"repo":    "codeskyblue",
		"name":    "gosuv",
		"tag":     tag.Version,
		"version": cleanVersion,
		"os_arch": osArch,
	})
	fmt.Println("Not finished yet. download from:", downloadURL)

	return nil

}
