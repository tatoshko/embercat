package core

type Board struct {
    Threads []Thread `json:"threads"`
}

type Thread struct {
    Subject string `json:"subject"`
    Comment string `json:"comment"`
    Num     string `json:"num"`
    Posts   []Post `json:"posts"`
}

type Post struct {
    Files []File `json:"files`
}

type File struct {
    FullName  string   `json:"fullname"`
    Md5       string   `json:"md5"`
    Path      string   `json:"path"`
    Thumbnail string   `json:"thumbnail"`
    Type      FileType `json:"type"`
}

type FileType int

const (
    t0 FileType = iota
    t1
    t2
    t3
    t4
    t5
    WEBM
    t7
    t8
    t9
    MP4
)
