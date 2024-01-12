package main

import (
	"fmt"
	"sync"
)

const (
	dbDriver = "pgx"
	dbSource = "postgresql://root:mysecret@localhost:5433/simple_bank?sslmode=disable"
)

type Person struct {
	ID      uint   `gorm:"primaryKey"`
	Name    string `gorm:"column:name"`
	Age     int    `gorm:"column:age"`
	Address string `gorm:"column:adrress"`
}

var wg sync.WaitGroup

func main() {
	//var wg sync.WaitGroup
	go sayHello()

	for i := 0; i < 5; i++ {
		fmt.Println("World")

	}
	//time.Sleep(1 * time.Millisecond)
	wg.Wait()
}

func sayHello() {

	wg.Add(1)
	for i := 0; i < 5; i++ {
		fmt.Println("Hello")
		//time.Sleep(100 * time.Millisecond)
	}
	wg.Done()
}
