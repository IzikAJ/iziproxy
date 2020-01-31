package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

// write failed request response
func writeFailResponse(writerPointer *http.ResponseWriter, status int, msg string) {
	writer := *writerPointer
	writer.WriteHeader(status)
	writer.Write([]byte(msg))
}

var stats = new(runtime.MemStats)

type statBlock struct {
	title string
	desc  string
	items []statItem
}
type statItem struct {
	title string
	value string
	desc  string
}

// write server stats response
func writeStatsResponse(writerPointer *http.ResponseWriter, server *Server) {
	writer := *writerPointer
	runtime.ReadMemStats(stats)

	rawStats, _ := json.Marshal(server.Stats)
	spaces := []statItem{}
	for space, signal := range server.space {
		spaces = append(spaces, statItem{
			title: space,
			value: fmt.Sprintf("<a href=\"http://%v.proxy.me\" target=\"_blank\"><b>%v</b> - %v</a>", space, space, signal),
		})
	}

	blocks := []statBlock{
		statBlock{
			title: "Spaces",
			items: spaces,
		},

		statBlock{
			title: "Mem stats",
			items: []statItem{
				statItem{"Alloc", fmt.Sprintf("%.2fMB", float64(stats.Alloc>>10)/1024), "bytes allocated and not yet freed"},
				statItem{"Sys", fmt.Sprintf("%.2fMB", float64(stats.Sys>>10)/1024), "bytes obtained from system"},
				statItem{"TotalAlloc", fmt.Sprintf("%.2fMB", float64(stats.TotalAlloc>>10)/1024), "bytes allocated (even if freed)"},
				statItem{"NumGC", fmt.Sprintf("%d", stats.NumGC), ""},
				statItem{"GCCPUFraction", fmt.Sprintf("%.6f%%", stats.GCCPUFraction*100), ""},
			},
		},

		statBlock{
			title: "Heap stats",
			items: []statItem{
				statItem{"HeapAlloc", fmt.Sprintf("%.2fMB", float64(stats.HeapAlloc>>10)/1024), ""},
				statItem{"HeapIdle", fmt.Sprintf("%.2fMB", float64(stats.HeapIdle>>10)/1024), ""},
				statItem{"HeapSys", fmt.Sprintf("%.2fMB", float64(stats.HeapSys>>10)/1024), ""},
				statItem{"HeapInuse", fmt.Sprintf("%.2fMB", float64(stats.HeapInuse>>10)/1024), ""},
				statItem{"HeapReleased", fmt.Sprintf("%.2fMB", float64(stats.HeapReleased>>10)/1024), ""},
			},
		},

		statBlock{
			title: "Other",
			items: []statItem{
				statItem{"Server Stats", fmt.Sprintf("%v", string(rawStats)), ""},
				statItem{"Current Time", fmt.Sprintf("%v", time.Now().String()), ""},
				statItem{"NumCPU", fmt.Sprintf("%d", runtime.NumCPU()), ""},
				statItem{"NumGoroutine", fmt.Sprintf("%d", runtime.NumGoroutine()), ""},
				statItem{"GO Version", fmt.Sprintf("%v", runtime.Version()), ""},
			},
		},
	}

	writer.WriteHeader(http.StatusOK)
	fmt.Fprintln(writer, "<html>")
	fmt.Fprintln(writer, "<head>")
	fmt.Fprintln(writer, "<meta http-equiv=\"refresh\" content=\"3;url=/stats\" />")
	fmt.Fprintln(writer, "</head>")
	fmt.Fprintln(writer, "<body>")

	for _, block := range blocks {
		fmt.Fprintf(writer, "<h3>%v</h3>", block.title)
		fmt.Fprintf(writer, "<p>%v</p>", block.desc)
		for _, item := range block.items {
			fmt.Fprintf(writer, "<b>%v:</b> %v<br/><i>%v</i></p>", item.title, item.value, item.desc)
		}
	}

	fmt.Fprintln(writer, "</body>")
	fmt.Fprintln(writer, "</html>")
}
