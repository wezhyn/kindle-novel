package mysql

import (
	"database/sql"
	"fmt"
	"sync"
)

var once sync.Once
var mysql *BookMysql

//数据库连接信息
const (
	USERNAME        = "kindle"
	PASSWORD        = "kindle"
	NETWORK         = "tcp"
	SERVER          = "www.wezhyn.com"
	PORT            = 3306
	DATABASE        = "kindle"
	MAX_LIFE_TIME   = 60
	MAX_CONNECT_NUM = 5
)

func init() {
	InitTables(NewInstance().db)
}

func InitTables(DB *sql.DB) {
	createSql := `CREATE TABLE IF NOT EXISTS book(
	id INT(8) PRIMARY KEY AUTO_INCREMENT NOT NULL,
	book_name VARCHAR(64) CHARACTER SET utf8 COLLATE utf8_unicode_ci not null ,
	web_name VARCHAR(64) CHARACTER SET utf8 COLLATE utf8_unicode_ci,
	title varchar(64) CHARACTER SET utf8 COLLATE utf8_unicode_ci not null ,
	number int not null ,
	url varchar(256) CHARACTER SET utf8 COLLATE utf8_unicode_ci not null ,
	url_hash bigint ,
	unique index book_name_number (book_name, number)
	)  ENGINE=InnoDB character set utf8 ; `

	if _, err := DB.Exec(createSql); err != nil {
		fmt.Println("create table failed:", err)
		return
	}
	fmt.Println("create table success")
}
