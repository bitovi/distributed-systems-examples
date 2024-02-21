package app

import (
	"bitovi/distributed-systems-examples/schemas"
	"errors"
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

func (cfg WorkflowConfig) OrderWorkflow(ctx workflow.Context, input schemas.WorkflowInput) (string, error) {
	logger := workflow.GetLogger(ctx)

	// Save Order
	ao := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 3,
		},
	})

	var orderId string
	if err := workflow.ExecuteActivity(ao, "SaveOrderActivity", input).
		Get(ctx, &orderId); err != nil {
		logger.Info("Failed to save", "error", err.Error())
		return "", err
	}

	logger.Info("Saved Order", "orderId", orderId)

	// Create a timer
	timerFuture := workflow.NewTimer(ctx, 1*time.Minute)
	// Start Transmission
	c := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    30 * time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    5 * time.Minute,
			MaximumAttempts:    3,
		},
	})
	transmissionFuture := workflow.ExecuteActivity(c, "TransmitOrderActivity", input)

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
		c := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: 30 * time.Second,
			RetryPolicy: &temporal.RetryPolicy{
				MaximumAttempts: 3,
			},
		})

		if err := workflow.ExecuteActivity(c, "UpdateOrderActivity", "CONFIRMED", orderId).
			Get(ctx, nil); err != nil {
			return orderId, err
		}
	case "failed":
		// Call UpdateOrder with "FAILED" status
		c := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: 30 * time.Second,
			RetryPolicy: &temporal.RetryPolicy{
				MaximumAttempts: 3,
			},
		})

		if err := workflow.ExecuteActivity(c, "UpdateOrderActivity", "FAILED", orderId).
			Get(ctx, nil); err != nil {
			return orderId, err
		}
		return orderId, errors.New("order was failed")
	}

	return orderId, nil
}
