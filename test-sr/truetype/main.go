package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomedium"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

type Editable interface {
	image.Image
	Set(x int, y int, c color.Color)
}

func main() {
	const (
		width  = 144 * 20
		height = 72 * 15
	)

	f, err := opentype.Parse(gomedium.TTF)
	if err != nil {
		log.Fatalf("Parse: %v", err)
	}

	imageFile, err := os.Open("out.png")
	if err != nil {
		log.Fatal("can't open image")
	}

	img, err := png.Decode(imageFile)
	if err != nil {
		log.Fatal("can't decode img")
	}

	dst, ok := img.(Editable)
	if !ok {
		log.Fatal("not editable img")
	}

	// dst := image.NewGray(image.Rect(0, 0, width, height))

	drawTextBox(f, "jelly", dst, image.Rect(0, 0, width/2, height/2))

	out, err := os.Create("out.png")
	if err != nil {
		log.Fatal("can't open output image")
	}

	err = png.Encode(out, dst)
	if err != nil {
		log.Fatal("can't encode png")
	}

	err = out.Close()
	if err != nil {
		log.Fatal("error closing output")
	}
}

func drawTextBox(f *sfnt.Font, text string, dst draw.Image, rect image.Rectangle) {
	fmt.Printf("input rect: %+v\n", rect)
	fmt.Printf("image bounds: %+v\n", dst.Bounds())
	textBoxBounds := rect.Bounds().Intersect(dst.Bounds())
	fmt.Printf("text box bounds: %+v\n", textBoxBounds)

	drawTextboxBackground(dst, textBoxBounds)
	_, drawer, drawnBounds := scaleFontFaceSize(f, text, dst, textBoxBounds)
	centerTextboxDrawer(&drawer, textBoxBounds, drawnBounds)
	drawer.DrawString(text)
}

func centerTextboxDrawer(d *font.Drawer, imgBounds image.Rectangle, drawnBounds fixed.Rectangle26_6) {
	topPad := drawnBounds.Min.Y.Round() - imgBounds.Min.Y
	bottomPad := drawnBounds.Max.Y.Round() - imgBounds.Max.Y
	leftPad := drawnBounds.Min.X.Round() - imgBounds.Min.X
	rightPad := drawnBounds.Max.X.Round() - imgBounds.Max.X

	d.Dot = fixed.P(
		d.Dot.X.Round()+int(math.Abs(float64(leftPad)+float64(rightPad))/2),
		d.Dot.Y.Round()-int(math.Abs(float64(topPad)+float64(bottomPad))/2),
	)
}

func drawTextboxBackground(dst draw.Image, rect image.Rectangle) {
	// Apply a gray rectangle background for text
	bounds := rect.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel := dst.At(x, y)
			r, g, b, _ := pixel.RGBA()
			dst.Set(x, y, color.NRGBA{uint8(r >> 9), uint8(g >> 9), uint8(b >> 9), uint8(200)})
		}
	}
}

func scaleFontFaceSize(f *sfnt.Font, text string, dst draw.Image, rect image.Rectangle) (face font.Face, drawer font.Drawer, bounds fixed.Rectangle26_6) {
	imgBounds := rect.Bounds()

	startingDotX := imgBounds.Min.X
	startingDotY := imgBounds.Max.Y

	fontsize := 1.
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    fontsize,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Fatalf("NewFace: %v", err)
	}

	drawer = font.Drawer{
		Dst:  dst,
		Src:  image.White,
		Face: face,
		Dot:  fixed.P(startingDotX, startingDotY),
	}
	bounds, _ = drawer.BoundString(text)
	fmt.Printf("Measured bounds: %+v\n", bounds)

	for math.Abs(float64(bounds.Max.X.Ceil())-float64(bounds.Min.X.Floor())) < float64(imgBounds.Dx())*0.9 &&
		math.Abs(float64(bounds.Max.Y.Ceil())-float64(bounds.Min.Y.Floor())) < float64(imgBounds.Dy())*0.9 {

		fontsize = fontsize + 1

		face, err = opentype.NewFace(f, &opentype.FaceOptions{
			Size:    fontsize,
			DPI:     72,
			Hinting: font.HintingNone,
		})
		if err != nil {
			log.Fatalf("NewFace: %v", err)
		}

		drawer = font.Drawer{
			Dst:  dst,
			Src:  image.White,
			Face: face,
			Dot:  fixed.P(startingDotX, startingDotY),
		}
		bounds, _ = drawer.BoundString(text)
	}

	fmt.Printf("fontsize: %f\n", fontsize)

	return face, drawer, bounds
}

func scaleFontFaceSizeAnalytical(f *sfnt.Font, text string, dst draw.Image) {
	imgBounds := dst.Bounds()

	startingDotX := imgBounds.Min.X
	startingDotY := imgBounds.Max.Y

	fontsize := 64.
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    fontsize,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Fatalf("NewFace: %v", err)
	}

	d := font.Drawer{
		Dst:  dst,
		Src:  image.White,
		Face: face,
		Dot:  fixed.P(startingDotX, startingDotY),
	}
	bounds, _ := d.BoundString(text)
	fmt.Printf("Measured bounds: %+v\n", bounds)

	xScale := imgBounds.Dx() / (bounds.Max.X.Ceil() - bounds.Min.X.Floor())
	yScale := imgBounds.Dy() / (bounds.Max.Y.Ceil() - bounds.Min.Y.Floor())

	if xScale < yScale {
		fontsize *= float64(xScale)
	} else {
		fontsize *= float64(yScale)
	}
	face, err = opentype.NewFace(f, &opentype.FaceOptions{
		Size:    fontsize,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Fatalf("NewFace: %v", err)
	}

	d = font.Drawer{
		Dst:  dst,
		Src:  image.White,
		Face: face,
		Dot:  fixed.P(startingDotX, startingDotY),
	}

	fmt.Printf("fontsize: %f\n", fontsize)

	bounds, _ = d.BoundString(text)
	fmt.Printf("Measured bounds: %+v\n", bounds)

	// Centre text

	fmt.Printf("drawer: %+v\n", bounds)
	fmt.Printf("img: %+v\n", imgBounds)

	fmt.Printf("The dot is at %v\n", d.Dot)

	topPad := bounds.Min.Y.Round() - imgBounds.Min.Y
	bottomPad := bounds.Max.Y.Round() - imgBounds.Max.Y
	leftPad := bounds.Min.X.Round() - imgBounds.Min.X
	rightPad := bounds.Max.X.Round() - imgBounds.Max.X

	fmt.Printf("topPad: %d\n", topPad)
	fmt.Printf("bottomPad: %d\n", bottomPad)
	fmt.Printf("leftPad: %d\n", leftPad)
	fmt.Printf("rightPad: %d\n", rightPad)

	d.Dot = fixed.P(
		d.Dot.X.Round()+int(math.Abs(float64(leftPad)+float64(rightPad))/2),
		d.Dot.Y.Round()-int(math.Abs(float64(topPad)+float64(bottomPad))/2),
	)

	fmt.Printf("The dot is at %v\n", d.Dot)
	d.DrawString(text)
	fmt.Printf("The dot is at %v\n", d.Dot)
}
