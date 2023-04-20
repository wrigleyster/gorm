package util

func Log(v ...interface{}) {
	if v[0] != nil {
		panic(v)
		//log.Fatal(v...)
	}
}
