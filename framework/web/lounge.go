package web

type Lounge struct {
	tables map[string]*Table
}

func NewLounge() *Lounge {
	return &Lounge{
		tables: make(map[string]*Table),
	}
}
func (l *Lounge) GetTable(name string) *Table {
	return l.tables[name]
}

func (l *Lounge) AddTable(tableId string, table *Table) {
	l.tables[tableId] = table
}
