package image

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"os"

	"golang.org/x/image/draw"
)

func CompressJPEG(srcPath string, quality int) ([]byte, error) {
	file, err := os.Open(srcPath)
	if err != nil {
		return nil, fmt.Errorf("Ошибка открытия файла %s: %w", srcPath, err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("Ошибка декодирования изображения %s: %w", srcPath, err)
	}

	bounds := img.Bounds()
	if bounds.Dx() > 1920 {
		newWidth := 1920
		ratio := float64(newWidth) / float64(bounds.Dx())
		newHeight := int(float64(bounds.Dy()) * ratio)

		resized := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
		draw.ApproxBiLinear.Scale(resized, resized.Bounds(), img, bounds, draw.Over, nil)
		img = resized
	}

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality}); err != nil {
		return nil, fmt.Errorf("Ошибка сжатия изображения %s: %w", srcPath, err)
	}

	return buf.Bytes(), nil
}
