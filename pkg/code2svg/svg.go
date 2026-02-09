package code2svg

import (
	_ "embed"
	"fmt"
	"regexp"
	"strings"
)

//go:embed code_preview.svg
var svgTemplate string

func GenerateSVG(code string, transparent bool) (string, error) {
	code = strings.ReplaceAll(code, "\r\n", "\n")
	lines := strings.Split(code, "\n")
	lineCount := len(lines)
	if lineCount == 0 {
		lineCount = 1
	}

	totalHeight := 80 + 20*lineCount

	const charWidth = 8.5
	maxWidth := 800.0
	for _, line := range lines {
		displayLine := strings.ReplaceAll(line, "\t", "    ")
		lineWidth := 40.0 + (float64(len(displayLine)) * charWidth)
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
	if transparent {
		svgStr = reBgRect.ReplaceAllString(svgStr, "")
	} else {
		svgStr = reBgRect.ReplaceAllString(svgStr, fmt.Sprintf(`<rect width="%d" height="%d" rx="8" class="bg"/>`, totalWidth, totalHeight))
	}

	var codeContent strings.Builder
	codeContent.WriteString("            <text class=\"base\" xml:space=\"preserve\">\n")
	for i, line := range lines {
		y := i * 20
		displayLine := strings.ReplaceAll(line, "\t", "    ")
		highlighted := HighlightCode(displayLine)
		codeContent.WriteString(fmt.Sprintf(`                <tspan x="0" y="%d">%s</tspan>`+"\n", y, highlighted))
	}
	codeContent.WriteString("            </text>")

	groupStartTag := `<g transform="translate(20, 40)">`
	startIdx := strings.Index(svgStr, groupStartTag)
	if startIdx != -1 {
		contentStart := startIdx + len(groupStartTag)
		endIdx := strings.Index(svgStr[contentStart:], "</g>")
		if endIdx != -1 {
			endIdx += contentStart
			// Actually placing the code inside the template and removing the current builer plate code
			// This is done by splitting it into a start and ending sequence
			svgStr = svgStr[:contentStart] + "\n" + codeContent.String() + "        " + svgStr[endIdx:]
		}
	}

	return svgStr, nil
}
