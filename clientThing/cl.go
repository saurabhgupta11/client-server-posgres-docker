package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	_ "debug/gosym"
	"fmt"
	"io/ioutil"
	"log"

	"./db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// load certificate of CA who signed web server's certificate
	pemServerCA, err := ioutil.ReadFile("cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}

	// create the credentials
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server's ca certificate")
	}

	config := &tls.Config{
		RootCAs:            certPool,
		InsecureSkipVerify: true,
	}

	return credentials.NewTLS(config), nil
}

func main() {
	var conn *grpc.ClientConn

	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("FAiled to load TLS credentials", err)
	}

	conn, err = grpc.Dial(":9000", grpc.WithTransportCredentials(tlsCredentials))
	if err != nil {
		log.Fatal("Error")
	}

	defer conn.Close()

	c := db.NewDatabaseServiceClient(conn)

	// infinte loop to make all the requests until client wants to exit
	var option int64
	var scanStruct = db.SingleRow{
		Id:        0,
		Age:       0,
		Firstname: "",
		Lastname:  "",
		Email:     "",
	}
	var scanId = db.Id{
		Id: 0,
	}
	for true {
		fmt.Println("\nEnter the number to get data/query you want")
		fmt.Println("1. Display Database")
		fmt.Println("2. Insert")
		fmt.Println("3. DeleteById")
		fmt.Println("4. FindByID")
		fmt.Println("5. UpdateByID")
		fmt.Print("Enter the option: ")
		_, err = fmt.Scan(&option)
		if err != nil {
			fmt.Println(err)
		}
		if option == 1 {
			message := db.Empty{}
			response, er := c.GetDB(context.Background(), &message)
			if er != nil {
				log.Fatal("option 1 failed", er)
			}
			fmt.Println("Database : ")
			for _, s := range response.Rows {
				fmt.Println(s)
			}
			fmt.Println()
		} else if option == 2 {
			fmt.Print("Enter Age: ")
			fmt.Scan(&scanStruct.Age)
			fmt.Print("Enter your First Name: ")
			fmt.Scan(&scanStruct.Firstname)
			fmt.Print("Enter your Last Name: ")
			fmt.Scan(&scanStruct.Lastname)
			fmt.Print("Enter your email: ")
			fmt.Scan(&scanStruct.Email)

			response, er := c.Insert(context.Background(), &scanStruct)
			if er != nil {
				log.Fatal("option 2 failed: ", er)
			}
			fmt.Println("User data saved: ", response)
		} else if option == 3 {
			fmt.Print("Enter Id: ")
			fmt.Scan(&scanId.Id)

			_, er := c.DeleteByID(context.Background(), &scanId)
			if er != nil {
				log.Fatal("option 3 failed: ", er)
			}
			fmt.Println("User data Deleted")
		} else if option == 4 {
			fmt.Print("Enter Id: ")
			fmt.Scan(&scanId.Id)

			response, er := c.FindByID(context.Background(), &scanId)
			if er != nil {
				fmt.Printf("No user with id: %d found\n", scanId.Id)
			} else {
				fmt.Println("User data : ", response)
			}
		} else if option == 5 {
			fmt.Print("Enter Id: ")
			fmt.Scan(&scanStruct.Id)
			fmt.Print("Enter Age: ")
			fmt.Scan(&scanStruct.Age)
			fmt.Print("Enter your First Name: ")
			fmt.Scan(&scanStruct.Firstname)
			fmt.Print("Enter your Last Name: ")
			fmt.Scan(&scanStruct.Lastname)
			fmt.Print("Enter your email: ")
			fmt.Scan(&scanStruct.Email)

			response, er := c.UpdateByID(context.Background(), &scanStruct)
			if er != nil {
				log.Fatal("option 5 failed: ", er)
			}
			fmt.Println("User data Updated: ", response)
		} else {
			fmt.Println("You entered an invalid value.....so exiting")
			break
		}
	}
}
