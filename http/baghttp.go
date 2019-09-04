package baghttp

import (
	"net/http"
	"sync/atomic"
)

func NewInstrumentedResponseWriter(w http.ResponseWriter) (http.ResponseWriter, *InstrumentedResponseWriter) {
	if closeNotifier, ok := w.(http.CloseNotifier); ok {
		iwcn := &instrumentedResponseWriterCloseNotify{
			InstrumentedResponseWriter: &InstrumentedResponseWriter{
				ResponseWriter: w,
			},
			CloseNotifier: closeNotifier,
		}

		return iwcn, iwcn.InstrumentedResponseWriter
	}

	iw := &InstrumentedResponseWriter{
		ResponseWriter: w,
	}

	return iw, iw
}

type instrumentedResponseWriterCloseNotify struct {
	*InstrumentedResponseWriter
	http.CloseNotifier
}

type InstrumentedResponseWriter struct {
	http.ResponseWriter

	status  int32
	written int64
}

func (iw *InstrumentedResponseWriter) Write(p []byte) (int, error) {
	n, err := iw.ResponseWriter.Write(p)

	atomic.AddInt64(&iw.written, int64(n))
	if atomic.LoadInt32(&iw.status) == 0 {
		atomic.StoreInt32(&iw.status, http.StatusOK)
	}

	return n, err
}

func (iw *InstrumentedResponseWriter) Info() ResponseInstrumentationInfo {
	return ResponseInstrumentationInfo{
		Written:     atomic.LoadInt64(&iw.written),
		Status:      atomic.LoadInt32(&iw.status),
		ContentType: iw.Header().Get("Content-Type"),
	}
}

func (iw *InstrumentedResponseWriter) WriteHeader(status int) {
	iw.ResponseWriter.WriteHeader(status)
	atomic.StoreInt32(&iw.status, int32(status))
}

func (iw *InstrumentedResponseWriter) Flush() {
	if flusher, ok := iw.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

type ResponseInstrumentationInfo struct {
	Written     int64
	Status      int32
	ContentType string
}
