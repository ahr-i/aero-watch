package handler

func (h *Handler) Close() {
	if h.db != nil {
		h.db.Close()
	}
}
