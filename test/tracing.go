package main

type TelemetryOptions struct {
	Name     string  `json:"name"`
	Endpoint string  `json:"endpoint"`
	Sampler  float64 `json:"sampler"`
	Batcher  string  `json:"batcher""`
}

func NewTelemetryOptions(name string, endpoint string, sampler float64, batcher string) *TelemetryOptions {
	return &TelemetryOptions{
		Name:     name,
		Endpoint: endpoint,
		Sampler:  sampler,
		Batcher:  batcher,
	}
}
