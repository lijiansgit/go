package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	addr        string
	user        string
	password    string
	maxOpenConn int
	maxIdleConn int
	loops       int
	sqlFile     string
	timeStr     string
	// ExecCount 执行总计
	ExecCount int64
	// ErrCount 执行失败总计
	ErrCount int64
	wg       *sync.WaitGroup
)

func init() {
	flag.StringVar(&addr, "addr", "10.200.150.202:3306", "mysql address: ip:port")
	flag.StringVar(&user, "user", "root", "mysql user")
	flag.StringVar(&password, "password", "root", "mysql user's password")
	flag.IntVar(&maxOpenConn, "maxOpenConn", 5000, "连接池中的最大的连接数")
	flag.IntVar(&maxIdleConn, "maxIdleConn", 100, "连接池中的闲置的连接数")
	flag.IntVar(&loops, "loops", 1, "循环执行sql文件的次数")
	//
	flag.StringVar(&sqlFile, "sqlFile", "./test.sql", "test sql file")
	flag.StringVar(&timeStr, "t", "1s", "每隔多长时间启动一个协程，默认为1秒，eg: 1us, 1ms")
}

func main() {
	flag.Parse()

	timeDur, err := time.ParseDuration(timeStr)
	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf(`%s:%s@tcp(%s)/test?charset=utf8`, user, password, addr)
	db, err := sql.Open(`mysql`, dsn)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	db.SetMaxOpenConns(maxOpenConn)
	db.SetMaxIdleConns(maxIdleConn)
	if err = db.Ping(); err != nil {
		panic(err)
	}

	sqls, err := openFile(sqlFile)
	if err != nil {
		panic(err)
	}

	for i := 1; i < loops; i++ {
		sqls = append(sqls, sqls...)
	}

	startTime := time.Now()
	go result()

	wg = new(sync.WaitGroup)
	for _, v := range sqls {
		wg.Add(1)
		go exec(db, v)

		time.Sleep(timeDur)
	}

	wg.Wait()

	log.Printf("ExecCount:%d ErrCount:%d \n", ExecCount, ErrCount)
	log.Println("Run Time: ", time.Since(startTime))
}

func exec(db *sql.DB, sqls string) {
	defer wg.Done()

	_, err := db.Exec(sqls)
	if err != nil {
		log.Printf("Exec error(%v) \n", err)
		atomic.AddInt64(&ErrCount, 1)
	}

	atomic.AddInt64(&ExecCount, 1)
}

func openFile(fileName string) (res []string, err error) {
	f, err := os.Open(fileName)
	if err != nil {
		return res, err
	}

	defer f.Close()

	br := bufio.NewReader(f)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		line := string(a)
		if line != "" {
			res = append(res, line)
		}
	}

	return res, nil
}

func result() {
	var (
		lastTimes int64
		diff      int64
		nowCount  int64
		errCount  int64
		timer     = int64(1)
	)

	for {
		nowCount = atomic.LoadInt64(&ExecCount)
		diff = nowCount - lastTimes
		lastTimes = nowCount
		errCount = atomic.LoadInt64(&ErrCount)
		log.Printf("ExecCount:%d ErrCount:%d exec/s:%d\n", nowCount, errCount, diff/timer)

		time.Sleep(time.Duration(timer) * time.Second)
	}
}
