package handler

import (
	"bytes"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/ahr-i/aero-watch/streaming/setting"
	"github.com/ahr-i/aero-watch/streaming/utils/logging"
	pkgerrors "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/yutopp/go-flv"
	flvtag "github.com/yutopp/go-flv/tag"
	"github.com/yutopp/go-rtmp"
	rtmpmsg "github.com/yutopp/go-rtmp/message"
)

var _ rtmp.Handler = (*rtmpStreamHandler)(nil)

type rtmpStreamHandler struct {
	rtmp.DefaultHandler

	app   string
	group string
	code  string

	mu     sync.Mutex
	ffmpeg *exec.Cmd
	stdin  io.WriteCloser
	flvEnc *flv.Encoder
}

func StartRTMPServer() {
	rtmpPort := setting.Setting.RTMPPort
	if rtmpPort == "" {
		logging.Warn("RTMP server skipped. rtmp_port is empty.")
		return
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":"+rtmpPort)
	if err != nil {
		logging.Error(err)
		return
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		logging.Error(err)
		return
	}

	server := rtmp.NewServer(&rtmp.ServerConfig{
		OnConnect: func(conn net.Conn) (io.ReadWriteCloser, *rtmp.ConnConfig) {
			return conn, &rtmp.ConnConfig{
				Handler: &rtmpStreamHandler{},
				ControlState: rtmp.StreamControlStateConfig{
					DefaultBandwidthWindowSize: 6 * 1024 * 1024 / 8,
				},
				Logger: newRTMPLogger(),
			}
		},
	})

	logging.Info("RTMP server start.")
	if err := server.Serve(listener); err != nil {
		logging.Error(err)
	}
}

func (h *rtmpStreamHandler) OnConnect(_ uint32, cmd *rtmpmsg.NetConnectionConnect) error {
	h.app = strings.Trim(cmd.Command.App, "/")
	return nil
}

func (h *rtmpStreamHandler) OnPublish(_ *rtmp.StreamContext, timestamp uint32, cmd *rtmpmsg.NetStreamPublish) error {
	group, code, ok := parseRTMPPublishPath(h.app, cmd.PublishingName)
	if !ok || !isValidStreamPart(group) || !isValidStreamPart(code) {
		logging.Warn("RTMP publish rejected. expected path: /live/{group}/{code}, received path: " + rtmpPublishPath(h.app, cmd.PublishingName))
		return stderrors.New("invalid stream path")
	}

	if !validateDrone(group, code) {
		logging.Warn("RTMP publish rejected. unauthorized drone: " + streamKey(group, code))
		return stderrors.New("unauthorized drone")
	}

	if err := h.startHLSWriter(group, code); err != nil {
		return err
	}

	h.group = group
	h.code = code
	markStreamLive(group, code)

	logging.Info("RTMP publish accepted: " + streamKey(group, code))
	return nil
}

func rtmpPublishPath(app string, publishingName string) string {
	app = strings.Trim(app, "/")
	publishingName = strings.Trim(publishingName, "/")
	if publishingName == "" {
		return "/" + app
	}

	return "/" + app + "/" + publishingName
}

func parseRTMPPublishPath(app string, publishingName string) (string, string, bool) {
	appParts := strings.Split(strings.Trim(app, "/"), "/")
	if len(appParts) == 0 || appParts[0] != "live" {
		return "", "", false
	}

	publishingName = strings.Trim(publishingName, "/")
	if len(appParts) == 1 {
		return splitStreamPath(publishingName)
	}

	if len(appParts) == 2 && !strings.Contains(publishingName, "/") {
		return appParts[1], publishingName, true
	}

	return "", "", false
}

func (h *rtmpStreamHandler) OnSetDataFrame(timestamp uint32, data *rtmpmsg.NetStreamSetDataFrame) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.flvEnc == nil {
		return nil
	}

	reader := bytes.NewReader(data.Payload)
	var script flvtag.ScriptData
	if err := flvtag.DecodeScriptData(reader, &script); err != nil {
		return nil
	}

	return h.flvEnc.Encode(&flvtag.FlvTag{
		TagType:   flvtag.TagTypeScriptData,
		Timestamp: timestamp,
		Data:      &script,
	})
}

func (h *rtmpStreamHandler) OnAudio(timestamp uint32, payload io.Reader) error {
	var audio flvtag.AudioData
	if err := flvtag.DecodeAudioData(payload, &audio); err != nil {
		return err
	}

	flvBody := new(bytes.Buffer)
	if _, err := io.Copy(flvBody, audio.Data); err != nil {
		return err
	}
	audio.Data = flvBody

	h.mu.Lock()
	defer h.mu.Unlock()

	if h.flvEnc == nil {
		return nil
	}

	return h.flvEnc.Encode(&flvtag.FlvTag{
		TagType:   flvtag.TagTypeAudio,
		Timestamp: timestamp,
		Data:      &audio,
	})
}

func (h *rtmpStreamHandler) OnVideo(timestamp uint32, payload io.Reader) error {
	var video flvtag.VideoData
	if err := flvtag.DecodeVideoData(payload, &video); err != nil {
		return err
	}

	flvBody := new(bytes.Buffer)
	if _, err := io.Copy(flvBody, video.Data); err != nil {
		return err
	}
	video.Data = flvBody

	h.mu.Lock()
	defer h.mu.Unlock()

	if h.flvEnc == nil {
		return nil
	}

	return h.flvEnc.Encode(&flvtag.FlvTag{
		TagType:   flvtag.TagTypeVideo,
		Timestamp: timestamp,
		Data:      &video,
	})
}

func (h *rtmpStreamHandler) OnClose() {
	h.mu.Lock()
	stdin := h.stdin
	ffmpeg := h.ffmpeg
	group := h.group
	code := h.code

	h.stdin = nil
	h.flvEnc = nil
	h.ffmpeg = nil
	h.mu.Unlock()

	if stdin != nil {
		_ = stdin.Close()
	}
	if ffmpeg != nil {
		_ = ffmpeg.Wait()
	}
	if group != "" && code != "" {
		markStreamOffline(group, code)
		logging.Info("RTMP publish closed: " + streamKey(group, code))
	}
}

func (h *rtmpStreamHandler) startHLSWriter(group string, code string) error {
	streamDir := hlsStreamDir(group, code)
	if err := os.MkdirAll(streamDir, 0755); err != nil {
		return pkgerrors.Wrap(err, "failed to create hls directory")
	}

	ffmpegPath := setting.Setting.FFmpegPath
	if ffmpegPath == "" {
		ffmpegPath = "ffmpeg"
	}

	segmentPath := filepath.Join(streamDir, "%03d.ts")
	indexPath := filepath.Join(streamDir, "index.m3u8")
	hlsTimeSeconds := setting.Setting.HLSTimeSeconds
	if hlsTimeSeconds <= 0 {
		hlsTimeSeconds = 1
	}
	hlsListSize := setting.Setting.HLSListSize
	if hlsListSize <= 0 {
		hlsListSize = 4
	}

	cmd := exec.Command(
		ffmpegPath,
		"-hide_banner",
		"-loglevel", "error",
		"-fflags", "+genpts",
		"-use_wallclock_as_timestamps", "1",
		"-f", "flv",
		"-i", "pipe:0",
		"-c", "copy",
		"-f", "hls",
		"-hls_time", strconv.Itoa(hlsTimeSeconds),
		"-hls_list_size", strconv.Itoa(hlsListSize),
		"-hls_flags", "delete_segments+omit_endlist+program_date_time",
		"-hls_segment_filename", segmentPath,
		indexPath,
	)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return pkgerrors.Wrap(err, "failed to open ffmpeg stdin")
	}

	if err := cmd.Start(); err != nil {
		_ = stdin.Close()
		return pkgerrors.Wrap(err, "failed to start ffmpeg")
	}

	encoder, err := flv.NewEncoder(stdin, flv.FlagsAudio|flv.FlagsVideo)
	if err != nil {
		_ = stdin.Close()
		_ = cmd.Process.Kill()
		return pkgerrors.Wrap(err, "failed to create flv encoder")
	}

	h.mu.Lock()
	h.ffmpeg = cmd
	h.stdin = stdin
	h.flvEnc = encoder
	h.mu.Unlock()

	return nil
}

func validateDrone(group string, code string) bool {
	if !setting.Setting.DroneValidateEnabled {
		return true
	}

	if setting.Setting.DroneOperationService == "" || setting.Setting.DroneValidatePath == "" {
		return false
	}

	validateURL, err := url.Parse(strings.TrimRight(setting.Setting.DroneOperationService, "/"))
	if err != nil {
		return false
	}

	validateURL.Path = path.Join(validateURL.Path, setting.Setting.DroneValidatePath)

	body := map[string]string{
		"group": group,
		"code":  code,
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		logging.Error(fmt.Sprintf("drone validation marshal failed: %v", err))
		return false
	}

	resp, err := http.Post(validateURL.String(), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logging.Error(fmt.Sprintf("drone validation failed: %v", err))
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices
}

func newRTMPLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	logger.SetLevel(logrus.PanicLevel)
	return logger
}
