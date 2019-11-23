package main

import (
	"fmt"
	"time"
)

func main() {

	tStr := "2019-02-27T09:32:09Z"

	t, err := time.Parse(time.RFC3339, tStr)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		fmt.Errorf("err: %v", err.Error())
		return
	}

	fmt.Printf("%v\n", t)
}
