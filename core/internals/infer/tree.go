package infer

import (
	"sort"
	"strings"
)

func BuildSchemaTree(flat map[string]*FieldStats) *SchemaNode {
	root := NewNode("", "")
	root.Kind = NodeObject

	paths := make([]string, 0, len(flat))
	for p := range flat {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	for _, path := range paths {
		fs := flat[path]
		attachPath(root, fs)
	}

	return root
}

func attachPath(root *SchemaNode, fs *FieldStats) {
	rawSegments := splitPath(fs.Path)

	curr := root
	currPath := ""

	for i := 0; i < len(rawSegments); i++ {
		seg := rawSegments[i]

		if strings.HasSuffix(seg, "[]") {
			base := strings.TrimSuffix(seg, "[]")

			curr = ensureChild(curr, base, joinPath(currPath, base))
			currPath = curr.Path

			elemName := base + "[]"
			curr = ensureChild(curr, elemName, joinPath(currPath, elemName))
			curr.Kind = NodeObject
			currPath = curr.Path

			if i == len(rawSegments)-1 {
				applyFieldStats(curr, fs)
			}
			continue
		}

		curr = ensureChild(curr, seg, joinPath(currPath, seg))
		currPath = curr.Path

		if i == len(rawSegments)-1 {
			applyFieldStats(curr, fs)
		}
	}
}

func ensureChild(parent *SchemaNode, name, path string) *SchemaNode {
	child, ok := parent.Children[name]
	if !ok {
		child = NewNode(path, name)
		parent.Children[name] = child
	}
	return child
}

func joinPath(prefix, seg string) string {
	if prefix == "" {
		return seg
	}
	return prefix + "." + seg
}

func applyFieldStats(n *SchemaNode, fs *FieldStats) {

	n.Count = maxInt(n.Count, fs.Count)

	for k, v := range fs.Types {
		n.Types[k] += v
	}
	for k, v := range fs.ArrayTypes {
		n.ArrayElemTypes[k] += v
	}

	if hasKind(fs.Types, KindArray) {
		n.Kind = NodeArray

		if hasKind(fs.ArrayTypes, KindObject) {
			elemName := n.Name + "[]"
			elemPath := n.Path + "[]"
			if _, ok := n.Children[elemName]; !ok {
				elem := NewNode(elemPath, elemName)
				elem.Kind = NodeObject
				n.Children[elemName] = elem
			}
		}
		return
	}

	if hasKind(fs.Types, KindObject) {
		n.Kind = NodeObject
		return
	}

	if len(fs.Types) == 0 {
		n.Kind = NodeUnknown
	} else {
		n.Kind = NodePrimitive
	}
}

func splitPath(path string) []string {
	if path == "" {
		return nil
	}
	return strings.Split(path, ".")
}

func hasKind(m map[Kind]int, k Kind) bool {
	_, ok := m[k]
	return ok
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
