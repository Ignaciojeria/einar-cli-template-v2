package storjbucket

import (
	"archetype/app/shared/infrastructure/storj"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"storj.io/uplink"
	"storj.io/uplink/edge"
)

type StorJBucket struct {
	fileExpiration  time.Duration
	sharedLinkCreds *edge.Credentials
	bucketName      string
	upLink          *storj.Uplink
}

func init() {
	ioc.Registry(
		NewStorJBucket,
		storj.NewUplink)
}

func NewStorJBucket(ul *storj.Uplink) (storj.UplinkManager, error) {
	sharedLinkExpiration := 10 * time.Minute
	fileExpiration := 7 * 24 * time.Hour
	bucketName := "insert-your-bucket-name"
	bucketFolderName := ""
	sharedLinkRestrictedAccess, err := ul.Access.Share(
		uplink.Permission{
			// only allow downloads
			AllowDownload: true,
			// this allows to automatically cleanup the access grants
			NotAfter: time.Now().Add(sharedLinkExpiration),
		}, uplink.SharePrefix{
			Bucket: bucketName,
			Prefix: bucketFolderName,
		},
	)
	if err != nil {
		return StorJBucket{}, fmt.Errorf("could not restrict access grant: %w", err)
	}

	// RegisterAccess registers the credentials to the linksharing and s3 sites.
	// This makes the data publicly accessible, see the security implications in https://docs.storj.io/dcs/concepts/access/access-management-at-the-edge.
	ctx := context.Background()
	credentials, err := ul.Config.RegisterAccess(ctx,
		sharedLinkRestrictedAccess,
		&edge.RegisterAccessOptions{Public: true})
	if err != nil {
		return StorJBucket{}, fmt.Errorf("could not register access: %w", err)
	}
	// Create Bucket
	_, err = ul.Project.CreateBucket(ctx, bucketName)
	if err != nil && !errors.Is(err, uplink.ErrBucketAlreadyExists) {
		return StorJBucket{}, fmt.Errorf("error creating bucket: %w", err)
	}

	// Ensure the desired Bucket within the Project is created.
	_, err = ul.Project.EnsureBucket(ctx, bucketName)
	if err != nil {
		return StorJBucket{}, fmt.Errorf("could not ensure bucket: %v", err)
	}

	bucket := StorJBucket{
		fileExpiration:  fileExpiration,
		sharedLinkCreds: credentials,
		bucketName:      bucketName,
		upLink:          ul,
	}
	return bucket, nil
}

func (b StorJBucket) CreatePublicSharedLink(ctx context.Context, objectKey string) (string, error) {
	// Create a public link that is served by linksharing service.
	url, err := edge.JoinShareURL("https://link.storjshare.io",
		b.sharedLinkCreds.AccessKeyID,
		b.bucketName, objectKey, nil)
	if err != nil {
		return "", fmt.Errorf("could not create a shared link: %w", err)
	}
	return url, nil
}

func (b StorJBucket) Upload(ctx context.Context, objectKey string, dataToUpload []byte) error {
	// Intitiate the upload of our Object to the specified bucket and key.
	upload, err := b.upLink.Project.UploadObject(ctx, b.bucketName, objectKey, &uplink.UploadOptions{
		// It's possible to set an expiration date for data.
		Expires: time.Now().Add(b.fileExpiration),
	})
	if err != nil {
		return fmt.Errorf("could not initiate upload: %v", err)
	}
	// Copy the data to the upload.
	buf := bytes.NewBuffer(dataToUpload)
	_, err = io.Copy(upload, buf)
	if err != nil {
		_ = upload.Abort()
		return fmt.Errorf("could not upload data: %v", err)
	}

	// Commit the uploaded object.
	err = upload.Commit()
	if err != nil {
		return fmt.Errorf("could not commit uploaded object: %v", err)
	}
	return nil
}

func (b StorJBucket) ListFiles(ctx context.Context) ([]string, error) {
	options := &uplink.ListObjectsOptions{
		Recursive: true,
	}
	objectList := b.upLink.Project.ListObjects(ctx, b.bucketName, options)

	var files []string
	for objectList.Next() {
		item := objectList.Item()
		files = append(files, item.Key)
	}

	if err := objectList.Err(); err != nil {
		return nil, fmt.Errorf("error listing files: %w", err)
	}

	return files, nil
}

func (b StorJBucket) Download(ctx context.Context, objectKey string) ([]byte, error) {

	// Start the download of the specified object from the bucket.
	download, err := b.upLink.Project.DownloadObject(ctx, b.bucketName, objectKey, nil)
	if err != nil {
		return nil, fmt.Errorf("could not initiate download: %w", err)
	}
	defer download.Close()

	// Read the data from the downloaded object.
	var data bytes.Buffer
	_, err = io.Copy(&data, download)
	if err != nil {
		return nil, fmt.Errorf("could not read data from download object: %w", err)
	}

	return data.Bytes(), nil
}
