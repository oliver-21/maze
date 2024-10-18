package main

import (
	"embed"
	"errors"
	"image"
	"image/png"
	"io/fs"
	"log"
	fsPath "path"
	"path/filepath"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"
	"github.com/hajimehoshi/ebiten/v2"
)

func getEbiten(file fs.File, name string) *ebiten.Image {
	// data, openable := file.Open(name)
	// if openable != nil {
	// 	panic("could not read: " + openable.Error())
	// }
	var (
		res image.Image
		err error
	)
	switch fsPath.Ext(name) {
	case ".png":
		res, err = png.Decode(file)
	default:
		res, _, err = image.Decode(file)
	}
	if err != nil {
		panic("could not decode with: " + err.Error())
	}
	return ebiten.NewImageFromImage(res)
}

//go:embed images
var images embed.FS

var cachFiles = map[string]*ebiten.Image{}

func get(file string) *ebiten.Image {
	prev := cachFiles[file]
	if prev != nil {
		return prev
	}
	data, err := images.Open("images/" + file)
	if err != nil {
		panic(err.Error())
	}
	defer data.Close()

	image := getEbiten(data, file)
	cachFiles[file] = image
	return image
}

func decode(file fs.File, name string) (s beep.StreamSeekCloser, format beep.Format, err error) {
	//TODO: support more extensions, mabey even allow taking audio from videos
	switch filepath.Ext(name) {
	case ".wav":
		return wav.Decode(file)
	case ".mp3":
		return mp3.Decode(file)
	case "":
		return nil, beep.Format{}, errors.New("no extension")
	default:
		return nil, beep.Format{}, errors.New("unknown extension: " + filepath.Ext(name))
	}
}

//go:embed theme.mp3
var embedThemeFile embed.FS

func getStreamer() (beep.StreamSeekCloser, beep.Format) {
	name := "theme.mp3"
	data, _ := embedThemeFile.Open(name)
	streamer, format, err := decode(data, name)

	if err != nil {
		log.Fatal(err)
	}
	return streamer, format
}

// func soundtrack() {
// 	done := make(chan bool)
// 	for {
// 		streamer, format := getStreamer()
// 		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

// 		speaker.Play(beep.Seq(streamer, beep.Callback(func() {
// 			done <- true
// 		})))
// 		<-done // Wait for sound to finish
// 		streamer.Close()
// 	}
// }

func getSound() *beep.Buffer {
	name := "theme.mp3"
	data, _ := embedThemeFile.Open(name)
	streamer, format, err := decode(data, name)

	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()
	return buffer
}

func soundtrack() {
	buffer := getSound()
	done := make(chan bool)
	for {
		segment := buffer.Streamer(0, buffer.Len())
		speaker.Play(beep.Seq(segment, beep.Callback(func() {
			done <- true
		})))
		<-done
	}
}
