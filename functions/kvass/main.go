package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	apex "github.com/apex/go-apex"
	"github.com/apex/go-apex/proxy"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"github.com/rs/cors"
)

var (
	svc    *s3.S3
	bucket *string
)

func newS3() {
	sess, err := session.NewSession()
	errHandler(err)
	svc = s3.New(sess)

	b := os.Getenv("AWS_S3_BUCKET")
	if b == "" {
		errHandler(errors.New("AWS_S3_BUCKET must be set"))
	}

	bucket = aws.String(b)
}

func getAllowedOrigins() []string {
	origins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if origins == "" {
		return []string{}
	}

	parts := strings.Split(origins, ",")
	for i, v := range parts {
		parts[i] = strings.TrimSpace(v)
	}

	return parts
}

var debugMode bool

func init() {
	debugMode, _ = strconv.ParseBool(os.Getenv("DEBUG_MODE"))
	if debugMode {
		for _, pair := range os.Environ() {
			l.Println(pair)
		}

		l.Printf("AllowedOrigins: %s\n", strings.Join(getAllowedOrigins(), ","))
	}

	newS3()
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			getHandler(w, req)
		case http.MethodPut:
			putHandler(w, req)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprint(w, "Only GET/PUT methods are allowed")
		}
	})

	handler := cors.New(cors.Options{
		AllowedOrigins: getAllowedOrigins(),
	}).Handler(mux)

	apex.Handle(proxy.Serve(handler))
}

func putHandler(w http.ResponseWriter, req *http.Request) {
	key := req.URL.Path

	bs, err := ioutil.ReadAll(req.Body)
	errHandler(err)
	if debugMode {
		l.Printf("key: %s\n", key)
		l.Printf("value: %s\n", string(bs))
	}

	input := &s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(bytes.NewReader(bs)),
		Bucket: bucket,
		Key:    aws.String(key),
	}
	_, err = svc.PutObject(input)
	errHandler(errors.Wrap(err, "svc.PutObject"))

	w.WriteHeader(http.StatusCreated)
}

func getHandler(w http.ResponseWriter, req *http.Request) {
	key := req.URL.Path
	output, err := svc.GetObject(&s3.GetObjectInput{Bucket: bucket, Key: aws.String(key)})
	errHandler(err)

	bs, err := ioutil.ReadAll(output.Body)
	errHandler(err)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(bs))
}

var l = log.New(os.Stderr, "", 0)

func errHandler(err error) {
	if err != nil {
		l.Fatal(err)
	}
}
