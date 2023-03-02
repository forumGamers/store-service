package helper

import "regexp"

func SlugGenerator(slug string) string {
		r := regexp.MustCompile(" ")
		res := r.ReplaceAllString(slug,"_")

	return res
}