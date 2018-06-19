package util

import (
	"time"
	"sync"
	"os"
	"bytes"
	"fmt"
	"log"
	"encoding/json"
	"math"
	"io"
	"path"
)

const (
	_VERSION_          = "1.0.1"               //版本
	DATEFORMAT         = "2006-01-02"          //日期格式(用于文件命名)
	TIMEFORMAT         = "2006/01/02 15:04:05" //时间格式(日志时间格式)
	_SPACE             = " "                   //参数分割
	_TABLE             = "\t"                  //日志文件行分隔符
	_JOIN              = "&"                   //参数连接符
	_FILE_OPERAT_MODE_ = 0644                  //文件操作权限模式
	_FILE_CREAT_MODE_  = 0666                  //文件建立权限模式
	_LABEL_            = "[_loggor_]"          //标签
)
const (
	//日志文件存储模式
	LOG_FILE_SAVE_USUAL = 1 //普通模式,不分割
	LOG_FILE_SAVE_SIZE  = 2 //大小分割
	LOG_FILE_SAVE_DATE  = 3 //日期分割
)
const (
	//文件大小单位
	_        = iota
	KB int64 = 1 << (iota * 10)
	MB
	GB
	TB
)
const (
	_EXTEN_NAME_               = ".log"                  //日志文件后缀名
	_CHECK_TIME_ time.Duration = 900 * time.Millisecond  //定时检测是否分割检测周期
	_WRITE_TIME_ time.Duration = 1300 * time.Millisecond //定时写入文件周期
)

var (
	IS_DEBUG     = false //调试模式
	TIMEER_WRITE = false //定时写入文件
)

type LOGGER interface {
	SetDebug(bool)                                      //设置日志文件路径及名称
	SetType(uint)                                       //设置日志类型
	SetRollingFile(string, string, int32, int64, int64) //按照文件大小分割
	SetRollingDaily(string, string)                     //按照日期分割
	SetRollingNormal(string, string)                    //设置普通模式
	Close()                                             //关闭
	Println(a ...interface{})                           //打印日志
	Printf(format string, a ...interface{})             //格式化输出
}

//==================================================================日志记录器
type LLogger struct {
	log_type     uint          //日志类型
	path         string        //日志文件路径
	dir          string        //目录
	filename     string        //文件名
	maxFileSize  int64         //文件大小
	maxFileCount int32         //文件个数
	dailyRolling bool          //日分割
	sizeRolling  bool          //大小分割
	nomalRolling bool          //普通模式(不分割)
	_suffix      int           //大小分割文件的当前序号
	_date        *time.Time    //文件时间
	mu_buf       *sync.Mutex   //缓冲锁
	mu_file      *sync.Mutex   //文件锁
	logfile      *os.File      //文件句柄
	timer        *time.Timer   //监视定时器
	writeTimer   *time.Timer   //批量写入定时器
	buf          *bytes.Buffer //缓冲区(公用buf保证数据写入的顺序性)
}

/**获取日志对象**/
func New() *LLogger {
	this := &LLogger{}
	this.buf = &bytes.Buffer{}
	this.mu_buf = new(sync.Mutex)
	this.mu_file = new(sync.Mutex)

	return this
}

/**格式行输出**/
func (this *LLogger) Printf(format string, a ...interface{}) {
	defer func() {
		if !TIMEER_WRITE {
			go this.bufWrite()
		}
	}()

	tp := fmt.Sprintf(format, a...)
	this.mu_buf.Lock()
	defer this.mu_buf.Unlock()
	this.buf.WriteString(
		fmt.Sprintf(
			"%s\t%d\t%s\n",
			time.Now().Format(TIMEFORMAT),
			this.log_type,
			tp,
		),
	)
}

/**逐行输出**/
func (this *LLogger) Println(a ...interface{}) {
	defer func() {
		if !TIMEER_WRITE {
			go this.bufWrite()
		}
	}()

	tp := fmt.Sprint(a...)
	this.mu_buf.Lock()
	defer this.mu_buf.Unlock()
	this.buf.WriteString(
		fmt.Sprintf(
			"%s\t%d\t%s\n",
			time.Now().Format(TIMEFORMAT),
			this.log_type,
			tp,
		),
	)
}

/**测试模式**/
func (this *LLogger) SetDebug(is_debug bool) {
	IS_DEBUG = is_debug
}

/**定时写入**/
func (this *LLogger) SetTimeWrite(time_write bool) *LLogger {
	TIMEER_WRITE = time_write

	return this
}

/**日志类型**/
func (this *LLogger) SetType(tp uint) {
	this.log_type = tp
}

/**大小分割**/
func (this *LLogger) SetRollingFile(dir, _file string, maxn int32, maxs int64, _u int64) {
	//0.输入合法性
	if this.sizeRolling ||
		this.dailyRolling ||
		this.nomalRolling {
		log.Println(_LABEL_, "mode can't be changed!")
		return
	}

	//1.设置各模式标志符
	this.sizeRolling = true
	this.dailyRolling = false
	this.nomalRolling = false

	//2.设置日志器各参数
	this.maxFileCount = maxn
	this.maxFileSize = maxs * int64(_u)
	this.dir = dir
	this.filename = _file
	for i := 1; i <= int(maxn); i++ {
		sizeFile := fmt.Sprintf(
			dir,
			_file,
			_EXTEN_NAME_,
			".",
			fmt.Sprintf("%05d", i),
		)
		if pathIsExist(sizeFile) {
			this._suffix = i
		} else {
			break
		}
	}
	//3.实时文件写入
	this.path = fmt.Sprint(
		dir,
		_file,
		_EXTEN_NAME_,
	)
	this.startLogger(this.path)
}

/**日期分割**/
func (this *LLogger) SetRollingDaily(dir, _file string) {
	//0.输入合法性
	if this.sizeRolling ||
		this.dailyRolling ||
		this.nomalRolling {
		log.Println(_LABEL_, "mode can't be changed!")
		return
	}

	//1.设置各模式标志符
	this.sizeRolling = false
	this.dailyRolling = true
	this.nomalRolling = false

	//2.设置日志器各参数
	this.dir = dir
	this.filename = _file
	this._date = getNowFormDate(DATEFORMAT)
	this.startLogger(
		fmt.Sprint(
			this.dir,
			this.filename,
			_EXTEN_NAME_,
			".",
			this._date.Format(DATEFORMAT),
		),
	)
}

/**普通模式**/
func (this *LLogger) SetRollingNormal(dir, _file string) {
	//0.输入合法性
	if this.sizeRolling ||
		this.dailyRolling ||
		this.nomalRolling {
		log.Println(_LABEL_, "mode can't be changed!")
		return
	}

	//1.设置各模式标志符
	this.sizeRolling = false
	this.dailyRolling = false
	this.nomalRolling = true

	//2.设置日志器各参数
	this.dir = dir
	this.filename = _file
	this.startLogger(
		fmt.Sprint(
			dir,
			_file,
			_EXTEN_NAME_,
		),
	)
}

/**关闭日志器**/
func (this *LLogger) Close() {
	//0.获取锁
	this.mu_buf.Lock()
	defer this.mu_buf.Unlock()
	this.mu_file.Lock()
	defer this.mu_file.Unlock()

	//1.关闭
	if nil != this.timer {
		this.timer.Stop()
	}
	if nil != this.writeTimer {
		this.writeTimer.Stop()
	}
	if this.logfile != nil {
		err := this.logfile.Close()

		if err != nil {
			log.Println(_LABEL_, "file close err", err)
		}
	} else {
		log.Println(_LABEL_, "file has been closed!")
	}

	//2.清理
	this.sizeRolling = false
	this.dailyRolling = false
	this.nomalRolling = false
}

//==================================================================内部工具方法
//初始日志记录器(各日志器统一调用)
func (this *LLogger) startLogger(tp string) {
	defer func() {
		if e, ok := recover().(error); ok {
			log.Println(_LABEL_, "WARN: panic - %v", e)
			log.Println(_LABEL_, string("panic"))
		}
	}()

	//1.初始化空间
	var err error
	this.buf = &bytes.Buffer{}
	this.mu_buf = new(sync.Mutex)
	this.mu_file = new(sync.Mutex)
	this.path = tp
	checkFileDir(tp)
	this.logfile, err = os.OpenFile(
		tp,
		os.O_RDWR|os.O_APPEND|os.O_CREATE,
		_FILE_OPERAT_MODE_,
	)
	if nil != err {
		log.Println(_LABEL_, "OpenFile err!")
	}

	//2.开启监控线程
	go func() {
		this.timer = time.NewTimer(_CHECK_TIME_)
		this.writeTimer = time.NewTimer(_WRITE_TIME_)
		if !TIMEER_WRITE {
			this.writeTimer.Stop()
		}

		for {
			select {
			//定时检测是否分割
			case <-this.timer.C:
				this.fileCheck()
				if IS_DEBUG && false {
					log.Printf("*") //心跳
				}
				break
				//定时写入文件(定时写入,会导致延时)
			case <-this.writeTimer.C:
				this.bufWrite()
				if IS_DEBUG && false {
					log.Printf(".") //心跳
				}
				break
			}
		}
	}()

	if IS_DEBUG {
		jstr, err := json.Marshal(this)
		if nil == err {
			log.Println(_LABEL_, _VERSION_, string(jstr))
		}
	}
}

//文件检测(会锁定文件)
func (this *LLogger) fileCheck() {
	//0.边界判断
	if nil == this.mu_file ||
		nil == this.logfile ||
		"" == this.path {

		return
	}
	defer func() {
		if e, ok := recover().(error); ok {
			log.Println(_LABEL_, "WARN: panic - %v", e)
			log.Println(_LABEL_, string("panic"))
		}
	}()

	//1.重命名判断
	var RENAME_FLAG bool = false
	var CHECK_TIME time.Duration = _CHECK_TIME_
	this.timer.Stop()
	defer this.timer.Reset(CHECK_TIME)
	if this.dailyRolling {
		//日分割模式
		now := getNowFormDate(DATEFORMAT)
		if nil != now &&
			nil != this._date &&
			now.After(*this._date) {
			//超时重名
			RENAME_FLAG = true
		} else {
			//检测间隔动态调整
			du := this._date.UnixNano() - now.UnixNano()
			abs := math.Abs(float64(du))
			CHECK_TIME = CHECK_TIME * time.Duration(abs/abs)
		}
	} else if this.sizeRolling {
		//文件大小模式
		if "" != this.path &&
			this.maxFileCount >= 1 &&
			fileSize(this.path) >= this.maxFileSize {
			//超量重名
			RENAME_FLAG = true
		}
	} else if this.nomalRolling {
		//普通模式
		RENAME_FLAG = false
	}

	//2.重名操作
	if RENAME_FLAG {
		this.mu_file.Lock()
		defer this.mu_file.Unlock()
		if IS_DEBUG {
			log.Println(_LABEL_, this.path, "is need rename.")
		}
		this.fileRename()
	}

	return
}

//重命名文件
func (this *LLogger) fileRename() {
	//1.生成文件名称
	var err error
	var newName string
	var oldName string
	defer func() {
		if IS_DEBUG {
			log.Println(
				_LABEL_,
				oldName,
				"->",
				newName,
				":",
				err,
			)
		}
	}()

	if this.dailyRolling {
		//日期分割模式(文件不重命名)
		oldName = this.path
		newName = this.path
		this._date = getNowFormDate(DATEFORMAT)
		this.path = fmt.Sprint(
			this.dir,
			this.filename,
			_EXTEN_NAME_,
			".",
			this._date.Format(DATEFORMAT),
		)
	} else if this.sizeRolling {
		//大小分割模式(1,2,3....)
		suffix := int(this._suffix%int(this.maxFileCount) + 1)
		oldName = this.path
		newName = fmt.Sprint(
			this.path,
			".",
			fmt.Sprintf("%05d", suffix),
		)
		this._suffix = suffix
		this.path = this.path
	} else if this.nomalRolling {
		//常规模式
	}

	//2.处理旧文件
	this.logfile.Close()
	if "" != oldName && "" != newName && oldName != newName {
		if pathIsExist(newName) {
			//删除旧文件
			err := os.Remove(newName)
			if nil != err {
				log.Println(_LABEL_, "remove file err", err.Error())
			}
		}
		err = os.Rename(oldName, newName)
		if err != nil {
			//重名旧文件
			log.Println(_LABEL_, "rename file err", err.Error())
		}
	}

	//3.创建新文件
	this.logfile, err = os.OpenFile(
		this.path,
		os.O_RDWR|os.O_APPEND|os.O_CREATE,
		_FILE_OPERAT_MODE_,
	)
	if err != nil {
		log.Println(_LABEL_, "creat file err", err.Error())
	}

	return
}

//缓冲写入文件
func (this *LLogger) bufWrite() {
	//0.边界处理
	if nil == this.buf ||
		"" == this.path ||
		nil == this.logfile ||
		nil == this.mu_buf ||
		nil == this.mu_file ||
		this.buf.Len() <= 0 {
		return
	}

	//1.数据写入
	var WRITE_TIME time.Duration = _WRITE_TIME_
	if nil != this.writeTimer {
		this.writeTimer.Stop()
		defer this.writeTimer.Reset(WRITE_TIME)
	}
	this.mu_file.Lock()
	defer this.mu_file.Unlock()
	this.mu_buf.Lock()
	defer this.mu_buf.Unlock()
	defer this.buf.Reset()
	n, err := io.WriteString(this.logfile, this.buf.String())
	if nil != err {
		//写入失败,校验文件,不存在则创建
		checkFileDir(this.path)
		this.logfile, err = os.OpenFile(
			this.path,
			os.O_RDWR|os.O_APPEND|os.O_CREATE,
			_FILE_OPERAT_MODE_,
		)
		if nil != err {
			log.Println(_LABEL_, "log bufWrite() err!")
		}
	}
	//根据缓冲压力进行动态设置写入间隔
	if n == 0 {
		WRITE_TIME = _WRITE_TIME_
	} else {
		WRITE_TIME = WRITE_TIME * time.Duration(n/n)
	}
}

//==================================================================辅助方法
//获取文件大小
func fileSize(file string) int64 {
	this, e := os.Stat(file)
	if e != nil {
		if IS_DEBUG {
			log.Println(_LABEL_, e.Error())
		}
		return 0
	}

	return this.Size()
}

//判断路径是否存在
func pathIsExist(path string) bool {
	_, err := os.Stat(path)

	return err == nil || os.IsExist(err)
}

//检查文件路径文件夹,不存在则创建
func checkFileDir(tp string) {
	p, _ := path.Split(tp)
	d, err := os.Stat(p)
	if err != nil || !d.IsDir() {
		if err := os.MkdirAll(p, _FILE_CREAT_MODE_); err != nil {
			log.Println(_LABEL_, "CheckFileDir() Creat dir faile!")
		}
	}
}

//获取当前指定格式的日期
func getNowFormDate(form string) *time.Time {
	t, err := time.Parse(form, time.Now().Format(form))
	if nil != err {
		log.Println(_LABEL_, "getNowFormDate()", err.Error())
		t = time.Time{}

		return &t
	}

	return &t
}

//==================================================================测试用例
func Test() {
	logg := New()
	logg.SetType(1)
	logg.SetRollingNormal("./logs", "logg")
	logg.Println("hello world!")
}
