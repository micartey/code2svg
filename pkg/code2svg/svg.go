package code2svg

import (
	_ "embed"
	"fmt"
	"regexp"
	"strings"
)

//go:embed code_preview.svg
var svgTemplate string

func GenerateSVG(code string) (string, error) {
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
		indent, content := CalculateIndent(line)
		lineWidth := 40.0 + float64(indent*20) + (float64(len(content)) * charWidth)
		if lineWidth > maxWidth {
			maxWidth = lineWidth
		}
	}
	totalWidth := int(maxWidth)
	if totalWidth > 800 {
		totalWidth += 20
	}

	svgStr := svgTemplate

	reSvgOpen := regexp.MustCompile(`<svg width="\d+" height="\d+" viewBox="0 0 \d+ \d+"`)
	svgStr = reSvgOpen.ReplaceAllString(svgStr, fmt.Sprintf(`<svg width="%d" height="%d" viewBox="0 0 %d %d"`, totalWidth, totalHeight, totalWidth, totalHeight))

	reBgRect := regexp.MustCompile(`<rect width="\d+" height="\d+" rx="8" class="bg"/>`)
	svgStr = reBgRect.ReplaceAllString(svgStr, fmt.Sprintf(`<rect width="%d" height="%d" rx="8" class="bg"/>`, totalWidth, totalHeight))

	var codeContent strings.Builder
	for i, line := range lines {
		indent, content := CalculateIndent(line)
		x := indent * 20
		y := i * 20
		highlighted := HighlightCode(content)

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
			svgStr = svgStr[:contentStart] + "\n" + codeContent.String() + "        " + svgStr[endIdx:]
		}
	}

	return svgStr, nil
}
