package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func RandString() string {
	char := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	arr := strings.Split(char, "")
	length := len(arr)
	ran := rand.New(rand.NewSource(time.Now().Unix()))
	randStr := ""
	for i := 0; i < 15; i++ {
		randStr = randStr + arr[ran.Intn(length)]
	}

	return randStr
}

// GetDays
// date2 大于 date1
// -1 失败
func GetDays(format, date1, date2 string) int {
	d1, err := time.ParseInLocation(format, date1, time.Local)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	d2, err := time.ParseInLocation(format, date2, time.Local)
	if err != nil {
		return -1
	}
	return int(d2.Sub(d1).Hours() / 24)
}

func main() {
	fmt.Println(GetDays("2006-01-02 15:04:05", "2022-02-23 11:00:34", time.Now().Format("2006-01-02 15:04:05")))
}