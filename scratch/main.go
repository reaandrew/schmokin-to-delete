package main

import (
	"fmt"
	"regexp"
)

func main() {
	data := "Average Transaction Rate (requests/sec)......: 1469.85"
	matched, err := regexp.Match(`Average Transaction Rate \(requests/sec\)[^\s]+\s[^0][\d\.]+`, []byte(data))
	fmt.Println(matched, err)
}
