package dbdetails

import (
	"fmt"
)

const (
	host     = "bhlx5eun6rbklyjklsjx-postgresql.services.clever-cloud.com"
	port     = 5432
	user     = "uaska2g7hzuc4dhy7non"
	password = "Ssl2e8p8woFbM7Y0YU8Z"
	dbname   = "bhlx5eun6rbklyjklsjx"
)

// SQLInfo is a function that returns a string to connect to the database
func SQLInfo() string {
	// psqlInfo is a string containing details to access database
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	return psqlInfo
}
