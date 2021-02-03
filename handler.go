package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const TelegramToken = `1523560914:AAFTl19_as8di5CjPaYeMkOFdEDcWL1p1hE`
const TelegramUrl = `https://api.telegram.org/bot`
const TelegramFileUrl = `https://api.telegram.org/file/bot`

func FilterMessage(message string) {

}

func SaveImage(id string) error {
	res, err := http.Get(fmt.Sprintf("%s%s/getFile?file_id=%s", TelegramUrl, TelegramToken, id))
	if err != nil {
		return err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	result := &GetFile{}
	if err := json.Unmarshal(body, result); err != nil {
		return err
	}

	path := result.Result.FilePath

	resImg, err := http.Get(fmt.Sprintf("%s%s/%s", TelegramFileUrl, TelegramToken, path))
	if err != nil {
		return err
	}

	defer resImg.Body.Close()

	img, err := ioutil.ReadAll(resImg.Body)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile("image.jpeg", img, 0666); err != nil {
		return err
	}
	return nil
}