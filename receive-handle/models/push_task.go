package models

type FileInfo struct {
	FileId    string `json:"id"`
	FileUrl   string `json:"url"`
	Host      string `json:"host"`
	FileType  string `json:"type"`
	FileSize  uint64 `json:"fsize"`
	PieceSize uint32 `json:"psize"`
	PPC       uint32 `json:"ppc"`
	CPPC      uint32 `json:"cppc"`
	Expire    int64  `json:"expire"`
}

type FileTask struct {
	FileInfo
	Operation string `json:"operation"`
}

type DistTask struct {
	FileId  string `json:"fid"`
	DistNum int    `json:"dist_num"`
}
