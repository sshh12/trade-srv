package events

type Symbol struct {
	ID       int
	Sym      string `pg:"type:'varchar',unique"`
	Name     string `pg:"type:'varchar'"`
	Industry string `pg:"type:'varchar'"`
	Sector   string `pg:"type:'varchar'"`
}
