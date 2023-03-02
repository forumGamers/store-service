package helper

import (
	"math"
	"regexp"
	"strconv"
)

func SlugGenerator(slug string) string {
	var p int
	for i := 0 ;i < len(slug);i++ {
		if slug[i] == ' ' {
			p++
		}else if i == 0 {
			p++
		}else if i == len(slug) - 1 {
			p++
		}
	}

	pre := "%"+strconv.FormatFloat(math.Floor(float64( 100 / p)),'f',2,64)+"%"

	r := regexp.MustCompile(" ")
	res := r.ReplaceAllString(slug,pre)

	return pre+res+pre
}