package resolver

import (
	"context"

	gmodel "github.com/DaisukeMatsumoto0925/backend/graph/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (r *mutationResolver) UploadFile(ctx context.Context, input gmodel.UploadFileInput) (*gmodel.UploadFilePayload, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)

	res, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String("bucketName"),
		Key:         aws.String("uploadPath/" + input.File.Filename),
		Body:        input.File.File,
		ContentType: aws.String(input.File.ContentType),
	})
	if err != nil {
		return nil, err
	}

	return &gmodel.UploadFilePayload{
		UploadedPath: res.Location,
	}, nil
}
