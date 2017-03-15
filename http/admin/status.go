package admin

import (
	"bytes"
	"fmt"
	"github.com/dustin/go-humanize"
	"go.rls.moe/nyx/http/errw"
	"go.rls.moe/nyx/http/middle"
	"net/http"
	"runtime"
	"sync"
	"time"
)

var memStat = map[string]interface{}{}
var memStatLock = new(sync.RWMutex)

var startTime = time.Now().UTC()

func init() {
	go func() {
		update()
		ticker := time.Tick(time.Second * 10)
		for _ = range ticker {
			update()
		}
	}()
}

func update() {
	memStatLock.Lock()
	defer memStatLock.Unlock()
	memStat["Uptime"] = uptime()
	memStat["GC"] = gcStat()
}
func uptime() map[string]interface{} {
	return map[string]interface{}{
		"human":   humanize.Time(startTime),
		"hours":   time.Now().Sub(startTime).Hours(),
		"seconds": time.Now().Sub(startTime).Seconds(),
	}
}

func gcStat() map[string]interface{} {
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)

	mem := map[string]interface{}{}
	mem["alloc"] = fmt.Sprintf("%.5f GiB", float64(m.Alloc)/1024/1024/1024)
	mem["calloc"] = fmt.Sprintf("%.5f GiB", float64(m.TotalAlloc)/1024/1024/1024)
	mem["sysmem"] = fmt.Sprintf("%.5f MiB", float64(m.Sys)/1024/1024)
	mem["lookups"] = fmt.Sprintf("× %d", m.Lookups)
	mem["mallocs"] = fmt.Sprintf("× %d", m.Mallocs)
	mem["frees"] = fmt.Sprintf("× %d", m.Frees)
	mem["liveobj"] = fmt.Sprintf("× %d", m.Mallocs-m.Frees)
	mem["heapalloc"] = fmt.Sprintf("%.5f MiB", float64(m.HeapSys)/1024/1024)
	mem["heaprelease"] = fmt.Sprintf("%.5f MiB", float64(m.HeapReleased)/1024/1024)
	mem["gcmeta"] = fmt.Sprintf("%.5f MiB", float64(m.GCSys)/1024/1024)
	mem["pause"] = fmt.Sprintf("%.5f min", float64(m.PauseTotalNs)/1000/1000/1000/60)
	mem["gctimes"] = fmt.Sprintf("× %d", m.NumGC)
	mem["fgctimes"] = fmt.Sprintf("× %d", m.NumForcedGC)
	mem["cpufrac"] = fmt.Sprintf("%.5f %%", m.GCCPUFraction*100)
	return map[string]interface{}{
		"numcpu":   runtime.NumCPU(),
		"numgor":   runtime.NumGoroutine(),
		"version":  runtime.Version(),
		"arch":     runtime.GOARCH,
		"os":       runtime.GOOS,
		"compiler": runtime.Compiler,
		"memory":   mem,
	}
}
func serveStatus(w http.ResponseWriter, r *http.Request) {
	memStatLock.RLock()
	defer memStatLock.RUnlock()
	ctx := middle.GetBaseCtx(r)
	ctx["Uptime"] = memStat["Uptime"]
	ctx["GC"] = memStat["GC"]
	dat := bytes.NewBuffer([]byte{})
	err := statusTmpl.Execute(dat, ctx)
	if err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}
	http.ServeContent(w, r, "panel.html", time.Now(),
		bytes.NewReader(dat.Bytes()))
}
