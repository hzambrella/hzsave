package record

import (
	"fmt"
	"testing"

	"golang.org/x/net/context"
)

var (
	ctx = context.Background()
	c   = New()
)

func TestAddRecord(t *testing.T) {
	err := c.AddRecord(ctx, "15111111111", "15122222222")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestDeleteRecord(t *testing.T) {
	err := c.DeleteRecord(ctx, 1)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestGetRecordList(t *testing.T) {
	result, err := c.GetRecordList(ctx, "15111111111")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)
}
