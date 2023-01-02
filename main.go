/*
   Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
   This file is licensed under the Apache License, Version 2.0 (the "License").
   You may not use this file except in compliance with the License. A copy of
   the License is located at
    http://aws.amazon.com/apache2.0/
   This file is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
   CONDITIONS OF ANY KIND, either express or implied. See the License for the
   specific language governing permissions and limitations under the License.
*/

// snippet-start:[s3.go.upload_object]
package main

// snippet-start:[s3.go.upload_object.imports]
import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// snippet-end:[s3.go.upload_object.imports]

// PutFile uploads a file to a bucket
// Inputs:
//     sess is the current session, which provides configuration for the SDK's service clients
//     bucket is the name of the bucket
//     filename is the name of the file
// Output:
//     If success, nil
//     Otherwise, an error from the call to Open or Upload
func PutFile(sess *session.Session, bucket *string, filepath *string, filename *string) error {
	// snippet-start:[s3.go.upload_object.open]
	file, err := os.Open(*filepath)
	// snippet-end:[s3.go.upload_object.open]
	if err != nil {
		fmt.Println("Unable to open file " + *filename)
		return err
	}

	defer file.Close()

	// snippet-start:[s3.go.upload_object.call]
	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: bucket,
		Key:    filename,
		Body:   file,
	})
	// snippet-end:[s3.go.upload_object.call]
	if err != nil {
		return err
	}

	return nil
}

func main() {
	is_disabled := strings.ToLower(os.Getenv("AWS_DISABLED"))
	if is_disabled == "1" || is_disabled == "true" {
		fmt.Println("AWS things are is disabled. good bye!")
		os.Exit(0)
	}

	// Endpoint setup
	endpoint := strings.ToLower(os.Getenv("AWS_ENDPOINT_URL"))
	endpointResolver := func(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
		if endpoint != "" {
			return endpoints.ResolvedEndpoint{
				URL: endpoint,
			}, nil
		}
	
		return endpoints.DefaultResolver().EndpointFor(service, region, optFns...)
	}

	// snippet-start:[s3.go.upload_object.args]
	bucket := flag.String("b", "", "The bucket to upload the file to")
	fpath := flag.String("f", "", "The file to upload")
	fname := flag.String("d", "", "The file name to dest")
	flag.Parse()

	if *bucket == "" || *fpath == "" {
		fmt.Println("You must supply a bucket name (-b BUCKET) and file name (-f FILE)")
		return
	}
	// snippet-end:[s3.go.upload_object.args]

	if *fname == "" {
		*fname = filepath.Base(*fpath)
	}

	// snippet-start:[s3.go.upload_object.session]
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			EndpointResolver: endpoints.ResolverFunc(endpointResolver),
		},
	}))
	// snippet-end:[s3.go.upload_object.session]

	f err := PutFile(sess, bucket, fpath, fname); err != nil {
		fmt.Println("Got error uploading file:")
		fmt.Println(err)
		return
	}
}
