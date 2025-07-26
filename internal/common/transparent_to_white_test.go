package common

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"
)

const TestImagePath = "test_data/1525666415627-logo_red_800x800.png"

func TestConvertTransparentToWhite(t *testing.T) {
	inputPath := TestImagePath
	input, err := os.ReadFile(inputPath)
	if err != nil {
		t.Fatalf("Failed to read input file: %v", err)
	}

	output, err := ConvertTransparentToWhite(input)
	if err != nil {
		t.Fatalf("ConvertTransparentToWhite failed: %v", err)
	}

	outputPath := fmt.Sprintf("%s_transparent_to_white.png", strings.TrimSuffix(inputPath, path.Ext(inputPath)))
	if err := os.WriteFile(outputPath, output, 0644); err != nil {
		t.Fatalf("Failed to write output file: %v", err)
	}
	t.Logf("Output written to %s", outputPath)
}

func TestPNGToJPG(t *testing.T) {
	inputPath := TestImagePath
	input, err := os.ReadFile(inputPath)
	if err != nil {
		t.Fatalf("Failed to read input file: %v", err)
	}

	output, err := PNGToJPG(input, 85)
	if err != nil {
		t.Fatalf("PNGToJPG failed: %v", err)
	}

	outputPath := fmt.Sprintf("%s_png_to_jpg.jpg", strings.TrimSuffix(inputPath, path.Ext(inputPath)))
	if err := os.WriteFile(outputPath, output, 0644); err != nil {
		t.Fatalf("Failed to write output file: %v", err)
	}
	t.Logf("Output written to %s", outputPath)
}
