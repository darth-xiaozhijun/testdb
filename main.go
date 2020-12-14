package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type Dbconn struct {
	Dsn string
	Db  *sql.DB
}

type UserTable struct {
	Uid        int
	Username   string
	Departname string
	Created    string
}

func main() {
	var err error
	dbConn := Dbconn{
		Dsn: "root:123456@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4",
	}
	dbConn.Db, err = sql.Open("mysql", dbConn.Dsn)
	if err != nil {
		panic(err)
		return
	}
	defer dbConn.Db.Close()
	//execData(&dbConn, "INSERT user_info(username,departname,created) VALUES('xzj','huawei','2020-12-14')")
	//preExecData(&dbConn, "INSERT user_info(username,departname,created) VALUES(?,?,?)", "tj", "youzan", "2020-12-14")
	userTable := dbConn.QueryRowData("select * from user_info where uid=1")
	fmt.Println(userTable)
	userMap := dbConn.QueryData("select * from user_info where uid>0")
	fmt.Println(userMap)
	userMap = dbConn.PreQueryData("select * from user_info where uid=?", 1)
	fmt.Println(userMap)
}

func execData(dbConn *Dbconn, sqlString string) {
	count, id, err := dbConn.ExecData(sqlString)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("受影响的行数：", count)
		fmt.Println("新添加数据的ID：", id)
	}
}

func preExecData(dbConn *Dbconn, sqlString string, args ...interface{}) {
	count, id, err := dbConn.PreExecData(sqlString, args...)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("受影响的行数：", count)
		fmt.Println("新添加数据的ID：", id)
	}
}

func (dbConn *Dbconn) ExecData(sqlString string) (count, id int64, err error) {
	result, err := dbConn.Db.Exec(sqlString)
	if err != nil {
		panic(err)
		return
	}
	if id, err = result.LastInsertId(); err != nil {
		panic(err)
		return
	}
	if count, err = result.RowsAffected(); err != nil {
		panic(err)
		return
	}
	return count, id, nil
}

func (dbConn *Dbconn) PreExecData(sqlString string, args ...interface{}) (count, id int64, err error) {
	stmt, err := dbConn.Db.Prepare(sqlString)
	defer stmt.Close()
	if err != nil {
		panic(err)
		return
	}
	result, err := stmt.Exec(args...)
	if err != nil {
		panic(err)
		return
	}
	if id, err = result.LastInsertId(); err != nil {
		panic(err)
		return
	}
	if count, err = result.RowsAffected(); err != nil {
		panic(err)
		return
	}
	return count, id, nil
}

func (dbConn *Dbconn) QueryRowData(sqlString string) (data UserTable) {
	user := new(UserTable)
	err := dbConn.Db.QueryRow(sqlString).Scan(&user.Uid, &user.Username, &user.Departname, &user.Created)
	if err != nil {
		panic(err)
		return
	}
	return *user
}

func (dbConn *Dbconn) QueryData(sqlString string) (resultSet map[int]UserTable) {
	rows, err := dbConn.Db.Query(sqlString)
	defer rows.Close()
	if err != nil {
		panic(err)
		return
	}
	resultSet = make(map[int]UserTable)
	user := new(UserTable)
	for rows.Next() {
		err := rows.Scan(&user.Uid, &user.Username, &user.Departname, &user.Created)
		if err != nil {
			panic(err)
			return
		}
		resultSet[user.Uid] = *user
	}
	return resultSet
}

func (dbConn *Dbconn) PreQueryData(sqlString string, args ...interface{}) (resultSet map[int]UserTable) {
	stmt, err := dbConn.Db.Prepare(sqlString)
	defer stmt.Close()
	if err != nil {
		panic(err)
		return
	}
	rows, err := stmt.Query(args...)
	defer rows.Close()
	if err != nil {
		panic(err)
		return
	}
	resultSet = make(map[int]UserTable)
	user := new(UserTable)
	for rows.Next() {
		err := rows.Scan(&user.Uid, &user.Username, &user.Departname, &user.Created)
		if err != nil {
			panic(err)
			return
		}
		resultSet[user.Uid] = *user
	}
	return resultSet
}

func (dbConn *Dbconn) PreQueryData2(sqlString string, args ...interface{}) {
	stmt, err := dbConn.Db.Prepare(sqlString)
	defer stmt.Close()
	if err != nil {
		panic(err)
		return
	}
	rows, err := stmt.Query(args...)
	defer rows.Close()
	if err != nil {
		panic(err)
		return
	}
	user := new(UserTable)
	for rows.Next() {
		err := rows.Scan(&user.Uid, &user.Username, &user.Departname, &user.Created)
		if err != nil {
			panic(err)
			return
		}
		fmt.Println(*user)
	}
}
