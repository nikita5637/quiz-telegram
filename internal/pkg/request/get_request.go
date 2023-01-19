package request

import (
	"context"
	"fmt"
)

// GetRequest ...
func (f *Facade) GetRequest(ctx context.Context, uuid string) ([]byte, error) {
	v, ok := f.cache[uuid]
	if !ok {
		request, err := f.requestStorage.GetRequestByUUID(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("get request error: %w", err)
		}

		f.cache[uuid] = request
		return request.Body, nil
	}

	return v.Body, nil
}
