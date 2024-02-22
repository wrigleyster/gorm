package gorm

import "database/sql"

type DbConsumer func(db *sql.DB)
type DataSource interface {
	From(table string) Stmt
	Orm(dbConsumer DbConsumer)
	DropTable(table string)
	With(columns string) TblDef
}

type TblDef interface {
	Key(primaryKey string) TblDef
	CreateTable(table string)
}
type Stmt interface {
	OrderAscendingBy(col string) Stmt
	OrderDescendingBy(col string) Stmt
	InnerJoin(table string, predicate string) Stmt
	LeftJoin(table string, predicate string) Stmt
	RightJoin(table string, predicate string) Stmt
	Where(predicate string, params ...any) Stmt
	Select(cols ...string) *sql.Rows
	Update(stmt string, params ...any) sql.Result
	Replace(values ...any) sql.Result
	Delete() sql.Result
}


func Qs(n int) []string {
	ss := make([]string, n)
	for i := range ss {
		ss[i] = "?"
	}
	return ss
}
