package mysql

import (
	"database/sql"
	"fmt"
	"talk/common"
	"talk/config"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db  *sql.DB
	err error
)

func Connect() bool {
	conf := config.InitConfig(common.RootPath() + "/config/config.config")
	link := conf["Username"] + ":" + conf["Passwd"] + "@tcp(" + conf["Host"] + ")/" + conf["Database"] + "?" + conf["Charset"]
	if common.Debug {
		fmt.Println(conf)
		fmt.Println(conf["Username"])
		fmt.Println(link)
	}
	db, err = sql.Open("mysql", link)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
func Select(sqlstmt string, params ...interface{}) []map[string]interface{} {
	stmt, e := db.Prepare(sqlstmt)
	if e != nil {
		fmt.Println("请检查sql语句->", sqlstmt, e)
	}
	defer stmt.Close()
	rows, e := stmt.Query(params...)
	if e != nil {
		fmt.Println("请检查参数->", params, e)
	}
	cols, e := rows.Columns()
	if e != nil {
		panic(e)
	}
	defer rows.Close()
	count := len(cols)
	Datas := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range cols {
			var v interface{}
			val := values[i]
			if b, ok := val.([]byte); ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		Datas = append(Datas, entry)
	}
	return Datas
}
func DML(sqlstmt string, params ...interface{}) bool {
	stmt, e := db.Prepare(sqlstmt)
	if e != nil {
		fmt.Println("请检查sql语句->", sqlstmt, e)
		return false
	}
	_, err = stmt.Exec(params...)
	if err != nil {
		fmt.Println("数据操作失败->", err)
		return false
	}
	stmt.Close()
	return true
}
func DDL(sqlstmt string) bool {
	stmt, e := db.Prepare(sqlstmt)
	if e != nil {
		fmt.Println("请检查sql语句->", sqlstmt, e)
		return true
	}
	if stmt != nil {
		stmt.Exec()
	}
	stmt.Close()
	return true
}
func Close() {
	db.Close()
}
