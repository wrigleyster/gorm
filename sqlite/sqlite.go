package sqlite

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wrigleyster/gorm"
	"github.com/wrigleyster/gorm/util"
)

type DS struct {
	dataSource string
}
type Stmt struct {
	db        DS
	cols      []string
	table     string
	where     string
	params    []any
	sortOrder string
	join      string
}

type TblDef struct {
	db      DS
	columns string
	key     string
}

func New(dataSource string) gorm.DataSource {
	return DS{dataSource: dataSource}
}

func (_ds DS) From(table string) gorm.Stmt {
	return Stmt{db: _ds, table: table}
}
func (_ds DS) Orm(dbConsumer gorm.DbConsumer) {
	db, err := sql.Open("sqlite3", _ds.dataSource)
	util.Log(err)
	defer db.Close()
	dbConsumer(db)
}
func (_ds DS) With(columns string) gorm.TblDef {
	return TblDef{db: _ds, columns: columns}
}
func (_ds DS) DropTable(table string) {
	_ds.Orm(func(db *sql.DB) {
		db.Exec("drop table " + table)
	})
}
func (_tblDef TblDef) CreateTable(table string) {
	_tblDef.db.Orm(func(db *sql.DB) {
		var key string
		if _tblDef.key != "" {
			key = ", primary key (" + _tblDef.key + ")"
		}
		_, err := db.Exec("create table if not exists `" + table + "` (" + _tblDef.columns + key + ")")
		util.Log(err)
	})
}
func (_tblDef TblDef) Key(primaryKey string) gorm.TblDef {
	_tblDef.key = primaryKey
	return _tblDef
}
func (_stmt Stmt) Where(predicate string, params ...any) gorm.Stmt {
	_stmt.where = "WHERE " + predicate
	_stmt.params = params
	return _stmt
}
func (_stmt Stmt) OrderAscendingBy(col string) gorm.Stmt {
	_stmt.sortOrder = "ORDER BY "+col+" ASC"
	return _stmt
}
func (_stmt Stmt) OrderDescendingBy(col string) gorm.Stmt {
	_stmt.sortOrder = "ORDER BY "+col+" DESC"
	return _stmt
}
func (_stmt Stmt) InnerJoin(table, predicate string) gorm.Stmt {
	_stmt.join = fmt.Sprintf("INNER JOIN %s ON %s", table, predicate)
	return _stmt
}
func (_stmt Stmt) LeftJoin(table, predicate string) gorm.Stmt {
	_stmt.join = fmt.Sprintf("LEFT OUTER JOIN %s ON %s", table, predicate)
	return _stmt
}
func (_stmt Stmt) RightJoin(table, predicate string) gorm.Stmt {
	_stmt.join = fmt.Sprintf("RIGHT OUTER JOIN %s ON %s", table, predicate)
	return _stmt
}
func (_stmt Stmt) Select(cols ...string) *sql.Rows {
	var rows *sql.Rows
	if len(cols) == 0 {
		cols = append(cols, "*")
	}
	_stmt.db.Orm(func(db *sql.DB) {
		stmt, err := db.Prepare(
			strings.Join([]string{
				"select", strings.Join(cols, ","),
				"from", _stmt.table, _stmt.join, _stmt.where,
				_stmt.sortOrder,
			}, " "))
		util.Log(err)
		rows, err = stmt.Query(_stmt.params...)
		util.Log(err)
	})
	return rows
}
func (_stmt Stmt) Update(stmt string, params ...any) sql.Result {
	var result sql.Result
	_stmt.db.Orm(func(db *sql.DB) {
		stmt, err := db.Prepare(
			strings.Join([]string{
				"update", _stmt.table, "set", stmt, _stmt.where,
			}, " "))
		util.Log(err)
		result, err = stmt.Exec(append(params, _stmt.params...)...)
	})
	return result
}
func (_stmt Stmt) Replace(values ...any) sql.Result {
	var result sql.Result
	_stmt.db.Orm(func(db *sql.DB) {
		stmt, err := db.Prepare(
			strings.Join([]string{
				"replace into", _stmt.table, "values", "(", strings.Join(gorm.Qs(len(values)), ","), ")",
			}, " "))
		util.Log(err)
		result, err = stmt.Exec(values...)
		util.Log(err)
	})
	return result
}
func (_stmt Stmt) Delete() sql.Result {
	var result sql.Result
	_stmt.db.Orm(func(db *sql.DB) {
		stmt, err := db.Prepare(
			strings.Join([]string{
				"delete from", _stmt.table, _stmt.where,
			}, " "))
		util.Log(err)
		result, err = stmt.Exec(_stmt.params...)
		util.Log(err)
	})
	return result
}

// From("github_request").Where("a=? and b=?", a, b).Select("col1", "col2")
// From("review").Where("hash=?", hash).Delete()
// Into("review").Where("hash=?", hash).insert()
//
