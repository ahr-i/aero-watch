package handler

func (h *Handler) Close() {
	h.gpsStore.Close()
}
