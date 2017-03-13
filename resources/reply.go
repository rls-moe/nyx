package resources

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/buntdb"
	"golang.org/x/crypto/blake2b"
)

type Reply struct {
	ID       int      `json:"id"`
	Text     string   `json:"text"`
	Image    []byte   `json:"image"`
	Thread   int      `json:"thread"`
	Board    string   `json:"board"`
	Metadata Metadata `json:"meta"`
}

func NewReply(tx *buntdb.Tx, host, board string, thread *Thread, in *Reply, noId bool) error {
	var err error

	if !noId {
		in.ID, err = getID()
		if err != nil {
			return err
		}
	} else {
	}

	dat, err := json.Marshal(in)
	if err != nil {
		return err
	}

	err = TestThread(tx, host, in.Board, in.Thread)
	if err != nil {
		return err
	}

	_, replaced, err := tx.Set(
		fmt.Sprintf(replyPath, escapeString(host), escapeString(board), thread.ID, in.ID),
		string(dat),
		nil)
	if err != nil {
		return err
	}
	if replaced {
		return errors.New("Admin already exists")
	}
	return nil
}

func UpdateReply(tx *buntdb.Tx, host, board string, r *Reply) error {
	if err := TestBoard(tx, host, board); err != nil {
		return err
	}
	if err := TestThread(tx, host, board, r.Thread); err != nil {
		return err
	}

	dat, err := json.Marshal(r)
	if err != nil {
		return err
	}

	_, replaced, err := tx.Set(
		fmt.Sprintf(replyPath, escapeString(host), escapeString(board), r.Thread, r.ID),
		string(dat),
		nil)
	if err != nil {
		return err
	}
	if !replaced {
		return fmt.Errorf("Reply %d/%d does not exist", r.Thread, r.ID)
	}
	return nil
}

func GetReply(tx *buntdb.Tx, host, board string, thread, id int) (*Reply, error) {
	var ret = &Reply{}
	dat, err := tx.Get(
		fmt.Sprintf(replyPath, escapeString(host), escapeString(board), thread, id),
	)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal([]byte(dat), ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func DelReply(tx *buntdb.Tx, host, board string, thread, id int) error {
	if _, err := tx.Delete(
		fmt.Sprintf(replyPath, escapeString(host), escapeString(board), thread, id),
	); err != nil {
		return err
	}
	return nil
}

func ListReplies(tx *buntdb.Tx, host, board string, thread int) ([]*Reply, error) {
	var replyList = []*Reply{}
	var err error

	err = TestThread(tx, host, board, thread)
	if err != nil {
		return nil, err
	}

	tx.AscendKeys(
		fmt.Sprintf(
			replySPath,
			escapeString(host),
			escapeString(board),
			thread,
		),
		func(key, value string) bool {
			var reply = &Reply{}
			err = json.Unmarshal([]byte(value), reply)
			if err != nil {
				return false
			}
			replyList = append(replyList, reply)
			return true
		})

	return replyList, err
}

func CalcTripCode(trip string) string {
	fullTrip := blake2b.Sum256([]byte(trip))
	return base64.RawStdEncoding.EncodeToString(fullTrip[:8])
}
