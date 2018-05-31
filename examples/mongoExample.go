package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"encoding/json"
)

type Person struct {
	Name  string
	Phone string
}

func main() {
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("test").C("people")
	err = c.Insert(&Person{"zuston", "+233639"},
		&Person{"ann", "+2323"})
	if err != nil {
		panic(err)
	}

	result := Person{}
	err = c.Find(bson.M{"name": "zuston"}).One(&result)
	if err != nil {
		panic(err)
	}

	fmt.Println("Phone:", result.Phone)

	var res []bson.M
	err = c.Find(bson.M{"name": "Ale"}).All(&res)
	if err != nil {
		panic(err)
	}
	fmt.Println(res[0]["name"])

	var ress bson.M
	iter := c.Find(bson.M{"name": "Ale"}).Iter()
	fmt.Println("---------GGGGGG---------")
	for iter.Next(&ress){
		l, _:= json.Marshal(ress)
		fmt.Println(string(l))
	}
}
