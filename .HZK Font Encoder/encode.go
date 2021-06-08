package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"strconv"
)

func main() {
	fmt.Println("Running, Please wait...")
	widths := []int{6, 8, 10, 12, 14, 16, 20, 24, 28, 32}
	var charImage image.Image
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	fontFile, err := os.Open("0T5UIC1.HZK")
	if err != nil {
		log.Fatal(err)
	}
	newFontFile, err := os.Create("0T5UIC1_NEW.HZK")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("")
	for size, width := range widths {
		height := 2 * width
		var numBytes int
		if width%8 > 0 {
			numBytes = ((width / 8) + 1)
		} else {
			numBytes = (width / 8)
		}
		picWidth := ((width + 1) * 16) + 1
		picHeight := ((height + 1) * 8) + 1
		imgPath := fmt.Sprintf("0x%02d_%dx%d_0-127.png", size, width, height)
		imgExists := true
		if _, err := os.Stat(imgPath); os.IsNotExist(err) {
			imgExists = false
		} else {
			imgFile, err := os.Open(imgPath)
			if err != nil {
				log.Fatal(err)
			}
			defer imgFile.Close()
			charImage, _, err = image.Decode(imgFile)
			if err != nil {
				log.Fatal(err)
			}
			bounds := charImage.Bounds()
			if picWidth != bounds.Max.X || picHeight != bounds.Max.Y {
				log.Fatal(fmt.Sprintf("Picture 0x0%d is not the correct dimensions. Picture needs to be %dx%d", size, picWidth, picHeight))
			}
		}
		charBytes := make([]byte, numBytes)
		for imgRow := 0; imgRow < 8; imgRow++ {
			for imgColumn := 0; imgColumn < 16; imgColumn++ {
				if size == 9 && imgRow == 7 && imgColumn == 15 {
					continue
				}
				for vert := 1; vert <= height; vert++ {
					bits := ""
					_, err := fontFile.Read(charBytes)
					if err != nil {
						log.Fatal(err)
					}
					if imgExists {
						for horz := 1; horz <= width; horz++ {
							R, G, B, A := charImage.At((imgColumn*(width+1))+horz, (imgRow*(height+1))+vert).RGBA()
							if A < 65000 {
								break
							}
							Pixel := int((R + G + B) / 3)
							if Pixel == int(color.Black.Y) {
								bits += "1"
							} else if Pixel == int(color.White.Y) {
								bits += "0"
							} else {
								log.Fatal("Picture could not be scanned correctly")
							}
						}
					}
					if bits == "" {
						newFontFile.Write(charBytes)
					} else {
						newCharBytes := BitsToBlocks(bits, numBytes)
						newFontFile.Write(newCharBytes)
					}
				}
			}
		}
	}

	rest := make([]byte, 2048)
	for {
		n, err := fontFile.Read(rest)
		if n == 0 {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if _, err := newFontFile.Write(rest[:n]); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Finished")
}

func BitsToBlocks(data string, byteLen int) []byte {
	var out []byte
	var str string
	for {
		if len(data) == (byteLen * 8) {
			break
		} else {
			data += "0"
			continue
		}
	}
	for i := len(data); i > 0; i -= 8 {
		if i-8 < 0 {
			str = string(data[0:i])
		} else {
			str = string(data[i-8 : i])
		}
		v, err := strconv.ParseUint(str, 2, 8)
		if err != nil {
			panic(err)
		}
		out = append([]byte{byte(v)}, out...)
	}
	return out
}
