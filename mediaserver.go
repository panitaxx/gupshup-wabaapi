package wabaapi

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path"
	"time"

	gcs "cloud.google.com/go/storage"
	"github.com/ansel1/merry"
)

//DefaultTimeout is the default timeout for the media server
var DefaultTimeout = 120 * time.Second

type MediaServer interface {
	//GetFile returns a file from the media server and its content type
	GetFile(requri string) (io.ReadCloser, string, error)
	//PutFile uploads a file to the media server and returns the path to the file
	PutFile(r io.Reader, contentType string) (string, error)
	//PutFileWithExt uploads a file to the media server and returns the path to the file
	PutFileWithExt(r io.Reader, ext string) (string, error)
}

type MediaServerMedia struct {
	Server      MediaServer
	Reader      io.ReadCloser
	ContentType string
}

func (media *MediaServerMedia) PutFile() (string, error) {
	url, err := media.Server.PutFile(media.Reader, media.ContentType)
	if err != nil {
		return "", merry.Wrap(err)
	}
	media.Reader.Close()
	return url, nil
}

type GCSMediaServer struct {
	Client     *gcs.Client
	Bucket     string
	PathPrefix string
	URLHost    string
}

func (ms *GCSMediaServer) GetFile(requri string) (io.ReadCloser, error) {
	if ms.Client == nil || ms.Bucket == "" {
		return nil, merry.New("GCSMediaServer not configured")
	}

	filename := path.Join(ms.PathPrefix, requri)
	obj := ms.Client.Bucket(ms.Bucket).Object(filename)
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	_, err := obj.Attrs(ctx)
	if err != nil {
		if errors.Is(err, gcs.ErrObjectNotExist) {
			return nil, merry.New("file not found").WithHTTPCode(http.StatusNotFound)
		}
		return nil, merry.Wrap(err)
	}

	return obj.NewReader(ctx)

}

func (ms *GCSMediaServer) PutFile(r io.Reader, contentType string) (string, error) {
	if ms.Client == nil || ms.Bucket == "" {
		return "", merry.New("GCSMediaServer not configured")
	}

	if contentType == "" {
		return "", merry.New("content type not specified")
	}

	exts, err := mime.ExtensionsByType(contentType)
	if err != nil {
		return "", merry.Wrap(err)
	}

	ext := ""
	if len(exts) > 0 {
		ext = "." + exts[0]
	}

	fn := fmt.Sprintf("%s-%d%s", createSecureRandomString(10), time.Now().Unix(), ext)

	filename := path.Join(ms.PathPrefix, fn)
	obj := ms.Client.Bucket(ms.Bucket).Object(filename)
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	w := obj.NewWriter(ctx)
	defer w.Close()

	w.ContentType = contentType

	_, err = io.Copy(w, r)
	if err != nil {
		return "", merry.Wrap(err)
	}

	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return "", merry.Wrap(err)
	}

	if ms.URLHost == "" {
		return attrs.MediaLink, nil
	}

	return fmt.Sprintf("%s/%s", ms.URLHost, filename), nil
}

func (ms *GCSMediaServer) PutFileWithExt(r io.Reader, ext string) (string, error) {
	if ms.Client == nil || ms.Bucket == "" {
		return "", merry.New("GCSMediaServer not configured")
	}

	if ext == "" {
		return "", merry.New("extension not specified")
	}

	ctype := mime.TypeByExtension(ext)
	if ctype == "" {
		return "", merry.New("content type not found")
	}

	return ms.PutFile(r, ctype)
}

func createSecureRandomString(length int) string {
	var b = make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	var alphabet = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	for i, v := range b {
		b[i] = alphabet[v%byte(len(alphabet))]
	}

	return string(b)
}
