package databases

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"log"
	_ "time"

	"shop.loadout.tf/src/server/config"

	_ "go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var cancelImagesConnect context.CancelFunc
var imagesBucket *gridfs.Bucket

func InitImagesDB(config config.Database) {
	var ctx context.Context
	ctx, cancelImagesConnect = context.WithCancel(context.Background())
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.ConnectURI))
	if err != nil {
		err := fmt.Errorf("error while initializing images DB %w", err)
		log.Println(err)
		panic(err)
	}

	defer closeImagesDB()

	imagesBucket, err = gridfs.NewBucket(client.Database(config.DBName), options.GridFSBucket().SetName(config.BucketName))
	if err != nil {
		err := fmt.Errorf("error while initializing bucket DB %w", err)
		log.Println(err)
		panic(err)
	}
}

func closeImagesDB() {
	if cancelImagesConnect != nil {
		cancelImagesConnect()
	}
}

func UploadImage(filename string, img image.Image) error {
	uploadStream, err := imagesBucket.OpenUploadStream(filename)
	if err != nil {
		return err
	}

	defer uploadStream.Close()

	buf := bytes.Buffer{}
	e := png.Encoder{
		CompressionLevel: png.BestSpeed,
	}
	err = e.Encode(&buf, img)
	if err != nil {
		return err
	}

	_, err = uploadStream.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("unable to write upload stream: %w", err)
	}

	return nil
}

func GetImage(filename string) ([]byte, error) {
	downloadStream, err := imagesBucket.OpenDownloadStreamByName(filename)
	if err != nil {
		return nil, err
	}
	defer downloadStream.Close()

	p := make([]byte, downloadStream.GetFile().Length)
	if _, err = downloadStream.Read(p); err != nil {
		return nil, err
	}

	return p, nil
}
