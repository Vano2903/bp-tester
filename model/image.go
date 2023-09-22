package model

type Image struct {
	Name        string `json:"name"`
	ImageID     string `json:"imageId"`
	BuildOutput []byte `json:"buildOutput"`
}
