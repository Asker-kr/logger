package logger

func doSomething(){
	Log(NOTICE, "Log Example. %d %s", 10, "param")
}

// it just a sample code.
func main() {
	// log file will be created at "./log/test_yyyymmdd.log"
	logPath := "./log"
	logPrefix := "test"

	// 1. make empty channel
	// 2. Start log with parameters
	// 3. use log
	// 4. call Close() for log close
	// 5. <-done wait channel to ensure that all logs have been completed.

	done := make(chan struct{})
	StartLog(logPath, logPrefix, 500, done)
	Log(NOTICE, "---------- Start  --------------")
	doSomething()
	Log(NOTICE, "----------   End  --------------")
	Close()
	<-done
}