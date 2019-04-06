package s3

import (
	"fmt"
	"github.com/apex/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/devster/tarreleaser/pkg/context"
	"github.com/devster/tarreleaser/pkg/pipe"
	"github.com/devster/tarreleaser/pkg/tmpl"
	"os"
	"path/filepath"
)

type Pipe struct{}

func (Pipe) String() string {
	return "s3 publishing"
}

func (Pipe) Default(ctx *context.Context) error {
	if ctx.Config.Publish.S3.ACL == "" {
		ctx.Config.Publish.S3.ACL = "private"
	}

	return nil
}

func (Pipe) Run(ctx *context.Context) error {
	if ctx.SkipPublish {
		return pipe.ErrSkipPublishEnabled
	}

	if ctx.Config.Publish.S3.Bucket == "" {
		log.Debug("skipping publish, no bucket configured")
		return nil
	}

	conf := ctx.Config.Publish.S3
	template := tmpl.New(ctx)
	bucket, err := template.Apply(conf.Bucket)
	if err != nil {
		return err
	}

	folder, err := template.Apply(conf.Folder)
	if err != nil {
		return err
	}

	key := filepath.Join(folder, ctx.Archive.Name)

	log.WithFields(log.Fields{
		"key":    fmt.Sprintf("%s:/%s", bucket, key),
		"region": conf.Region,
		"ACL":    conf.ACL,
	}).Info("uploading...")

	awsConfig := &aws.Config{
		Credentials: credentials.NewChainCredentials([]credentials.Provider{
			&credentials.EnvProvider{},
			&credentials.SharedCredentialsProvider{
				Profile: conf.Profile,
			},
		}),
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil
	}

	svc := s3.New(sess, &aws.Config{
		Region: aws.String(conf.Region),
	})

	f, err := os.Open(ctx.Archive.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   f,
		ACL:    aws.String(conf.ACL),
	})
	return err
}
