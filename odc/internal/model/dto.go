package model

type ConnectResponse struct {
	ID string `json:"id"`
}

type ErrorResponse struct {
	Code    string `json:"error_code"`
	Message string `json:"error_msg"`
}

type SetupResponse struct {
	Path     string `json:"path"`
	FileName string `json:"file_name"`
}

type ResultDataRequest struct {
	FileName string `json:"file_name"`
}
