package models

import "time"

type PatchRequestPayload struct {
	Id         string       `json:"id"`
	Schemas    []string     `json:"schemas"`
	Operations []Operations `json:"operations"`
}

type Operations struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

type Meta struct {
	ResourceType string    `json:"resourceType"`
	Created      time.Time `json:"created"`
	LastModified time.Time `json:"lastModified"`
	Version      string    `json:"version"`
}

