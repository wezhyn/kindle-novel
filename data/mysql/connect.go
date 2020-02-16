package mysql

import (
	"database/sql"
	"fmt"
	"kindle/data"
	log2 "log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Item struct {
	id int
	data.UpdateItem
}

type BookMysql struct {
	db          *sql.DB
	errorLogger *log2.Logger
}

func (m BookMysql) Last(bookName string) (int, error) {
	lastSql := ` select number from book where book_name=? order by number desc limit 1;`
	row := m.db.QueryRow(lastSql, bookName)
	var number int
	if err := row.Scan(&number); err != nil {
		return -1, err
	} else {
		return number, nil
	}
}

func (m BookMysql) Save(data data.UpdateItem) error {

	result, err := m.db.Exec("insert INTO book(book_name,web_name,title,"+
		"number,url,url_hash) values(?,?,?,?,?,?)", data.BookName, data.WebName, data.Title,
		data.Number, data.Url, data.UrlHash)
	if err != nil {
		m.errorLogger.Printf("Insert data failed,err: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected() //通过RowsAffected获取受影响的行数
	if err != nil {
		m.errorLogger.Printf("Get RowsAffected failed,err:%v", err)
		return err
	}
	if rowsAffected <= 0 {
		return fmt.Errorf("错误的插入数据")
	}
	return nil
}
func NewInstance() *BookMysql {
	once.Do(func() {
		url := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
		db, err := sql.Open("mysql", url)
		if err != nil {
			panic(err)
		}
		mysql = new(BookMysql)
		mysql.errorLogger = log2.New(os.Stdout, "ERROR: ", log2.Ldate|log2.Ltime|log2.Lshortfile)
		db.SetConnMaxLifetime(MAX_LIFE_TIME * time.Second)
		db.SetMaxOpenConns(MAX_CONNECT_NUM)
		mysql.db = db
	})
	return mysql
}
