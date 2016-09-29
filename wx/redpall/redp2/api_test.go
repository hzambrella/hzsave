package redp

import (
	"fmt"
	"testing"
)

func TestGetDrawCountAPI(t *testing.T) {
	num, err := GetDrawCount(ctx, "1")
	if err != nil {
		fmt.Println("api_err")
	}
	fmt.Println(num)
}

func TestUpdateDrawCountAPI(t *testing.T) {
	var add int64 = 1
	err := UpdateDrawCount(ctx, "1", add)
	if err != nil {
		fmt.Println("api_err")
	}
	num, err := GetDrawCount(ctx, "1")
	if err != nil {
		fmt.Println("api_err")
	}
	fmt.Println(num)
}

func TestFriendListOfRecord(t *testing.T) {
	fmt.Println("---1----")
	result, err := friendListOfRecord(ctx, "e0bbe405-5d14-475e-8767-d3f823d7363d", "1")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func TestIncrRecLiuLiang(t *testing.T) {
	fmt.Println("---1----")
	err := incrRecLiuLiang(ctx,1000 ,"e0bbe405-5d14-475e-8767-d3f823d7363d")
	if err != nil {
		fmt.Println(err)
	}
}

func TestUpdateParentId(t *testing.T) {
	fmt.Println("---1----")
	err := UpdateParentId(ctx,"b491092b-2f21-4dc4-a1f1-2f0a7a77f7b0","bf1d69b4-749f-4031-89a1-60f4722996fa","1")
	if err != nil {
		fmt.Println(err)
	}
}
