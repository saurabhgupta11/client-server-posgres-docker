package db

import (
	"context"
	"database/sql"

	"../dbconnection"
)

// Db global pointer for database access in API struct methods
var dbaccess *sql.DB

func init() {
	dbaccess = dbconnection.ConnectDB()
}

// Server struct
type Server struct{}

// GetDB gets data from the database and returns to the client
func (s *Server) GetDB(ctx context.Context, message *Empty) (*Rows, error) {

	var reply []*SingleRow

	sqlStatement := `SELECT * from users
		ORDER BY id
	`
	rows, err := dbaccess.Query(sqlStatement)
	if err != nil {
		// handle this error better than this
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		record := &SingleRow{
			Id:        0,
			Age:       0,
			Firstname: "",
			Lastname:  "",
			Email:     "",
		}
		err = rows.Scan(&record.Id, &record.Age, &record.Firstname, &record.Lastname, &record.Email)
		if err != nil {
			// handle this error
			panic(err)
		}
		reply = append(reply, record)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return &Rows{Rows: reply}, nil
}

// Insert adds a row into the database
func (s *Server) Insert(ctx context.Context, message *SingleRow) (*SingleRow, error) {
	sqlStatement := `
		INSERT INTO
		users (age, firstname, lastname, email)
		VALUES ($1, $2, $3, $4)
		RETURNING id,age,firstname,lastname,email
	`

	// data will be feeded into the reply with scan in QueryRow
	reply := &SingleRow{
		Id:        0,
		Age:       0,
		Firstname: "",
		Lastname:  "",
		Email:     "",
	}

	err := dbaccess.QueryRow(sqlStatement, message.Age, message.Firstname, message.Lastname, message.Email).Scan(&reply.Id, &reply.Age, &reply.Firstname, &reply.Lastname, &reply.Email)
	if err != nil {
		panic(err)
	}

	return reply, nil
}

// DeleteByID deletes a users data by the Id
func (s *Server) DeleteByID(ctx context.Context, message *Id) (*Empty, error) {
	sqlStatement := `
		DELETE FROM users
		WHERE id = $1;
	`
	_, err := dbaccess.Exec(sqlStatement, message.Id)
	if err != nil {
		return &Empty{}, err
	}

	return &Empty{}, nil
}

// FindByID finds a user by the id in the database
func (s *Server) FindByID(ctx context.Context, message *Id) (*SingleRow, error) {
	sqlStatement := `
		SELECT id, age, firstname, lastname, email 
		FROM users 
		WHERE id=$1;
	`
	reply := &SingleRow{
		Id:        0,
		Age:       0,
		Firstname: "",
		Lastname:  "",
		Email:     "",
	}
	row := dbaccess.QueryRow(sqlStatement, message.Id)
	switch err := row.Scan(&reply.Id, &reply.Age, &reply.Firstname, &reply.Lastname, &reply.Email); err {
	case sql.ErrNoRows:
		return reply, sql.ErrNoRows
	case nil:
		return reply, nil
	default:
		panic(err)
	}
}

// UpdateByID finds a user by the id in the database
func (s *Server) UpdateByID(ctx context.Context, message *SingleRow) (*SingleRow, error) {
	sqlStatement := `
		UPDATE users
		SET age = $2, firstname = $3, lastname = $4, email = $5
		WHERE id=$1
		RETURNING id,age,firstname,lastname,email;
	`
	reply := &SingleRow{
		Id:        0,
		Age:       0,
		Firstname: "",
		Lastname:  "",
		Email:     "",
	}
	row := dbaccess.QueryRow(sqlStatement, message.Id, message.Age, message.Firstname, message.Lastname, message.Email)
	switch err := row.Scan(&reply.Id, &reply.Age, &reply.Firstname, &reply.Lastname, &reply.Email); err {
	case sql.ErrNoRows:
		return reply, sql.ErrNoRows
	case nil:
		return reply, nil
	default:
		panic(err)
	}
}
