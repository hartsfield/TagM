package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

var fexts map[string]error = map[string]error{
	"image/png":     errors.New("png"),
	"image/svg+xml": errors.New("svg"),
	"image/jpeg":    errors.New("jpg"),
	"image/gif":     errors.New("gif"),
	"video/mp4":     errors.New("mp4"),
	"video/webm":    errors.New("webm"),
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	post, err := parseForm(r)
	if err != nil {
		log.Println(status(w, "Invalid Form", err))
	}
	b, err := json.Marshal(post)
	if err != nil {
		log.Println(status(w, "Invalid Form", err))
	}
	if err = zhPost(post); err == nil {
		ajaxResponse(w, map[string]string{
			"status":     "success",
			"replyID":    post.ID,
			"itemString": string(b),
		})
		cache()
		return
	}
	log.Println(status(w, "Database Error", err))
}
func parseForm(r *http.Request) (*post, error) {
	mr, err := r.MultipartReader()
	if err != nil {
		return nil, err
	}
	var c_ *credentials = r.Context().Value(ctxkey).(*credentials)
	var post *post = &post{
		ID:         genID(15),
		TS:         time.Now(),
		TimeString: time.Now().Format(time.RFC822),
		Author:     c_.User.ID,
	}
	c_.User.Posts = append(c_.User.Posts, post.ID)
	for {
		part, err_part := mr.NextPart()
		if err_part == io.EOF {
			break
		}
		if part.FormName() == "Media" {
			post.Type = "Media"
			if handleFile(part, post) != nil {
				return nil, err
			}
		}
		if part.FormName() == "ProfilePic" {
			post.Type = "ProfilePic"
			if handleFile(part, post) != nil {
				return nil, err
			}
		}
		if part.FormName() == "ProfileBG" {
			post.Type = "ProfileBG"
			if handleFile(part, post) != nil {
				return nil, err
			}
		}
		if part.FormName() == "uptext" {
			txt, err := readPart(part)
			if err != nil {
				return nil, err
			}
			post.Text = template.HTML(txt)
		}
		if part.FormName() == "work" {
			if c_.User.Work, err = readPart(part); err != nil {
				return nil, err
			}
		}
		if part.FormName() == "location" {
			if c_.User.Location, err = readPart(part); err != nil {
				return nil, err
			}
		}
		if part.FormName() == "about" {
			if c_.User.About, err = readPart(part); err != nil {
				return nil, err
			}
		}
	}
	return post, nil
}

func readPart(part *multipart.Part) (string, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(part); err != nil {
		return "", err
	}
	return buf.String(), nil
}
func handleFile(part *multipart.Part, data *post) error {
	fileBytes, err := io.ReadAll(io.LimitReader(part, 10<<20))
	if err != nil {
		return err
	}
	ex := fexts[http.DetectContentType(fileBytes)].Error()
	tempFile, err := os.CreateTemp("public/temp", "u-*."+fmt.Sprint(ex))
	if err == nil {
		defer tempFile.Close()
		if _, err = tempFile.Write(fileBytes); err == nil {
			data.TempFileName = tempFile.Name()
			switch ex {
			case "png", "jpg", "gif", "svg":
				data.MediaType = template.HTML("<img class='post-img post-media' src='/" + data.TempFileName + "' />")
			case "mp4", "webm":
				data.MediaType = template.HTML("<video class='post-video post-media' src='/" + data.TempFileName + "' />")
			}
			return nil
		}
		return err
	}
	return err
}
