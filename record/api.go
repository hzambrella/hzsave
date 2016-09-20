package record

import (
	"golang.org/x/net/context"
)

//go:generate ndp_tools_gen -pkg .

// @trace
func addRecord(ctx context.Context, caller, callee string) error {
	return db.AddRecord(ctx, caller, callee)
}

// @trace
func getRecordList(ctx context.Context, caller string) ([]RecordList, error) {
	return db.GetRecordList(ctx, caller)
}

// @trace
func deleteRecord(ctx context.Context, id int) error {
	return db.DeleteRecord(ctx, id)
}
