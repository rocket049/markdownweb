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

type MyDB struct {
	db *sql.DB
}

func dbOpen() *MyDB {
	var db *sql.DB
	var err error
	db, err = sql.Open("sqlite3", "file:./data/log.db?cache=shared&_mutex=full&_busy_timeout=9999999")
	if err != nil {
		return nil
	}
	_, err = db.Exec("create table if not exists log(path text UNIQUE,num int DEFAULT 0);")
	if err != nil {
		return nil
	}
	return &MyDB{db}
}

func (p *MyDB) Close() {
	p.db.Close()
}

func (p *MyDB) InsertPath(v string) bool {
	_, err := p.db.Exec("insert into log(path,num) values(?,1);", v)
	if err != nil {
		return false
	} else {
		return true
	}
}

func (p *MyDB) UpdatePathOnce(v string) {
	p.db.Exec("update log set num=(select num from log where path=?)+1 where path=?;", v, v)
}

func InsertOrUpdatePath(p string) {
	db := dbOpen()
	if db == nil {
		return
	}
	defer db.Close()
	ok := db.InsertPath(p)
	//log.Println(ok)
	if ok == false {
		db.UpdatePathOnce(p)
	}
	//log.Println("Insert or Update path")
}

/*
func main() {
	flag.Parse()
	p := flag.Arg(0)
	InsertOrUpdatePath(p)
}
*/
