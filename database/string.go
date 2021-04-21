package database

import "os"

var DATABASE_URL string = os.Getenv("DATABASE_URL_FOR_CRAWLER")
