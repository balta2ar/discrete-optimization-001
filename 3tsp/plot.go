package main

import (
    "fmt"
    "math"
    "code.google.com/p/draw2d/draw2d"
    "image"
    "image/png"
    "image/color"
    "log"
)

type Map struct {
    city    []Point
}

type Point struct {
    x, y float64
}

func main() {
    i := image.NewRGBA(image.Rect(0, 0, 1000, 1000))
    gc := draw2d.NewGraphicContext(i)
    fc := draw2d.NewGraphicContext(i)
    fc.SetFillColor(color.RGBA{0, 0, 255, 0x80})
    fc.SetStrokeColor(color.RGBA{0, 0, 255, 0x80})


    for i := 0; i < n-1; i++ {
        x0 := (m.city[path[i]].x)/10 + 10
        x1 := (m.city[path[i+1]].x)/10 +10
        y0 := (m.city[path[i]].y)/10 +10
        y1 := (m.city[path[i+1]].y)/10 +10
        gc.MoveTo(x0, y0)
        gc.LineTo(x1,y1)
        fc.ArcTo(x0, y0, 2, 2, 0, 2*math.Pi)
        fc.SetLineWidth(1.0)
        fc.FillStroke()
    }
    x0 := (m.city[path[0]].x)/10 +10
    x1 := (m.city[path[n-1]].x)/10 +10
    y0 := (m.city[path[0]].y)/10 +10
    y1 := (m.city[path[n-1]].y)/10 +10
    gc.MoveTo(x0, y0)
    gc.LineTo(x1,y1)

    gc.Stroke()
    saveToPngFile("TestPath.png", i)
}

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
        //fmt.Printf("Wrote %s OK.\n", filePath)
}

