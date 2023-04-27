package models

type Log struct {
	Path         string `json:"path"`
	UserId       int    `json:"UserId"`
	Method       string `json:"method"`
	StatusCode   int    `json:"statusCode"`
	ResponseTime int    `json:"responseTime"`
	Origin       string `json:"origin"`
}