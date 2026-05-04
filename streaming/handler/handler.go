package handler

import "github.com/gorilla/mux"

func CreateHandler() *Handler {
	mux := mux.NewRouter()
	handler := &Handler{
		Handler: mux,
	}

	// 서버 상태 확인용 API
	mux.HandleFunc("/ping", handler.pingHandler).Methods("GET")

	// 특정 드론 스트림의 HLS 재생 주소를 반환
	mux.HandleFunc("/api/v1/streams/hls", handler.hlsStreamHandler).Methods("POST")
	// 특정 드론 스트림의 현재 화면을 캡쳐해서 base64로 반환
	mux.HandleFunc("/api/v1/streams/capture", handler.captureStreamHandler).Methods("POST")
	// 현재 RTMP로 들어오고 있는 모든 드론 스트림 목록을 반환
	mux.HandleFunc("/api/v1/streams/live", handler.liveStreamsHandler).Methods("GET")
	// 특정 드론의 스트리밍 상태를 반환
	mux.HandleFunc("/api/v1/streams/{group}/{code}", handler.streamStatusHandler).Methods("GET")

	// 현재 스트리밍 중인 드론 목록을 웹 화면으로 보여줌
	mux.HandleFunc("/admin/streams", handler.adminStreamsHandler).Methods("GET")
	// 선택한 드론의 HLS 영상을 웹에서 바로 재생하는 테스트 페이지
	mux.HandleFunc("/admin/streams/{group}/{code}", handler.adminStreamPlayerHandler).Methods("GET")

	// 브라우저 또는 HLS 플레이어가 직접 접근하는 HLS 파일 경로
	mux.PathPrefix("/hls/").HandlerFunc(handler.hlsFileHandler).Methods("GET")

	return handler
}
