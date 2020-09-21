package service

import (
	"image"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"os"
	"unicode/utf8"

	"github.com/golang/freetype"
)

type Poster struct {
	*Content
	Avatar *Avatar
	*Rect
	Qr *Qr
}

type Content struct {
	Title   string
	Author  string
	Company string
	BgPath  string
	DstPath string
	DstFile *os.File
}

type Rect struct {
	X0 int
	Y0 int
	X1 int
	Y1 int
}

type Qr struct {
	Path string
	X    int
	Y    int
}

type DrawText struct {
	JPG draw.Image

	Title string
	X0    int
	Y0    int
	Size0 float64

	TitleNext string
	X3        int
	Y3        int

	Author string
	X1     int
	Y1     int
	Size1  float64

	Company string
	X2      int
	Y2      int
	Size2   float64
}

func NewPoster(content Content, rect *Rect, avatar Avatar, qr Qr) *Poster {
	return &Poster{
		Content: &content,
		Rect:    rect,
		Avatar:  &avatar,
		Qr:      &qr,
	}
}

func (p *Poster) Generate() (err error) {
	p.DstFile, err = os.OpenFile(p.DstPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	bgFile, err := os.Open(p.BgPath)
	if err != nil {
		return err
	}
	defer bgFile.Close()

	bgImage, err := jpeg.Decode(bgFile)
	if err != nil {
		return err
	}

	err = p.Avatar.Thumb()
	if err != nil {
		return
	}

	avatarFile, err := os.Open(p.Avatar.ThumbPath)
	if err != nil {
		return err
	}
	defer avatarFile.Close()

	avatarImage, err := jpeg.Decode(avatarFile)
	if err != nil {
		return err
	}

	qrFile, err := os.Open(p.Qr.Path)
	if err != nil {
		return err
	}
	defer qrFile.Close()

	qrImage, err := jpeg.Decode(qrFile)
	if err != nil {
		return err
	}

	jpg := image.NewRGBA(image.Rect(p.Rect.X0, p.Rect.Y0, p.Rect.X1, p.Rect.Y1))
	draw.Draw(jpg, jpg.Bounds(), bgImage, bgImage.Bounds().Min, draw.Over)
	draw.Draw(jpg, jpg.Bounds(), avatarImage, avatarImage.Bounds().Min.Sub(image.Pt(p.Avatar.X, p.Avatar.Y)), draw.Over)
	draw.Draw(jpg, jpg.Bounds(), qrImage, avatarImage.Bounds().Min.Sub(image.Pt(p.Qr.X, p.Qr.Y)), draw.Over)

	drawText := &DrawText{
		JPG: jpg,

		Title: p.Title,
		X0:    94,
		Y0:    892,
		Size0: 40,

		Author: p.Author,
		X1:     500,
		Y1:     892,
		Size1:  35,

		Company: p.Company,
		X2:      500,
		Y2:      948,
		Size2:   35,
	}

	if utf8.RuneCountInString(p.Title) > 10 {
		drawText.Title = substring(p.Title, 0, 10)
		drawText.TitleNext = substring(p.Title, 10, utf8.RuneCountInString(p.Title))
		drawText.X3 = 94
		drawText.Y3 = 948
	}

	err = p.DrawPoster(drawText, "msyhbd.ttc")

	if err != nil {
		return err
	}

	return nil
}

func substring(source string, start int, end int) string {
	var r = []rune(source)
	length := len(r)

	if start < 0 || end > length || start > end {
		return ""
	}

	if start == 0 && end == length {
		return source
	}

	var substring = ""
	for i := start; i < end; i++ {
		substring += string(r[i])
	}
	return substring
}

func (p *Poster) DrawPoster(d *DrawText, fontName string) error {
	fontSource := "./data/fonts/" + fontName
	fontSourceBytes, err := ioutil.ReadFile(fontSource)
	if err != nil {
		return err
	}

	trueTypeFont, err := freetype.ParseFont(fontSourceBytes)
	if err != nil {
		return err
	}
	fc := freetype.NewContext()
	fc.SetDPI(72)
	fc.SetFont(trueTypeFont)
	fc.SetFontSize(d.Size0)
	fc.SetClip(d.JPG.Bounds())
	fc.SetDst(d.JPG)
	fc.SetSrc(image.Black)

	_, err = fc.DrawString(d.Title, freetype.Pt(d.X0, d.Y0))
	if err != nil {
		return err
	}

	_, err = fc.DrawString(d.TitleNext, freetype.Pt(d.X3, d.Y3))
	if err != nil {
		return err
	}

	fc.SetFontSize(d.Size1)
	_, err = fc.DrawString(d.Author, freetype.Pt(d.X1, d.Y1))
	if err != nil {
		return err
	}

	fc.SetFontSize(d.Size2)
	_, err = fc.DrawString(d.Company, freetype.Pt(d.X2, d.Y2))
	if err != nil {
		return err
	}

	err = jpeg.Encode(p.DstFile, d.JPG, nil)
	if err != nil {
		return err
	}

	return nil
}
