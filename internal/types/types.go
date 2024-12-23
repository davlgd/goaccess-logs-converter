package types

type LogEntry struct {
	Date   string `json:"date"`
	Source struct {
		IP string `json:"ip"`
	} `json:"source"`
	HTTP *struct {
		Request struct {
			Method, Path string `json:"method,path"`
		} `json:"request"`
		Response struct {
			StatusCode int `json:"statusCode"`
		} `json:"response"`
	} `json:"http"`
	BytesOut int `json:"bytesOut"`
}
