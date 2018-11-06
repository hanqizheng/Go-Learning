package main

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"
	bson "gopkg.in/mgo.v2/bson"
)

// User is the collection of users
type User struct {
	UserName string
	Password string
}

func main() {
	session, err := mgo.Dial("localhost")

	if err != nil {
		fmt.Println(err.Error())
	}

	db := session.DB("hqz")
	collection := db.C("users")

	// 插入
	err = collection.Insert(&User{"xiaoming", "root"}, &User{"lina", "long"})
	if err != nil {
		fmt.Println(err.Error())
	}

	// 删除
	err = collection.Remove(bson.M{"username": "lina"})

	// 更新
	err = collection.Update(bson.M{"_id": bson.ObjectIdHex("5be105b200e82a2b3f909e7b")}, bson.M{"$set": bson.M{"password": "123455677"}})

	// 单个查询
	result := User{}
	err = collection.Find(bson.M{"username": "hanqizheng"}).One(&result)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(result)

	// 多项查询
	result := User{}
	iter := collection.Find(nil).Iter()

	for iter.Next(&result) {
		fmt.Println(result)
	}

}
