package request

import (
	"context"
	"fmt"

	"github.com/go-xorm/builder"
)

// UnregisterRequest ...
func (f *Facade) UnregisterRequest(ctx context.Context, uuid string) error {
	request, err := f.requestStorage.GetRequestByUUID(ctx, uuid)
	if err != nil {
		return fmt.Errorf("unregister request error: %w", err)
	}

	records, err := f.requestStorage.Find(ctx, builder.NewCond().And(
		builder.Eq{
			"group_uuid": request.GroupUUID,
		},
	))
	if err != nil {
		return fmt.Errorf("unregister request error: %w", err)
	}

	for _, record := range records {
		err := f.requestStorage.Delete(ctx, record.ID)
		if err != nil {
			return fmt.Errorf("unregister request error: %w", err)
		}

		delete(f.cache, record.UUID)
	}

	return nil
}
