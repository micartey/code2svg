package code2svg

import (
	"fmt"
	"html"
	"regexp"
	"sort"
	"strings"
)

const (
	patternComment  = `//.*`
	patternString   = `".*?"`
	patternKeyword  = `\b(fn|let|return|if|else|while|for|match|type|struct|enum|impl|use|mod|pub|crate|async|await|static|mut|const|ref|move|where|dyn|trait|package|import|func|var|range|chan|go|select|defer|interface|map|switch|case|default|break|continue|fallthrough)\b`
	patternFunction = `\b(\w+!?)[\s]*[\(!]`
	patternType     = `\b[A-Z]\w*\b`
	patternVariable = `\b[a-z_]\w*\b`
)

var (
	reComment  = regexp.MustCompile(patternComment)
	reString   = regexp.MustCompile(patternString)
	reKeyword  = regexp.MustCompile(patternKeyword)
	reFunction = regexp.MustCompile(patternFunction)
	reType     = regexp.MustCompile(patternType)
	reVariable = regexp.MustCompile(patternVariable)
)

type replacement struct {
	start, end int
	class      string
	priority   int
}

func HighlightCode(line string) string {
	if line == "" {
		return ""
	}

	// Handle comments first as they override everything else
	if reComment.MatchString(line) {
		idx := reComment.FindStringIndex(line)
		before := line[:idx[0]]
		comment := line[idx[0]:]
		return highlightTokens(before) + fmt.Sprintf(`<tspan class="comment">%s</tspan>`, html.EscapeString(comment))
	}

	return highlightTokens(line)
}

func highlightTokens(line string) string {
	if line == "" {
		return ""
	}

	// Handle strings first recursively
	if reString.MatchString(line) {
		idx := reString.FindStringIndex(line)
		before := line[:idx[0]]
		str := line[idx[0]:idx[1]]
		after := line[idx[1]:]
		return highlightTokens(before) + fmt.Sprintf(`<tspan class="string">%s</tspan>`, html.EscapeString(str)) + highlightTokens(after)
	}

	var replacements []replacement

	// Keywords: 1
	for _, match := range reKeyword.FindAllStringIndex(line, -1) {
		replacements = append(replacements, replacement{match[0], match[1], "keyword", 1})
	}
	// Types: 2
	for _, match := range reType.FindAllStringIndex(line, -1) {
		replacements = append(replacements, replacement{match[0], match[1], "type", 2})
	}
	// Functions: 3
	for _, match := range reFunction.FindAllSubmatchIndex([]byte(line), -1) {
		if len(match) >= 4 {
			replacements = append(replacements, replacement{match[2], match[3], "function", 3})
		}
	}
	// Variables: 4
	for _, match := range reVariable.FindAllStringIndex(line, -1) {
		replacements = append(replacements, replacement{match[0], match[1], "variable", 4})
	}

	if len(replacements) == 0 {
		return html.EscapeString(line)
	}

	// Sort by start position, then priority
	sort.Slice(replacements, func(i, j int) bool {
		if replacements[i].start == replacements[j].start {
			if replacements[i].end == replacements[j].end {
				return replacements[i].priority < replacements[j].priority
			}
			return (replacements[i].end - replacements[i].start) > (replacements[j].end - replacements[j].start)
		}
		return replacements[i].start < replacements[j].start
	})

	// Filter overlaps
	var filtered []replacement
	lastEnd := -1
	for _, r := range replacements {
		if r.start >= lastEnd {
			filtered = append(filtered, r)
			lastEnd = r.end
		}
	}

	var finalResult strings.Builder
	lastPos := 0
	for _, r := range filtered {
		finalResult.WriteString(html.EscapeString(line[lastPos:r.start]))
		finalResult.WriteString(fmt.Sprintf(`<tspan class="%s">%s</tspan>`, r.class, html.EscapeString(line[r.start:r.end])))
		lastPos = r.end
	}
	finalResult.WriteString(html.EscapeString(line[lastPos:]))

	return finalResult.String()
}
