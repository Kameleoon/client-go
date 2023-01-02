package utils

import (
	"fmt"
	"strings"
)

func ArrayToString(a []int, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func Remove(sequence interface{}, index int) interface{} {
	if seq, ok := sequence.([]interface{}); ok {
		seq[index] = seq[len(seq)-1]
		return seq[:len(seq)-1]
	}
	return nil
}
