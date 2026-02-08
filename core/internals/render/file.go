package render

import (
	"strings"

	"github.com/deformal/mongo-auto-type-gen/internals/infer"
)

type FileComposer struct {
	opt      TSOptions
	wroteHdr bool
	blocks   []string
}

func NewFileComposer(opt TSOptions) *FileComposer {
	return &FileComposer{opt: opt}
}

func (f *FileComposer) AddCollection(tree *infer.SchemaNode, totalDocs int, rootTypeName string) {
	opt := f.opt
	opt.RootTypeName = rootTypeName
	block := RenderTypeScript(tree, totalDocs, opt)
	block = strings.TrimSpace(block)
	if block != "" {
		f.blocks = append(f.blocks, block)
	}
}

func (f *FileComposer) String() string {
	var b strings.Builder

	if f.opt.ObjectIDAs == "ObjectId" {
		b.WriteString("export type ObjectId = string;\n\n")
	}
	if f.opt.DateAs == "string" {
		b.WriteString("export type ISODateString = string;\n\n")
	}

	for i, blk := range f.blocks {
		if i > 0 {
			b.WriteString("\n\n")
		}
		b.WriteString(blk)
	}

	out := strings.TrimSpace(b.String())
	if out == "" {
		return ""
	}
	return out + "\n"
}
