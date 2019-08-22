package main

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

type LogProcess struct {
	r chan []byte
	w chan []byte
	read Reader
	write Writer
}


type Reader interface {
	Read(rr chan []byte)
}

type Writer interface {
	Write(ww chan []byte)
}

type ReadFromFile struct {
	path string
}

type WriteToDb struct {
	dsn string
}


func (rr *ReadFromFile) Read(r chan []byte)  {

	file,err:=os.Open(rr.path)
	if err!=nil{
		panic(err)
	}

	//如果要从文件末尾开始收集日志则：
	//file.Seek(0,2)

	bufFile := bufio.NewReader(file)
	count := 0
	for{
		count++
		if count>5{
			break
		}
		line,err:=bufFile.ReadBytes('\n')
		if err==io.EOF{
			time.Sleep(100*time.Millisecond)
			continue
		}else if err!=nil{
			panic(err)
		}
		r <- line[:len(line)-1]
	}
}

func (l *LogProcess) ProcessLog()  {
	for log:=range l.r{
		l.w <- log
	}
}

type LogInfo struct {
	timeLocal time.Time
	ip string
	method string
	path string
	http string
	status string
}
func (ww *WriteToDb) Write(w chan []byte)  {

	//120.216.207.220 - - [22/Aug/2019:14:58:35 +0800] "GET http://finance.sina.com.cn/ HTTP/1.1" 200 194 "http://finance.sina.com.cn/" "Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0)"

	var exp  = `([^ ]*) ([^ ]*) ([^ ]*) (\[(.*)\]) (\"(.*?)\") (-|[0-9]*) (-|[0-9]*) (\".*?\") (\".*?\")`

	re := regexp.MustCompile(exp)

	loc,_:=time.LoadLocation("Asia/Shanghai")

	for log:=range w{

		matches:=re.FindStringSubmatch(string(log))

		t,err:=time.ParseInLocation("02/Jan/2006:15:04:05 +0800",matches[5],loc)

		info:=LogInfo{}

		if err!=nil{
			fmt.Println("time err",err.Error(),matches[5])
			continue
		}

		info.timeLocal = t

		info.ip = matches[1]

		//"GET /lnmp.gif HTTP/1.1"
		getInfo := strings.Split(matches[7], " ")
		if len(getInfo) != 3{
			fmt.Println("split err",matches[7])
			continue
		}

		info.method = getInfo[0]

		u,err:= url.Parse(getInfo[1])

		if err!=nil{
			fmt.Println("Parse err",err)
		}

		info.path = u.Path

		info.status = matches[8]


		fmt.Println("info:",info)

	}
}


func main(){

	rr1:=&ReadFromFile{
		path:"./access.log",
	}

	ww1:=&WriteToDb{
		dsn:"111",//TODO
	}

	var l = &LogProcess{
		r : make(chan []byte),
		w : make(chan []byte),
		read:rr1,
		write:ww1,
	}

	go 	l.read.Read(l.r)
	go	l.ProcessLog()
	go	l.write.Write(l.w)

	time.Sleep(1*time.Second)
}