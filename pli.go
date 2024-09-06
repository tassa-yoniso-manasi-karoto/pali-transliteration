package pli

import (
	"strings"
	"unicode/utf8"
	"regexp"

	"github.com/gookit/color"
	"github.com/k0kubun/pp"
)

	//TODO Yamakkan ๎
var (
	wantDebug = false
	p = map[string]any{}
	reSpace = regexp.MustCompile(`\p{Z}+|[\n\r]+`)
	m = map[string]string{
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
		"ึ":  "iṃ",
		"สฺม": "sm",
		"สฺว": "sv",
		"ทฺว": "dv",
	}
)

func ThaiToLatin(str string, mode int) (out string) {
	after := ""
	ToAddafter := []string{"เ", "โ"}
	if mode == 1 { //PHONETIC (NORMAL) STYLE THAI PALI
		for str != "" {
			r, _ := utf8.DecodeRuneInString(str)
			char := string(r)
			str = strings.TrimPrefix(str, char)
			if corresp, ok := m[char]; ok {
				if !contains(ToAddafter, char) {
					r, _ := utf8.DecodeRuneInString(str)
					nextchar := string(r)
					if _, ok := m[nextchar]; char == "ง" && !ok {
						corresp = "ṁ"
					}
					out += corresp + after
					after = "" 
				} else {
					after = corresp
				}
			} else {
				out += char
			}
		}
	} else if mode == 2 { // PINTU STYLE THAI PALI
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



func placeholder() {
	color.Redln(" 𝒻*** 𝓎ℴ𝓊 𝒸ℴ𝓂𝓅𝒾𝓁ℯ𝓇")
	pp.Println("𝓯*** 𝔂𝓸𝓾 𝓬𝓸𝓶𝓹𝓲𝓵𝓮𝓻")
}

