package controllers

import (
	"MyApp/global"
	"MyApp/models"
	"fmt"
	"io/ioutil"
	"net/http"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"crypto/md5"
    "encoding/hex"
	"log"
	"github.com/gin-gonic/gin"
)


// 初始化日志文件
func init() {
    logFile, err := os.OpenFile("video_upload.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatalf("Failed to open log file: %v", err)
    }
    log.SetOutput(logFile)
}

// 接收前端日志数据的接口
func ReceiveFrontendLog(ctx *gin.Context) {
    var logData struct {
        Log string `json:"log"`
    }
    if err := ctx.ShouldBindJSON(&logData); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid log data"})
        return
    }
    log.Printf("Frontend Log: %s", logData.Log)
    ctx.JSON(http.StatusOK, gin.H{"message": "Log received successfully"})
}


func UploadVideoL(ctx *gin.Context) {
	// 移除 ctx.ShouldBind(&video)

	file, err := ctx.FormFile("video")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No video uploaded"})
		return
	}

	videoName := file.Filename
	videoSize := file.Size

	// 记录上传时间
	_ = time.Now()

	// 自动迁移模型
	if err := global.DB.AutoMigrate(&models.Video{}); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error migrating video model"})
		return
	}

	// 创建视频对象
	video := models.Video{
		Name:  videoName,
		Size:  videoSize,
		Title: videoName,
	}

	// 保存视频信息到数据库
	if err := global.DB.Create(&video).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating video"})
		return
	}

	// 创建上传目录（如果不存在）
	if err := os.MkdirAll("D:/NewBack/scripts/4C2025/video_action_cls/video_test", os.ModePerm); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating uploads directory"})
		return
	}
	 // 获取上传的文件
	 file, err = ctx.FormFile("video")
	 if err != nil {
		 ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting video from request"})
		 return
	 }

	 // 新文件名
	 newFileName := "test.mp4"
	 newFilePath := filepath.Join("D:/NewBack/scripts/4C2025/video_action_cls/video_test", newFileName)
 
	 // 打开上传的文件
	 src, err := file.Open()
	 if err != nil {
		 ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error opening uploaded video"})
		 return
	 }
	 defer src.Close()

	 // 创建新文件
	 dst, err := os.Create(newFilePath)
	 if err != nil {
		 ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating new video file"})
		 return
	 }
	 defer dst.Close()

	 // 复制文件内容
	 if _, err := io.Copy(dst, src); err != nil {
		 ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error copying video file"})
		 return
	 }

	 // 更新数据库中的视频名
	 video.Name = newFileName
	 if err := global.DB.Save(&video).Error; err != nil {
		 ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating video name in database"})
		 return
	 }

 	fmt.Println("Processing video: ", videoName)

	 // 调用Python脚本处理视频
	pythonPath := "D:/ForAnaconda/python.exe"
	scriptPath := "./scripts/4C2025/video_action_cls/test_fuse_final.py"
	//moduleDir := "D:/NewBack/scripts/4C2025/video_action_cls/ByteTrack/exps/example/mot"

	// 设置环境变量
	//env := os.Environ()
	//env = append(env, fmt.Sprintf("PYTHONPATH=%s", moduleDir))

	processVideo := exec.Command(pythonPath, scriptPath)
	
	//processVideo.Env = env
	 // 获取标准错误输出
	 
	 stderr, err := processVideo.StderrPipe()
	 if err != nil {
		 ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting stderr pipe"})
		 return
	 }
    
	 // 启动命令
	 if err := processVideo.Start(); err != nil {
		 ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error starting video processing script"})
		 return
	 }

	 // 读取标准错误输出
	 stderrBytes, err := ioutil.ReadAll(stderr)
	 if err != nil {
		 ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading stderr"})
		 return
	 }
 
	 // 等待命令执行完成
	 if err := processVideo.Wait(); err != nil {
		 ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error processing video: %s", string(stderrBytes))})
		 return
	 }





	 // 记录日志
	 // Declare and initialize filePairs
	 filePairs := []FilePair{}

	 // Define uploadTime as the current time
	 uploadTime := time.Now()

	 for _, pair := range filePairs {
		log.Printf("Video Name: %s, Upload Time: %s, Text Summary: %s", videoName, uploadTime.Format(time.RFC3339), pair.Text)
	}



 	fmt.Println("Process video successfully: ", videoName)
	 ctx.JSON(http.StatusOK, gin.H{
		 "message":     "Video uploaded and processed successfully",
		 "video":       newFileName,
		 "size":        video.Size,
		 "uploaded_at": time.Now(),
		 "title":       video.Title,
	 })
}

type FilePair struct {
    VideoURL string `json:"video_url"`
    Text     string `json:"text"`
    MD5Hash  string `json:"md5_hash"`
}

// 获取匹配的视频和文本文件对
func GetMatchedFiles(c *gin.Context) {
	
    var filePairs []FilePair
    videoFiles := make(map[string]string)
    textFiles := make(map[string]string)

    // 遍历textoutput文件夹
    err := filepath.Walk("D:/NewBack/scripts/4C2025/video_action_cls/video_test_output/shear_test", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }

        // 分离文件名和扩展名
        fileName := info.Name()
        ext := filepath.Ext(fileName)

        if ext == ".mp4" {
            videoFiles[fileName[:len(fileName)-4]] = "/videos/" + fileName // 返回相对应的网络 URL
        } else if ext == ".txt" {
            textFiles[fileName[:len(fileName)-4]] = path
        }
        return nil
    })

    if err != nil {
        c.JSON(500, gin.H{"error": "遍历文件夹失败"})
        return
    }


	
    // 匹配视频和文本文件
    for key, videoPath := range videoFiles {
        if textPath, ok := textFiles[key]; ok {
            textContent, err := os.ReadFile(textPath)
            if err != nil {
                c.JSON(500, gin.H{"error": fmt.Sprintf("读取文本文件 %s 失败", textPath)})
                return
            }
            // 确保文本内容不超过100字
            if len(textContent) > 100 {
                textContent = textContent[:100]
            }

			originalVideoPath := "D:/NewBack/scripts/4C2025/video_action_cls/video_test_output/shear_test/" + key + ".mp4"
            convertedVideoPath := "D:/NewBack/scripts/4C2025/video_action_cls/video_test_output/shear_test/" + key + "_h264.mp4"

			fmt.Println("Converting video: ", originalVideoPath)
            // 使用FFmpeg进行编码转换
			ffmpegCmd := exec.Command("ffmpeg", "-i", originalVideoPath, "-c:v", "libx264", "-preset", "medium", "-crf", "23", "-c:a", "aac", "-b:a", "128k", convertedVideoPath)
			if err := ffmpegCmd.Run(); err != nil {
				c.JSON(500, gin.H{"error": fmt.Sprintf("FFmpeg command failed: %v", err)})
				return
			}

            // 调用 calculateFileMD5 函数计算视频文件的 MD5 哈希值
            md5Hash, err := calculateFileMD5("D:/NewBack/scripts/4C2025/video_action_cls/video_test_output/shear_test/" + key + ".mp4")
            if err != nil {
                c.JSON(500, gin.H{"error": fmt.Sprintf("计算视频文件 %s 的 MD5 哈希值失败", videoPath)})
                return
            }
            
			
			
			filePair := FilePair{
                VideoURL: "/videos/" + filepath.Base(convertedVideoPath),
                Text:     string(textContent),
                MD5Hash:  md5Hash,
            }
            filePairs = append(filePairs, filePair)
        }
    }


// 设置完整的响应头
c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
c.Writer.Header().Set("Pragma", "no-cache")
c.Writer.Header().Set("Expires", "0")
c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")


    c.JSON(200, gin.H{"file_pairs": filePairs})
}


func calculateFileMD5(filePath string) (string, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return "", err
    }
    defer file.Close()

    hash := md5.New()
    if _, err := io.Copy(hash, file); err != nil {
        return "", err
    }

    return hex.EncodeToString(hash.Sum(nil)), nil
}

// GetProcessedVideo 函数用于返回处理后的 text.mp4 视频文件
func GetProcessedVideo(c *gin.Context) {
    // 定义视频文件的路径
    videoPath := "D:/NewBack/scripts/4C2025/video_action_cls/video_test_output/test1.mp4"

    // 检查文件是否存在
    if _, err := os.Stat(videoPath); os.IsNotExist(err) {
        c.JSON(http.StatusNotFound, gin.H{"error": "视频文件未找到"})
        return
    }

	md5Hash, err := calculateFileMD5(videoPath)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "无法计算文件的 MD5 哈希值"})
        return
    }

    // 设置响应头，包含 MD5 哈希值
    c.Writer.Header().Set("Content-MD5", md5Hash)

    // 打开视频文件
    file, err := os.Open(videoPath)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "无法打开视频文件"})
        return
    }
    defer file.Close()

    // 获取文件信息
    fileInfo, err := file.Stat()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取文件信息"})
        return
    }

    // 设置响应头
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", fileInfo.Name()))
	c.Writer.Header().Set("Content-Type", "video/mp4")
	c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	c.Writer.Header().Set("Accept-Ranges", "bytes")

    // 将文件内容写入响应体
	http.ServeContent(c.Writer, c.Request, fileInfo.Name(), fileInfo.ModTime(), file)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "无法发送视频文件"})
        return
    }
}