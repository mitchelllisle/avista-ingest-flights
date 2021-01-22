package gcp

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/mitchelllisle/avista-ingest-flights/src/utils"
	"google.golang.org/api/iterator"
	"io/ioutil"
	"time"
)

type GCS struct {
	Project string
	Bucket string
}

type File struct {
	File string
	Created time.Time
	ContentType string
}

func InitGCS(project, bucket string) *GCS {
	return &GCS{
		Project: project,
		Bucket: bucket,
	}
}

func (g *GCS) ListFilesInBucket(bucket string, prefix string) []File {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	utils.PanicOnError(err, "storage.NewClient: %v")
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	it := client.Bucket(bucket).Objects(ctx, &storage.Query{Prefix: prefix})
	var files []File
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}

		if attrs != nil {
			file := File{
				File: attrs.Name,
				ContentType: attrs.ContentType,
				Created: attrs.Created,
			}
			files = append(files, file)
		}
	}
	return files
}

func (g *GCS) DownloadFile(object string) []byte {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	utils.PanicOnError(err, "storage.NewClient: %v")
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	rc, err := client.Bucket(g.Bucket).Object(object).NewReader(ctx)
	utils.PanicOnError(err, fmt.Sprintf("Unable to create object %s", object))
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	utils.PanicOnError(err, "Unable to read object")
	return data
}

func (g *GCS) UploadFile(payload []byte, fileName string, contentType string) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	utils.PanicOnError(err, "storage.NewClient: %v")
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	wc := client.Bucket(g.Bucket).Object(fileName).NewWriter(ctx)
	wc.ContentType = contentType

	_, err = wc.Write(payload)
	utils.PanicOnError(err, fmt.Sprintf(
"createFile: unable to write data to bucket %q, file %q: %v", g.Bucket, fileName, err))

	err = wc.Close()
	utils.PanicOnError(err, "unable to upload data")
}