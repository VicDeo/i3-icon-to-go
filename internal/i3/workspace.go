package i3

import (
	"i3-icon-to-go/internal/config"
	"strings"

	"github.com/mdirkse/i3ipc"
)

// We need both Workspace and I3Node methods FTW so we join them by name
type UnifiedWorkspace struct {
	Workspace i3ipc.Workspace
	Node      *i3ipc.I3Node
	MergeWith *i3ipc.Workspace
}

func (w UnifiedWorkspace) GetNewName() string {
	apps := []string{}
	for _, leaf := range w.Node.Leaves() {
		apps = append(apps, w.getAppTitle(leaf))
	}
	return config.BuildName(w.Workspace.Num, apps)
}

func (w UnifiedWorkspace) IsCustom() bool {
	idx := strings.IndexByte(w.Workspace.Name, ':')
	return idx >= 0
}

func (w UnifiedWorkspace) IsDuplicateOf(ww UnifiedWorkspace) bool {
	return w.getIdByName() == ww.getIdByName()
}

func (w UnifiedWorkspace) getAppTitle(leaf *i3ipc.I3Node) string {
	name := ""
	for _, attr := range []string{
		leaf.Name,
		leaf.Window_Properties.Instance,
		leaf.Window_Properties.Class,
		leaf.Window_Properties.Title,
	} {
		if attr == "" {
			continue
		}
		name = config.GetAppIcon(attr)
		if !config.IsNoMatchIcon(name) {
			return name
		}
		break
	}
	if leaf.Name != "" {
		return config.TrimAppName(leaf.Name)
	}
	return "?"
}

func (w UnifiedWorkspace) getIdByName() string {
	if idx := strings.IndexByte(w.Workspace.Name, ':'); idx >= 0 {
		return w.Workspace.Name[:idx]
	}
	return w.Workspace.Name
}
