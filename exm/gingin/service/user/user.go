package user

import (
	"database/sql"
)

//User is a struct 对应解析JSON
type User struct {
	UserName string `json: "username"`
	Password string `json: "password"`
}

// AddOne is 添加一个用户
func AddOne(username string, password string) {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test")
	checkErr(err)

	defer db.Close()

	opt, err := db.Prepare("INSERT users SET username=?, password=?")
	checkErr(err)

	_, err = opt.Exec(username, password)
	checkErr(err)
}

// DelOne is 删除一个用户
func DelOne(username string) {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test")
	checkErr(err)

	defer db.Close()

	opt, err := db.Prepare("DELETE FROM users WHERE username=?")
	checkErr(err)

	_, err = opt.Exec(username)
	checkErr(err)

}

// UpdateOne is 更新一个用户
func UpdateOne(info *User) {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test")

	checkErr(err)

	defer db.Close()

	opt, err := db.Prepare("UPDATE FROM users SET password=? WHERE username=?")
	checkErr(err)

	_, err = opt.Exec(info)
	checkErr(err)
}

// FindAll is 查询整张User表
func FindAll() []User {
	var users []User

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test")

	checkErr(err)

	defer db.Close()

	results, err := db.Query("SELECT * FROM users")

	checkErr(err)

	for results.Next() {
		var user User

		err = results.Scan(&user.UserName, &user.Password)

		checkErr(err)

		users = append(users, user)
	}

	return users
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
