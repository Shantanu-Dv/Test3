package utils

import (
	"doc-reco-go/internal/constants"
	"fmt"
	fuzzy "github.com/paul-mannino/go-fuzzywuzzy"
	"regexp"
	"strings"
)

func GetDoubtMessageId(context int) int {
	if context == constants.RecommendationContextConstant.PreChat {
		return 0
	}
	return -1
}

/*
	\u0400-\u04FF -> Russian
	\u0900-\u097F -> Devnagri
	\u0A00-\u0AFF -> Punjabi + Gujrati
	\u0C00-\u0C7F -> Telgue
	\u0B80-\u0BFF -> Tamil
	\u0D00–\u0D7F -> Malyalam
*/
var languageDetectRegex = regexp.MustCompile("[\u0400-\u04FF\u0900-\u097F\u0A00-\u0AFF\u0C00-\u0C7F\u0B80-\u0BFF\u0D00–\u0D7F]")

func CheckBadQueryError(searchTerm string, badQuery, isOcr bool) error {

	if !badQuery {
		return nil
	}

	if isOcr && len(strings.Split(searchTerm, " ")) < 3 {
		return fmt.Errorf(constants.BadQueryErrorMessage.InsufficientOcr)
	}

	if found := languageDetectRegex.Find([]byte(searchTerm)); found != nil {
		return fmt.Errorf(constants.BadQueryErrorMessage.UnsupportedLanguage)
	}

	return nil
}

var mathsSymbolRegex *regexp.Regexp

func ReplaceMathsTextToSymbol(text string) string {

	if mathsSymbolRegex == nil {
		var mathsTextKey []string
		for k, _ := range constants.MathsTextSymbolMap {
			mathsTextKey = append(mathsTextKey, k)
		}
		mathsSymbolRegex = regexp.MustCompile(fmt.Sprintf(`(?mi)\b(%s)\b`, strings.Join(mathsTextKey, "|")))
	}

	return mathsSymbolRegex.ReplaceAllStringFunc(text, func(s string) string {
		return constants.MathsTextSymbolMap[strings.ToLower(s)]
	})
}

func GetPercentageQuestionMatch(searchText, questionText string) int {
	// opts: forceAscii: false, cleanse: true
	return fuzzy.TokenSortRatio(searchText, questionText, false, true)
}
