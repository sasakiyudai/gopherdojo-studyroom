package imgconv

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
)

func Converter(dir, src, dst string) error {
	imgPaths, err := getFiles(dir, src)
	if err != nil {
		return err
	}

	for _, path := range imgPaths {
		if err := convert(path, dst); err != nil {
			return err
		}
	}
	return nil
}

func getFiles(dir, src string) ([]string, error) {
	var imgPaths []string

	if f, err := os.Stat(dir); err != nil {
		return nil, err
	} else if !f.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", dir)
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == "."+src {
			imgPaths = append(imgPaths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return imgPaths, nil
}

func convert(filePath, dst string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Error closing file: %s\n", err)
		}
	}()

	img, _, err := image.Decode(f)
	if err != nil {
		return err
	}

	outputPath := renameExt(filePath, dst)
	output, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := output.Close(); err != nil {
			log.Printf("Error closing file: %s\n", err)
		}
	}()

	switch dst {
	case "jpg", "jpeg":
		return convertJPG(img, output)
	case "png":
		return convertPNG(img, output)
	case "gif":
		return convertGIF(img, output)
	default:
		return fmt.Errorf("%s is not supported", dst)
	}
}

func renameExt(filePath, dst string) string {
	return filePath[:len(filePath)-len(filepath.Ext(filePath))] + "." + dst
}

func convertJPG(img image.Image, output *os.File) error {
	if err := jpeg.Encode(output, img, nil); err != nil {
		return err
	}
	return nil
}

func convertPNG(img image.Image, output *os.File) error {
	if err := png.Encode(output, img); err != nil {
		return err
	}
	return nil
}

func convertGIF(img image.Image, output *os.File) error {
	if err := gif.Encode(output, img, nil); err != nil {
		return err
	}
	return nil
}