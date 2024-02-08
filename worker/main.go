package main

import (
	app "bitovi/distributed-systems-examples"
	"bitovi/distributed-systems-examples/activities"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.uber.org/zap"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func main() {
	temporalClient, err := client.Dial(client.Options{
		Namespace: app.Namespace,
		HostPort:  "temporal-dev-server:7233",
	})
	if err != nil {
		log.Fatalln("Unable to start Temporal client, is the server running? ", zap.Error(err))
	}
	defer temporalClient.Close()

	db, err := connectDB()
	if err != nil {
		panic(err)
	}

	w := worker.New(temporalClient, "demo", worker.Options{})

	wf := app.NewWorkflow()
	w.RegisterWorkflow(wf.Workflow)
	w.RegisterActivity(activities.NewSaveOrder(db).SaveOrderActivity)
	w.RegisterActivity(activities.NewUpdateOrder(db).UpdateOrderActivity)
	w.RegisterActivity(activities.NewTransmitOrder().TransmitOrderActivity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}

func connectDB() (*sql.DB, error) {
	loadEnv()

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s sslmode=disable",
		host, port, user, password)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully connected!")
	return db, nil
}
