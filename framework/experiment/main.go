package main

import (
	"expvar"
	"fmt"
)

func main() {
	fmt.Println("vim-go")

	point := expvar.NewString("root_path").Init()
	point.Set("/")
	fmt.Println("root_path is:", root_path)
	fmt.Println("point is:", point)

}
