package utils

import (
	"doc-reco-go/internal/constants"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"regexp"
	"strings"
)

func FindEndOfMath(delimiter string, text string, startIndex int) int {
	/*
		helper function to find latex delimiters Adapted from
		https://github.com/Khan/perseus/blob/master/src/perseus-markdown.jsx
	*/

	index := startIndex
	braceLevel := 0

	delimLen := len(delimiter)

	for ; index < len(text) - 1; {
		character := string(text[index])

		if braceLevel <= 0 && text[index:index+delimLen] == delimiter {
			return index
		} else if character == "\\" {
			index += 1
		} else if character == "{" {
			braceLevel += 1
		} else if character == "}" {
			braceLevel -= 1
		}

		index += 1
	}
	return -1
}

func SplitAtDelimiters(startData []map[string]interface{}, leftDelim, rightDelim string, display bool) []map[string]interface{} {
	/*
		helper function to find latex delimiters and split at it
		:param startData:
		:param leftDelim:
		:param rightDelim:
		:param display:
		:return:
	*/

	var finalData []map[string]interface{}

	for _, data := range startData {
		if data["type"] != "text" {
			finalData = append(finalData,data)
			//continue
		}
		text := data["data"].(string)

		lookingForLeft := true
		currIndex := 0

		nextIndex := strings.Index(text, leftDelim)
		if nextIndex != -1 {
			currIndex = nextIndex
			finalData = append(finalData, map[string]interface{}{
				"type":  "text",
				"data":  text[0:currIndex],
				"start": 0,
				"end":   nextIndex,
			})
			lookingForLeft = false
		}

		for ; true; {
			if lookingForLeft {
				nextIndex = strings.Index(text[currIndex:], leftDelim)
				if nextIndex == -1 {
					break
				}
				nextIndex += currIndex
				finalData = append(finalData, map[string]interface{}{
					"type":  "text",
					"data":  text[currIndex:nextIndex],
					"start": currIndex,
					"end":   nextIndex,
				})
				currIndex = nextIndex
			} else {
				nextIndex = FindEndOfMath(rightDelim, text, currIndex+len(leftDelim))
				if nextIndex == -1 {
					break
				}
				finalData = append(finalData, map[string]interface{}{
					"type": "math",
					"data": text[currIndex+len(leftDelim):nextIndex],
					"rawData": text[currIndex:nextIndex+len(rightDelim)],
					"display":    display,
					"start":      currIndex,
					"end":        nextIndex + len(rightDelim),
					"leftDelim":  leftDelim,
					"rightDelim": rightDelim,
				})
				currIndex = nextIndex + len(rightDelim)
			}

			lookingForLeft = !lookingForLeft
		}

		finalData = append(finalData, map[string]interface{}{
			"type":  "text",
			"data":  text[currIndex:],
			"start": currIndex,
			"end":   len(text),
		})
	}
	return finalData
}


func SeparateMath(text_ string) ([]string, []string, error){
	/*
		Returns separate lists for normal text and latex
		:param text_:
		:param delimiter:
		:param return_tokens:
		:return:
	*/
	delimiter := map[string]interface{}{"left": "$$", "right": "$$", "display": true}
	data := []map[string]interface{}{{"type": "text", "data": text_}}

	var text []string
	var math []string
	for _, token := range SplitAtDelimiters(data, delimiter["left"].(string), delimiter["right"].(string), delimiter["display"].(bool)){
		if token["type"] == "text" {
			text = append(text, token["data"].(string))
		}else if token["type"] == "math" {
			math = append(math, token["data"].(string))
		} else {
			return text, math, fmt.Errorf("invalid Token Type")
		}
	}
	return text, math, nil
}

func SplitWithDelimiters(text string, delimiters []map[string]interface{}) []map[string]interface{} {
	/*
	function to seperate latex and text into a list of phrases
	:param text:
	:param delimiters:
	:return:
	*/
	data := []map[string]interface{}{{"type": "text", "data": text}}

	for _, delimiter := range delimiters {
		data = SplitAtDelimiters(data, delimiter["left"].(string), delimiter["right"].(string), false)
	}
	return data
}

func UseShingles(text string, minWordCount int) bool {
	if len(strings.Split(text, " ")) > minWordCount {
		return false
	}
	return true
}

func StripHtmlTags(str string) string {
	policy := bluemonday.StripTagsPolicy()
	result := policy.Sanitize(str)
	return result
}



func handleLatexSymbols(text string) string {
	var latexSymbolRegex *regexp.Regexp

	var latexSymbolKey []string

	for k,_ := range constants.LatexMap {
		latexSymbolKey = append(latexSymbolKey,k)
	}

	latexSymbolRegex = regexp.MustCompile(fmt.Sprintf(`(\%s)\b`,strings.Join(latexSymbolKey,`|\`)))

	return latexSymbolRegex.ReplaceAllStringFunc(text, func(s string) string{
		return constants.LatexMap[s]
	})
}

func handleCaptureGroups(text string) string {
	for k,v:=range constants.LatexCaptureGroupMap {
		var re = regexp.MustCompile(fmt.Sprintf(`(?mU)\%v\s*{(.*)}\s*{(.*)}`,k))
		sub := v
		text = re.ReplaceAllString(text,sub)
	}
	return text
}

func handleMatrix(text string) string {
	var re = regexp.MustCompile(`\\begin{(?:\w?|small)matrix}(.*)\\end{(?:\w?|small)matrix}`)
	return re.ReplaceAllString(text, "|$1|")
}

func LatexToPlainText (str string) string {
	str = handleLatexSymbols(str)

	//handle matrix
	str = handleMatrix(str)

	// capture groups
	str = handleCaptureGroups(str)

	// remove escape chars
	for k,v := range constants.LatexEscapeSymbols {
		str = strings.Replace(str,k,v,-1)
	}

	// two symbols handled separately as they are used in escape chars
	str = strings.Replace(str,`&`,``,-1)
	str = strings.Replace(str,`#`,``,-1)
	str = strings.Replace(str,`\`,``,-1)

	return str
}