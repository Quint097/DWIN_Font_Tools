package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("======  HZK Font Encoder  ======")
	fmt.Println("")
	srcpath, err := os.Executable()
	srcpath = filepath.Dir(srcpath)
	if err != nil {
		fmt.Println("==!  " + err.Error())
		os.Exit(1)
	}
	path := strings.TrimSuffix(srcpath, "src")
	if _, err := os.Stat(path + "\\Images"); os.IsNotExist(err) {
		fmt.Println("==!  Can't find the image folder.")
		os.Exit(1)
	}
	var filename string
	if len(os.Args) < 2 {
		fmt.Print("==>  Please enter the name of your new font file ending with .HZK: ")
		fmt.Scanln(&filename)
		filename = path + "\\" + filename
	}
	if !strings.HasSuffix(filename, ".HZK") {
		fmt.Println("==!  Filename can not be used.")
		os.Exit(1)
	}
	fmt.Println("")
	widths := []int{6, 8, 10, 12, 14, 16, 20, 24, 28, 32}
	var charImage image.Image
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	fontFile, err := os.Open(srcpath + "\\0T5UIC1.HZK")
	if err != nil {
		fmt.Println("==!  " + err.Error())
		os.Exit(1)
	}
	newFontFile, err := os.Create(filename)
	if err != nil {
		fmt.Println("==!  " + err.Error())
		os.Exit(1)
	}
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
		imgPath := fmt.Sprintf(path+"\\Images\\0x%02d_%dx%d_0-127.png", size, width, height)
		fmt.Println("===  Processing: " + imgPath)
		imgExists := true
		if _, err := os.Stat(imgPath); os.IsNotExist(err) {
			imgExists = false
			fmt.Println("===  Skipping: Image not found.")
		} else {
			imgFile, err := os.Open(imgPath)
			if err != nil {
				fmt.Println("==!  " + err.Error())
				os.Exit(1)
			}
			defer imgFile.Close()
			charImage, _, err = image.Decode(imgFile)
			if err != nil {
				fmt.Println("==!  " + err.Error())
				os.Exit(1)
			}
			bounds := charImage.Bounds()
			if picWidth != bounds.Max.X || picHeight != bounds.Max.Y {
				fmt.Printf("==!  Picture 0x0%d is not the correct dimensions. Picture needs to be %dx%d\n", size, picWidth, picHeight)
				os.Exit(1)
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
						fmt.Println("==!  " + err.Error())
						os.Exit(1)
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
								fmt.Println("==!  Picture could not be scanned correctly")
								os.Exit(1)
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
		fmt.Println("")
	}

	rest := make([]byte, 2048)
	for {
		n, err := fontFile.Read(rest)
		if n == 0 {
			break
		}
		if err != nil {
			fmt.Println("==!  " + err.Error())
			os.Exit(1)
		}
		if _, err := newFontFile.Write(rest[:n]); err != nil {
			fmt.Println("==!  " + err.Error())
			os.Exit(1)
		}
	}
	fmt.Println("======  Finished  ======")
	fmt.Println("")
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
			fmt.Println("==!  " + err.Error())
			os.Exit(1)
		}
		out = append([]byte{byte(v)}, out...)
	}
	return out
}
