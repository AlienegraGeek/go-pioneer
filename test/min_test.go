package test

import (
	"go-pioneer/min"
	"testing"
)

func TestUpload(t *testing.T) {
	client := min.Init()
	min.Upload(client, "my-bucket")
}

func TestDownload(t *testing.T) {
	client := min.Init()
	min.Download(client, "/Users/yuvan/Documents/github/go-pioneer/file/data2.json")
}

func TestCreateFolder(t *testing.T) {
	client := min.Init()
	err := min.CreateFolder(client, "my-bucket", "test")
	if err != nil {
		t.Fatalf("CreateFolder error: %v", err)
	}
}

func TestQuota(t *testing.T) {
	min.BucketQuotaApi("my-bucket")
}
