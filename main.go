package main

import (
	"github.com/go-xorm/xorm"
	"github.com/k0kubun/pp"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type User struct {
	Id        int       `xorm:"INTEGER NOT NULL"`
	Name      string    `xorm:"TEXT NOT NULL"`
	Age       int       `xorm:"INTEGER NOT NULL"`
	CreatedAt time.Time `xorm:"NUMERIC"`
}

func main() {

	db := func() *xorm.Engine {
		engine, err := xorm.NewEngine("sqlite3", "./db.sqlite")
		engine.ShowSQL = true
		if err != nil {
			panic(err)
		}
		return engine
	}()

	pp.Println(func(engine *xorm.Engine) *User {
		obj := User{}
		_, err := engine.Id(1).Get(&obj)
		if err != nil {
			panic(err)
		}
		return &obj
	}(db))

	pp.Println(func(engine *xorm.Engine) []User {
		obj := []User{}
		err := engine.Cols("id").Find(&obj)
		if err != nil {
			panic(err)
		}
		return obj
	}(db))

	pp.Println(func(engine *xorm.Engine) map[int]int {
		obj := []User{}
		err := engine.Cols("id").Find(&obj)
		if err != nil {
			panic(err)
		}
		dst := make(map[int]int)
		for _, o := range obj {
			id := o.Id
			dst[id] = id
		}
		return dst
	}(db))

	pp.Println(func(engine *xorm.Engine) []int {
		// Numbers of Concurrency
		concurrency := 0
		done := make(chan bool)

		// Get All IDs
		allIDs := make(map[int]int)
		concurrency++
		go func(engine *xorm.Engine, dst *map[int]int, d chan bool) {
			obj := []User{}
			err := engine.Cols("id").Find(&obj)
			if err != nil {
				panic(err)
			}
			for _, o := range obj {
				id := o.Id
				(*dst)[id] = id
			}
			done <- true
		}(db, &allIDs, done)

		// Get IDs to ignore
		concurrency++
		ignoringIDs := make(map[int]int)
		go func(engine *xorm.Engine, dst *map[int]int, d chan bool) {
			obj := []User{}
			err := engine.Cols("id").Where("id % 2 = 0").Find(&obj)
			if err != nil {
				panic(err)
			}
			for _, o := range obj {
				id := o.Id
				(*dst)[id] = id
			}
			done <- true
		}(db, &ignoringIDs, done)

		// Waiting
		for i := 0; i < concurrency; i++ {
			<-done
		}

		// Delete
		for _, id := range ignoringIDs {
			delete(allIDs, id)
		}

		// Be Slice from Map
		dst := func(ids map[int]int) []int {
			dst := make([]int, len(ids))
			i := 0
			for _, id := range ids {
				dst[i] = id
				i++
			}
			return dst
		}(allIDs)
		return dst
	}(db))

}
