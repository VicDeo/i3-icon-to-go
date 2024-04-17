package i3

import (
	"fmt"
	"strings"

	"github.com/mdirkse/i3ipc"
)

type UnifiedWorkspaces []UnifiedWorkspace

func ProcessWorkspaces() error {
	var commands []string

	i3, err := i3ipc.GetIPCSocket()
	if err != nil {
		return err
	}
	defer i3.Close()

	unifiedWorkspaces, err := getWorkspaces(i3)
	if err != nil {
		return err
	}
	// i3 placement rules create a workspace with a name that has just digits only
	// Thus a newly created window will be placed to the workspace '2'
	// and we end up with the duplicated workspaces '2' and '2:icon'
	// Merge such workspaces on the first pass
	for _, w := range unifiedWorkspaces {
		if w.Workspace.Num < 0 || w.MergeWith == nil {
			continue
		}
		targetName := w.MergeWith.Name
		for _, l := range w.Node.Leaves() {
			commands = append(commands, getMoveAppToWorkspaceCommand(l.ID, targetName))
		}
		if w.Workspace.Focused {
			commands = append(commands, getFocusCommand(targetName))
		}
	}
	for _, w := range unifiedWorkspaces {
		if w.Workspace.Num < 0 || w.MergeWith != nil {
			continue
		}
		newName := w.GetNewName()
		if w.Workspace.Name != newName {
			commands = append(commands, getRenameWorkspaceCommand(w.Workspace.Name, newName))
		}
	}

	// Send all commands at once as there could be a mess otherwise
	command := strings.Join(commands, ";")
	fmt.Println(command)
	if _, error := i3.Command(command); error != nil {
		fmt.Println(error)
		return error
	}
	return nil
}

func getWorkspaces(i3 *i3ipc.IPCSocket) (UnifiedWorkspaces, error) {
	unifiedWorkspaces := []UnifiedWorkspace{}
	tree, err := i3.GetTree()
	if err != nil {
		return nil, err
	}
	workspaces, err := i3.GetWorkspaces()
	if err != nil {
		return nil, err
	}

	workspaceNodes := tree.Workspaces()
	for _, w := range workspaces {
		unifiedWorkspaces = append(unifiedWorkspaces, UnifiedWorkspace{
			Workspace: w,
			Node:      getNodeByName(workspaceNodes, w.Name),
			MergeWith: nil,
		})
	}
	// Searching for duplicates e.g. '2' and '2:'
	for i := 0; i < len(unifiedWorkspaces); i++ {
		outer := unifiedWorkspaces[i]
		for j := i + 1; j < len(unifiedWorkspaces); j++ {
			inner := unifiedWorkspaces[j]
			if outer.IsDuplicateOf(inner) {
				if !outer.IsCustom() {
					unifiedWorkspaces[i].MergeWith = &inner.Workspace
				} else {
					unifiedWorkspaces[j].MergeWith = &outer.Workspace
				}
			}
		}
	}
	return unifiedWorkspaces, nil
}

func getNodeByName(nodes []*i3ipc.I3Node, name string) *i3ipc.I3Node {
	for _, n := range nodes {
		if n.Name == name {
			return n
		}
	}
	return nil
}
