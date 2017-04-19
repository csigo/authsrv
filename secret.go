package authsrv

import "crypto/rand"

const readable = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenSalt(length int) string {
	var bytes = make([]byte, length)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = readable[v%byte(len(readable))]
	}
	return string(bytes)
}
