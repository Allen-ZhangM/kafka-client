package models

const (
	UrlPath_Register       = "register"
	UrlPath_GetTasks       = "get_tasks"
	UrlPath_PrefetchReport = "prefetch_report"
)

const (
	GetTasksLimit = 10
)

const (
	PushTask_File = "FileTask"
	PushTask_Dist = "DistTask"
)
const (
	PushFile_Big  = "Big"
	PushFile_M3U8 = "M3U8"
)
const (
	FileOperation_Download = "Download"
	FileOperation_Update   = "Update"
	FileOperation_Delete   = "Delete"
)

type PrefetchStatus string

const (
	Prefetch_Init        PrefetchStatus = "init"
	Prefetch_Probed                     = "probed"
	Prefetch_ToPush                     = "topush"
	Prefetch_Started                    = "started"
	Prefetch_Downloading                = "downloading"
	Prefetch_Finished                   = "finished"
	Prefetch_Failed                     = "failed"
	Prefetch_Disable                    = "disable"
)

type RegisterResp struct {
	IP string `json:"ip"`
}

type PushTasksResp struct {
	Error string     `json:"error"`
	Tasks []PushTask `json:"tasks"`
}

type PushTask struct {
	PushID string `json:"push_id"`
	Type   string `json:"type"`
	Data   string `json:"data"`
}

type PrefetchReport struct {
	FileId string         `json:"fid"`
	Status PrefetchStatus `json:"status"`
}
