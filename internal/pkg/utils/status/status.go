package status

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

// GetErrorInfoFromStatus ...
func GetErrorInfoFromStatus(status *status.Status) *errdetails.ErrorInfo {
	for _, detail := range status.Details() {
		switch t := detail.(type) {
		case *errdetails.ErrorInfo:
			return t
		}
	}

	return nil
}
