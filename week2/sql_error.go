package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)


//dao遇到 sql.ErrNoRows 不应改wrap，应该直接处理，不返回数据即可

/* 答案
dao:

 return errors.Wrapf(errors.NotFound, fmt.Sprintf("sql: %s error: %v", sql, err))


biz:

if errors.Is(err, errors.NotFound} {

}

if errors.Reason(err, xxxx) == xxxx {

}
 */

func main() {
	db, err := sql.Open("mysql", "root:123456@tcp(192.168.33.111:3306)/wxziroom");
	defer db.Close()

	if err != nil {
		panic(err)
	}
	id := 1
	nickname, err := queryDataWithId(db, id)
	fmt.Println(nickname, err)
}

func queryDataWithId(db *sql.DB, id int) (string, error) {
	var nickname string
	err := db.QueryRow("SELECT nickname FROM hd_join WHERE id=?", id).Scan(&nickname)
	switch {
	case err == sql.ErrNoRows:
		return "", nil
	case err != nil:
		return "", err
	default:
		return nickname, nil
	}
}

