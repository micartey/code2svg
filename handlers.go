package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func handleSVG(w http.ResponseWriter, r *http.Request) {
	codeBase64 := r.URL.Query().Get("code")
	if codeBase64 == "" {
		body, _ := io.ReadAll(r.Body)
		codeBase64 = string(body)
	}

	if codeBase64 == "" {
		http.Error(w, "Missing code parameter", http.StatusBadRequest)
		return
	}

	// Refactored: Base64 decoding and SVG generation logic extracted to helper functions
	decoded, err := decodeBase64(codeBase64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid base64 string: %v", err), http.StatusBadRequest)
		return
	}

	svg, err := generateSVG(string(decoded))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write([]byte(svg))
}

func generateSVG(code string) (string, error) {
	code = strings.ReplaceAll(code, "\r\n", "\n")
	lines := strings.Split(code, "\n")
	lineCount := len(lines)
	if lineCount == 0 {
		lineCount = 1
	}

	totalHeight := 60 + 20*lineCount

	const charWidth = 8.5
	maxWidth := 800.0
	for _, line := range lines {
		indent, content := calculateIndent(line)
		lineWidth := 40.0 + float64(indent*20) + (float64(len(content)) * charWidth)
		if lineWidth > maxWidth {
			maxWidth = lineWidth
		}
	}
	totalWidth := int(maxWidth)
	if totalWidth > 800 {
		totalWidth += 20
	}

	svgTemplate, err := os.ReadFile("code_preview.svg")
	if err != nil {
		return "", fmt.Errorf("could not read SVG template")
	}

	svgStr := string(svgTemplate)

	reSvgOpen := regexp.MustCompile(`<svg width="\d+" height="\d+" viewBox="0 0 \d+ \d+"`)
	svgStr = reSvgOpen.ReplaceAllString(svgStr, fmt.Sprintf(`<svg width="%d" height="%d" viewBox="0 0 %d %d"`, totalWidth, totalHeight, totalWidth, totalHeight))

	reBgRect := regexp.MustCompile(`<rect width="\d+" height="\d+" rx="8" class="bg"/>`)
	svgStr = reBgRect.ReplaceAllString(svgStr, fmt.Sprintf(`<rect width="%d" height="%d" rx="8" class="bg"/>`, totalWidth, totalHeight))

	var codeContent strings.Builder
	for i, line := range lines {
		indent, content := calculateIndent(line)
		x := indent * 20
		y := i * 20
		highlighted := highlightCode(content)

		if x == 0 {
			codeContent.WriteString(fmt.Sprintf(`            <text y="%d" class="base">%s</text>`+"\n", y, highlighted))
		} else {
			codeContent.WriteString(fmt.Sprintf(`            <text x="%d" y="%d" class="base">%s</text>`+"\n", x, y, highlighted))
		}
	}

	groupStartTag := `<g transform="translate(20, 40)">`
	startIdx := strings.Index(svgStr, groupStartTag)
	if startIdx != -1 {
		contentStart := startIdx + len(groupStartTag)
		endIdx := strings.Index(svgStr[contentStart:], "</g>")
		if endIdx != -1 {
			endIdx += contentStart
			// This line effectively "removes" the boilerplate code content from the template
			// by slicing the string and injecting the newly generated code content.
			svgStr = svgStr[:contentStart] + "\n" + codeContent.String() + "        " + svgStr[endIdx:]
		}
	}

	return svgStr, nil
}
