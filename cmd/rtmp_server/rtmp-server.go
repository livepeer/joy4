package main

import (
	"fmt"
	"strings"

	"github.com/livepeer/joy4/av/avutil"
	"github.com/livepeer/joy4/format"
	"github.com/livepeer/joy4/format/rtmp"
)

func init() {
	format.RegisterAll()
}

func main() {
	server := &rtmp.Server{}

	server.HandlePlay = func(conn *rtmp.Conn) {
		segs := strings.Split(conn.URL.Path, "/")
		url := fmt.Sprintf("%s://%s", segs[1], strings.Join(segs[2:], "/"))
		fmt.Printf("===> got play conn %+v\n", conn)
		src, _ := avutil.Open(url)
		avutil.CopyFile(conn, src)
	}

	server.HandlePublish = func(conn *rtmp.Conn) {
		fmt.Printf("===> got publish conn %s\n", conn.URL)
		/*
			ch := channels[conn.URL.Path]
			if ch == nil {
				ch = &Channel{}
				ch.que = pubsub.NewQueue()
				query := conn.URL.Query()
				if q := query.Get("cachegop"); q != "" {
					var n int
					fmt.Sscanf(q, "%d", &n)
					ch.que.SetMaxGopCount(n)
				}
				channels[conn.URL.Path] = ch
			} else {
				ch = nil
			}
			l.Unlock()
			if ch == nil {
				return
			}

			avutil.CopyFile(ch.que, conn)

			l.Lock()
			delete(channels, conn.URL.Path)
			l.Unlock()
			ch.que.Close()
		*/
		go func() {
			for {
				pkt, err := conn.ReadPacket()
				if err != nil {
					fmt.Printf("conn %s err=%v\n", conn.URL.String(), err)
					break
				}
				if false {
					fmt.Printf("Got packet idx=%d dts=%s\n", pkt.Idx, pkt.Time)
				}
			}
		}()
	}

	fmt.Println("Listening")

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("err=%v\n", err)
	}

	// ffplay rtmp://localhost/rtsp/192.168.1.1/camera1
	// ffplay rtmp://localhost/rtmp/live.hkstv.hk.lxdns.com/live/hks
}
