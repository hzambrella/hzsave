package redp

import (
	"golang.org/x/net/context"

	"fmt"
	"testing"
)

var (
	ctx  = context.Background()
	redp = New()
)

func TestCreateAct(t *testing.T) {
	result, err := redp.CreateAct(ctx, &Activity{
		Title:   "红包",
		HuafeiR: 100,
		Jifen:   100,
		Mb:      100,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}

func TestActById(t *testing.T) {
	result, err := redp.ActById(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}

func TestCreateRec(t *testing.T) {
	result, err := redp.CreateRec(ctx, &Record{
		RetCount: 666,
		Huafei:   666,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}

func TestRecById(t *testing.T) {
	result, err := redp.RecById(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}

func TestRecByUser(t *testing.T) {
	result, err := redp.RecByUser(ctx, 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}

func TestMarkDrawed(t *testing.T) {
	rec, err := redp.RecById(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.MarkDrawed(ctx, rec); err != nil {
		t.Fatal(err)
	}
	fmt.Println(rec.DrawStatus)
}

func TestDrawNum(t *testing.T) {
	num, err := redp.DrawNum(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(num)
}

func TestGetDrawCount(t *testing.T) {
	num, err := redp.GetDrawCount(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(num)
}

func TestUpdateDrawCount(t *testing.T) {
	num, err := redp.GetDrawCount(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(num)

	err = redp.UpdateDrawCount(ctx, 1, num+3)
	if err != nil {
		t.Fatal(err)
	}
	num, err = redp.GetDrawCount(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(num)
}
