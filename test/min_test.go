package test

import (
	"AlienegraGeek/go-pioneer/min"
	"testing"
)

func TestUpload(t *testing.T) {
	client := min.Init()
	min.Upload(client, "my-bucket")
}

func TestDownload(t *testing.T) {
	client := min.Init()
	min.Download(client, "/Users/yuvan/Documents/github/go-pioneer/file/data1.json")
}
