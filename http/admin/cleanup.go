package admin

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/buntdb"
	"go.rls.moe/nyx/http/errw"
	"go.rls.moe/nyx/http/middle"
	"go.rls.moe/nyx/resources"
	"net/http"
	"strings"
	"time"
)

func handleCleanup(w http.ResponseWriter, r *http.Request) {
	sess := middle.GetSession(r)
	if sess == nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}
	if sess.CAttr("mode") != "admin" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	fmt.Println("Beginning cleanup...")
	db := middle.GetDB(r)
	var delKeys = []string{}
	err := db.View(func(tx *buntdb.Tx) error {
		var err error
		tx.AscendKeys("*", func(key, value string) bool {
			keyType := detectType(key)
			if keyType == "thread" {
				var host string
				host, err = resources.GetHostnameFromKey(key)
				if err != nil {
					fmt.Printf("Error: %s", err)
					return false
				}
				var thread = &resources.Thread{}
				err = json.Unmarshal([]byte(value), thread)
				if err != nil {
					fmt.Printf("Error: %s", err)
					return false
				}
				threadTime := resources.DateFromId(thread.ID)
				dur := threadTime.Sub(time.Now())
				if dur > time.Hour*24*7 {
					fmt.Printf("Sched %s for deletion: expired\n", key)
					delKeys = append(delKeys, key)
					return true
				}
				err = resources.FillReplies(tx, host, thread)
				if err != nil {
					fmt.Printf("Error: %s", err)
					return false
				}
				if len(thread.GetReplies()) == 0 {
					fmt.Printf("Sched %s for deletion: empty\n", key)
					delKeys = append(delKeys, key)
					return true
				}
				if _, err := resources.GetReply(tx, host, thread.Board, thread.ID, thread.StartReply); err == buntdb.ErrNotFound {
					fmt.Printf("Sched %s for delection: main reply dead\n", key)
					delKeys = append(delKeys, key)
					return true
				}
			} else if keyType == "reply" {
				var host string
				host, err = resources.GetHostnameFromKey(key)
				if err != nil {
					fmt.Printf("Error: %s", err)
					return false
				}
				var reply = &resources.Reply{}
				err = json.Unmarshal([]byte(value), reply)
				if err != nil {
					fmt.Printf("Error: %s", err)
					return false
				}
				replyTime := resources.DateFromId(reply.ID)
				dur := replyTime.Sub(time.Now())
				if dur > time.Hour*24*7 {
					fmt.Printf("Sched %s for deletion: expired\n", key)
					delKeys = append(delKeys, key)
					return true
				}
				if val, ok := reply.Metadata["deleted"]; ok && val == "yes" {
					fmt.Printf("Sched %s for deletion: deleted\n", key)
					delKeys = append(delKeys, key)
					return true
				}
				if err := resources.TestThread(tx, host, reply.Board, reply.Thread); err == buntdb.ErrNotFound {
					fmt.Printf("Sched %s for deletion: missing parent %d: %s\n", key, reply.Thread, err)
					delKeys = append(delKeys, key)
					return true
				}
			}
			return true
		})
		/* Insert cleanup codes here */
		return err
	})
	fmt.Println("Removing sched' entries")
	db.Update(func(tx *buntdb.Tx) error {
		for _, v := range delKeys {
			fmt.Printf("Deleting %s\n", v)
			tx.Delete(v)
		}
		return nil
	})
	fmt.Println("Shrinking DB")
	err = db.Shrink()
	if err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}
	fmt.Println("Finished Cleanup")

	http.Redirect(w, r, "/admin/panel.html", http.StatusSeeOther)
}

func detectType(key string) string {
	if strings.Contains(key, "/jack/") {
		if strings.HasSuffix(key, "/board-data") {
			return "board"
		}
		if strings.HasSuffix(key, "/thread") {
			return "thread"
		}
		if strings.HasSuffix(key, "/reply-data") {
			return "reply"
		}
		return "system"
	}
	return "none"
}
