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
	Where(predicate string, params ...string) Stmt
	Select(cols ...string) *sql.Rows
	Update(stmt string, params ...string) sql.Result
	Replace(values ...string) sql.Result
	Delete() sql.Result
}

func Qs(n int) []string {
	ss := make([]string, n)
	for i, _ := range ss {
		ss[i] = "?"
	}
	return ss
}
func GetParams(params ...string) []interface{} {
	all := make([]interface{}, len(params))
	for i, v := range params {
		all[i] = v
	}
	return all
}
