// This package was written by xojoc (http://xojoc.pw)
// and is in the Public Domain do what you want with it.

/*Package turtle implements basic primitives for turtle graphics.
 */
package turtle

import (
	"bufio"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"
	"path"
)

type line struct {
	x     float64
	y     float64
	x1    float64
	y1    float64
	width float64
	color color.Color
}

func min(a, b, c float64) float64 {
	return math.Min(math.Min(a, b), c)
}
func max(a, b, c float64) float64 {
	return math.Max(math.Max(a, b), c)
}

func adjustSize(l []line) {
	minx := math.Inf(1)
	miny := math.Inf(1)
	for i := 0; i < len(l); i++ {
		minx = min(minx, l[i].x, l[i].x1)
		miny = min(miny, l[i].y, l[i].y1)
	}
	for i := 0; i < len(l); i++ {
		l[i].x -= minx
		l[i].x1 -= minx
		l[i].y -= miny
		l[i].y1 -= miny
	}
}

func widthHeight(l []line) (int, int) {
	maxx := math.Inf(-1)
	maxy := math.Inf(-1)
	for i := 0; i < len(l); i++ {
		maxx = max(maxx, l[i].x, l[i].x1)
		maxy = max(maxy, l[i].y, l[i].y1)
	}
	return int(math.Ceil(maxx)), int(math.Ceil(maxy))
}

// Turtle keeps track of the state of the Turtle.
type Turtle struct {
	X float64
	Y float64
	// Angle in radiants.
	A    float64
	draw bool

	color color.Color
	width float64
	lines []line
}

func New() *Turtle {
	t := &Turtle{}
	t.Rotate(180)
	t.PenDown()
	t.SetColor(color.RGBA{0, 0, 0, 0xff})
	t.SetWidth(5.0)
	return t
}

func dot(image draw.Image, x, y int, c color.Color) {
	image.Set(x, y, c)
}

// Bresenham's algorithm
func drawLine(image draw.Image, fx0, fy0, fx1, fy1, fw float64, c color.Color) {
	x0 := int(fx0 + 0.5)
	x1 := int(fx1 + 0.5)
	y0 := int(fy0 + 0.5)
	y1 := int(fy1 + 0.5)
	w := int(fw + 0.5)

	dx := x1 - x0
	if dx < 0 {
		dx = -dx
	}
	dy := y1 - y0
	if dy < 0 {
		dy = -dy
	}
	var sx, sy int
	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}
	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}
	err := dx - dy

	for {
		h := int(fw/2 + 0.5)
		for i := 0; i < w; i++ {
			if dx > dy {
				dot(image, x0, y0-i+h, c)
			} else {
				dot(image, x0-i+h, y0, c)
			}
		}
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

func saveToPNG(name string, image image.Image) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	buf := bufio.NewWriter(f)
	err = png.Encode(buf, image)
	if err != nil {
		return err
	}
	err = buf.Flush()
	if err != nil {
		return err
	}
	return f.Close()
}

func (t *Turtle) Save(name string) error {
	adjustSize(t.lines)
	w, h := widthHeight(t.lines)
	image := image.NewRGBA(image.Rect(0, 0, w, h))
	for i := 0; i < len(t.lines); i++ {
		drawLine(image, t.lines[i].x, t.lines[i].y, t.lines[i].x1, t.lines[i].y1, t.lines[i].width, t.lines[i].color)
	}
	if path.Ext(name) == ".png" {
		return saveToPNG(name, image)
	} else {
		log.Fatal("unknown file extension: " + path.Ext(name))
	}
	return nil
}

func (t *Turtle) SetColor(c color.Color) {
	t.color = c
}

func (t *Turtle) SetWidth(w float64) {
	t.width = w
}

func degToRad(d float64) float64 {
	return d * math.Pi / 180.0
}

func radToDeg(r float64) float64 {
	return r * 180.0 / math.Pi
}

func (t *Turtle) Rotate(angle float64) {
	t.A += degToRad(angle)
}

func (t *Turtle) Move(d float64) {
	x := t.X
	y := t.Y
	t.X += math.Cos(t.A) * (-d)
	t.Y += math.Sin(t.A) * (-d)
	if t.draw {
		t.lines = append(t.lines, line{x, y, t.X, t.Y, t.width, t.color})
	}
}

func (t *Turtle) PenUp() {
	t.draw = false
}

func (t *Turtle) PenDown() {
	t.draw = true
}

/*
func (t *Turtle) Setx(x float64) {
	t.x = x
}
func (t *Turtle) Sety(y float64) {
	t.y = y
}
func (t *Turtle) Getx() float64 {
	return t.x
}
func (t *Turtle) Gety() float64 {
	return t.y
}

func (t *Turtle) SetAngle(a float64) {
	t.angle = degToRad(a)
}
func (t *Turtle) GetAngle() float64 {
	return radToDeg(t.angle)
}
*/
