package shop

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"time"
	_ "time"

	"github.com/baldurstod/randstr"
	_ "go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
)

/*
var cancelImagesConnect context.CancelFunc
var imagesBucket *gridfs.Bucket

func closeImagesDB() {
	if cancelImagesConnect != nil {
		cancelImagesConnect()
	}
}
*/

func InsertImage(img image.Image) (string, error) {
	if shopDb == nil {
		return "", errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	buf := bytes.Buffer{}
	e := png.Encoder{
		CompressionLevel: png.BestSpeed,
	}
	err := e.Encode(&buf, img)
	if err != nil {
		return "", fmt.Errorf("failed to encode image: <%w>", err)
	}

	filename := randstr.String(32, "0123456789abcdefghijklmnopqrstuvwxyz")
	err = writeImage(filename, buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("failed to insert image: <%w>", err)
	}

	_, err = shopDb.Exec(`INSERT INTO images (filename, date_created, date_updated)
	VALUES ($1, $2, $3)`,
		filename,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return "", fmt.Errorf("failed to insert image: <%w>", err)
	}

	return filename, nil
}

func GetImage(filename string) ([]byte, error) {
	return readImage(filename)
}

func getFilePath(filename string) string {
	// Ensure id is at least 4 char long
	i := fmt.Sprintf("%04v", filename)

	return "./images/" + i[0:2] + "/" + i[2:4] + "/" + filename
}

func getFileDir(filename string) string {
	// Ensure id is at least 4 char long
	i := fmt.Sprintf("%04v", filename)

	return "./images/" + i[0:2] + "/" + i[2:4] + "/"
}

func readImage(filename string) ([]byte, error) {
	file, err := os.Open(getFilePath(filename))
	if err != nil {
		return nil, fmt.Errorf("error while opening file %w", err)
	}
	defer file.Close()

	buf, err := io.ReadAll(bufio.NewReader(file))
	if err != nil {
		return nil, fmt.Errorf("failed to read file "+filename+": <%w>", err)
	}
	return buf, nil
}

func writeImage(filename string, buf []byte) error {
	path := getFileDir(filename)
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("failed to create folders for path "+path+": <%w>", err)
	}

	err = os.WriteFile(getFilePath(filename), buf, 0666)
	if err != nil {
		return fmt.Errorf("failed to write file "+filename+": <%w>", err)
	}
	return nil

	/*
		file, err := os.Open(getFilePath(filename))
		if err != nil {
			return nil, fmt.Errorf("error while opening file %w", err)
		}
		defer file.Close()

		buf, err := io.ReadAll(bufio.NewReader(file))
		if err != nil {
			return nil, fmt.Errorf("failed to read file "+filename+": <%w>", err)
		}
		return buf, nil
	*/
}

/*


func readImage(filename string) (image.Image, error) {
	file, err := os.Open(getFilePath(filename))
	if err != nil {
		return nil, fmt.Errorf("error while opening file %w", err)
	}
	defer file.Close()

	img, err := png.Decode(bufio.NewReader(file))
	if err != nil {
		return nil, fmt.Errorf("error while decoding image "+filename+" %w", err)
	}

	return img, nil
}
*/
