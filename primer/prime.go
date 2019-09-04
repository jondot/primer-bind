package primer

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/fogleman/primitive/primitive"
	"github.com/nfnt/resize"
)

type OnIterationDone interface {
	Do(int)
}
type JsonDummy struct {
	hello   string
	meaning int
}

func Bench() {
	for i := 0; i < 10000000; i++ {
		_, _ = json.Marshal(JsonDummy{hello: "world", meaning: 42})
	}
	fmt.Print("Done\n")
}

func Process(name string, iters int, size int, workers int, mode int, iterdone OnIterationDone) string {
	out := ""
	//files, err := ioutil.ReadDir(name)
	//if err != nil {
	//log.Fatal(err)
	//}

	//for _, file := range files {
	//println(file.Name())
	//out = out + file.Name()
	//}
	if workers <= 0 {
		workers = 1
	}
	if size <= 0 {
		size = 250
	}
	if mode <= 0 {
		mode = 1
	}

	OutputSize := 1024

	Outputs := []string{name + ".out.png"}
	Alpha := 128
	Number := iters

	name = strings.Replace(name, "file://", "", -1)
	input, err := primitive.LoadImage(name)
	if err != nil {
		fmt.Printf("Error: %v", err)
		iterdone.Do(-1)
	}

	input = resize.Thumbnail(uint(size), uint(size), input, resize.Bilinear)

	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	bg := primitive.MakeColor(primitive.AverageImageColor(input))
	fmt.Printf("Input: ok")

	model := primitive.NewModel(input, bg, OutputSize, workers)
	fmt.Printf("iteration %d, time %.3f, score %.6f\n", 0, 0.0, model.Score)
	start := time.Now()
	for i := 1; i <= Number; i++ {
		// find optimal shape and add it to the model
		model.Step(primitive.ShapeType(mode), Alpha, 1)
		elapsed := time.Since(start).Seconds()
		fmt.Printf("iteration %d, time %.3f, score %.6f\n", i, elapsed, model.Score)
		iterdone.Do(i)

		// write output image(s)
		for _, output := range Outputs {
			ext := strings.ToLower(filepath.Ext(output))
			saveFrames := strings.Contains(output, "%") && ext != ".gif"
			if saveFrames || i == Number {
				path := output
				if saveFrames {
					path = fmt.Sprintf(output, i)
				}
				fmt.Printf("writing %s\n", path)
				switch ext {
				default:
					fmt.Errorf("unrecognized file extension: %s", ext)
				case ".png":
					fmt.Printf("saving: %s", path)
					primitive.SavePNG(path, model.Context.Image())
					//case ".jpg", ".jpeg":
					//check(primitive.SaveJPG(path, model.Context.Image(), 95))
					//case ".svg":
					//check(primitive.SaveFile(path, model.SVG()))
					//case ".gif":
					//frames := model.Frames(0.001)
					//check(primitive.SaveGIFImageMagick(path, frames, 50, 250))
					//}
				}
			}
		}
	}

	return out
}
