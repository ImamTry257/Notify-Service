package validator

import (
	"strings"

	notifypbv2 "github.com/ImamTry257/lms-proto-notify/gen/go/notify"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var validNotifyTypes = map[string]bool{
	"OTP":            true,
	"RESET_PASSWORD": true,
	"ACTIVATION":     true,
}

// ValidateNotifyRequest validates the required fields of a NotifyRequest.
// Returns a gRPC InvalidArgument error listing all violations if any field is invalid.
func ValidateNotifyRequest(req *notifypbv2.NotifyRequest) error {
	var violations []string

	if strings.TrimSpace(req.Email) == "" {
		violations = append(violations, "email is required")
	}
	if strings.TrimSpace(req.Type) == "" {
		violations = append(violations, "type is required")
	} else if !validNotifyTypes[req.Type] {
		violations = append(violations, "type must be one of: OTP, RESET_PASSWORD, ACTIVATION")
	}
	if strings.TrimSpace(req.Data) == "" {
		violations = append(violations, "data is required")
	}

	if len(violations) > 0 {
		return status.Errorf(codes.InvalidArgument, "invalid request: %s", strings.Join(violations, "; "))
	}

	return nil
}
