package main

import (
	"os"
	"bufio"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"strings"
)

var imagePath = flag.String("image", "", "Image path. JPG or PNG (example -image=path/to/image).")
var outputPath = flag.String("output", "", "Output path (example -output=path/to/html/file).")

func init() {
	flag.Parse()
	if (*imagePath == "") {
		log.Fatal("image path not specified")
	}
	if (*outputPath == "") {
		log.Fatal("output path not specified")
	}
	if !strings.HasSuffix(*outputPath, ".html"){
		*outputPath += ".html"
	}
}

func main() {
	reader, err := os.Open(*imagePath);
	defer reader.Close()
	if err != nil {
		log.Fatal("file open error", err)
	}

	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal("image decode error", err)
	}

	outputFile, err := os.OpenFile(*outputPath, os.O_CREATE, 0660)
	if err != nil {
		log.Fatal("output file create fail", err)
	}

	writer := bufio.NewWriter(outputFile);
	bounds := img.Bounds()
	writer.WriteString(fmt.Sprintf(`
	<html>
		<head>
			<style>
				.img {
					width:%dpx;
					height:%dpx;
					background-image:
	`,
		bounds.Max.X, bounds.Max.Y))

	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			writer.WriteString(fmt.Sprintf("radial-gradient(1px circle at %dpx %dpx, rgba(%d,%d,%d,%d), transparent 1px)",
				x, y, (r >> 8) & 0xFF, (g >> 8) & 0xFF, (b >> 8) & 0xFF, (a >> 8) & 0xFF))

			if x != bounds.Max.X - 1 || y != bounds.Max.Y - 1 {
				writer.WriteByte(',')
				writer.WriteByte('\n')
			}
		}
	}

	writer.WriteString(`;
			</style>
		</head>
		<body>
			<div class="img"></div>
		</body>
	</html>
	`)
	writer.Flush()

	log.Println("Done: Image [", *imagePath, "] -> [", *outputPath, "]")
}