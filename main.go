package main

// #cgo LDFLAGS: -lharfbuzz
// #include <harfbuzz/hb.h>
// #include <harfbuzz/hb-ot.h>
// #include <stdlib.h>
// #include <string.h>
// #include <dfuncs.h>
import "C"
import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
	"unsafe"

	"golang.org/x/image/vector"
)

var ptrMap map[unsafe.Pointer]interface{}

func SavePointer(data interface{}) unsafe.Pointer {
	if ptrMap == nil {
		ptrMap = map[unsafe.Pointer]interface{}{}
	}
	var ptr unsafe.Pointer = C.malloc(C.size_t(1))
	if ptr == nil {
		panic("cannot save pointer")
	}
	ptrMap[ptr] = data
	return ptr
}

func GoDataPointer(ptr unsafe.Pointer) interface{} {
	if data, ok := ptrMap[ptr]; ok {
		return data
	}
	return nil
}

func FreePointer(ptr unsafe.Pointer) {
	if ptr == nil {
		return
	}
	delete(ptrMap, ptr)
	C.free(ptr)
}

type blob_t *C.hb_blob_t
type face_t *C.hb_face_t
type font_t *C.hb_font_t
type buffer_t *C.hb_buffer_t

func CreateBlobFromFile(path string) blob_t {
	chptr := C.CString(path)
	defer C.free(unsafe.Pointer(chptr))
	blob := C.hb_blob_create_from_file(chptr)
	return blob
}

func CreateFace(blob blob_t, p1 uint) face_t {
	return C.hb_face_create(blob, C.uint(p1))
}

func DestroyBlob(blob blob_t) {
	C.hb_blob_destroy(blob)
}

func CreateFont(face face_t) font_t {
	return C.hb_font_create(face)
}

func OpenTypeSetFuncs(font font_t) {
	C.hb_ot_font_set_funcs(font)
}

func FontSetScale(font font_t, size int) {
	C.hb_font_set_scale(font, C.int(size)*64, C.int(size)*64)
}

func CreateBufferFromUtf8(text string) buffer_t {
	buffer := C.hb_buffer_create()
	chptr := C.CString(text)
	defer C.free(unsafe.Pointer(chptr))
	C.hb_buffer_add_utf8(buffer, chptr, -1, 0, -1)
	C.hb_buffer_guess_segment_properties(buffer)
	return buffer
}

//export moveTo
func moveTo(data unsafe.Pointer, to_x, to_y C.float) {
	drawData, _ := (GoDataPointer(data)).(*Data)
	var (
		cx = drawData.px + float32(to_x/64.)
		cy = drawData.py - float32(to_y/64.)
	)
	if drawData.rs != nil {
		drawData.rs.MoveTo(cx-drawData.minX, cy-drawData.minY)
	} else {
		drawData.setMinX(cx)
		drawData.setMinY(cy)
	}

}

//export lineTo
func lineTo(data unsafe.Pointer, to_x, to_y C.float) {
	drawData, _ := (GoDataPointer(data)).(*Data)
	var (
		cx = drawData.px + float32(to_x/64.)
		cy = drawData.py - float32(to_y/64.)
	)
	if drawData.rs != nil {
		drawData.rs.LineTo(cx-drawData.minX, cy-drawData.minY)
	} else {
		drawData.setMinX(cx)
		drawData.setMinY(cy)
	}
}

//export quadraticTo
func quadraticTo(data unsafe.Pointer, control_x, control_y, to_x, to_y C.float) {
	drawData, _ := (GoDataPointer(data)).(*Data)
	var (
		ax = drawData.px + float32(control_x/64.)
		ay = drawData.py - float32(control_y/64.)
		cx = drawData.px + float32(to_x/64.)
		cy = drawData.py - float32(to_y/64.)
	)
	if drawData.rs != nil {
		drawData.rs.QuadTo(ax-drawData.minX, ay-drawData.minY, cx-drawData.minX, cy-drawData.minY)
	} else {
		drawData.setMinX(ax, cx)
		drawData.setMinY(ay, cy)
	}
}

//export cubeTo
func cubeTo(data unsafe.Pointer, control1_x, control1_y, control2_x, control2_y, to_x, to_y C.float) {
	drawData, _ := (GoDataPointer(data)).(*Data)
	var (
		ax = drawData.px + float32(control1_x/64.)
		ay = drawData.py - float32(control1_y/64.)
		bx = drawData.px + float32(control2_x/64.)
		by = drawData.py - float32(control2_y/64.)
		cx = drawData.px + float32(to_x/64.)
		cy = drawData.py - float32(to_y/64.)
	)
	if drawData.rs != nil {
		drawData.rs.CubeTo(ax-drawData.minX, ay-drawData.minY, bx-drawData.minX,
			by-drawData.minY, cx-drawData.minX, cy-drawData.minY)
	} else {
		drawData.setMinX(ax, bx, cx)
		drawData.setMinY(ay, by, cy)
	}
}

//export closePath
func closePath(data unsafe.Pointer) {
	drawData, _ := (GoDataPointer(data)).(*Data)
	if drawData.rs != nil {
		drawData.rs.ClosePath()
	}
}

type Data struct {
	rs                             *vector.Rasterizer
	px, py, minX, minY, maxX, maxY float32
}

func (data *Data) setMinX(numbers ...float32) {
	for _, x := range numbers {
		if data.minX > x {
			data.minX = x
		}
		if x > data.maxX {
			data.maxX = x
		}
	}
}

func (data *Data) setMinY(numbers ...float32) {
	for _, y := range numbers {
		if data.minY > y {
			data.minY = y
		}
		if y > data.maxY {
			data.maxY = y
		}
	}
}

func colorizeImage(img *image.NRGBA, color color.Color) {
	for i := 0; i < img.Rect.Dx(); i++ {
		for j := 0; j < img.Rect.Dy(); j++ {
			img.Set(i, j, color)
		}
	}
}

var (
	width  int
	height int
)

func main() {

	blob := CreateBlobFromFile("./fonts/IranNastaliq.ttf")
	face := CreateFace(blob, 0)
	DestroyBlob(blob)
	font := CreateFont(face)
	OpenTypeSetFuncs(font)
	FontSetScale(font, 0.75*64)
	buffer := CreateBufferFromUtf8("نیوشا گوگولی")

	C.hb_shape(font, buffer, nil, 0)

	len := C.hb_buffer_get_length(buffer)
	glyphInfos := unsafe.Slice(C.hb_buffer_get_glyph_infos(buffer, nil), int(len))
	glyphPositions := unsafe.Slice(C.hb_buffer_get_glyph_positions(buffer, nil), int(len))
	data := Data{
		rs:   nil,
		minX: math.MaxFloat32,
		minY: math.MaxFloat32,
		maxX: -math.MaxFloat32,
		maxY: -math.MaxFloat32,
	}
	dataPtr := SavePointer(&data)
	defer FreePointer(dataPtr)
	dfuncs := C.hb_draw_funcs_create()
	C.hb_draw_funcs_set_move_to_func(dfuncs, (C.hb_draw_move_to_func_t)(C.move_to), nil, nil)
	C.hb_draw_funcs_set_line_to_func(dfuncs, (C.hb_draw_line_to_func_t)(C.line_to), nil, nil)
	C.hb_draw_funcs_set_quadratic_to_func(dfuncs, (C.hb_draw_quadratic_to_func_t)(C.quadratic_to), nil, nil)
	C.hb_draw_funcs_set_cubic_to_func(dfuncs, (C.hb_draw_cubic_to_func_t)(C.cube_to), nil, nil)
	C.hb_draw_funcs_set_close_path_func(dfuncs, (C.hb_draw_close_path_func_t)(C.close_path), nil, nil)

	x_cursor, y_cursor := 0.0, 0.0
	for i := 0; i < int(len); i++ {
		gid := glyphInfos[i].codepoint
		// cluster := glyphInfos[i].cluster
		xadd, yadd := float64(float64(glyphPositions[i].x_advance)/64.), float64(float64(glyphPositions[i].y_advance)/64.)
		xoff, yoff := float64(float64(glyphPositions[i].x_offset)/64.), float64(float64(glyphPositions[i].y_offset)/64.)
		data.px = (float32(x_cursor) + float32(xoff))
		data.py = -(float32(y_cursor) + float32(yoff))
		chptr := unsafe.SliceData(make([]C.char, 32))
		C.hb_font_get_glyph_name(font, gid, chptr, 32)
		// glyphname := C.GoString(chptr)
		// fmt.Printf("glyph='%s'	cluster=%d	advance=(%g,%g)	offset=(%g,%g)\n",
		// 	glyphname, cluster, xadd, yadd, xoff, yoff)
		C.hb_font_draw_glyph(font, gid, dfuncs, dataPtr)

		x_cursor += xadd
		y_cursor += yadd
	}
	width, height = int(math.Ceil(float64(data.maxX-data.minX))), int(math.Ceil(float64(data.maxY-data.minY)))

	data.rs = vector.NewRasterizer(width, height)
	x_cursor, y_cursor = 0.0, 0.0
	for i := 0; i < int(len); i++ {
		gid := glyphInfos[i].codepoint
		// cluster := glyphInfos[i].cluster
		xadd, yadd := float64(float64(glyphPositions[i].x_advance)/64.), float64(float64(glyphPositions[i].y_advance)/64.)
		xoff, yoff := float64(float64(glyphPositions[i].x_offset)/64.), float64(float64(glyphPositions[i].y_offset)/64.)
		data.px = (float32(x_cursor) + float32(xoff))
		data.py = -(float32(y_cursor) + float32(yoff))
		chptr := unsafe.SliceData(make([]C.char, 32))
		C.hb_font_get_glyph_name(font, gid, chptr, 32)
		// glyphname := C.GoString(chptr)
		// fmt.Printf("glyph='%s'	cluster=%d	advance=(%g,%g)	offset=(%g,%g)\n",
		// 	glyphname, cluster, xadd, yadd, xoff, yoff)
		C.hb_font_draw_glyph(font, gid, dfuncs, dataPtr)

		x_cursor += xadd
		y_cursor += yadd
	}
	data.rs.DrawOp = draw.Src
	file, _ := os.OpenFile("./out.png", os.O_WRONLY|os.O_APPEND, 0666)
	//dst := image.NewNRGBA(image.Rect(0, 0, width, height))
	//src := image.NewNRGBA(image.Rect(0, 0, width, height))
	//colorizeImage(dst, color.Black)
	//colorizeImage(src, color.RGBA{
	//	R: 255,
	//})
	dst := image.NewAlpha(image.Rect(0, 0, width, height))
	data.rs.Draw(dst, dst.Bounds(), image.Opaque, image.Point{})
	png.Encode(file, dst)
	const asciiArt = ".++8"
	buf := make([]byte, 0, height*(width+1))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			a := dst.AlphaAt(x, y).A
			buf = append(buf, asciiArt[a>>6])
		}
		buf = append(buf, '\n')
	}
	os.Stdout.Write(buf)
	C.hb_buffer_destroy(buffer)
	C.hb_font_destroy(font)
	C.hb_face_destroy(face)
}
