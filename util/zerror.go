package util

func CheckPanic(err error){
	if err!=nil {
		panic(err)
	}
}

func CheckErr(err error) bool{
	if err!=nil {
		return false
	}
	return true
}
