package mariadb

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/wrigleyster/gorm"
	"github.com/wrigleyster/gorm/util"
	"strings"
)

type DS struct {
	dbname, host   string
	port           int
	user, password string
}
type Stmt struct {
	db        DS
	cols      []string
	table     string
	predicate string
	params    []any
}
type TblDef struct {
	db         DS
	columnDefs string
	key        string
}

func New(dbname, host string, port int) gorm.DataSource {
	return DS{dbname: dbname, host: host, port: port}
}

func (_ds DS) From(table string) gorm.Stmt {
	return Stmt{db: _ds, table: table}
}
func (_ds DS) Orm(dbConsumer gorm.DbConsumer) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", _ds.user, _ds.password, _ds.host, _ds.port, _ds.dbname)
	db, err := sql.Open("mysql", dsn)
	util.Log(err)
	defer db.Close()
	dbConsumer(db)
}
func (_ds DS) With(columns string) gorm.TblDef {
	return TblDef{db: _ds, columnDefs: columns}
}
func (_ds DS) DropTable(table string) {
	_ds.Orm(func(db *sql.DB) {
		_, err := db.Exec("drop table " + table)
		util.Log(err)
	})
}
func (_tblDef TblDef) Key(primaryKey string) gorm.TblDef {
	_tblDef.key = primaryKey
	return _tblDef
}
func (_tblDef TblDef) CreateTable(table string) {
	_tblDef.db.Orm(func(db *sql.DB) {
		var key string
		if _tblDef.key != "" {
			key = ", primary key (" + _tblDef.key + ")"
		}
		_, err := db.Exec("create table if not exists `" + table + "` (" + _tblDef.columnDefs + key + ")")
		util.Log(err)
	})
}
func (_stmt Stmt) Where(predicate string, params ...any) gorm.Stmt {
	_stmt.predicate = predicate
	_stmt.params = params
	return _stmt
}
func (_stmt Stmt) Select(cols ...string) *sql.Rows {
	var rows *sql.Rows
	if len(cols) == 0 {
		cols = append(cols, "*")
	}
	_stmt.db.Orm(func(db *sql.DB) {
		var where string
		if _stmt.predicate != "" {
			where = "where"
		}
		stmt, err := db.Prepare(
			strings.Join([]string{
				"select", strings.Join(cols, ","),
				"from", _stmt.table, where, _stmt.predicate,
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
		var where string
		if _stmt.predicate != "" {
			where = "where"
		}
		stmt, err := db.Prepare(
			strings.Join([]string{
				"update", _stmt.table, "set", stmt, where, _stmt.predicate,
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
				"delete from", _stmt.table, "where", _stmt.predicate,
			}, " "))
		util.Log(err)
		result, err = stmt.Exec(_stmt.params...)
		util.Log(err)
	})
	return result
}
