package main

import (
		"fmt"
		"database/sql"
		_ "github.com/go-sql-driver/mysql"
)

type User struct {
		Name string `json:"Name"`
}

func main()  {
		fmt.Println("Go MySQL")

		db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test")

		if err != nil {
				panic(err.Error())
		}

		defer db.Close()

		// add 
		insert, err := db.Query("INSERT INTO users VALUES('huanan')")
		if err != nil {
				panic(err.Error())
		}
		defer insert.Close()
		fmt.Println("Successfully inserted into users")

		// update
		update, err := db.Query("UPDATE users SET name = 'hqz' WHERE name = 'hanqizheng'")
		if err != nil {
				panic(err.Error())
		}
		defer update.Close()

		// select
		results, err := db.Query("SELECT name FROM users")

		if err != nil {
				panic(err.Error())
		}

		for results.Next() {
				var user User

				err = results.Scan(&user.Name)

				if err != nil {
						panic(err.Error())
				}

				fmt.Println(user.Name)
		}

		// delete
		delete, err := db.Query("DELETE FROM users WHERE name = 'lisi'")
		if err != nil {
				panic(err.Error())
		}
		defer delete.Close()
}