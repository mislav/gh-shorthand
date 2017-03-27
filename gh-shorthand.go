package main

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/zerowidth/gh-shorthand/alfred"
	"github.com/zerowidth/gh-shorthand/config"
	"github.com/zerowidth/gh-shorthand/parser"
	"os"
	"strings"
)

func main() {
	var input string
	var items = []alfred.Item{}

	if len(os.Args) < 2 {
		input = ""
	} else {
		input = strings.Join(os.Args[1:], " ")
	}

	path, _ := homedir.Expand("~/.gh-shorthand.yml")
	cfg, err := config.LoadFromFile(path)
	if err != nil {
		items = []alfred.Item{errorItem("when loading ~/.gh-shorthand.yml", err.Error())}
		printItems(items)
		return
	}

	printItems(generateItems(cfg, input))
}

var repoIcon = octicon("repo")
var issueIcon = octicon("git-pull-request")

func generateItems(cfg *config.Config, input string) []alfred.Item {
	items := []alfred.Item{}

	if input == "" {
		return items
	}

	// input includes leading space or leading mode char followed by a space
	if len(input) > 0 {
		if input[0:1] != " " {
			return items
		}
		input = input[1:]
	}

	result := parser.Parse(cfg.RepoMap, input)
	icon := repoIcon
	usedDefault := false

	if result.Repo == "" && cfg.DefaultRepo != "" {
		result.Repo = cfg.DefaultRepo
		usedDefault = true
	}

	if result.Repo != "" {
		uid := "gh:" + result.Repo
		title := "Open " + result.Repo
		arg := "open https://github.com/" + result.Repo

		if result.Issue != "" {
			uid += "#" + result.Issue
			title += "#" + result.Issue
			arg += "/issues/" + result.Issue
			icon = issueIcon
		}

		if result.Match != "" {
			title += " (" + result.Match
			if result.Issue != "" {
				title += "#" + result.Issue
			}
			title += ")"
		}

		if usedDefault {
			title += " (default repo)"
		}

		items = append(items, alfred.Item{
			UID:   uid,
			Title: title + " on GitHub",
			Arg:   arg,
			Valid: true,
			Icon:  &icon,
		})
	}

	if !strings.ContainsAny(input, " /") {
		for key, repo := range cfg.RepoMap {
			if strings.HasPrefix(key, input) && key != result.Match && repo != result.Repo {
				items = append(items, alfred.Item{
					UID:          "gh:" + repo,
					Title:        fmt.Sprintf("Open %s (%s) on GitHub", repo, key),
					Arg:          "open https://github.com/" + repo,
					Valid:        true,
					Autocomplete: " " + key,
					Icon:         &repoIcon,
				})
			}
		}

		if input != "" {
			items = append(items, alfred.Item{
				Title:        fmt.Sprintf("Open %s... on GitHub", input),
				Autocomplete: " " + input,
				Valid:        false,
			})
		}
	}
	return items
}

func errorItem(context, msg string) alfred.Item {
	return alfred.Item{
		Title:    fmt.Sprintf("Error %s", context),
		Subtitle: msg,
		Valid:    false,
	}
}

func printItems(items []alfred.Item) {
	doc := alfred.Items{Items: items}
	if err := json.NewEncoder(os.Stdout).Encode(doc); err != nil {
		panic(err.Error())
	}
}

// octicon is relative to the alfred workflow, so this tells alfred to retrieve
// icons from there rather than relative to this go binary.
func octicon(name string) alfred.Icon {
	return alfred.Icon{
		Path: fmt.Sprintf("octicons-%s.png", name),
	}
}
