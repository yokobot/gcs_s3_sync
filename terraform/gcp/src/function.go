package gcs_s3_sync

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/functions/metadata"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// GCSEvent is the payload of a GCS event.
type GCSEvent struct {
	Kind                    string                 `json:"kind"`
	ID                      string                 `json:"id"`
	SelfLink                string                 `json:"selfLink"`
	Name                    string                 `json:"name"`
	Bucket                  string                 `json:"bucket"`
	Generation              string                 `json:"generation"`
	Metageneration          string                 `json:"metageneration"`
	ContentType             string                 `json:"contentType"`
	TimeCreated             time.Time              `json:"timeCreated"`
	Updated                 time.Time              `json:"updated"`
	TemporaryHold           bool                   `json:"temporaryHold"`
	EventBasedHold          bool                   `json:"eventBasedHold"`
	RetentionExpirationTime time.Time              `json:"retentionExpirationTime"`
	StorageClass            string                 `json:"storageClass"`
	TimeStorageClassUpdated time.Time              `json:"timeStorageClassUpdated"`
	Size                    string                 `json:"size"`
	MD5Hash                 string                 `json:"md5Hash"`
	MediaLink               string                 `json:"mediaLink"`
	ContentEncoding         string                 `json:"contentEncoding"`
	ContentDisposition      string                 `json:"contentDisposition"`
	CacheControl            string                 `json:"cacheControl"`
	Metadata                map[string]interface{} `json:"metadata"`
	CRC32C                  string                 `json:"crc32c"`
	ComponentCount          int                    `json:"componentCount"`
	Etag                    string                 `json:"etag"`
	CustomerEncryption      struct {
		EncryptionAlgorithm string `json:"encryptionAlgorithm"`
		KeySha256           string `json:"keySha256"`
	}
	KMSKeyName    string `json:"kmsKeyName"`
	ResourceState string `json:"resourceState"`
}

// HelloGCS consumes a GCS event.
func HelloGCS(ctx context.Context, e GCSEvent) error {
	meta, err := metadata.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("metadata.FromContext: %v", err)
	}
	log.Printf("Event ID: %v\n", meta.EventID)
	log.Printf("Event type: %v\n", meta.EventType)
	log.Printf("Bucket: %v\n", e.Bucket)
	log.Printf("File: %v\n", e.Name)
	log.Printf("Metageneration: %v\n", e.Metageneration)
	log.Printf("Created: %v\n", e.TimeCreated)
	log.Printf("Updated: %v\n", e.Updated)

	switch {
	case (meta.EventType == "google.storage.object.finalize"):
		Finalized(ctx, e)
	case (meta.EventType != "google.storage.object.delete"):
		Delete(ctx, e)
	}

	return nil
}

// get aws credentail from gcp secrets
func GetSecret(s string) (string, error) {
    log.Printf("GetSecret start.")
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
        log.Printf("failed to create secretmanager client: %v", err)
        return fmt.Errorf("GetSecret error: %v", err)
	}

    projectId := "yokobot-dev"
    secret := s

    name := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectId, secret)

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
        log.Printf("failed to access secret verion: %v", err)
        return fmt.Errorf("GetSecret error: %v", err)
	}
	value := string(result.Payload.Data)
    log.Printf("GetSecret end.")
	return value, nil
}

// Return S3 client
func S3Client() *s3.S3 {
    log.Printf("S3Client start.")
	aws_access_key_id, _ := GetSecret("aws_access_key_id")
	aws_secret_access_key, _ := GetSecret("aws_secret_access_key")

	// s3 client作る
	sess := session.Must(session.NewSession(&aws.Config{
        Region: aws.String("ap-northeast-1"),
        Credentials: credentials.NewStaticCredentials(aws_access_key_id, aws_secret_access_key, ""),
    }))
	svc := s3.New(sess)

    log.Printf("S3Client end.")
	return svc
}

// finalized event
func Finalized(ctx context.Context, e GCSEvent) error {
    log.Printf("Finalized start.")
	meta, err := metadata.FromContext(ctx)
	if err != nil {
        log.Printf("Finalized get metadata failed.")
		return fmt.Errorf("Finalized error: %v", err)
	}

	// 同名ファイルがs3に存在しているか確認して、存在していれば何もしない、存在しなければファイルをs3にコピーする
	svc := S3Client()

	input := &s3.ListObjectsInput{
		Bucket: aws.String(e.Bucket),
		Prefix: aws.String(e.Name),
	}

	resp, err := svc.ListObjects(input)

	if err != nil {
		log.Printf("s3 ListObjects error: %v", err)
        return fmt.Errorf("Finalized error: %v", err)
	}

    //gcsのファイルをtmpにコピーしてs3にpushする TODO

	for _, item := range resp.Contents {
		object_name := *item.Key
		if object_name != e.Name {
			input := &s3.PutObjectInput{
				Bucket: aws.String(e.Bucket),
				Key:    aws.String(e.Name),
			}
			_, err := svc.PutObject(input)
			if err != nil {
				log.Printf("s3 PutObject error: %v", err)
                return fmt.Errorf("Finalized error: %v", err)
			}
		}
        log.Printf("s3 PutObject: Key is exist.")
	}

    log.Printf("Finalize end.")
	return nil
}

// delete event
func Delete(ctx context.Context, e GCSEvent) error {
	meta, err := metadata.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("metadata.FromContext: %v", err)
	}

	if meta.EventType != "google.storage.object.delete" {
		return fmt.Errorf("Event is not delete. %v")
	}

	//s3に同名ファイルが存在しているか確認して、存在していればファイルを削除する
	svc := S3Client()

	// s3に同名ファイルがあるかを調べる
	input := &s3.ListObjectsInput{
		Bucket: aws.String(e.Bucket),
		Prefix: aws.String(e.Name),
	}
	resp, err := svc.ListObjects(input)
	if err != nil {
		return fmt.Errorf("s3 ListObjects error: %v", err)
	}

	// 同名のファイルが存在していればs3にファイルを削除する
	for _, item := range resp.Contents {
		object_name := *item.Key
		if object_name == e.Name {
			input := &s3.DeleteObjectInput{
				Bucket: aws.String(e.Bucket),
				Key:    aws.String(e.Name),
			}
			_, err := svc.DeleteObject(input)
			if err != nil {
				return fmt.Errorf("s3 ListObjects error: %v", err)
			}
		}
	}

	return nil
}
