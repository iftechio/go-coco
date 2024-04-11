package profile

import (
	"net/http"
	"net/http/pprof"
	"os"
	"strings"

	"github.com/google/gops/agent"
	xLogger "github.com/xieziyu/go-coco/utils/logger"
)

const (
	pprofAddress = ":25110"
	gopsAddress  = ":25111"
)

func init() {
	// 避免在测试环境载入
	// 在测试环境载入，会直接在两个端口上启动服务器
	// 如果需要并行得启动另外测试，就是会发生端口冲突
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.") {
			return
		}
	}
	go startPprofServer()
	go startGopsServer()
}

func startPprofServer() {
	xLogger.Infof("pprof started on %s", pprofAddress)
	// 手动挂载 router, http.DefaultServerMux 不安全
	s := http.NewServeMux()
	s.HandleFunc("/debug/pprof/", pprof.Index)
	s.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	s.HandleFunc("/debug/pprof/profile", pprof.Profile)
	s.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	s.HandleFunc("/debug/pprof/trace", pprof.Trace)
	if err := http.ListenAndServe(pprofAddress, s); err != nil {
		xLogger.Fatal(err)
	}
}

func startGopsServer() {
	xLogger.Infof("gops started on %s", gopsAddress)
	if err := agent.Listen(agent.Options{
		Addr:                   gopsAddress,
		ConfigDir:              "",
		ShutdownCleanup:        false,
		ReuseSocketAddrAndPort: false,
	}); err != nil {
		xLogger.Fatal(err)
	}
}
