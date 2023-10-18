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
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

type urgencyColor image.Image

var urgencyColors = struct {
	low        urgencyColor
	mediumLow  urgencyColor
	mediumHigh urgencyColor
	high       urgencyColor
}{
	low:        image.NewUniform(color.RGBA{52, 235, 64, 255}),
	mediumLow:  image.NewUniform(color.RGBA{241, 236, 0, 255}),
	mediumHigh: image.NewUniform(color.RGBA{235, 88, 52, 255}),
	high:       image.NewUniform(color.RGBA{235, 0, 0, 255}),
}

type textboxPosition int

const (
	topLeft = iota
	topRight
	bottomLeft
	bottomRight
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

	f, err := opentype.Parse(goregular.TTF)
	if err != nil {
		log.Fatalf("Parse: %v", err)
	}

	// imageFile, err := os.Open("out.png")
	imageFile, err := os.Open("../../demo-heatmaps-updated/bronchus1_sr.png")
	if err != nil {
		log.Fatal("can't open image")
	}

	img, err := png.Decode(imageFile)
	if err != nil {
		log.Fatal("can't decode img")
	}

	// dst, ok := img.(Editable)
	// if !ok {
	// log.Fatal("not editable img")
	// }

	dst := image.NewRGBA(img.Bounds())
	draw.Draw(dst, dst.Bounds(), img, image.Point{}, draw.Src)

	textBox := image.Rect(0, 0,
		int(float32(dst.Rect.Max.X)*0.75),
		int(float32(dst.Rect.Max.Y)*0.1),
	)

	drawTextBox(f,
		[]string{
			"NGT Malposition risk: HIGH (41.0%)",
			"Risk of Bronchial NGT: LOW (5.0%)",
			"Risk of Rain: 20%",
		},
		dst,
		textBox,
	)

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

func drawTextBox(f *sfnt.Font, lines []string, dst draw.Image, rect image.Rectangle) {
	fmt.Printf("input rect: %+v\n", rect)
	fmt.Printf("image bounds: %+v\n", dst.Bounds())
	textBoxBounds := rect.Bounds().Intersect(dst.Bounds())
	fmt.Printf("text box bounds: %+v\n", textBoxBounds)

	// Divide the textbox space equally between lines
	textLineBounds := splitRectangleLines(rect, len(lines))
	fmt.Printf("split lines bounds: %+v\n", textLineBounds)

	// Create a scaled drawer for each line of text
	// Track the boundaries used by each drawer
	var fontSizes []float64
	var lineDrawers []font.Drawer
	var drawnBounds []image.Rectangle
	for i := range lines {
		fontsize, lineDrawer, lineBounds := scaleFontFaceSize(f, lines[i], dst, textLineBounds[i])

		fontSizes = append(fontSizes, fontsize)
		lineDrawers = append(lineDrawers, lineDrawer)
		drawnBounds = append(drawnBounds, fixedRectToRect(lineBounds))
	}

	// Each drawer fontsize scaling was done independently
	// Align them so all text is rendered at same fontsize
	smallestFontsize := math.MaxFloat64
	for _, fontsize := range fontSizes {
		if fontsize < smallestFontsize {
			smallestFontsize = fontsize
		}
	}

	// Update drawers and their bounds
	for i := range lineDrawers {
		lineDrawers[i], drawnBounds[i] = setDrawerFontsize(f, smallestFontsize, lines[i], dst, drawnBounds[i])
	}

	// Create textbox from union of all line boundaries
	textBoxBounds = unionRects(drawnBounds)
	fmt.Printf("textbox set to drawn bounds: %+v\n", textBoxBounds)
	// Add padding
	textBoxBounds = scaleRect(textBoxBounds, 1.1, 2)
	fmt.Printf("textbox scaled: %+v\n", textBoxBounds)

	// TODO: margin would need to be adjusted if allowed
	// textBoxBounds = positionTextBox(textBoxBounds, dst.Bounds(), topLeft)

	// Margin
	marginTranslate := image.Point{
		int(float64(dst.Bounds().Dx()) * 0.07),
		int(float64(dst.Bounds().Dy()) * 0.07),
	}

	textBoxBounds = textBoxBounds.Add(marginTranslate)

	// If applying margins and padding now, need to update drawer positions again
	for i := range lineDrawers {
		lineDrawers[i].Dot = lineDrawers[i].Dot.Add(
			fixed.Point26_6{
				X: fixed.I(marginTranslate.X),
				Y: fixed.I(marginTranslate.Y),
			},
		)
	}

	// Draw textbox first, background element
	drawTextboxBackground(dst, textBoxBounds)

	// Draw each line on top
	for i := range lineDrawers {
		lineDrawers[i].Src = urgencyColors.mediumHigh
		lineDrawers[i].DrawString(lines[i])
	}
}

// TODO this kind of works, but after textbox is shrunk to fit text looks strange
func positionTextBox(box image.Rectangle, imgBounds image.Rectangle, position textboxPosition) image.Rectangle {

	var targetMin image.Point

	switch position {
	case topLeft:
		targetMin = image.Point{0, 0}
	case topRight:
		targetMin = image.Point{
			imgBounds.Max.X - box.Dx(),
			0,
		}
	case bottomLeft:
		targetMin = image.Point{
			0,
			imgBounds.Max.Y - box.Dy(),
		}
	case bottomRight:
		targetMin = image.Point{
			imgBounds.Max.X - box.Dx(),
			imgBounds.Max.Y - box.Dy(),
		}
	default:
		log.Fatal("unrecognized position option %d", position)
	}

	return box.Add(targetMin.Sub(box.Min))
}

func unionRects(rects []image.Rectangle) image.Rectangle {
	union := rects[0]

	for i := 1; i < len(rects); i++ {
		union = union.Union(rects[i])
	}

	return union
}

func alignTextLinesDrawers(f *sfnt.Font, text string, dst draw.Image, textBoxBounds image.Rectangle) (font.Drawer, image.Rectangle) {
	_, drawer, drawnBounds := scaleFontFaceSize(f, text, dst, textBoxBounds)

	return drawer, fixedRectToRect(drawnBounds)
}

func splitRectangleLines(rect image.Rectangle, n int) []image.Rectangle {
	// If not exact, just leave a gap - we will center later
	newHeight := rect.Dy() / n

	var lineBounds []image.Rectangle

	for i := rect.Min.Y; i <= rect.Max.Y-newHeight; i += newHeight {
		lineBounds = append(lineBounds,
			image.Rect(
				rect.Min.X,
				i,
				rect.Max.X,
				i+newHeight,
			),
		)
	}

	return lineBounds
}

func fixedRectToRect(fixedRect fixed.Rectangle26_6) image.Rectangle {
	return image.Rect(fixedRect.Min.X.Floor(),
		fixedRect.Min.Y.Floor(),
		fixedRect.Max.X.Ceil(),
		fixedRect.Max.Y.Ceil(),
	)
}

// scaleRect scales the provided rect in place.
// Height and width are both multipliers and are scaled separately.
func scaleRect(rect image.Rectangle, width, height float32) image.Rectangle {
	scaledDx := int(float32(rect.Dx()) * width)
	scaledDy := int(float32(rect.Dy()) * height)

	diffX := scaledDx - rect.Dx()
	diffY := scaledDy - rect.Dy()

	return image.Rect(
		rect.Min.X-diffX/2,
		rect.Min.Y-diffY/2,
		rect.Max.X+diffX/2,
		rect.Max.Y+diffY/2,
	)
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
	mask := image.NewUniform(color.Alpha{150})
	draw.DrawMask(dst, rect, image.Black, image.Point{}, mask, image.Point{}, draw.Over)
}

func setDrawerFontsize(f *sfnt.Font, fontsize float64, text string, dst draw.Image, rect image.Rectangle) (drawer font.Drawer, bounds image.Rectangle) {
	imgBounds := rect.Bounds()
	startingDotX := imgBounds.Min.X
	startingDotY := imgBounds.Max.Y

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
	boundFixed, _ := drawer.BoundString(text)
	bounds = fixedRectToRect(boundFixed)

	return drawer, bounds
}

func scaleFontFaceSize(f *sfnt.Font, text string, dst draw.Image, rect image.Rectangle) (fontsize float64, drawer font.Drawer, bounds fixed.Rectangle26_6) {
	imgBounds := rect.Bounds()

	startingDotX := imgBounds.Min.X
	startingDotY := imgBounds.Max.Y

	fontsize = 1.
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

	return fontsize, drawer, bounds
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
