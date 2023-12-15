package common

import (
	"bytes"
	"context"
	"fmt"
	gomime "github.com/cubewise-code/go-mime"
	"github.com/dustin/go-humanize"
	"github.com/minio/minio-go/v7"
	"github.com/pkg/errors"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Client struct {
	Primary *MinioClient
	logger  *zap.SugaredLogger
}

type MinioClient struct {
	Client       *minio.Client
	filesBucket  string
	s3PublicHost string
}

type UploadFile struct {
	ContentType string
	Path        string
	File        *[]byte // don't used now
	FileStream  io.Reader
	Size        int64
}

func (u *UploadFile) FromMultipart(file multipart.File, header *multipart.FileHeader) {
	extFile := filepath.Ext(header.Filename)
	u.ContentType = gomime.TypeByExtension(extFile)
	u.Path = header.Filename
	u.Size = header.Size
	u.FileStream = file
}

type S3Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	PublicHost      string
	UseSSL          bool
}

func NewS3Client(logger *zap.SugaredLogger) (*S3Client, error) {
	primary := S3Config{
		Endpoint:        os.Getenv("S3_HOST"),
		AccessKeyID:     os.Getenv("S3_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("S3_SECRET_ACCESS_KEY"),
		PublicHost:      os.Getenv("S3_PUBLIC_HOST"),
		UseSSL:          false,
	}

	client, err := minio.New(primary.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(primary.AccessKeyID, primary.SecretAccessKey, ""),
		Secure: primary.UseSSL,
	})

	if err != nil {
		return nil, err
	}

	return &S3Client{
		Primary: &MinioClient{
			Client:       client,
			s3PublicHost: os.Getenv("S3_PUBLIC_HOST"),
			filesBucket:  "files",
		},
		logger: logger,
	}, nil
}

func (m *MinioClient) Upload(ctx context.Context, bucket string, filePath string, file *UploadFile, stream io.Reader, meta map[string]string) error {
	var reader io.Reader

	if file.FileStream != nil {
		reader = stream
	} else if file.File != nil {
		reader = bytes.NewReader(*file.File)
	} else {
		return fmt.Errorf("отсутствует контент файла %s", filePath)
	}

	opts := minio.PutObjectOptions{ContentType: file.ContentType, UserMetadata: meta}
	if file.Size == -1 {
		opts.PartSize = 5 * humanize.MiByte
	}

	infoFiles, err := m.Client.PutObject(
		ctx,
		bucket,
		filePath,
		reader,
		file.Size,
		opts,
	)

	if err != nil {
		return err
	}

	log.Printf("Successfully uploaded file: %s of size %d\n", filePath, infoFiles.Size)
	return nil
}

func (s *S3Client) UploadToFiles(ctx context.Context, filePath string, file *UploadFile) error {
	userMetaData := map[string]string{"x-amz-acl": "public-read"}

	reader := Reader{
		Reader: nil,
		s:      make([]byte, 0),
		mux:    sync.Mutex{},
	}

	if file.FileStream != nil {
		reader.Reader = file.FileStream
	}

	if err := s.Primary.Upload(ctx, s.Primary.filesBucket, filePath, file, reader.View(), userMetaData); err != nil {
		return errors.Wrapf(err, "failed to upload s3 file: %s", file.Path)
	}

	return nil
}

func (s *S3Client) RemoveFile(ctx context.Context, filePath string) error {
	g, waitCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return s.Primary.Client.RemoveObject(waitCtx, s.Primary.filesBucket, filePath, minio.RemoveObjectOptions{
			ForceDelete: true,
		})
	})

	return g.Wait()
}

func (s *S3Client) GetPublicURL(bucket string, filePath string) string {
	bucketPath := bucket

	return fmt.Sprintf("http://%s/%s/%s", s.Primary.s3PublicHost, bucketPath, filePath)
}
