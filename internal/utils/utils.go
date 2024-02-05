package utils

func MakeUrlPathValid(str *string) {
	if len(*str) == 0 {
		return
	}

	if (*str)[0] != '/' {
		*str = "/" + *str
	}

	for (*str)[len(*str)-1] == '/' {
		*str = (*str)[:len(*str)-1]
	}
}
