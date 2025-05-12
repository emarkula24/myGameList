package utils

import "database/sql"

type Env struct {
	DB       *sql.DB
	FrontUrl string
}
