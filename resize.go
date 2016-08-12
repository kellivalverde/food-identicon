package main

import (
	"fmt"
	"image/jpeg"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

// exists returns whether the given file or directory exists or not
// from http://stackoverflow.com/questions/10510691/how-to-check-whether-a-file-or-directory-denoted-by-a-path-exists-in-golang
func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func resizeEverything() {
	searchDir := "./ingredients/"

	fileList := []string{}
	err := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		os.MkdirAll(filepath.Join("resized", path, "../"), os.ModePerm)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range fileList {
		if exists(path.Join("resized", f)) {
			continue
		}
		fmt.Println(f)

		if !strings.Contains(f, ".jpg") {
			continue
		}

		file, err := os.Open(f)
		if err != nil {
			continue
		}

		// decode jpeg into image.Image
		img, err := jpeg.Decode(file)
		if err != nil {
			continue
		}
		file.Close()

		// resize to width 100 using Lanczos resampling
		// and preserve aspect ratio
		m := resize.Resize(100, 100, img, resize.Lanczos3)

		out, err := os.Create(path.Join("resized", f))
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		// write new image to file
		jpeg.Encode(out, m, nil)

	}
}

func getFileNames(ingredientImages []string{}) {
	for _, ingredient := range ingredients {
		ingredientFolder := strings.Join(strings.Split(strings.TrimSpace(ingredient), " "), "-")
		if !exists(path.Join("resized", "ingredients", ingredientFolder)) {
			continue
		}
		fileList := []string{}
		err := filepath.Walk(path.Join("resized", "ingredients", ingredientFolder), func(path string, f os.FileInfo, err error) error {
			fileList = append(fileList, path)
			return nil
		})
		ingredientImages = append(ingredientImages, fileList[rand.Intn(len(fileList))])
	}
}

func main() {
	// ingredients := []string{"olive oil", "butter", "flour", "baking soda"}
	resizeEverything()
}
