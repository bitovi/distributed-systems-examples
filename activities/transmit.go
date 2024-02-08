package activities

import (
	"os"
	"bitovi/distributed-systems-examples/schemas"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.temporal.io/sdk/activity"
)

type TransmitConfig struct{}

func NewTransmitOrder() *TransmitConfig {
	return &TransmitConfig{}
}

func (cgf TransmitConfig) TransmitOrderActivity(ctx context.Context, i schemas.WorkflowInput) error {
	jsonData, err := json.Marshal(i)
	if err != nil {
		activity.GetLogger(ctx).Error("Error marshaling JSON", "Error", err)
		return err
	}

	url := os.Getenv("TRANSMIT_URL")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		activity.GetLogger(ctx).Error("Error creating request", "Error", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		activity.GetLogger(ctx).Error("Error sending request", "Error", err)
		return err
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	return nil
}
