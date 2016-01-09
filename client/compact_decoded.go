package client

// CompactData is the common compact decoded structure
type CompactData struct {
	ID, Date, Version string
	Columns           []string
	Rows              [][]string
}

// Indexer provices cached lookup for CompactData
type Indexer func(col string, row int) string

// Indexer create the cache
func (m *CompactData) Indexer() Indexer {
	index := make(map[string]int)
	for i, c := range m.Columns {
		index[c] = i
	}
	return func(col string, row int) string {
		return m.Rows[row][index[col]]
	}
}
