package infer

type NodeKind int

const (
	NodePrimitive NodeKind = iota
	NodeObject
	NodeArray
	NodeUnknown
)

type SchemaNode struct {
	Path string

	Name string

	Kind NodeKind

	Types map[Kind]int

	ArrayElemTypes map[Kind]int

	Count int

	Children map[string]*SchemaNode
}

func NewNode(path, name string) *SchemaNode {
	return &SchemaNode{
		Path:           path,
		Name:           name,
		Kind:           NodeUnknown,
		Types:          map[Kind]int{},
		ArrayElemTypes: map[Kind]int{},
		Children:       map[string]*SchemaNode{},
	}
}
