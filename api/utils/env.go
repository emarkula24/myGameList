package utils

import (
	"database/sql"

	giantbomb "example.com/mygamelist/interfaces"
)

type Env struct {
	DB       *sql.DB
	FrontUrl string
	API      giantbomb.Client
}
