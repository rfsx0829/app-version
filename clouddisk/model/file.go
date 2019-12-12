package model

// File means file
type File struct {
	FID      int    `json:"fid"`
	UID      int    `json:"uid"`
	Filename string `json:"filename"`
	Filepath string `json:"filepath"`
	MD5Value string `json:"md5value"`
}
