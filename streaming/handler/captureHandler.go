package handler

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/ahr-i/aero-watch/streaming/setting"
)

func (h *Handler) captureStreamHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req streamRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rend.JSON(w, http.StatusBadRequest, nil)
		return
	}

	if !isValidStreamPart(req.Group) || !isValidStreamPart(req.Code) {
		rend.JSON(w, http.StatusBadRequest, nil)
		return
	}

	if _, exists := getStream(req.Group, req.Code); !exists {
		rend.JSON(w, http.StatusNotFound, nil)
		return
	}

	capturedAt := time.Now().UTC()
	imagePath, err := captureStreamFrame(req.Group, req.Code, capturedAt)
	if err != nil {
		rend.JSON(w, http.StatusInternalServerError, nil)
		return
	}

	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		rend.JSON(w, http.StatusInternalServerError, nil)
		return
	}

	rend.JSON(w, http.StatusOK, streamCaptureResponseBody{
		Group:       req.Group,
		Code:        req.Code,
		Status:      "okay",
		ImageBase64: base64.StdEncoding.EncodeToString(imageData),
		ImageType:   "image/jpeg",
	})
}

func captureStreamFrame(group string, code string, capturedAt time.Time) (string, error) {
	key := streamKey(group, code)
	streamDir := hlsStreamDir(group, code)
	segmentPath, err := latestHLSSegmentPath(streamDir)
	if err != nil {
		return "", err
	}

	captureRoot := setting.Setting.CaptureRoot
	if captureRoot == "" {
		captureRoot = "./captures"
	}

	captureDir := filepath.Join(captureRoot, key)
	if err := os.MkdirAll(captureDir, 0755); err != nil {
		return "", err
	}

	imagePath := filepath.Join(captureDir, capturedAt.Format("20060102T150405Z")+".jpg")
	ffmpegPath := setting.Setting.FFmpegPath
	if ffmpegPath == "" {
		ffmpegPath = "ffmpeg"
	}

	cmd := exec.Command(
		ffmpegPath,
		"-hide_banner",
		"-loglevel", "error",
		"-y",
		"-i", segmentPath,
		"-frames:v", "1",
		"-q:v", "2",
		imagePath,
	)

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return imagePath, nil
}

func latestHLSSegmentPath(streamDir string) (string, error) {
	indexPath := filepath.Join(streamDir, "index.m3u8")
	file, err := os.Open(indexPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	latestSegment := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		latestSegment = line
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	if latestSegment == "" {
		return "", os.ErrNotExist
	}

	if filepath.IsAbs(latestSegment) {
		return latestSegment, nil
	}

	return filepath.Join(streamDir, latestSegment), nil
}
