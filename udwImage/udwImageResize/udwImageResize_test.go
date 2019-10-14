package udwImageResize

import (
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwImage"
	"github.com/tachyon-protocol/udw/udwTest"
	"testing"
)

var pngContent16x16 = []byte("\x89PNG\r\n\n\x00\x00\x00\rIHDR\x00\x00\x00\x00\x00\x00\b\x00\x00\x00\xf3\xffa\x00\x00\xcdIDAT8\x8d\x95\x91\xbf\x8aSA\xc6g\xe6\x92MD\x92\xa0\xf2\b\"\x8ao 6\t\xf8\f\xf1\t\xb4l,m\xb1\xd4N\xb1\xb0\xb4\xbc\x90Z\x84\xb2`%h\nu\x9b%H&7\xf7\xdf\xcc\vsw6+\xeb\a3\xcc|\xbf\xf9\xce9\xc2F\xd3\xe9\xf4.p\x87(\xc6\xf8y4}\xd8{\x99\xa6\xe9Ӳ,\xd59\xa7Y\x96\xa9sN\x9ds\xbaZ\xad\xd49\xa7\xcb\xe5Rg\xb3YH\xd3\xf4\xfe\xb6/\xd9>\xd4uM\x9e\xe7\x00\xa8*\x00\"\x82\xaa\xa2\xaa\f3ߤi\xcax<~\v`\xb6\"\x82\xb5k-Ƙ\x9d%\",\v\xaa\xaa2\xfd~\xffu\xe3\xddI\xb0\xfdؘS\xb6\xaab\xade0\xa0\xaa\xf4z\xbd\xdf\xa0I\xb0O\"B\x92\xfc}^\x96\xe5\xfec\xce\x9c\xa7\v'\xb80`\xbbv\x80\xc3\xefq\xf5W\x96\xd91\xf7n\xbe\xf8\x809\xedK\xae(\xf8\xf2\xee1\x97\xaf\f\xf9\xe5Nmg\xc6\xd8\xcc=\xc6H\x8cUED\xa0m\v\xaeuS\xe7n\xe0\xc4\xec\xbd\xdf\xec\xc1\x97\x8d\"\xea=\x87y\x87g\x9f\xbe\x85\x8f\xdf>2[Fc\xa4\xaek\xaa\xaa\xa2\xaa*\xbc\xbe \xa8'\xfa@\\\xe7*\xe4ݫ\x00ݤ\xd3\xe9km{\xbd^\xb7\xe7\xf3\xf9\x99&\x95e\xc9\xda\xe4\xf8P\xf5j\r\xc02\x8f\xa8\xaaM\xac\xb50\x93\xc9乪\xbeܔ\xd5$3O^\xdd~x\xe3\xd6\xf5\ae^\xab\x8cDL+Q\xd5*;>ji\xf0H\xf3\xd3\xc1\xc1\xadV\xcb\x00\x84H\x92\xa4))\xda\"b6ਪQD|\x8c\xb1\xf8Z\xb0\xfb\xcf<qb\x98\x00\x00\x00\x00IEND\xaeB`\x82\n")

func TestMustResizePngFileToHeightAndWidth(ot *testing.T) {

	udwFile.MustDelete("testFile")
	defer udwFile.MustDelete("testFile")
	udwFile.MustWriteFileWithMkdir("testFile/1.png", pngContent16x16)
	MustResizePngFileToHeightAndWidth(MustResizePngFileToHeightAndWidthRequest{
		InFilePath:  "testFile/1.png",
		OutFilePath: "testFile/2.png",
		Height:      16,
		Width:       16,
	})
	udwTest.Equal(udwFile.MustReadFile("testFile/2.png"), pngContent16x16)

	MustResizePngFileToHeightAndWidth(MustResizePngFileToHeightAndWidthRequest{
		InFilePath:  "testFile/1.png",
		OutFilePath: "testFile/1.png",
		Height:      16,
		Width:       16,
	})
	udwTest.Equal(udwFile.MustReadFile("testFile/1.png"), pngContent16x16)

	MustResizePngFileToHeightAndWidth(MustResizePngFileToHeightAndWidthRequest{
		InFilePath:  "testFile/1.png",
		OutFilePath: "testFile/1.png",
		Height:      8,
		Width:       8,
	})
	config := udwImage.MustPngDecodeConfigFromFile("testFile/1.png")
	udwTest.Equal(config.Height, 8)
	udwTest.Equal(config.Width, 8)
}
