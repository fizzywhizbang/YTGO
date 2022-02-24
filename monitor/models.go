package monitor

type YTCM struct {
	Db_name     string `json:"db_name"`
	Db_user     string `json:"db_user"`
	Db_password string `json:"db_password"`
	Db_host     string `json:"db_host"`
	Db_port     string `json:"db_port"`
	BaseDL      string `json:"base_download_dir"`
	Defbrowser  string `json:"defbrowser"`
	JD          string `json:"jd"`
	FolderWatch string `json:"folderwatch"`
}

type Status struct {
	ID   int
	Name string
}

type Tags struct {
	ID   int
	Name string
}

type Search struct {
	ID   int
	Name string
	Link string
}

type Video struct {
	id           int
	yt_videoid   string
	title        string
	description  string
	publisher    string
	publish_date int
	watched      int
}

type Channel struct {
	id              int
	displayname     string
	dldir           string
	yt_channelid    string
	lastpub         int
	lastcheck       int
	archive         int
	notes           string
	date_added      int
	last_feed_count int
}
