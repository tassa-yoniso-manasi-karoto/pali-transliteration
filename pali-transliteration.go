package pli

import (
	"strings"
	"unicode/utf8"
	"regexp"
)

	
var m = map[string]string{
	"อ": "",
	"ภ": "bh",
	"ม": "m",
	"ล": "l",
	"พ": "b",
	"ก": "k",
	"ข": "kh",
	"ค": "g",
	"ฆ": "gh",
	"ง": "ṅ",
	"จ": "c",
	"ฉ": "ch",
	"ช": "j",
	"ฌ": "jh",
	"ญ": "ñ",
	"ฏ": "ṭ",
	"ฐ": "ṭh",
	"ฑ": "ḍ",
	"ฒ": "ḍh",
	"ณ": "ṇ",
	"ต": "t",
	"ถ": "th",
	"ท": "d",
	"ธ": "dh",
	"น": "n",
	"ป": "p",
	"ผ": "ph",
	"ย": "y",
	"ร": "r",
	"ว": "v",
	"ส": "s",
	"ห": "h",
	"ฬ": "ḷ",
	"ะ": "a",
	"ั":  "a",
	"ุ":  "u",
	"า": "ā",
	"ิ":  "i",
	"ี":  "ī",
	"ู":  "ū",	
	"เ": "e",	
	"โ": "o",	
	"์":  "-",	
	"ํ":  "ṃ", //"ṁ",
	"๐": "0",
	"๑": "1",
	"๒": "2",
	"๓": "3",
	"๔": "4",
	"๕": "5",
	"๖": "6",
	"๗": "7",
	"๘": "8",
	"๙": "9",
	"ฺ":  "",
	"ฯ": ".",
	// unicode committee be like: You thought "ṃ" and "ṃ" were the same? How simple minded.
	// You think we wouldn't throw at you visually indistinguishable characters just be cause we can? ...Haha.
	"ึ":  "iṃ",
	"สฺม": "sm",
	"สฺว": "sv",
	"ทฺว": "dv",
}

func ThaiToRoman(str string, mode int) (out string) {
	after := ""
	ToAddafter := []string{"เ", "โ"}
	if mode == 1 {
		for str != "" {
			r, _ := utf8.DecodeRuneInString(str)
			char := string(r)
			if corresp, ok := m[char]; ok {
				if !contains(ToAddafter, char) {
					out += corresp + after
					after = "" 
				} else {
					after = corresp
				}
			} else {
				out += char
			}
			str = strings.TrimPrefix(str, char)
		}
	} else if mode > 1 {
		cons   := strings.Split("มลพอกขคฆงภจฉชฌญฏฐฑฒณตถทธปผยรวสหฬน", "")
		vowels := strings.Split("ะัิีึืุูา", "")
		combinations := []string{"ทฺว", "สฺว", "สฺม"}
		for str != "" {
			found := false
			for _, comb := range combinations {
				if str, found = strings.CutPrefix(str, comb); found {
					out += m[comb] + after
					after = ""
					continue
				}
			}
			r, size := utf8.DecodeRuneInString(str)
			char := string(r)
			nxt, _ := utf8.DecodeRuneInString(str[size:])
			charnxt := string(nxt)
			if corresp, ok := m[char]; ok {
				if !contains(ToAddafter, char) {
					if after == "" && contains(cons, char) && charnxt != "ฺ" && !contains(vowels, charnxt) {
						after = "a"
					}
					out += corresp + after
					after = "" 
				} else {
					after = corresp
				}
			} else {
				out += char
			}
			str = strings.TrimPrefix(str, char)
		}
		re := regexp.MustCompile(` +\.`)
		out = re.ReplaceAllString(out, ".")
	} else {
		return str
	}
	return
}


func contains[T comparable](arr []T, i T) bool {
	for _, a := range arr {
		if a == i {
			return true
		}
	}
	return false
}
