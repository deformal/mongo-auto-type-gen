package render

import (
	"fmt"
	"sort"
	"strings"

	"github.com/deformal/mongo-auto-type-gen/core/internals/infer"
)

type TSOptions struct {
	RequiredThreshold float64
	DateAs            string
	ObjectIDAs        string
	NullPolicy        string
	UseInterface      bool
	RootTypeName      string
	AllUsedTypeNames  map[string]bool // All type names used across all collections (root + embedded)
}

func RenderTypeScript(root *infer.SchemaNode, totalDocs int, opt TSOptions) string {
	if opt.RootTypeName == "" {
		opt.RootTypeName = "Root"
	}
	if opt.NullPolicy == "" {
		opt.NullPolicy = "optional"
	}

	var objectNodes []*infer.SchemaNode
	collectObjectNodes(root, &objectNodes)

	sort.Slice(objectNodes, func(i, j int) bool {
		if len(objectNodes[i].Path) != len(objectNodes[j].Path) {
			return len(objectNodes[i].Path) < len(objectNodes[j].Path)
		}
		return objectNodes[i].Path < objectNodes[j].Path
	})

	typeNames := map[string]string{}
	typeNames[""] = opt.RootTypeName

	// Mark root type name as used
	if opt.AllUsedTypeNames != nil {
		opt.AllUsedTypeNames[opt.RootTypeName] = true
	}

	for _, n := range objectNodes {
		typeName := pascalFromPath(n.Path)

		// If this type name is already used, prefix it with the root collection name
		if opt.AllUsedTypeNames != nil && opt.AllUsedTypeNames[typeName] {
			typeName = opt.RootTypeName + typeName
		}

		typeNames[n.Path] = typeName

		// Mark this type name as used for future collections
		if opt.AllUsedTypeNames != nil {
			opt.AllUsedTypeNames[typeName] = true
		}
	}

	var b strings.Builder

	if opt.ObjectIDAs == "ObjectId" {
		b.WriteString("export type ObjectId = string;\n\n")
	}
	if opt.DateAs == "string" {
		b.WriteString("export type ISODateString = string;\n\n")
	}

	b.WriteString(renderObjectDef(root, totalDocs, opt, typeNames))
	b.WriteString("\n\n")

	for _, n := range objectNodes {
		def := renderObjectDef(n, totalDocs, opt, typeNames)
		if strings.TrimSpace(def) == "" {
			continue
		}
		b.WriteString(def)
		b.WriteString("\n\n")
	}

	return strings.TrimSpace(b.String()) + "\n"
}

func collectObjectNodes(n *infer.SchemaNode, out *[]*infer.SchemaNode) {
	if n.Path != "" && n.Kind == infer.NodeObject {
		*out = append(*out, n)
	}

	keys := make([]string, 0, len(n.Children))
	for k := range n.Children {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		collectObjectNodes(n.Children[k], out)
	}
}

func renderObjectDef(n *infer.SchemaNode, totalDocs int, opt TSOptions, typeNames map[string]string) string {
	name := typeNames[n.Path]
	kw := "type"
	if opt.UseInterface {
		kw = "interface"
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("export %s %s ", kw, name))
	if kw == "type" {
		b.WriteString("= ")
	}
	b.WriteString("{\n")

	childNames := make([]string, 0, len(n.Children))
	for k := range n.Children {
		childNames = append(childNames, k)
	}
	sort.Strings(childNames)
	printable := 0
	for _, childKey := range childNames {
		child := n.Children[childKey]

		if strings.HasSuffix(child.Name, "[]") {
			continue
		}

		printable++

		fieldName := child.Name

		isOptional := isFieldOptional(child, totalDocs, opt.RequiredThreshold, opt.NullPolicy)

		tsType := tsTypeForNode(child, opt, typeNames)

		if isOptional {
			b.WriteString(fmt.Sprintf("  %s?: %s;\n", fieldName, tsType))
		} else {
			b.WriteString(fmt.Sprintf("  %s: %s;\n", fieldName, tsType))
		}
	}
	if printable == 0 && n.Path != "" {
		return ""
	}

	b.WriteString("}")
	if kw == "interface" {

	} else {
		b.WriteString(";")
	}
	return b.String()
}

func isFieldOptional(n *infer.SchemaNode, totalDocs int, threshold float64, nullPolicy string) bool {
	if totalDocs <= 0 {
		return true
	}

	if float64(n.Count)/float64(totalDocs) < threshold {
		return true
	}

	if nullPolicy == "optional" {
		if _, ok := n.Types[infer.KindNull]; ok {
			return true
		}
	}
	return false
}

func tsTypeForNode(n *infer.SchemaNode, opt TSOptions, typeNames map[string]string) string {
	switch n.Kind {
	case infer.NodeObject:

		if tn, ok := typeNames[n.Path]; ok && tn != "" {
			return tn
		}
		return "Record<string, unknown>"

	case infer.NodeArray:

		elemNode, ok := n.Children[n.Name+"[]"]
		if ok && elemNode != nil && elemNode.Kind == infer.NodeObject {
			return fmt.Sprintf("%s[]", typeNames[elemNode.Path])
		}

		elemUnion := tsUnionFromKinds(n.ArrayElemTypes, opt)
		if elemUnion == "" {
			elemUnion = "unknown"
		}
		// Only wrap in parentheses if it's a union type (contains |)
		if strings.Contains(elemUnion, "|") {
			return fmt.Sprintf("(%s)[]", elemUnion)
		}
		return fmt.Sprintf("%s[]", elemUnion)

	case infer.NodePrimitive:
		u := tsUnionFromKinds(n.Types, opt)
		if u == "" {
			return "unknown"
		}

		if opt.NullPolicy == "optional" {
			u = stripNullFromUnion(u)
			if u == "" {
				return "unknown"
			}
		}
		return u

	default:

		u := tsUnionFromKinds(n.Types, opt)
		if u == "" {
			return "unknown"
		}
		if opt.NullPolicy == "optional" {
			u = stripNullFromUnion(u)
			if u == "" {
				return "unknown"
			}
		}
		return u
	}
}

func tsUnionFromKinds(m map[infer.Kind]int, opt TSOptions) string {
	parts := make([]string, 0, len(m))
	for k := range m {
		switch k {
		case infer.KindObject, infer.KindArray:

			continue
		case infer.KindString:
			parts = append(parts, "string")
		case infer.KindNumber:
			parts = append(parts, "number")
		case infer.KindBoolean:
			parts = append(parts, "boolean")
		case infer.KindNull:
			parts = append(parts, "null")
		case infer.KindDate:
			if opt.DateAs == "Date" {
				parts = append(parts, "Date")
			} else {
				parts = append(parts, "ISODateString")
			}
		case infer.KindObjectID:
			if opt.ObjectIDAs == "ObjectId" {
				parts = append(parts, "ObjectId")
			} else {
				parts = append(parts, "string")
			}
		default:
			parts = append(parts, "unknown")
		}
	}

	if len(parts) == 0 {
		return ""
	}

	sort.Strings(parts)
	parts = unique(parts)
	parts = deduplicateStringAliases(parts)
	return strings.Join(parts, " | ")
}

// deduplicateStringAliases removes redundant base string types from unions when semantic aliases exist.
// For example: "ISODateString | string" becomes "ISODateString" to preserve semantic meaning
func deduplicateStringAliases(parts []string) []string {
	hasString := false
	hasISODateString := false
	hasObjectId := false

	for _, p := range parts {
		switch p {
		case "string":
			hasString = true
		case "ISODateString":
			hasISODateString = true
		case "ObjectId":
			hasObjectId = true
		}
	}

	// If we have semantic type aliases along with base "string", prefer the semantic types
	if hasString && (hasISODateString || hasObjectId) {
		result := make([]string, 0, len(parts))
		for _, p := range parts {
			// Skip base "string" when we have semantic aliases
			if p == "string" {
				continue
			}
			result = append(result, p)
		}
		return result
	}

	return parts
}

func stripNullFromUnion(u string) string {
	parts := strings.Split(u, " | ")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if strings.TrimSpace(p) != "null" {
			out = append(out, strings.TrimSpace(p))
		}
	}
	return strings.Join(out, " | ")
}

func unique(in []string) []string {
	if len(in) == 0 {
		return in
	}
	out := []string{in[0]}
	for i := 1; i < len(in); i++ {
		if in[i] != in[i-1] {
			out = append(out, in[i])
		}
	}
	return out
}

func pascalFromPath(path string) string {
	if path == "" {
		return "Root"
	}
	segs := strings.Split(path, ".")
	var b strings.Builder
	for _, s := range segs {
		s = strings.TrimSuffix(s, "[]")
		b.WriteString(pascal(s))
	}
	return b.String()
}

func pascal(s string) string {
	if s == "" {
		return s
	}

	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})
	for i := range parts {
		if len(parts[i]) == 0 {
			continue
		}
		parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
	}
	return strings.Join(parts, "")
}
