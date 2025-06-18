package onboarding

import (
    "context"
    "fmt"
    "time"

<<<<<<< Updated upstream
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "google.golang.org/protobuf/types/known/durationpb"
=======
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
>>>>>>> Stashed changes

    "go.temporal.io/api/workflowservice/v1"
    "go.temporal.io/sdk/client"
)

func RegisterNamespace(ctx context.Context, hostPort, namespace string, retentionDays int) error {
    nsClient, err := client.NewNamespaceClient(client.Options{HostPort: hostPort})
    if err != nil {
        return fmt.Errorf("unable to create namespace client: %w", err)
    }
    defer nsClient.Close()

<<<<<<< Updated upstream
    if retentionDays < 1 {
        retentionDays = 1
    }
    retention := &durationpb.Duration{
        Seconds: int64(retentionDays) * int64(24*time.Hour) / int64(time.Second),
    }

    req := &workflowservice.RegisterNamespaceRequest{
        Namespace:                        namespace,
        WorkflowExecutionRetentionPeriod: retention,
    }
    if err := nsClient.Register(ctx, req); err != nil {
        if s, ok := status.FromError(err); ok && s.Code() == codes.AlreadyExists {
            return nil
        }
        return fmt.Errorf("failed to register namespace %q: %w", namespace, err)
    }
=======
	if _, err := nsClient.Describe(ctx, namespace); err == nil {
		return nil
	} else if s, ok := status.FromError(err); !ok || s.Code() != codes.NotFound {
		return fmt.Errorf("failed to describe namespace %q: %w", namespace, err)
	}

	if retentionDays < 1 {
		retentionDays = 1
	}
	retention := &durationpb.Duration{
		Seconds: int64(retentionDays) * int64(24*time.Hour) / int64(time.Second),
	}

	req := &workflowservice.RegisterNamespaceRequest{
		Namespace:                        namespace,
		WorkflowExecutionRetentionPeriod: retention,
	}
	if err := nsClient.Register(ctx, req); err != nil {
		return fmt.Errorf("failed to register namespace %q: %w", namespace, err)
	}
>>>>>>> Stashed changes

    return nil
}
