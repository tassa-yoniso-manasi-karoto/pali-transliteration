package pli

import (
	"strings"
	"unicode/utf8"
	"regexp"
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"runtime"
	"path/filepath"
	"sort"
	
	"github.com/rs/zerolog/log"
	"github.com/gookit/color"
	"github.com/k0kubun/pp"
	"github.com/rivo/uniseg"
)

	//TODO Yamakkan ๎
var (
	wantDebug = false
	kanaPatterns []string
	kanaScheme = make(map[string]string)
	p = map[string]any{}
	atm = []string{} // atm = at the moment
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

func ThaiToRoman(str string, mode int) (out string) {
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


func RomanToKana(str string) (out string) {
	if len(kanaPatterns) == 0 {
		if ok := initKana(); !ok {
			kanaPatterns = kanaPatternsBackup
			kanaScheme = kanaSchemeBackup
		}
	}
	str = strings.ToLower(str)
	str = strings.ReplaceAll(str, "ṃ", "ṁ")
Outerloop:
	for str != "" {
		for _, s := range kanaPatterns {
			if _, found := strings.CutPrefix(str, s); found {
				if wantDebug { color.Greenln("MATCHED:", s)}
				out += kanaScheme[s]
				str = strings.TrimPrefix(str, s)
				continue Outerloop
			}
		}
		r, _ := utf8.DecodeRuneInString(str)
		c := string(r)
		out += c
		if wantDebug {
			if !reSpace.MatchString(c) {
				color.Redf("MATCH FAILED: \"%s\"\n", c)
			} else {
				fmt.Print("\n")
			}
		}
		str = strings.TrimPrefix(str, c)
	}
	return
}


func initKana() bool {
	ex, err := os.Executable()
	if err != nil {
		log.Error().Err(err).Msg("An error occurred")
		return false
	}
	path := filepath.Dir(ex)
	if runtime.GOOS == "windows" {
		path += `\`
	} else {
		path += "/"
	}
	data, err := ioutil.ReadFile(path+"transliteration.json")
	if err != nil {
		log.Error().Err(err).Msg("Can't access transliteration.json")
		return false
	}
	err = json.Unmarshal(data, &p)
	if err != nil {
		log.Error().Err(err).Msg("Unmarshal error")
		return false
	}
	parseKanaTree(0, p)
	for k, _ := range kanaScheme {
		kanaPatterns = append(kanaPatterns, k)
	}
	sort.Slice(kanaPatterns, func(i, j int) bool {
		l1, l2 := uniseg.GraphemeClusterCount(kanaPatterns[i]), uniseg.GraphemeClusterCount(kanaPatterns[j])
		if l1 != l2 {
			return l1 > l2
		}
		return kanaPatterns[i] > kanaPatterns[j]
	})
	return true
}



func parseKanaTree(depth int, m map[string]interface{}) {
	// atm = at the moment
	for k, v := range m {
		if len(atm) > depth {
			//fmt.Println(atm, "===>", atm[:depth+1])
			atm = atm[:depth]
		}
		switch v.(type) {
		case string:
			tmp := strings.Join(append(atm, strings.ToLower(k)), "")
			//fmt.Println(tmp, "=", strings.ToLower(v.(string)))
			kanaScheme[tmp] = strings.ToLower(v.(string))
		case map[string]interface{}:
			atm = append(atm, strings.ToLower(k))
			parseKanaTree(depth+1, v.(map[string]interface{}))
		default:
			panic("wrong house you fool")
		}
	}
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

