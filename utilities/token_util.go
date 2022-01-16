// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package utilities

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type NameChange struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// @me Discord Patch request to change Username
func (in *Instance) NameChanger(name string) (http.Response, error) {

	url := "https://discord.com/api/v9/users/@me"

	data := NameChange{
		Username: name,
		Password: in.Password,
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return http.Response{}, err
	}

	req, err := http.NewRequest("PATCH", url, strings.NewReader(string(bytes)))

	if err != nil {
		return http.Response{}, err
	}
	cookie, err := in.GetCookieString()
	if err != nil {
		return http.Response{}, fmt.Errorf("error while getting cookie %v", err)
	}

	req.Header.Add("Authorization", in.Token)
	req.Header.Add("cookie", cookie)

	resp, err := in.Client.Do(CommonHeaders(req))
	if err != nil {
		return http.Response{}, err
	}
	defer resp.Body.Close()

	return *resp, nil

}

type AvatarChange struct {
	Avatar string `json:"avatar"`
}

// @me Discord Patch request to change Avatar
func (in *Instance) AvatarChanger(avatar string) (http.Response, error) {

	url := "https://discord.com/api/v9/users/@me"

	avatar = "data:image/png;base64," + avatar

	data := AvatarChange{
		Avatar: avatar,
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return http.Response{}, err
	}
	req, err := http.NewRequest("PATCH", url, strings.NewReader(string(bytes)))

	if err != nil {
		return http.Response{}, err
	}
	cookie, err := in.GetCookieString()
	if err != nil {
		return http.Response{}, fmt.Errorf("error while getting cookie %v", err)
	}

	req.Header.Add("Authorization", in.Token)
	req.Header.Add("cookie", cookie)

	resp, err := http.DefaultClient.Do(CommonHeaders(req))
	if err != nil {
		return http.Response{}, err
	}

	return *resp, nil

}

// Encoding images to b64
func EncodeImg(pathToImage string) (string, error) {

	image, err := os.Open(pathToImage)

	if err != nil {
		return "", err
	}

	defer image.Close()

	reader := bufio.NewReader(image)
	imagebytes, _ := ioutil.ReadAll(reader)

	extension := http.DetectContentType(imagebytes)

	switch extension {
	default:
		return "", fmt.Errorf("unsupported image type: %s", extension)
	case "image/jpeg":
		img, err := jpeg.Decode(bytes.NewReader(imagebytes))
		if err != nil {
			return "", err
		}
		buf := new(bytes.Buffer)
		if err := png.Encode(buf, img); err != nil {
			return "", err
		}
		return base64.StdEncoding.EncodeToString(buf.Bytes()), nil

	case "image/png":
		return base64.StdEncoding.EncodeToString(imagebytes), nil
	}
}

// Get all file paths in a directory
func GetFiles(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		paths = append(paths, dir+"/"+file.Name())
	}
	return paths, nil
}
