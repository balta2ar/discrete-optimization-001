package main

import (
	"bufio"
	"fmt"
	"math"
	"code.google.com/p/draw2d/draw2d"
    "image"
    "image/png"
	"image/color"
	"log"
	"os"
	"strings"
)

type Map struct {
	city    []Point
}

type Point struct {
	x, y float64
}


func main() {

	file, _ := os.Open(os.Args[1])
	reader := bufio.NewReader(file)
	// read number of cities
	firstLine, _ := reader.ReadString('\n')
	pointN := readInt(firstLine)
	n := int(pointN[0])

	m := new(Map)
	m.city = make([]Point, n)

	var xb, yb float64 = 0, 0

	// read city
	for i := 0; i < n; i++ {
		line, _ := reader.ReadString('\n')
		point := readInt(line)
		p := Point{float64(point[0]), float64(point[1])}
		m.city[i] = p
		if point[0] > xb {
			xb = point[0]
		}
		if point[1] >  yb {
			yb = point[1]

		}
	}


	// read path
        var path []float64
	in := bufio.NewReader(os.Stdin)
	jj := 0
	for input, _ := in.ReadString('\n'); jj <2; input, _ = in.ReadString('\n') {
		path = readInt(input)
		jj++
	}

	img := image.NewRGBA(image.Rect(0, 0, int(xb) + 100, int(yb) + 100))
        gc := draw2d.NewGraphicContext(img)
        fc := draw2d.NewGraphicContext(img)
	fc.SetFillColor(color.RGBA{0, 0, 255, 0x80})
	fc.SetStrokeColor(color.RGBA{0, 0, 255, 0x80})

	for i:= 0; i < n-1; i++ {
		x0 := (m.city[int(path[i])].x) + 10
		x1 := (m.city[int(path[i+1])].x) +10
		y0 := (m.city[int(path[i])].y) +10
		y1 := (m.city[int(path[i+1])].y) +10
		gc.MoveTo(x0, y0)
		gc.LineTo(x1,y1)
		fc.ArcTo(x0, y0, 2, 2, 0, 2*math.Pi)
        fc.SetLineWidth(1.0)
		fc.FillStroke()
	}
	x0 := (m.city[int(path[0])].x) +10
	x1 := (m.city[int(path[n-1])].x) +10
	y0 := (m.city[int(path[0])].y) +10
	y1 := (m.city[int(path[n-1])].y) +10
	gc.MoveTo(x0, y0)
	gc.LineTo(x1,y1)

	gc.Stroke()
    saveToPngFile("TestPath.png", img)


}
//////////////////////////////////////////////////////////
// return int from line
/////////////////////////////////////////////////////////
func readInt(line string) []float64 {

	words := strings.Split(line, " ")
	ints := []float64{}
	for _, word := range words {
		if word == "" {
			continue
		}
		var f float64
		fmt.Sscanf(word, "%g", &f)
		
		ints = append(ints, f)
	}
	return ints
}
/////////////////////////////////////////////////////
//Save png file
/////////////////////////////////////////////////////
func saveToPngFile(filePath string, m image.Image) {
        f, err := os.Create(filePath)
        if err != nil {
                log.Println(err)
                os.Exit(1)
        }
        defer f.Close()
        b := bufio.NewWriter(f)
        err = png.Encode(b, m)
        if err != nil {
                log.Println(err)
                os.Exit(1)
        }
        err = b.Flush()
        if err != nil {
                log.Println(err)
                os.Exit(1)
        }
        fmt.Printf("Wrote %s OK.\n", filePath)
}
