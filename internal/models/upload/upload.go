package upload

import (
	"fmt"
	"log/slog"
	"mime/multipart"
)

type Upload struct {
	Title []string
	Cover []*multipart.FileHeader
	Album []*multipart.FileHeader
}

func New(title []string, cover []*multipart.FileHeader, album []*multipart.FileHeader) (*Upload, error) {
	slog.Info("Files accepted", "title", title, "cover", len(cover), "album", len(album))
	//title
	if len(title[0]) > 30 || len(title[0]) == 0 {
		return nil, fmt.Errorf("title length is longer than 30 characters or title missing")
	}
	//cover
	if len(cover) > 1 || len(cover) == 0 {
		return nil, fmt.Errorf("more than 1 cover or cover missing")
	}
	if cover[0].Header["Content-Type"][0] != "image/jpeg" {
		return nil, fmt.Errorf("file %s is not a jpeg format", cover[0].Filename)
	}
	//album image
	if len(album) == 0 {
		return nil, fmt.Errorf("album missing")
	}
	for i, _ := range album {
		if album[i].Header["Content-Type"][0] != "image/jpeg" {
			return nil, fmt.Errorf("file %s is not a jpeg format", album[i].Filename)
		}
	}
	slog.Info("Validate album successfully")
	return &Upload{
		Title: title,
		Cover: cover,
		Album: album}, nil
}
