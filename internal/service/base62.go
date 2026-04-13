package service

const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func encodeBase62(num int64) string {
	if num == 0 {
		return "0"
	}

	var encoded []byte

	for num > 0 {
		remainder := num % 62
		encoded = append([]byte{base62Chars[remainder]}, encoded...)
		num = num / 62
	}

	return string(encoded)
}
