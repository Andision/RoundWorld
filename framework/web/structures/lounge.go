package structures

import "fmt"

// Lounge is an interface that defines the methods for managing game tablesMap.
type Lounge interface {
	GetTableById(name string) Table
	AddTable(table Table) error
}

type LoungeImpl struct {
	tablesMap map[string]Table
}

func NewLounge() Lounge {
	return &LoungeImpl{
		tablesMap: make(map[string]Table),
	}
}

func (l *LoungeImpl) GetTableById(name string) Table {
	return l.tablesMap[name]
}

func (l *LoungeImpl) AddTable(table Table) error {
	// Check if the table already exists
	if _, exists := l.tablesMap[table.GetTableId()]; exists {
		return fmt.Errorf("table %s already exists", table.GetTableId())
	}

	l.tablesMap[table.GetTableId()] = table
	return nil
}
