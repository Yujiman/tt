package worker

func (h *Handler) echo() {
	h.logger.Debugw("Echo worker")
}
