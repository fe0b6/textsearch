package textsearch

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
	"strings"
)

var (
	mystemBin string
)

// Init - инициализация
func Init(path string) {
	mystemBin = path
}

// GetIndex - строку а возвращаем индекс для нее
func GetIndex(word string) (ind []string, err error) {
	ans, err := GetIndexes([]string{word})
	if err != nil {
		log.Println("[error]", err)
		return
	}

	ind = ans[0]
	return
}

// GetIndexes - строки а возвращаем индекс для нее
func GetIndexes(words []string) (ind [][]string, err error) {
	var inputBuffer, outBuffer bytes.Buffer

	for _, s := range words {
		s = strings.Replace(strings.Replace(s, "\n", " ", -1), "\r", "", -1)
		inputBuffer.Write([]byte(s + "\n"))
	}

	proc := exec.Command(mystemBin, "--format", "json", "--generate-all", "--weight")
	proc.Stdin = &inputBuffer
	proc.Stdout = &outBuffer

	err = proc.Start()
	if err != nil {
		log.Println("[error]", err)
		return
	}

	err = proc.Wait()
	if err != nil {
		log.Println("[error]", err)
		return
	}

	ind = [][]string{}
	for _, b := range strings.Split(string(outBuffer.Bytes()), "\n") {
		if b == "" {
			continue
		}
		var ans []answer
		err = json.Unmarshal([]byte(b), &ans)
		if err != nil {
			log.Println("[error]", err)
			return
		}

		arr := []string{}
		for _, a := range ans {
			arr = append(arr, getBest(a.Analysis))
		}

		ind = append(ind, arr)
	}

	return
}

func getBest(arr []analys) (ans string) {
	var best float64
	for _, a := range arr {
		if a.Wt > best {
			ans = a.Lex
			best = a.Wt
		}
	}

	return
}
