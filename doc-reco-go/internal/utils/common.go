package utils

import (
	"encoding/json"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/go-playground/validator"
	"net/http"
	"reflect"
	"strings"
	"unicode"
)

func FindElementInSlice(element interface{}, slice interface{}) bool {
	s := reflect.ValueOf(slice)

	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Interface() == element {
			return true
		}
	}
	return false
}

func GetSliceDifference(a, b []int) []int {
	var diff []int
	for _, elm := range a {
		if !FindElementInSlice(elm, b) {
			diff = append(diff, elm)
		}
	}
	return diff
}

func LimitWords(text string, limit int) string {
	if limit > len(text) {
		limit = len(text)
	}
	return text[:limit]
}

func ValidateRequestStructTag(request interface{}) error {
	validate := validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	err := validate.Struct(request)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return fmt.Errorf("%s %s", err.Field(), err.Tag())
		}
	}
	return nil
}

func ParseRequestBody(r *http.Request, t interface{}) error {
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		return err
	}

	return r.Body.Close()
}

func SentryCaptureError(r *http.Request, err error) {
	ctx := r.Context()
	hub := sentry.GetHubFromContext(ctx)
	hub.Scope().SetTags(map[string]string{
		"url":   r.URL.Path,
		"error": err.Error(),
	})
	hub.CaptureMessage(err.Error())
}

func FilterSearchText(text string) (filteredText string) {
	for _, r := range text {
		if unicode.IsGraphic(r) {
			filteredText += string(r)
		}
	}

	replacer := strings.NewReplacer("\t", " ", "\n", " ", "\r", "", "\a", "", "\"", "", `\`, `\\`)
	filteredText = replacer.Replace(text)
	filteredText = strings.Trim(filteredText, " ")
	return
}

func ChainArr(array1 []string, array2 string, array3 []string, array4[]string) []string {
	jointStringArr := array1
	jointStringArr = append(jointStringArr, array2)
	jointStringArr = append(jointStringArr, array3...)
	jointStringArr = append(jointStringArr, array4...)

	return jointStringArr
}