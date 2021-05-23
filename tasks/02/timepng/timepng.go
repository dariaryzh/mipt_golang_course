package timepng

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"time"
)

// TimePNG записывает в `out` картинку в формате png с текущим временем
func TimePNG(out io.Writer, t time.Time, c color.Color, scale int) {
	_ = png.Encode(out, buildTimeImage(t, c, scale))
}

// buildTimeImage создает новое изображение с временем `t`
func buildTimeImage(t time.Time, c color.Color, scale int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, (X*5+4)*scale, Y*scale))
	T := t.Format(Layout)
	var runes [5]rune
	for i, s := range T {
		runes[i] = s
	}
	fillWithMask(img.SubImage(image.Rect(0, 0, X*scale, Y*scale)).(*image.RGBA), nums[runes[0]], c, scale)
	for i := 1; i < 5; i++ {
		fillWithMask(img.SubImage(image.Rect(i*(X+1)*scale, 0, (X+i*(X+1))*scale, Y*scale)).(*image.RGBA), nums[runes[i]], c, scale)
	}
	return img
}

// fillWithMask заполняет изображение `img` цветом `c` по маске `mask`. Маска `mask`
// должна иметь пропорциональные размеры `img` с учетом фактора `scale`
// NOTE: Так как это вспомогательная функция, можно считать, что mask имеет размер (3x5)
func fillWithMask(img *image.RGBA, mask []int, c color.Color, scale int) {
	for x := 0; x < X; x++ {
		for y := 0; y < Y; y++ {
			if mask[x+y*3] == 1 {
				for i := 0; i < scale; i++ {
					for j := 0; j < scale; j++ {
						img.Set((x*scale+i)+img.Bounds().Min.X, y*scale+j, c)
					}
				}
			}
		}
	}
}

const Layout = "15:04"

const X = 3
const Y = 5

var nums = map[rune][]int{
	'0': {
		1, 1, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1,
	},
	'1': {
		0, 1, 1,
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
	},
	'2': {
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
	},
	'3': {
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
	},
	'4': {
		1, 0, 1,
		1, 0, 1,
		1, 1, 1,
		0, 0, 1,
		0, 0, 1,
	},
	'5': {
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
	},
	'6': {
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1,
	},
	'7': {
		1, 1, 1,
		0, 0, 1,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
	},
	'8': {
		1, 1, 1,
		1, 0, 1,
		1, 1, 1,
		1, 0, 1,
		1, 1, 1,
	},
	'9': {
		1, 1, 1,
		1, 0, 1,
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
	},
	':': {
		0, 0, 0,
		0, 1, 0,
		0, 0, 0,
		0, 1, 0,
		0, 0, 0,
		0, 0, 0,
	},
}
