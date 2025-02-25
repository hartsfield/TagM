// Provided Under BSD (2 Clause)
//
// Copyright 2025 Johnathan A. Hartsfield
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
// this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
// this list of conditions and the following disclaimer in the documentation
// and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS “AS IS”
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.
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
	// NOTE: the output of part.FormName() must be a direct match with the
	// JSON key sent from the client. Check the JSON tags on the struct we
	// marshal from the clients JSON request, check the JSON its self, and
	// check the html form and it's child elements for the "name"
	// attributes, making sure these keys are exact matches.
	for {
		part, err_part := mr.NextPart()
		if err_part == io.EOF {
			break
		}
		///////////////////////////////////////////////////////////////
		/////////////////////////    MEDIA    /////////////////////////
		///////////////////////////////////////////////////////////////
		// Post media
		if part.FormName() == "Media" { // see: upload.html
			post.Type = "Media"
			if handleFile(part, post) != nil {
				return nil, err
			}
		}
		// Profile Pic
		if part.FormName() == "ProfilePic" { // see: profile.html
			post.Type = "ProfilePic"
			if handleFile(part, post) != nil {
				return nil, err
			}
		}
		// Profile Background
		if part.FormName() == "ProfileBG" { // see: profile.html
			post.Type = "ProfileBG"
			if handleFile(part, post) != nil {
				return nil, err
			}
		}
		///////////////////////////////////////////////////////////////
		/////////////////////////    TEXT     /////////////////////////
		///////////////////////////////////////////////////////////////
		// Post main text
		if part.FormName() == "uptext" { // see: upload.html
			txt, err := readPart(part)
			if err != nil {
				return nil, err
			}
			post.Text = template.HTML(txt)
		}
		// User profile "work" input
		if part.FormName() == "work" { // see: profile.html
			if c_.User.Work, err = readPart(part); err != nil {
				return nil, err
			}
		}
		// User profile "location" input
		if part.FormName() == "location" { // see: profile.html
			if c_.User.Location, err = readPart(part); err != nil {
				return nil, err
			}
		}
		// User profile "about" input
		if part.FormName() == "about" { // see: profile.html
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
	// Limit the file size we accept. 10<<20 is a bitwise operation,
	// meaning no one knows what the hell it means. We look it up to find
	// that we are "shifting bits", in this case the bits of the left
	// operand (10). We are shifting them (bitwise), to the left (<<), by
	// the number specified by the right operand (20). This is equivalent
	// to multiplying the left operand by 2 raised to the power of the
	// right operand. Therefore, 10<<20 is equivalent to: 10 * (2 ** 20)
	// which is approximately ~10 megabytes. (note: I'm not re-writing this
	// comment when we change this).
	fileBytes, err := io.ReadAll(io.LimitReader(part, 200<<20)) // 200mb
	if err != nil {
		return err
	}

	// fexts conception was aimed at reducing the code base, acting as a
	// type of "switch". If a mime type is supported it "errors" with its
	// proper extension.
	var fexts map[string]error = map[string]error{
		"image/png":  errors.New("png"),
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
			case "png", "jpg", "gif":
				data.MediaType = markupMedia(
					"img", data.TempFileName)
			case "mp4", "webm":
				data.MediaType = markupMedia(
					"vid", data.TempFileName)
			}
			return nil
		}
		return err
	}
	return err
}

// markupMedia() is used to wrap the media element displayed in a post with
// the appropriate html markup/tag. Images get the "img" tag, while videos
// get the "video" tag, along with attributes "controls", "autoplay", and
// "mute". Since it only supports images and videos, only the parameter "img"
// is recognized as the forst parameter, anything else defaults to the video
// tag. We also pass the path to the media as the second parameter.
func markupMedia(typ, path string) template.HTML {
	if typ == "img" {
		return template.HTML(
			"<img class='item-img item-media'" +
				" src='/" + path + "'/>")
	}
	return template.HTML(
		"<video class='item-video item-media' controls autoplay mute" +
			" src='/" + path + "' />")
}
