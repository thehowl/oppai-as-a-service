package main

type baseResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
}

var err500 = baseResponse{false, "An error occurred"}
