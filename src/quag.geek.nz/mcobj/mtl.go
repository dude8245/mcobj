package main

import (
	"fmt"
	"io"
	"os"
)

func printMtl(w io.Writer, blockId uint16) {
	if !noColor {
		fmt.Fprintln(w, "usemtl", MaterialNamer.NameBlockId(blockId))
	}
}

func writeMtlFile(filename string) os.Error {
	if noColor {
		return nil
	}

	var outFile, outErr = os.Create(filename)
	if outErr != nil {
		return outErr
	}
	defer outFile.Close()

	for _, color := range colors {
		color.Print(outFile)
	}

	return nil
}

type Vec2 struct {
	x float32
	y float32
}

type TexCoord struct {
	topLeft     Vec2
	bottomRight Vec2
}

func NullTexCoord() TexCoord {
	return TexCoord{Vec2{0, 0},
		Vec2{0, 0}}
}
func (t TexCoord) equals(u TexCoord) bool {
	return t.topLeft.x == u.topLeft.x &&
		t.topLeft.y == u.topLeft.y &&
		t.bottomRight.x == u.bottomRight.x &&
		t.bottomRight.y == u.bottomRight.y

}

func (t TexCoord) isNull() bool {
	return t.equals(NullTexCoord())
}

func NewTexCoord(v00 float64, v01 float64, v10 float64, v11 float64) TexCoord {
	return TexCoord{Vec2{float32(v00), float32(v01)},
		Vec2{float32(v10), float32(v11)}}

}

func (t TexCoord) TopLeft() Vec2 {
	return t.topLeft
}

func (t TexCoord) BottomRight() Vec2 {
	return t.bottomRight
}

func (t TexCoord) TopRight() Vec2 {
	return Vec2{t.bottomRight.x, t.topLeft.y}
}

func (t TexCoord) BottomLeft() Vec2 {
	return Vec2{t.topLeft.x, t.bottomRight.y}
}

func (t TexCoord) vertex(i int) Vec2 {
	switch i {
	case 0:
		return t.TopLeft()
	case 1:
		return t.TopRight()
	case 2:
		return t.BottomRight()
	case 3:
		return t.BottomLeft()
	}
	return Vec2{0, 0}
}


type MTL struct {
	blockId  byte
	metadata byte
	color    uint32
	name     string
	sideTex  TexCoord
	frontTex TexCoord
	topTex   TexCoord
	botTex   TexCoord
}

func (mtl *MTL) Print(w io.Writer) {
	var (
		r = mtl.color >> 24
		g = mtl.color >> 16 & 0xff
		b = mtl.color >> 8 & 0xff
		a = mtl.color & 0xff
	)

	fmt.Fprintf(w, "# %s\nnewmtl %s\nKd %.4f %.4f %.4f\nd %.4f\nillum 1\nmap_Kd terrain.png\n\n", mtl.name, MaterialNamer.NameBlockId(uint16(mtl.blockId)+uint16(mtl.metadata)*256), float64(r)/255, float64(g)/255, float64(b)/255, float64(a)/255)
}

func (mtl *MTL) colorId() uint16 {
	var id = uint16(mtl.blockId)
	if mtl.metadata != 255 {
		id += uint16(mtl.metadata) << 8
	}
	return id
}

func init() {
	colors = make([]MTL, 256)
	for i, _ := range colors {
		colors[i] = MTL{byte(i), 255, 0x800000ff, fmt.Sprintf("Unknown.%d", i), NullTexCoord(), NullTexCoord(), NullTexCoord(), NullTexCoord()}
	}

	extraData = make(map[byte]bool)
}

var (
	extraData map[byte]bool

	colors []MTL

	MaterialNamer BlockIdNamer
)

type BlockIdNamer interface {
	NameBlockId(blockId uint16) string
}

type NumberBlockIdNamer struct{}

func (n *NumberBlockIdNamer) NameBlockId(blockId uint16) (name string) {
	var idByte = byte(blockId & 0xff)
	var extraValue, extraPresent = extraData[idByte]
	if extraValue && extraPresent {
		name = fmt.Sprintf("%d_%d", idByte, blockId>>8)
	} else {
		name = fmt.Sprintf("%d", idByte)
	}
	return
}

type NameBlockIdNamer struct{}

func (n *NameBlockIdNamer) NameBlockId(blockId uint16) (name string) {
	var idByte = byte(blockId & 0xff)
	var extraValue, extraPresent = extraData[idByte]
	if extraValue && extraPresent {
		for _, color := range colors {
			if color.blockId == idByte && color.metadata == uint8(blockId>>8) {
				return color.name
			}
		}
	} else {
		for _, color := range colors {
			if color.blockId == idByte {
				return color.name
			}
		}
	}
	return
}