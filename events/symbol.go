package events

// Symbol is stock or fund
type Symbol struct {
	tableName struct{} `pg:",discard_unknown_columns"`
	ID       int
	Sym      string `pg:"type:'varchar',unique"`
	Name     string `pg:"type:'varchar'"`
	Industry string `pg:"type:'varchar'"`
	Sector   string `pg:"type:'varchar'"`
}
