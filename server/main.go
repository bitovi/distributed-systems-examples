package main

import (
	app "bitovi/distributed-systems-examples"
	"bitovi/distributed-systems-examples/schemas"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
)

var (
	workflowClient client.Client
)

func main() {
	var err error
	// The client is a heavyweight object that should be created once per process.
	workflowClient, err = client.Dial(client.Options{
		HostPort:  os.Getenv("TEMPORAL_HOST_PORT"),
		Namespace: app.Namespace,
	})
	if err != nil {
		panic(err)
	}
	defer workflowClient.Close()

	fmt.Println("Starting server 3000...")
	http.HandleFunc("/order", createHandler)
	_ = http.ListenAndServe(":3000", nil)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
	var data schemas.WorkflowInput
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Start the workflow
	workflowOptions := client.StartWorkflowOptions{
		ID:        uuid.NewString(),
		TaskQueue: "demo",
	}

	wf := app.NewWorkflow()

	// Assuming you have a workflow function defined as YourWorkflow
	we, err := workflowClient.ExecuteWorkflow(context.Background(), workflowOptions, wf.Workflow, data)
	if err != nil {
		log.Fatalln("Error starting workflow", err)
	}
	fmt.Printf("Created new order id:%s.\n", workflowOptions.ID)

	// Optionally, send back the WorkflowID and RunID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"orderID": we.GetID(),
	})
}
