
func (resp Response) SetCustomHeaders(w http.ResponseWriter) {
	for key, value := range resp.Headers {
		w.Header().Set(key, value)
	}
}
