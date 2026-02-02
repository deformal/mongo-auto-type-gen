package infer

type FieldStats struct {
	Path        string
	Count       int
	Types       map[Kind]int
	ArrayTypes  map[Kind]int
	ObjectStats map[string]*FieldStats
}

func NewFieldStats(path string) *FieldStats {
	return &FieldStats{
		Path:        path,
		Types:       map[Kind]int{},
		ArrayTypes:  map[Kind]int{},
		ObjectStats: map[string]*FieldStats{},
	}
}
