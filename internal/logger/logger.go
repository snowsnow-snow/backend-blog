// log 输出的文件夹为当天的日期，单个 log 文件的大小由 common.GlobalConfig.Log.MaxSize 决定
// log 文件名称由 common.GlobalConfig.Log.name 决定，每个文件会有编号
// 编号是指当天生成的第几个 log 文件，默认从 0 开始

package logger

import (
	"backend-blog/config"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	yearMonthDay              = ""
	currLogFileName           = ""
	currLogFolder             = ""
	currLogFilePathName       = ""
	logFileMaxSize            = config.GlobalConfig.Log.MaxSize * 1024 * 1024 // log 文件大小，单位为 byte
	currLogFileSize     int64 = 0                                             // 记录当前写入的 log 大小，单位为 byte
	Info                *log.Logger
	Error               *log.Logger
	Warn                *log.Logger
	before              = "->"
	after               = "<-"
)

func init() {
	if currYearMonthDay := time.Now().Format("2006-01-02"); yearMonthDay != currYearMonthDay {
		currLogFolder = createLogFolder(config.GlobalConfig.Log.Path)
		file := createLogFile(config.GlobalConfig.Log.Name)
		mountLog(file)
	}
}

// createLogFolder 当系统第一次启动，或者距离上次创建日志文件夹的日期已发生变化
// 就会新建一个新的以当前日期为名称的文件夹，后续的日志保存在此目录
func createLogFolder(logPath string) string {
	if folderNotExist(logPath) {
		if err := createFolder(logPath); err != nil {
			log.Fatalf("Error Create LogPath: %v\n", err)
		}
	}
	// 获取当前年月日
	yearMonthDay = time.Now().Format("2006-01-02")
	var currLogFolder string
	var createFolderErr error
	currLogFolder = filepath.Join(logPath, yearMonthDay)
	createFolderErr = createFolder(currLogFolder)
	if createFolderErr != nil {
		log.Fatalf("Error Create File: %v\n", createFolderErr)
	}
	return currLogFolder
}

// createLogFile 创建 log 文件，log 输出的内容会被记录在此文件中
func createLogFile(logFilename string) *os.File {
	files, _ := os.ReadDir(currLogFolder)
	logFile, maxSort := getLogFileInfo(files)
	// logFile 不为空
	if logFile != nil && logFile.Size() > logFileMaxSize {
		currLogFileSize = logFile.Size()
		maxSort++
	}
	currLogFileName = logFilename + "_" + strconv.FormatInt(maxSort, 10)
	currLogFileSize = 0
	currLogFilePathName = filepath.Join(currLogFolder, currLogFileName+".log")
	file, openFileErr := os.OpenFile(currLogFilePathName, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModeAppend|os.ModePerm)
	if openFileErr != nil {
		log.Fatalf("Error Open File: %v\n", openFileErr)
	}
	return file
}
func mountLog(file *os.File) {
	Info = log.New(io.MultiWriter(file, os.Stderr), "INFO ", log.Ldate|log.Ltime)
	Error = log.New(io.MultiWriter(file, os.Stderr), "ERROR ", log.Lshortfile|log.Ldate|log.Ltime)
	Warn = log.New(io.MultiWriter(file, os.Stderr), "WARN ", log.Lshortfile|log.Ldate|log.Ltime)
}

// getLogFileInfo 获取当前日志文件夹下最新的一个文件并返回文件信息和编号
// 根据文件编号确定最新的文件，如果当前文件夹下没有日志文件，返回 -1
func getLogFileInfo(files []os.DirEntry) (fs.FileInfo, int64) {
	var maxSort int64 = -1
	var maxSortFile os.DirEntry
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileNameAll := path.Base(file.Name())
		fileNameFix := fileNameAll[0 : len(fileNameAll)-len(path.Ext(fileNameAll))]
		split := strings.Split(fileNameFix, "_")
		if len(split) > 1 {
			sort, _ := strconv.Atoi(split[len(split)-1])
			currSort := int64(sort)
			if currSort > maxSort {
				maxSortFile = file
				maxSort = currSort
			}
		}
	}
	if maxSort == -1 {
		maxSort = 0
	}
	if maxSortFile != nil {
		info, _ := maxSortFile.Info()
		return info, maxSort

	}
	return nil, maxSort
}

// RequestBefore 输出 log, log 格式为: Time Method - Path - IP - StatusCode
func RequestBefore(c *fiber.Ctx) error {
	if currYearMonthDay := time.Now().Format("2006-01-02"); yearMonthDay != currYearMonthDay {
		currLogFolder = createLogFolder(config.GlobalConfig.Log.Path)
		file := createLogFile(config.GlobalConfig.Log.Name)
		mountLog(file)
	}
	if currLogFileSize > logFileMaxSize {
		file := createLogFile(config.GlobalConfig.Log.Name)
		mountLog(file)
	}
	requestID := c.Get("X-Request-Id", uuid.NewString())
	// 在此处使用请求ID进行日志记录等操作
	c.Set("X-Request-Id", requestID)
	// 将请求ID存储在ctx对象中，以便后续使用
	c.Locals("Request-ID", requestID)

	Info.Println(before)
	// 将日志记录到文件中
	printLog := c.Method() + " - " + c.Path() + " - " + c.IP() + " - " + requestID
	Info.Println(printLog)
	if c.Method() != fiber.MethodGet {
		if c.Get(fiber.HeaderContentType) == fiber.MIMEApplicationJSON {
			Info.Println("Request body: ", "\n", string(c.Request().Body()))
		}
	}
	return c.Next()
}
func RequestAfter(c *fiber.Ctx) error {
	if err := c.Next(); err != nil {
		return err
	}
	// 将日志记录到文件中
	printLog := "Response status code: " + strconv.Itoa(c.Response().StatusCode())
	if isFile, ok := c.Locals("isFile").(bool); isFile && ok {
		Info.Println(printLog, "\n", "                         File size: (byte)", len(c.Response().Body()))
	} else {
		Info.Println(printLog, "\n", "                         Result: ", string(c.Response().Body()))
	}
	Info.Println(after)
	ChangeCurrLogFileSize()
	return nil
}
func ChangeCurrLogFileSize() {
	file, err := os.Stat(currLogFilePathName)
	if err != nil {
		Error.Println("File size: %v", err)
	}
	currLogFileSize = file.Size()
}
func createFolder(folderPath string) error {
	var err error
	if folderNotExist(folderPath) {
		err = os.Mkdir(folderPath, 0755)
	}
	return err
}
func folderNotExist(folderPath string) bool {
	_, err := os.Stat(folderPath)
	return os.IsNotExist(err)
}
