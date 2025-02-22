package enum

type DownloadStatus string

const (
	DownloadCompleted  DownloadStatus = "COMPLETED"
	DownloadError      DownloadStatus = "ERROR"
	DownloadNotStarted DownloadStatus = "NOT_STARTED"
)

func (s DownloadStatus) String() string {
	return string(s)
}
