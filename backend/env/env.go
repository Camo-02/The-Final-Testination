package env

import (
	"backend/utils"
)

var DB_PORT = utils.GetEnv("DB_PORT")
var DB_HOST = utils.GetEnv("DB_HOST")
var DB_USERNAME = utils.GetEnv("DB_USERNAME")
var DB_PASSWORD = utils.GetEnv("DB_PASSWORD")
var DB_NAME = utils.GetEnv("DB_NAME")
var API_PORT = utils.GetEnv("PUBLIC_API_PORT")
