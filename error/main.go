package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

var db *sql.DB

func initDB() (err error) {
	dsn := "root:88wang@tcp(39.96.93.59)/todo?charset=utf8mb4&parseTime=True"
	// 不会校验账号密码是否正确
	// 注意！！！这里不要使用:=，我们是给全局变量赋值，然后在main函数中使用全局变量db
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return errors.Wrap(err, "[dao] sql.Open error")
	}

	// 尝试鱼数据库建立连接（校验dsn是否正确）
	err = db.Ping()
	if err != nil {
		return errors.Wrap(err, "[dao] db.Ping error")
	}

	return nil
}

func deinitDB() {
	db.Close()
}

type task struct {
	id        int
	title     string
	completed int
}

func queryRowFromTable(t *task, id int) error {
	sqlStr := "select id, title, completed from todos where id=?"
	// 非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放
	err := db.QueryRow(sqlStr, id).Scan(&t.id, &t.title, &t.completed)
	if err != nil {
		return errors.Wrap(err, "[dao] db.QueryRow error")
	}
	return nil
}

func main() {
	err := initDB()
	if err != nil {
		fmt.Printf("init db failed, err: %+v\n", err)
		return
	}
	defer deinitDB()

	db.SetMaxOpenConns(1000)
	db.SetMaxIdleConns(100)

	var tk task
	tk.id = 1
	err = queryRowFromTable(&tk, tk.id)
	if err != nil {
		fmt.Printf("[service] query DB err: %+v\n", err)
		return
	}
	fmt.Printf("table %d info [%v, %v, %v]\n", tk.id, tk.id, tk.title, tk.completed)

}
