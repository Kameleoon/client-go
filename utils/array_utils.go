package utils

import (
	"fmt"
	"strings"
)

func ArrayToString(a []int, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}

func Contains(s []interface{}, str string) bool {
	for _, v := range s {
		if _, ok := v.(string); ok {
			if v == str {
				return true
			}
		}
	}
	return false
}

func SliceContains[T comparable](s []T, x T) bool {
	for _, v := range s {
		if v == x {
			return true
		}
	}
	return false
}

func Map(array []interface{}, f func(interface{}) interface{}) []interface{} {
	mapArray := make([]interface{}, len(array))
	for i, v := range array {
		mapArray[i] = f(v)
	}
	return mapArray
}

func MapToStringDict(array []interface{}, f func(*map[string]interface{}, interface{})) map[string]interface{} {
	mapDict := make(map[string]interface{}, len(array))
	for _, v := range array {
		f(&mapDict, v)
	}
	return mapDict
}

func Remove(sequence interface{}, index int) interface{} {
	if seq, ok := sequence.([]interface{}); ok {
		seq[index] = seq[len(seq)-1]
		return seq[:len(seq)-1]
	}
	return nil
}
