package family

import (
	"fmt"

	"testing"

	"golang.org/x/net/context"
)

var (
	ctx = context.Background()
	c   = New()
)

func TestGetCallee(t *testing.T) {
	result, err := c.GetCallee(ctx, 1000005, "testing")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}

func TestGetFamilyNum(t *testing.T) {
	callee, max, err := GetFamilyNum(ctx, 1000005, "testing")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(callee, max)
}

/*
func TestDeleteCallee(t *testing.T) {
	fmt.Println("---------------1----------------")
	err := c.deleteCallee(ctx, "1000005", "135-1054-4459")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("delete complete")
}
*/
func TestDeleteFamilynum(t *testing.T) {
	fmt.Println("---------------1----------------")
	err := deleteFamilyNum(ctx, 1000005, 100, "testing", "740740740")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("delete complete")
}

func TestGetBindNumOfUser(t *testing.T) {
	fmt.Println("---------------1----------------")
	result, err := c.getBindNumOfUser(ctx, 1000005, "testing")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("bindnum:", result)
}

/*
func TestGetMaxBignum(t *testing.T) {
	fmt.Println("---------------1----------------")
	result, err := c.getMaxBignum(ctx)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("maxbignum:", result)
}

/*
func TestFamilyMaxNum(t *testing.T) {
	fmt.Println("---------------1----------------")
	result, err := FamilyMaxNum(ctx)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("maxbignum:", result)
}
*/
func TestGetCurrentBigNum(t *testing.T) {
	fmt.Println("---------------1----------------")
	result, err := c.getCurrentBigNum(ctx, "18310849050")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("bindbignumid:", result)
}

func TestGetNumId(t *testing.T) {
	fmt.Println("---------------1----------------")
	result, err := c.getNumId(ctx, "testing", 1000144, "15392474846")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("result:", result)

}

/*
func TestAlloctionNum(t *testing.T) {
	fmt.Println("---------------1----------------")
	err := c.alloctionNum(ctx, "1000144", "110110119")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("set complete")
}

func TestSetFamilynum(t *testing.T) {
	fmt.Println("---------------1----------------")
	err := SetFamilynum(ctx, "1000144", "7407407474740")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("set complete")
}
*/
func TestUpdateBigNumId(t *testing.T) {
	fmt.Println("---------------1----------------")
	err := updateBigNumId(ctx, 1000005, 10000000, "testing", "15670259258", "740740740", "075532990932")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("set complete")
}

func TestCheckBigNumId(t *testing.T) {
	fmt.Println("---------------1----------------")
	list, err := dc.checkBigNumId(ctx, "testing", 1000005, 1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("result:", list)
}
