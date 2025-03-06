package api

import (
	"encoding/base64"
	"errors"
	"image/png"
	"log"
	"net/url"
	"strings"
	"time"

	printfulsdk "github.com/baldurstod/go-printful-sdk"
	"github.com/baldurstod/randstr"
	"shop.loadout.tf/src/server/model"
	"shop.loadout.tf/src/server/mongo"
)

func initMockupTasks(tasks []*model.MockupTask) error {
	for task := range tasks {
		log.Println(task)
	}

	go RunTasks()

	return nil
}

var running = false

func RunTasks() {
	if running {
		return
	}
	running = true
	for {
		ProcessMockupTasks()
		time.Sleep(10 * time.Second)
	}
}

func ProcessMockupTasks() error {
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
		err = mongo.UploadImage(filename, mockup)
		if err != nil {
			return err
		}

		for _, productID := range task.ProductIDs {
			product, err := mongo.FindProduct(productID.Hex())
			if err != nil {
				return err
			}

			imageURL, err := url.JoinPath(imagesConfig.BaseURL, "/image/", filename)
			if err != nil {
				return errors.New("unable to create image url")
			}

			product.AddFile(task.Template.Placement, imageURL)

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
