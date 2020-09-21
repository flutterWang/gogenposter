package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"unicode/utf8"

	"github.com/boombuler/barcode/qr"
)

var url = "https://www.bagevent.com/event/gopherchina2020"

type Member struct {
	Title      string
	Author     string
	Company    string
	BgPath     string
	DstPath    string
	MemberPath string
	ThumbPath  string
}

func Generate() (err error) {
	os.MkdirAll("./data/gen/thumb", os.ModePerm)
	os.MkdirAll("./data/gen/dst", os.ModePerm)
	os.MkdirAll("./data/gen/qrcode", os.ModePerm)

	// 生成二维码
	qrcode := NewQrCode(url, 140, 140, qr.M, qr.Auto)
	filePath, err := qrcode.Encode("./data/gen/qrcode")
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadFile("./poster.json")
	if err != nil {
		return err
	}

	var members []Member
	err = json.Unmarshal(data, &members)
	if err != nil {
		return err
	}

	for _, value := range members {
		fmt.Println(utf8.RuneCountInString(value.Title))
		poster := NewPoster(
			Content{
				Title:   value.Title,
				Author:  value.Author,
				Company: value.Company,
				BgPath:  value.BgPath,
				DstPath: value.DstPath,
			},
			&Rect{
				X0: 0,
				Y0: 0,
				X1: 750,
				Y1: 1334,
			},
			Avatar{
				Path:      value.MemberPath,
				ThumbPath: value.ThumbPath,
				X:         59,
				Y:         192,
				Width:     632,
				Height:    627,
			},
			Qr{
				Path: filePath,
				X:    500,
				Y:    1058,
			},
		)
		err = poster.Generate()
	}

	return
}
