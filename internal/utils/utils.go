package utils

func MakeUrlPathValid(str *string) {
	val := *str
	if val[0] != '/' {
		*str = "/" + *str
	}
	for val[len(val)-1] == '/' {
		*str = (*str)[:len(*str)-1]
	}
}

