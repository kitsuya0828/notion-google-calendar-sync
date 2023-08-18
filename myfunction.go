package mycloudeventfunction

import (
	"context"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/Kitsuya0828/notion-googlecalendar-sync/run"
	"github.com/cloudevents/sdk-go/v2/event"
)

func init() {
    // Register a CloudEvent function with the Functions Framework
    functions.CloudEvent("MyCloudEventFunction", myCloudEventFunction)
}

// Function myCloudEventFunction accepts and handles a CloudEvent object
func myCloudEventFunction(ctx context.Context, e event.Event) error {
    // Your code here
    // Access the CloudEvent data payload via e.Data() or e.DataAs(...)
	run.Run()

    // Return nil if no error occurred
    return nil
}