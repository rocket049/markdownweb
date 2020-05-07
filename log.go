package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

/*
create table if not exists log(path text UNIQUE,num int DEFAULT 0);
insert into log(path,num) values('path1',1);
update log set num=(select num from log where path='path1')+1 where path='path1';
*/

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "file:./data/log.db?_mutex=full&_busy_timeout=9999999")
	if err != nil {
		panic(err)
	}
}

func DbClose() {
	db.Close()
}

func insertPath(p string) bool {
	_, err := db.Exec("insert into log(path,num) values(?,1);", p)
	if err != nil {
		return false
	} else {
		return true
	}
}

func updatePathOnce(p string) {
	db.Exec("update log set num=(select num from log where path=?)+1 where path=?;", p, p)
}

func InsertOrUpdatePath(p string) {
	ok := insertPath(p)
	if ok == false {
		updatePathOnce(p)
	}
}

/*
func main() {
	flag.Parse()
	p := flag.Arg(0)
	ok := insertPath(p)
	fmt.Println("insert:", ok)
	if ok == false {
		updatePathOnce(p)
		fmt.Println("update once")
	}
}
*/
