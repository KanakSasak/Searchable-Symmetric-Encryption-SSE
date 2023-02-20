package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"
	"strings"
)

func main() {
	plaintext := "Alice and Bob want to send message securely"
	fromIndexWord := strings.Index(plaintext, "message")
	toIndexWord := fromIndexWord + 7
	log.Println("Plain Text : ", plaintext)
	log.Println("Index Position : " + strconv.Itoa(fromIndexWord) + " to " + strconv.Itoa(toIndexWord) + "")
	log.Println("Index Word : ", plaintext[fromIndexWord:toIndexWord])
	log.Println("Index Word Length : ", len(plaintext[fromIndexWord:toIndexWord]))
	fmt.Println("-------------------------***-------------------------------")
	plainbin := binary(plaintext)
	log.Println("Binary :", plainbin)
	log.Println("Binary Length :", len(plainbin))
	log.Println("Plain Length :", len(plaintext))
	log.Println("Plain ASCII :", []byte(plaintext))
	fmt.Println("-------------------------***-------------------------------")

	keyplain, keybin, err := GenerateRandomString(len(plaintext))
	if err != nil {
		panic(err)
	}

	log.Println("Key in plain : ", keyplain)
	log.Println("Key in binary : ", keybin)
	log.Println("Key in binary length : ", len(keybin))
	log.Println("Key in plain length : ", len(keyplain))
	log.Println("key in ASCII : ", []byte(keyplain))
	log.Println("Index key Word : ", keyplain[fromIndexWord:toIndexWord])
	log.Println("Index key Length : ", len(keyplain[fromIndexWord:toIndexWord]))
	fmt.Println("-------------------------***-------------------------------")
	Titxt := produceTi(keyplain)
	TiBin := binary(Titxt)
	stringTI := Titxt

	log.Println("Ti ASCII : ", []byte(Titxt))
	log.Println("Ti in Plain : ", Titxt)
	log.Println("Ti length : ", len(Titxt))
	log.Println("Ti in binary : ", TiBin)
	log.Println("Ti in binary length : ", len(TiBin))
	log.Println("Index Ti Word : ", stringTI[fromIndexWord:toIndexWord])
	log.Println("Index Ti Length : ", len(stringTI[fromIndexWord:toIndexWord]))
	fmt.Println("-------------------------***-------------------------------")

	chiper := encryptStreamChiper(TiBin, plainbin)
	log.Println("Chiper : ", chiper)
	log.Println("Chiper in binary length : ", len(chiper))
	chiperdecode := binaryToText(chiper)
	log.Println("Chiper ASCII : ", []byte(chiperdecode))
	log.Println("Chiper text : ", chiperdecode)
	log.Println("Chiper text length: ", len(chiperdecode))
	log.Println("Index Chiper Word: ", chiperdecode[fromIndexWord:toIndexWord])
	log.Println("Index Chiper Length : ", len(chiperdecode[fromIndexWord:toIndexWord]))
	plain := decryptStreamChiper(TiBin, chiper)
	log.Println("Plain : ", plain)
	log.Println("Plain in binary length : ", len(plain))
	plaindecode := binaryToText(plain)
	log.Println("Plain ASCII : ", []byte(plaindecode))
	log.Println("Plain text : ", plaindecode)
	log.Println("Plain text length: ", len(plaindecode))

	if plain != plainbin {
		panic("failed to decrypt")
	}

	fmt.Println("-------------------------*** Start Scenario SSE ***-------------------------------")
	log.Println("Alice give Bob index word : ", plaintext[fromIndexWord:toIndexWord])
	log.Println("Alice give Bob index word position : from " + strconv.Itoa(fromIndexWord) + " to " + strconv.Itoa(toIndexWord) + "")
	prooftxt := proof(chiperdecode, plaintext[fromIndexWord:toIndexWord], fromIndexWord, toIndexWord)
	log.Println("Bob give Alice proof : ", prooftxt)
	log.Println("Alice compare proof : " + prooftxt + " and Ti : " + stringTI[fromIndexWord:toIndexWord] + "")
	if prooftxt != stringTI[fromIndexWord:toIndexWord] {
		panic("failed to proofing")
	} else {
		log.Println("Proofing Success!")
	}

}

func binary(s string) string {
	res := ""
	for _, c := range s {
		res = fmt.Sprintf("%s%.8b", res, c)
	}
	return res
}

func GenerateRandomString(n int) (string, string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", "", err
		}
		ret[i] = letters[num.Int64()]
	}

	token := binary(string(ret))

	return string(ret), token, nil
}

func encryptStreamChiper(key string, message string) string {

	// keyarr := []rune(key)
	// messagearr := []rune(message)
	keyarr := []byte(key)
	messagearr := []byte(message)
	chipertxt := ""
	for i := 0; i < len(messagearr); i++ {
		ciphertmp := messagearr[i] ^ keyarr[i]
		// log.Println(string(messagearr[i]))
		// log.Println(keyarr[i])
		chipertxt = fmt.Sprintf("%s%s", chipertxt, fmt.Sprint(ciphertmp))
	}
	return chipertxt

}

func produceTi(key string) string {
	//keyarr := []rune(key)
	keyarr := []byte(key) //convert to ASCII

	////Ti := (Si,Fki (Si)) == Fki(Si) = 1^1 mod 5
	Ti := make([]byte, len(keyarr))

	for i := 0; i < len(keyarr); i++ {
		// res := math.Pow(float64(5), float64(keyarr[i])) //pangkat
		//res := math.Mod(float64(keyarr[i]), 14) //modulo
		res := int(keyarr[i]) + 3
		// log.Println(int(keyarr[i]))
		// log.Println("*************")
		// log.Println(res)
		copy(Ti[i:], string(int(res)))
	}

	return string(Ti)
}

func decryptStreamChiper(key string, chiper string) string {

	keyarr := []byte(key)
	chiperarr := []byte(chiper)
	plaintext := ""

	for i := 0; i < len(chiperarr); i++ {
		plaintexttmp := chiperarr[i] ^ keyarr[i]
		plaintext = fmt.Sprintf("%s%s", plaintext, fmt.Sprint(plaintexttmp))

	}

	return plaintext

}

func binaryToText(binarystring string) string {
	plaintext := ""
	n := 8
	for n <= len(binarystring) {
		lass := n - 8
		// log.Println(valuex)
		// log.Println(n)

		x := binarystring[lass:n]
		y, e := strconv.Atoi(x)
		if e != nil {
			panic(e)
		}
		//fmt.Println(string(binaryToDecimal(y)))
		plaintext = fmt.Sprintf("%s%s", plaintext, string(binaryToDecimal(y)))
		n += 8
	}

	return plaintext

}

func binaryToDecimal(num int) int {
	var remainder int
	index := 0
	decimalNum := 0
	for num != 0 {
		remainder = num % 10
		num = num / 10
		decimalNum = decimalNum + remainder*int(math.Pow(2, float64(index)))
		index++
	}
	return decimalNum
}

func proof(chiper string, word string, fromindex int, toindex int) string {
	chiperword := chiper[fromindex:toindex]
	wordarr := []byte(word)
	chiperwordarr := []byte(chiperword)
	//plaintext := ""
	// fmt.Println(word)
	// fmt.Println(chiperword)
	prooftxt := make([]byte, len(wordarr))
	for i := 0; i < len(wordarr); i++ {
		plaintexttmp := chiperwordarr[i] ^ wordarr[i]
		// fmt.Println(plaintexttmp)
		//plaintext = fmt.Sprintf("%s%s", plaintext, fmt.Sprint(plaintexttmp))
		copy(prooftxt[i:], string(int(plaintexttmp)))
	}

	// fmt.Println(plaintext)
	// fmt.Println(prooftxt)
	// fmt.Println(string(prooftxt))

	return string(prooftxt)

}
