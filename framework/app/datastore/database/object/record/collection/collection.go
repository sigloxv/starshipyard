package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("vim-go")
}

type Collection struct {
	Name          string
	Records       []*Record
	LastUpdatedAt time.Time
	History       chan *Record
}
