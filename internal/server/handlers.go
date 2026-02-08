package server

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strconv"

	"micartey.dev/code2svg/pkg/code2svg"
)

func HandleSVG(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, false)
}

func HandlePNG(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, true)
}

func handleRequest(w http.ResponseWriter, r *http.Request, asPNG bool) {
	codeBase64 := r.URL.Query().Get("code")
	transparent := r.URL.Query().Get("transparent") == "true"

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

	svg, err := code2svg.GenerateSVG(string(decoded), transparent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if asPNG {
		png, err := convertSVGToPNG(svg)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to convert to PNG: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", strconv.Itoa(len(png)))
		w.Write(png)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Content-Length", strconv.Itoa(len(svg)))
	w.Write([]byte(svg))
}

func convertSVGToPNG(svg string) ([]byte, error) {
	cmd := exec.Command("rsvg-convert", "-f", "png", "-b", "none")
	cmd.Stdin = bytes.NewBufferString(svg)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("rsvg-convert error: %v, stderr: %s", err, stderr.String())
	}

	return out.Bytes(), nil
}
