// upload_handler.go houses the functions used for user post submissions and
// other multipart/form-data requests. This file may be split up in the future.
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

// uploadHandler() is the entry point for post uploads and step 1 of the upload
// process. We parse the form data sent by the client, marshal it so that we
// may return it to the client, and respond with the appropriate ajaxResponse.
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the form data sent by the client into a post{}
	post, err := parseForm(r)
	if err != nil {
		log.Println(status(w, "Invalid Form", err))
	}

	// Marshal the freshly parsed post{} into its JSON representation in
	// []byte form.
	b, err := json.Marshal(post)
	if err != nil {
		log.Println(status(w, "Invalid Form", err))
	}

	// Add the post to the database sets/maps.
	if err = zhPost(post); err == nil {
		// custom Ajax response returning the new posts ID and JSON
		// representation (if any).
		ajaxResponse(w, map[string]string{
			"status":     "success",
			"replyID":    post.ID,
			"itemString": string(b),
		})
		// We cache here in development, but for production we won't be
		// cache()ing the database after every submission, there is a
		// better way.
		cache()
		return
	}
	log.Println(status(w, "Database Error", err))
}

// parseForm() parses multipart/form-data sent by the client. This is used for
// every form except auth, but I may break it down into smaller functions
// eventually.
func parseForm(r *http.Request) (*post, error) {
	mr, err := r.MultipartReader()
	if err != nil {
		return nil, err
	}

	// obtain the users credentials from the context:
	var c_ *credentials = r.Context().Value(ctxkey).(*credentials)

	// create a post with default values:
	var post *post = &post{
		ID:         genID(15),
		TS:         time.Now(),
		TimeString: time.Now().Format(time.RFC822),
		Author:     c_.User.ID,
	}

	// Add the post ID to sorted set(s):
	_, err = zaddUsersPosts(c_, post)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Read the multipart/form-data, cycling through each form part,
	// checking the part.FormName(), and responding based on the output.
	// A switch didn't seem to work properly here.
	for {
		part, err_part := mr.NextPart()
		if err_part == io.EOF {
			break
		}

		// Post media
		if part.FormName() == "Media" {
			post.Type = "Media"
			if handleFile(part, post) != nil {
				return nil, err
			}
		}

		// Profile Pic
		if part.FormName() == "ProfilePic" {
			post.Type = "ProfilePic"
			if handleFile(part, post) != nil {
				return nil, err
			}
		}

		// Profile Background
		if part.FormName() == "ProfileBG" {
			post.Type = "ProfileBG"
			if handleFile(part, post) != nil {
				return nil, err
			}
		}

		// Post main text
		if part.FormName() == "uptext" {
			txt, err := readPart(part)
			if err != nil {
				return nil, err
			}
			post.Text = template.HTML(txt)
		}

		// User profile "work" input
		if part.FormName() == "work" {
			if c_.User.Work, err = readPart(part); err != nil {
				return nil, err
			}
		}

		// User profile "location" input
		if part.FormName() == "location" {
			if c_.User.Location, err = readPart(part); err != nil {
				return nil, err
			}
		}

		// User profile "about" input
		if part.FormName() == "about" {
			if c_.User.About, err = readPart(part); err != nil {
				return nil, err
			}
		}
	}
	return post, nil
}

// readPart(*multipart.Part) is used in the parseForm() function to reduce
// repeated code. It converts the form part to a string or returns an error if
// it can't.
func readPart(part *multipart.Part) (string, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(part); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// handleFile() is used to handle file uploads.
func handleFile(part *multipart.Part, data *post) error {

	// Limit the file size we accept.
	fileBytes, err := io.ReadAll(io.LimitReader(part, 10<<20))
	if err != nil {
		return err
	}

	// fexts is used to reduce the code base a little, acting as a type of
	// "switch". If a mime type is supported it "errors" with its proper
	// extension.
	var fexts map[string]error = map[string]error{
		"image/png": errors.New("png"),
		// "image/svg+xml": errors.New("svg"),
		"image/jpeg": errors.New("jpg"),
		"image/gif":  errors.New("gif"),
		"video/mp4":  errors.New("mp4"),
		"video/webm": errors.New("webm"),
	}

	ex := fexts[http.DetectContentType(fileBytes)].Error()
	tempFile, err := os.CreateTemp("public/temp", "u-*."+fmt.Sprint(ex))

	if err == nil {
		defer tempFile.Close()
		if _, err = tempFile.Write(fileBytes); err == nil {
			data.TempFileName = tempFile.Name()

			// create the HTML element based on the file type:
			switch ex {
			case "png", "jpg", "gif", "svg":
				data.MediaType = template.HTML(
					"<img class='item-img item-media' src='/" +
						data.TempFileName + "' />")
			case "mp4", "webm":
				data.MediaType = template.HTML(
					"<video class='item-video item-media' src='/" +
						data.TempFileName + "' />")
			}
			return nil
		}
		return err
	}
	return err
}
