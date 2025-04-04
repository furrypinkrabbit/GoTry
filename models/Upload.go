package models

 

import (

    "gorm.io/gorm"

)

 

type Video struct {

    gorm.Model
    Title    string "binding:\"required\""
    Name     string
    Size     int64
    VideoPath string // 存储视频文件的相对路径或URL
}
