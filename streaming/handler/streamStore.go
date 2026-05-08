package handler

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ahr-i/aero-watch/streaming/setting"
	"github.com/ahr-i/aero-watch/streaming/utils/logging"
)

const (
	streamStatusLive    = "LIVE"
	streamStatusOffline = "OFFLINE"
)

type streamInfo struct {
	Group      string
	Code       string
	StartedAt  time.Time
	LastSeenAt time.Time
}

var streams = struct {
	sync.RWMutex
	items      map[string]streamInfo
	offlineAts map[string]time.Time
}{
	items:      make(map[string]streamInfo),
	offlineAts: make(map[string]time.Time),
}

func streamKey(group string, code string) string {
	return group + "-" + code
}

func splitStreamPath(value string) (string, string, bool) {
	group, code, ok := strings.Cut(value, "/")
	return group, code, ok && group != "" && code != ""
}

func hlsStreamDir(group string, code string) string {
	return filepath.Join(setting.Setting.HLSRoot, group, code)
}

func hlsIndexPath(group string, code string) string {
	return filepath.Join(hlsStreamDir(group, code), "index.m3u8")
}

func markStreamLive(group string, code string) streamInfo {
	key := streamKey(group, code)
	now := time.Now().UTC()

	streams.Lock()
	defer streams.Unlock()

	info, exists := streams.items[key]
	if exists {
		info.LastSeenAt = now
		streams.items[key] = info
		delete(streams.offlineAts, key)
		return info
	}

	info = streamInfo{
		Group:      group,
		Code:       code,
		StartedAt:  now,
		LastSeenAt: now,
	}
	streams.items[key] = info
	delete(streams.offlineAts, key)

	return info
}

func markStreamOffline(group string, code string) {
	streams.Lock()
	defer streams.Unlock()

	key := streamKey(group, code)
	delete(streams.items, key)
	streams.offlineAts[key] = time.Now().UTC()
}

func getStream(group string, code string) (streamInfo, bool) {
	pruneStaleStreams()

	key := streamKey(group, code)

	streams.RLock()
	info, exists := streams.items[key]
	streams.RUnlock()
	if exists {
		return info, true
	}

	return streamFromHLS(group, code)
}

func listStreams() []streamInfo {
	pruneStaleStreams()

	streams.RLock()
	list := make([]streamInfo, 0, len(streams.items))
	for _, info := range streams.items {
		list = append(list, info)
	}
	streams.RUnlock()

	for _, info := range streamsFromHLS() {
		if _, exists := findStream(list, info.Group, info.Code); !exists {
			list = append(list, info)
		}
	}

	sort.Slice(list, func(i int, j int) bool {
		return list[i].StartedAt.After(list[j].StartedAt)
	})

	return list
}

func findStream(list []streamInfo, group string, code string) (streamInfo, bool) {
	for _, info := range list {
		if info.Group == group && info.Code == code {
			return info, true
		}
	}

	return streamInfo{}, false
}

func streamFromHLS(group string, code string) (streamInfo, bool) {
	indexPath := hlsIndexPath(group, code)
	fileInfo, err := os.Stat(indexPath)
	if err != nil || fileInfo.IsDir() {
		return streamInfo{}, false
	}
	if isStale(fileInfo.ModTime().UTC()) {
		return streamInfo{}, false
	}

	return streamInfo{
		Group:      group,
		Code:       code,
		StartedAt:  fileInfo.ModTime().UTC(),
		LastSeenAt: fileInfo.ModTime().UTC(),
	}, true
}

func pruneStaleStreams() {
	streams.Lock()
	defer streams.Unlock()

	for key, info := range streams.items {
		indexPath := hlsIndexPath(info.Group, info.Code)
		fileInfo, err := os.Stat(indexPath)
		if err == nil && !fileInfo.IsDir() {
			lastSeenAt := fileInfo.ModTime().UTC()
			info.LastSeenAt = lastSeenAt
			streams.items[key] = info
		}
	}

	for key, offlineAt := range streams.offlineAts {
		if isStale(offlineAt) {
			delete(streams.offlineAts, key)
			removeHLSDirectory(key)
		}
	}
}

func isStale(lastSeenAt time.Time) bool {
	return time.Since(lastSeenAt) > streamTimeout()
}

func streamTimeout() time.Duration {
	if setting.Setting.StreamTimeoutSeconds <= 0 {
		return 30 * time.Second
	}

	return time.Duration(setting.Setting.StreamTimeoutSeconds) * time.Second
}

func streamsFromHLS() []streamInfo {
	groupEntries, err := os.ReadDir(setting.Setting.HLSRoot)
	if err != nil {
		return nil
	}

	list := make([]streamInfo, 0, len(groupEntries))
	for _, groupEntry := range groupEntries {
		if !groupEntry.IsDir() || !isValidStreamPart(groupEntry.Name()) {
			continue
		}

		codeEntries, err := os.ReadDir(filepath.Join(setting.Setting.HLSRoot, groupEntry.Name()))
		if err != nil {
			continue
		}

		for _, codeEntry := range codeEntries {
			if !codeEntry.IsDir() || !isValidStreamPart(codeEntry.Name()) {
				continue
			}

			info, exists := streamFromHLS(groupEntry.Name(), codeEntry.Name())
			if exists {
				list = append(list, info)
				continue
			}

			indexPath := hlsIndexPath(groupEntry.Name(), codeEntry.Name())
			fileInfo, err := os.Stat(indexPath)
			if err == nil && !fileInfo.IsDir() && isStale(fileInfo.ModTime().UTC()) {
				removeHLSDirectory(streamKey(groupEntry.Name(), codeEntry.Name()))
			}
		}
	}

	return list
}

func removeHLSDirectory(key string) {
	group, code, ok := strings.Cut(key, "-")
	if !ok || !isValidStreamPart(group) || !isValidStreamPart(code) {
		return
	}

	dir := hlsStreamDir(group, code)
	err := os.RemoveAll(dir)
	if err != nil {
		logging.Error(err)
		return
	}

	logging.Info("Removed stale HLS directory: " + dir)
}
