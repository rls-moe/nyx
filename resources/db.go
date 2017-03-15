package resources

import (
	"errors"
	"fmt"
	"github.com/tidwall/buntdb"
	"regexp"
	"strings"
)

const (
	setup         = "/jack/setup"
	hostEnable    = "/jack/%s/enabled"
	boardPath     = "/jack/%s/board/%s/board-data"
	threadPath    = "/jack/%s/board/%s/thread/%032d/thread-data"
	threadSPath   = "/jack/%s/board/%s/thread/*/thread-data"
	replyPath     = "/jack/%s/board/%s/thread/%032d/reply/%032d/reply-data"
	replySPath    = "/jack/%s/board/%s/thread/%032d/reply/*/reply-data"
	modPassPath   = "/jack/%s/pass/mod/%s/mod-data"
	adminPassPath = "/jack/./pass/admin/%s/admin-data"
)

func GetHostnameFromKey(key string) (string, error) {
	regex := regexp.MustCompile(`^/jack/(.+)/(board|pass)`)
	res := regex.FindStringSubmatch(key)
	if len(res) != 3 {
		fmt.Printf("Found %d keys: %s", len(res), res)
		return "", errors.New("Could not find host in key")
	}
	return unescapeString(res[1]), nil
}

func InitialSetup(db *buntdb.DB) error {
	return db.Update(func(tx *buntdb.Tx) error {
		if _, err := tx.Get(setup); err != nil {
			fmt.Println("")
			if err != buntdb.ErrNotFound {
				fmt.Println("DB setup not known.")
				return err
			}
			fmt.Println("DB not setup.")
			tx.Set(setup, "yes", nil)
		} else {
			fmt.Println("DB setup.")
			return nil
		}

		fmt.Println("Creating Indices")
		err := tx.CreateIndex("board/short", "/jack/*/board/*/board-data", buntdb.IndexJSON("short"))
		if err != nil {
			return err
		}
		err = tx.CreateIndex("replies", "/jack/*/board/*/thread/*/reply/*/reply-data", buntdb.IndexJSON("thread"))
		if err != nil {
			return err
		}
		err = tx.CreateIndex("board/thread", "/jack/*/board/*/thread/*/thread-data", buntdb.IndexJSON("board"))
		if err != nil {
			return err
		}

		fmt.Println("Creating default admin")
		admin := &AdminPass{
			ID: "admin",
		}
		err = admin.HashLogin("admin")
		if err != nil {
			return err
		}
		fmt.Println("Saving default admin to DB")
		err = NewAdmin(tx, admin)
		if err != nil {
			return err
		}

		fmt.Println("Committing setup...")

		return nil
	})
}

func CreateHost(db *buntdb.DB, hostname string) error {
	return db.Update(func(tx *buntdb.Tx) error {
		hostname = escapeString(hostname)
		_, replaced, err := tx.Set(fmt.Sprintf(hostEnable, "hostname"), "", nil)
		if err != nil {
			tx.Rollback()
			return err
		}
		if replaced {
			tx.Rollback()
			return errors.New("Hostname already enabled")
		}

		board := &Board{
			ShortName: "d",
			LongName:  "default",
			Metadata: map[string]string{
				"locked":      "true",
				"description": "Default Board",
			},
		}
		err = NewBoard(tx, hostname, board)
		if err != nil {
			tx.Rollback()
			return err
		}

		return nil
	})
}

func escapeString(in string) string {
	in = strings.Replace(in, ".", ".dot.", -1)
	in = strings.Replace(in, "-", ".minus.", -1)
	in = strings.Replace(in, "\\", ".backslash.", -1)
	in = strings.Replace(in, "*", ".star.", -1)
	in = strings.Replace(in, "?", ".ask.", -1)
	in = strings.Replace(in, "/", ".slash.", -1)
	in = strings.Replace(in, "@", ".at.", -1)
	in = strings.Replace(in, ">>", ".quote.", -1)
	in = strings.Replace(in, ">", ".arrow-left.", -1)
	in = strings.Replace(in, "<", ".arrow-right.", -1)
	return in
}

func unescapeString(in string) string {
	in = strings.Replace(in, ".arrow-right.", "<", -1)
	in = strings.Replace(in, ".arrow-left.", ">", -1)
	in = strings.Replace(in, ".quote.", ">>", -1)
	in = strings.Replace(in, ".at.", "@", -1)
	in = strings.Replace(in, ".slash.", "/", -1)
	in = strings.Replace(in, ".ask.", "?", -1)
	in = strings.Replace(in, ".star.", "*", -1)
	in = strings.Replace(in, ".backslash.", "\\", -1)
	in = strings.Replace(in, ".minus.", "-", -1)
	in = strings.Replace(in, ".dot.", ".", -1)
	return in
}
