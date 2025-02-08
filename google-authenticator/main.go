package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func prefixRandom(otp string) string {
	if len(otp) == 6 {
		return otp
	}

	for i:=(6-len(otp)); i>0; i-- {
	randomNumber := strconv.Itoa(rand.Intn(10))
	 otp = randomNumber + otp
}
     return otp
}
	
func getHOTPToken(secret string, interval int64) string {
	key, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	check(err)
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, uint64(interval))

	hash := hmac.New(sha1.New, key)
	hash.Write(bs)
	h := hash.Sum(nil)

	o := h[19] & 15

	var header uint32
	r := bytes.NewReader(h[o : o+4])
	err = binary.Read(r, binary.BigEndian, &header)
	check(err)

	h12 := (int32(header) & 0x7fffffff) % 1000000
	otp := strconv.Itoa(int(h12))
	return prefixRandom(otp)

}

func getTOTPToken(secret string) string {
	return getHOTPToken(secret, time.Now().Unix()/10)
}

func main() {
	for {
		data,err := ioutil.ReadFile("secret.txt")
		check(err)
		otp := getTOTPToken(string(data))
		fmt.Println("Your OTP is:", otp)
		time.Sleep(10 * time.Second)
	}
}