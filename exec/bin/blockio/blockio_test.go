package main

import (
	"context"
	"io/ioutil"
	"testing"
)

func Test_getMajMin(t *testing.T) {
	ctx:=context.Background()
	majMin, err := getMajMin(ctx)
	if err != nil {
		t.Error(err.Error())
	}else {
		t.Log(majMin)
	}
}

func Test_updateBlkio(t *testing.T) {
	err := updateBlkio("read", "1024000000", context.TODO())
	if err != nil {
		t.Error("updateBlkio read error")
		t.Error(err)
	}
	file, err := ioutil.ReadFile(ReadFileName)
	if err != nil {
		t.Error(err)
	}
	t.Log("print"+ReadFileName)
	t.Log(string(file))
	err = updateBlkio("write", "1024000001", context.TODO())
	if err != nil {
		t.Error("updateBlkio write error")
		t.Error(err)
	}
	file, err = ioutil.ReadFile(WriteFileName)
	if err != nil {
		t.Error(err)
	}
	t.Log("print"+WriteFileName)
	t.Log(string(file))
	err = updateBlkio("read", "0", context.TODO())
	if err != nil {
		t.Error("updateBlkio read error")
		t.Error(err)
	}
	file, err = ioutil.ReadFile(ReadFileName)
	if err != nil {
		t.Error(err)
	}
	t.Log("print"+ReadFileName)
	t.Log(string(file))
	err = updateBlkio("write", "0", context.TODO())
	if err != nil {
		t.Error("updateBlkio write error")
		t.Error(err)
	}
	file, err = ioutil.ReadFile(WriteFileName)
	if err != nil {
		t.Error(err)
	}
	t.Log("print"+WriteFileName)
	t.Log(string(file))
	err = updateBlkio("undefine", "1024000000", context.TODO())
	if err !=UnsupportErr {
		t.Error("updateBlkio undefine error")
		t.Error(err)
	}
}
