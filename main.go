package main

import (
	"flag"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/fogleman/gg"
)

// calculateOffset calculates the offset needed to center image b within image a.
func calculateOffset(a image.Image, b image.Image) (int, int) {
	aw, ah := a.Bounds().Dx(), a.Bounds().Dy()
	bw, bh := b.Bounds().Dx(), b.Bounds().Dy()
	return (aw - bw) / 2, (ah - bh) / 2
}

// loadImage loads an image from the given file path.
func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// saveImage saves an image to the given file path.
func saveImage(img image.Image, path string) error {
	outFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outFile.Close()

	return png.Encode(outFile, img)
}

// processScreenshot processes a single screenshot by overlaying it onto the frame and saving the result.
func processScreenshot(frameImg image.Image, screenshotPath string, outputDir string) error {
	img, err := loadImage(screenshotPath)
	if err != nil {
		log.Printf("画像の読み込みに失敗しました: %v", err)
		return err
	}

	output := image.NewRGBA(frameImg.Bounds())

	offsetX, offsetY := calculateOffset(frameImg, img)

	dc := gg.NewContextForRGBA(output)
	dc.DrawRoundedRectangle(float64(offsetX), float64(offsetY), float64(img.Bounds().Dx()), float64(img.Bounds().Dy()), 80)
	dc.Clip()

	dc.DrawImage(img, offsetX, offsetY)
	draw.Draw(output, output.Bounds(), frameImg, image.Point{}, draw.Over)

	outputPath := filepath.Join(outputDir, filepath.Base(screenshotPath))
	err = saveImage(output, outputPath)
	if err != nil {
		log.Printf("出力ファイルの保存に失敗しました: %v", err)
		return err
	}

	log.Printf("画像を保存しました: %s", outputPath)
	return nil
}

// processScreenshotsInFolder processes all PNG screenshots in the given folder.
func processScreenshotsInFolder(folderPath string, frameImg image.Image, outputDir string) error {
	return filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".png" {
			return processScreenshot(frameImg, path, outputDir)
		}
		return nil
	})
}

func main() {
	folderPath := flag.String("screenshots", "integration_test/screenshots", "スクリーンショット画像のフォルダパス")
	framePath := flag.String("frame", "assets/internal/frame.png", "フレーム画像のパス")
	outputDir := flag.String("output", "assets/internal/release_screenshots", "出力ディレクトリ")

	flag.Parse()

	frameImg, err := loadImage(*framePath)
	if err != nil {
		log.Fatalf("フレーム画像の読み込みに失敗しました: %v", err)
	}

	err = os.MkdirAll(*outputDir, os.ModePerm)
	if err != nil {
		log.Fatalf("出力ディレクトリの作成に失敗しました: %v", err)
	}

	err = processScreenshotsInFolder(*folderPath, frameImg, *outputDir)
	if err != nil {
		log.Fatalf("画像の処理に失敗しました: %v", err)
	}
}
