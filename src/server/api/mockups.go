package api

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"
	"net/url"
	"time"

	printfulsdk "github.com/baldurstod/go-printful-sdk"
	"golang.org/x/image/draw"
	"shop.loadout.tf/src/server/databases/shop"
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
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered: ", r)
		}
	}()

	tasks, err := shop.FindMockupTasks()
	if err != nil {
		return err
	}

	for _, task := range tasks {
		b, err := shop.GetImage(task.SourceImage)
		if err != nil {
			return errors.New("error while decoding image")
		}

		img, err := png.Decode(bytes.NewReader(b))
		if err != nil {
			return errors.New("Error while decoding image")
		}

		mockup, err := printfulsdk.GenerateMockup(img, task.Template)
		if err != nil {
			return err
		}

		filename, err := shop.InsertImage(mockup)
		if err != nil {
			return err
		}

		filenameThumb, err := shop.InsertImage(createThumbnail(mockup, 100))
		if err != nil {
			return err
		}

		for _, productID := range task.ProductIDs {
			product, err := shop.FindProduct(productID)
			if err != nil {
				return err
			}

			imageURL, err := url.JoinPath(imagesConfig.BaseURL, "/image/", filename)
			if err != nil {
				return errors.New("unable to create image url")
			}

			imageURLThumb, err := url.JoinPath(imagesConfig.BaseURL, "/image/", filenameThumb)
			if err != nil {
				return errors.New("unable to create image url")
			}

			product.SetFile(task.Template.Placement, imageURL, imageURLThumb)

			err = shop.UpdateProduct(product)
			if err != nil {
				return err
			}
		}

		task.DateUpdated = time.Now()
		task.Status = "completed"
		err = shop.UpdateMockupTask(task)
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
