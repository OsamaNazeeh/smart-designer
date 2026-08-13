package api

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func getInputfileFormat(filename string) (string, error) {
	cmd := exec.Command(
		"ffprobe",
		"-v", "error",
		"-show_entries", "format=format_name",
		"-of", "default=noprint_wrappers=1:nokey=1",
		filename,
	)

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	ext := getFromat(string(output))
	return ext, nil  
}


func getFromat(file string) string {
	format := strings.TrimSpace(file)
	ext := ""
	switch format {
		case "jpeg_pipe", "mjpeg":
			ext = "jpg"
		case "png_pipe":
			ext = "png"
		case "webp_pipe":
			ext = "webp"
	}
	return ext 
}

func buildFFmpegOptions(
	filePath string,
	transformations Transformations,
	h *Handler, 
) ([]string, string, error) {
	inputFormat, err := getInputfileFormat(filePath)
	if err != nil {
		return nil, "", err
	}

	var filters []string

	// Grayscale
	if transformations.FiltersOptions.Grayscale {
		filters = append(filters, "format=gray")
	}

	// Dither
	if transformations.FiltersOptions.Dither {
		filters = append(filters,
			"palettegen[p]",
			"[0:v][p]paletteuse=dither=bayer:bayer_scale=3",
		)
	}

	// Resize
	if transformations.ResizeOption != (Resize{}) {
		r := transformations.ResizeOption

		filters = append(filters,
			"scale="+
				strconv.Itoa(r.Width)+":"+
				strconv.Itoa(r.Height),
		)
	}

	// Crop
	if transformations.CropOption != (Crop{}) {
		c := transformations.CropOption

		filters = append(filters,
			"crop="+
				strconv.Itoa(c.Width)+":"+
				strconv.Itoa(c.Height)+":"+
				strconv.Itoa(c.X)+":"+
				strconv.Itoa(c.Y),
		)
	}

	// Rotate
	if transformations.Rotate != 0 {
		filters = append(filters,
			"rotate="+strconv.Itoa(transformations.Rotate)+"*PI/180",
		)
	}

	// Output format
	outputFormat := inputFormat
	format := strings.ToLower(strings.TrimPrefix(transformations.Format, "."))
	if format != "" {
		switch format {
		case "jpg", "png", "webp":
			outputFormat = format
		default:
			return nil, "", fmt.Errorf("unsupported output format: %s", transformations.Format)
		}
	}


	
	outputPath := filepath.Join(os.TempDir(), "transformed." + outputFormat)

	args := []string{
		"-i", filePath,
	}

	if len(filters) > 0 {
		args = append(args,
			"-vf",
			strings.Join(filters, ","),
		)
	}

	args = append(args,
		"-y",
		outputPath,
	)
	fmt.Println(outputPath)
	return args, outputPath, nil
}