package redp

import (
	"fmt"
	"testing"
)

func TestGetDrawCountAPI(t *testing.T) {
	num, err := GetDrawCount(ctx, 1)
	if err != nil {
		fmt.Println("api_err")
	}
	fmt.Println(num)
}

func TestUpdateDrawCountAPI(t *testing.T) {
	var add int64 = 1
	err := UpdateDrawCount(ctx, 1, add)
	if err != nil {
		fmt.Println("api_err")
	}
	num, err := GetDrawCount(ctx, 1)
	if err != nil {
		fmt.Println("api_err")
	}
	fmt.Println(num)
}

func TestRedpDrawUserList(t *testing.T) {
	result,err:=redpDrawUserList(ctx,1 )
	if err != nil {
		fmt.Println("api_err")
	}
	fmt.Println(result)
}

func TestGetContactMobile(t *testing.T) {
	result,err:=getContactMobile(ctx,"15671682392" )
	if err != nil {
		fmt.Println("api_err")
	}
	fmt.Println(result)
}

func TestFriendListOfRecord(t *testing.T) {
	fmt.Println("---1----")
	result,err:=friendListOfRecord(ctx,1000163 ,1 )
	if err != nil {
		fmt.Println("api_err")
	}
	fmt.Println(result)
}
