package cache

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"

	"bytes"
	"io/ioutil"
	"net/url"
	"time"
)

type s3objects struct {
	bucket, prefix string
	s3             *s3.S3
}

func getS3ForRegion(region string) *s3.S3 {
	var sess *session.Session
	var err error

	if region == "" {
		sess, err = session.NewSession()
	} else {
		sess, err = session.NewSession(&aws.Config{Region: aws.String(region)})
	}
	if err != nil {
		log.Fatal("Could not initialize session", region, err)
		return nil
	}

	return s3.New(sess)
}

func NewS3ObjectCache(uri, region string) ICache {
	parsed, err := url.Parse(uri)
	if err != nil {
		log.Fatal("Invalid uri to NewS3ObjectCache", uri, err)
		return nil
	}
	instance := getS3ForRegion(region)
	return &s3objects{bucket: parsed.Host, prefix: parsed.Path, s3: instance}
}

func (cc *s3objects) Add(key string, value interface{}, expiration *time.Time) error {
	if value == nil {
		// delete the key then
		_, err := cc.s3.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(cc.bucket),
			Key:    aws.String(cc.prefix + key),
		})
		return err
	}

	bb, ok := value.([]byte)
	if !ok {
		log.Fatal("s3Objects allows values of type []byte only")
		return nil
	}

	var metadata map[string]*string
	if expiration != nil {
		formatted := expiration.Format(time.UnixDate)
		metadata = map[string]*string{"expires": &formatted}
	}

	_, err := cc.s3.PutObject(&s3.PutObjectInput{
		Bucket:   aws.String(cc.bucket),
		Key:      aws.String(cc.prefix + key),
		Body:     bytes.NewReader(bb),
		Metadata: metadata,
	})
	return err
}

func (cc *s3objects) Get(key string) (result interface{}, expires *time.Time, err error) {
	resp, err := cc.s3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(cc.bucket),
		Key:    aws.String(cc.prefix + key),
	})

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "NoSuchKey" {
				// got nothing. silently return nothing
				return nil, nil, nil
			}
		}
		log.Println("Failed to fetch", cc.prefix+key, "from s3", cc.bucket, " ", err)
		return
	}
	defer resp.Body.Close()
	result, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, nil, err
	}

	if resp.Metadata != nil {
		if expiration, ok := resp.Metadata["expires"]; ok && expiration != nil {
			if expiresResult, err := time.Parse(time.UnixDate, *expiration); err == nil {
				expires = &expiresResult
			}
		}
	}

	if expires != nil && expires.Before(time.Now()) {
		// object is expired.  delete it and refetch
		result = nil
		expires = nil
		err = cc.Add(key, nil, nil)
	}
	return
}
