package util

import "log"

func Log(v ...any) {
	if v[0] != nil {
		log.Println(v...)
		panic(v[len(v)-1])
	}
}
