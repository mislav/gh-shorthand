package completion

import (
	"path/filepath"

	"github.com/sahilm/fuzzy"
	"github.com/zerowidth/gh-shorthand/pkg/alfred"
)

type projectDirMode int

const (
	modeEdit projectDirMode = iota
	modeTerm
)

func projectDirItems(dirs map[string]string, search string, mode projectDirMode) (items alfred.Items) {
	projects := map[string]string{}
	projectNames := []string{}

	for base, expanded := range dirs {
		if dirs, err := findProjectDirs(expanded); err == nil {
			for _, dirname := range dirs {
				short := filepath.Join(base, dirname)
				full := filepath.Join(expanded, dirname)
				projects[short] = full
				projectNames = append(projectNames, short)
			}
		} else {
			items = append(items, ErrorItem("Invalid project directory: "+base, err.Error()))
		}
	}

	if len(search) > 0 {
		sorted := fuzzy.Find(search, projectNames)
		projectNames = []string{}
		for _, result := range sorted {
			projectNames = append(projectNames, result.Str)
		}
	}

	for _, short := range projectNames {
		var item = alfred.Item{
			Title: short,
			Valid: true,
			Text:  &alfred.Text{Copy: projects[short], LargeType: projects[short]},
			Mods: &alfred.Mods{
				Alt: &alfred.ModItem{
					Valid:     true,
					Arg:       projects[short],
					Subtitle:  "Open finder in " + short,
					Icon:      finderIcon,
					Variables: alfred.Variables{"action": "finder"},
				},
			},
		}

		if mode == modeEdit {
			item.UID = "ghe:" + short
			item.Subtitle = "Edit " + short
			item.Arg = projects[short]
			item.Variables = alfred.Variables{"action": "edit"}
			item.Icon = editorIcon
			item.Mods.Cmd = &alfred.ModItem{
				Valid:     true,
				Arg:       projects[short],
				Subtitle:  "Open terminal in " + short,
				Icon:      terminalIcon,
				Variables: alfred.Variables{"action": "term"},
			}
		} else {
			item.UID = "ght:" + short
			item.Subtitle = "Open terminal in " + short
			item.Arg = projects[short]
			item.Variables = alfred.Variables{"action": "term"}
			item.Icon = terminalIcon
			item.Mods.Cmd = &alfred.ModItem{
				Valid:     true,
				Arg:       projects[short],
				Subtitle:  "Edit " + short,
				Icon:      editorIcon,
				Variables: alfred.Variables{"action": "edit"},
			}
		}

		items = append(items, item)
	}

	return
}
