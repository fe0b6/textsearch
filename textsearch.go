package textsearch

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	mystemBin   string
	minWordLen  int
	ruReg       *regexp.Regexp
	strSplitReg *regexp.Regexp
)

// Init - инициализация
func Init(path string, ml int) {
	mystemBin = path
	minWordLen = ml
	ruReg = regexp.MustCompile("^[а-яА-Я]+$")
	strSplitReg = regexp.MustCompile("[^a-zA-Zа-яА-Я0-9]")
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

	tmpInd := [][]string{}
	for _, s := range words {
		s = strings.Replace(strings.Replace(s, "\n", " ", -1), "\r", "", -1)

		arr := []string{}
		for _, w := range strSplitReg.Split(s, -1) {
			if ruReg.MatchString(w) {
				continue
			}

			arr = append(arr, w)
		}
		tmpInd = append(tmpInd, arr)

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

	for i, b := range strings.Split(string(outBuffer.Bytes()), "\n") {
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
			w := getBest(a.Analysis)
			if utf8.RuneCountInString(w) < minWordLen {
				continue
			}
			arr = append(arr, w)
		}

		// ДОбавляем не русские слова и цифры
		if len(tmpInd[i]) > 0 {
			arr = append(arr, tmpInd[i]...)
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

// Search - Поиск совпадений
func (q Query) Search(ind []string) (ok bool) {
	for _, w := range q.Words {
		for _, w2 := range ind {
			if w2 == w {
				ok = true
				return
			}
		}
	}

	return
}
