package main

import (
	"flag"
	"fmt"
)

func main() {
	var des = flag.String("des", "Welcome", "Destination")
	flag.Parse()
	fmt.Println(*des)
}
