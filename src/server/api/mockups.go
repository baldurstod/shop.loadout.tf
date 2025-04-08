package api

import (
	"encoding/base64"
	"errors"
	"image"
	"image/png"
	"net/url"
	"strings"
	"time"

	printfulsdk "github.com/baldurstod/go-printful-sdk"
	"github.com/baldurstod/randstr"
	"golang.org/x/image/draw"
	"shop.loadout.tf/src/server/mongo"
)

func RunTasks() {
	go func() {
		for {
			processMockupTasks()
			time.Sleep(10 * time.Second)
		}
	}()
}

func processMockupTasks() error {
	tasks, err := mongo.FindMockupTasks()
	if err != nil {
		return err
	}

	for _, task := range tasks {
		b64data := task.SourceImage[strings.IndexByte(task.SourceImage, ',')+1:] // Remove data:image/png;base64,
		img, err := png.Decode(base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64data)))
		if err != nil {
			return errors.New("Error while decoding image")
		}

		mockup, err := printfulsdk.GenerateMockup(img, task.Template)
		if err != nil {
			return err
		}

		filename := randstr.String(32)
		filenameThumb := filename + "_thumb"

		err = mongo.UploadImage(filename, mockup)
		if err != nil {
			return err
		}

		err = mongo.UploadImage(filenameThumb, createThumbnail(mockup, 100))
		if err != nil {
			return err
		}

		for _, productID := range task.ProductIDs {
			product, err := mongo.FindProduct(productID)
			if err != nil {
				return err
			}

			imageURL, err := url.JoinPath(imagesConfig.BaseURL, "/image/", filename)
			if err != nil {
				return errors.New("unable to create image url")
			}

			product.SetFile(task.Template.Placement, imageURL, imageURL+"_thumb")

			err = mongo.UpdateProduct(product)
			if err != nil {
				return err
			}
		}

		task.Status = "completed"
		err = mongo.UpdateMockupTask(task)
		if err != nil {
			return err
		}
	}
	return nil
}

func createThumbnail(i image.Image, size int) image.Image {
	r := i.Bounds()
	width := r.Dx()
	height := r.Dy()
	thumbWidth := size
	thumbHeight := size
	if width > height {
		if width > 0 {
			thumbHeight = size * height / width
		}
	} else {
		if height > 0 {
			thumbWidth = size * width / height
		}
	}

	thumb := image.NewNRGBA(image.Rect(0, 0, thumbWidth, thumbHeight))
	draw.ApproxBiLinear.Scale(thumb, thumb.Bounds(), i, i.Bounds(), draw.Src, nil)

	return thumb

}
