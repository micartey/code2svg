package server

import (
	"fmt"
	"io"
	"net/http"

	"micartey.dev/code2svg/pkg/code2svg"
)

func HandleSVG(w http.ResponseWriter, r *http.Request) {
	codeBase64 := r.URL.Query().Get("code")
	if codeBase64 == "" {
		body, _ := io.ReadAll(r.Body)
		codeBase64 = string(body)
	}

	if codeBase64 == "" {
		http.Error(w, "Missing code parameter", http.StatusBadRequest)
		return
	}

	decoded, err := code2svg.DecodeBase64(codeBase64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid base64 string: %v", err), http.StatusBadRequest)
		return
	}

	svg, err := code2svg.GenerateSVG(string(decoded))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write([]byte(svg))
}
