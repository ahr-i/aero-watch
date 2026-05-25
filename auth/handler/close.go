package handler

func (h *Handler) Close() {
	if h.store != nil {
		h.store.Close()
	}
}
