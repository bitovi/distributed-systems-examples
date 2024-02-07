package app

import (
	"bitovi/distributed-systems-examples/schemas"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

var Namespace = "distributed-systems-examples"

type WorkflowConfig struct {
}

func NewWorkflow() *WorkflowConfig {
	return &WorkflowConfig{}
}

func (cfg WorkflowConfig) Workflow(ctx workflow.Context, input schemas.WorkflowInput) (string, error) {
	logger := workflow.GetLogger(ctx)

	// Save Order
	orderId, err := cfg.SaveOrder(ctx, input)
	if err != nil {
		logger.Info("Failed to save", "error", err.Error())

		return "", err
	}
	logger.Info("Saved Order", "orderId", orderId)

	// Create a timer
	timerFuture := workflow.NewTimer(ctx, 5*time.Minute)
	// Start Transmission
	transmissionFuture := cfg.TransmitOrderFuture(ctx, input)

	selector := workflow.NewSelector(ctx)
	var outcome string

	selector.AddFuture(transmissionFuture, func(f workflow.Future) {
		var result string
		err := f.Get(ctx, &result)
		if err != nil {
			outcome = "failed"
		} else {
			outcome = "completed"
		}
	}).AddFuture(timerFuture, func(f workflow.Future) {
		outcome = "failed"
	})

	// Wait for one of the futures to complete
	selector.Select(ctx)

	switch outcome {
	case "completed":
		// Call UpdateOrder with "CONFIRMED" status
		if err := cfg.UpdateOrder(ctx, "CONFIRMED", *orderId); err != nil {
			return *orderId, err
		}
	case "failed":
		// Call UpdateOrder with "FAILED" status
		if err := cfg.UpdateOrder(ctx, "FAILED", *orderId); err != nil {
			return *orderId, err
		}
	}

	return *orderId, nil
}

func (cfg WorkflowConfig) SaveOrder(ctx workflow.Context, input schemas.WorkflowInput) (*string, error) {
	c := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 3,
		},
	})

	var orderId string
	if err := workflow.ExecuteActivity(c, "SaveOrderActivity", input).
		Get(ctx, &orderId); err != nil {
		return nil, err
	}
	return &orderId, nil
}

func (cfg WorkflowConfig) UpdateOrder(ctx workflow.Context, status, orderId string) error {
	c := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 3,
		},
	})

	if err := workflow.ExecuteActivity(c, "UpdateOrderActivity", status, orderId).
		Get(ctx, nil); err != nil {
		return err
	}
	return nil
}

func (cfg WorkflowConfig) TransmitOrderFuture(ctx workflow.Context, input schemas.WorkflowInput) workflow.Future {
	c := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    2 * time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    5 * time.Minute,
			MaximumAttempts:    3,
		},
	})

	return workflow.ExecuteActivity(c, "TransmitOrderActivity", input)
}
