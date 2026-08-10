package shop

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"
	"time"
	_ "time"

	"github.com/baldurstod/randstr"
)

func InsertImage(img image.Image, thumb image.Image) (string, error) {
	if shopDb == nil {
		return "", errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	imgBuf := bytes.Buffer{}
	thumbBuf := bytes.Buffer{}
	e := png.Encoder{
		CompressionLevel: png.BestSpeed,
	}
	err := e.Encode(&imgBuf, img)
	if err != nil {
		return "", fmt.Errorf("failed to encode image: <%w>", err)
	}
	err = e.Encode(&thumbBuf, thumb)
	if err != nil {
		return "", fmt.Errorf("failed to encode thunb: <%w>", err)
	}

	filename := randstr.String(32, "0123456789abcdefghijklmnopqrstuvwxyz")

	_, err = shopDb.Exec(`INSERT INTO images (filename, image, thumb, date_created, date_updated)
	VALUES ($1, $2, $3, $4, $5)`,
		filename,
		imgBuf.Bytes(),
		thumbBuf.Bytes(),
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return "", fmt.Errorf("failed to insert image: <%w>", err)
	}

	return filename, nil
}

func GetImage(filename string) ([]byte, error) {
	if shopDb == nil {
		return nil, errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	query := `SELECT image FROM images WHERE filename = $1;`
	row := shopDb.QueryRow(query, filename)

	var image string

	err := row.Scan(&image)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row in GetImage: <%w>", err)
	}

	return []byte(image), nil
}
