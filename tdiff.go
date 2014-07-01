// +build ignore

package main

import (
	"fmt"
	"os"
	"time"
)

const (
	ssDate = "2006-01-02"
	ssTime = "15:04:05"
)

func main() {
	form := fmt.Sprintf("%s-%s", ssDate, ssTime)
	ts1, ts2 := os.Args[1], os.Args[2]

	t1, err := time.Parse(form, ts1)
	if err != nil {
		panic(err)
	}

	t2, err := time.Parse(form, ts2)
	if err != nil {
		panic(err)
	}

	fmt.Println(t1)

	d := t1.Sub(t2)
	if t1.Before(t2) {
		d = t2.Sub(t1)
	}

	fmt.Println(d)

	fmt.Println(t1.Add(d))

}
