package resources

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/buntdb"
)

type Thread struct {
	ID         int      `json:"id"`
	StartReply int      `json:"start"`
	Board      string   `json:"board"`
	Metadata   Metadata `json:"-"`

	intReply *Reply

	intReplies []*Reply
}

func (t *Thread) GetReplies() []*Reply {
	return t.intReplies
}

func (t *Thread) GetReply() *Reply {
	return t.intReply
}

func NewThread(tx *buntdb.Tx, host, board string, in *Thread, in2 *Reply) error {
	var err error

	err = TestBoard(tx, host, in.Board)
	if err != nil {
		return err
	}

	in.ID, err = getID()
	if err != nil {
		return err
	}
	in2.Thread = in.ID

	in2.ID, err = getID()
	if err != nil {
		return err
	}
	in.StartReply = in2.ID

	dat, err := json.Marshal(in)
	if err != nil {
		return err
	}

	_, replaced, err := tx.Set(
		fmt.Sprintf(threadPath, escapeString(host), escapeString(board), in.ID),
		string(dat),
		nil)

	if err != nil {
		return err
	}
	if replaced {
		return errors.New("Thread already exists")
	}

	return NewReply(tx, host, board, in, in2, true)
}

func TestThread(tx *buntdb.Tx, host, board string, id int) error {
	err := TestBoard(tx, host, board)
	if err != nil {
		return err
	}

	_, err = tx.Get(
		fmt.Sprintf(threadPath, escapeString(host), escapeString(board), id),
	)
	return err
}

func GetThread(tx *buntdb.Tx, host, board string, id int) (*Thread, error) {
	var ret = &Thread{}
	dat, err := tx.Get(
		fmt.Sprintf(threadPath, escapeString(host), escapeString(board), id),
	)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal([]byte(dat), ret); err != nil {
		return nil, err
	}

	ret.intReply, err = GetReply(tx, host, board, id, ret.StartReply)
	if err != nil && err == buntdb.ErrNotFound {
		ret.intReply = &Reply{
			Board:     ret.Board,
			Thread:    ret.ID,
			ID:        -1,
			Image:     nil,
			Thumbnail: nil,
			Metadata: map[string]string{
				"deleted": "not found",
			},
			Text: "[not found]",
		}
	} else if err != nil {
		return nil, err
	}
	return ret, nil
}

func DelThread(tx *buntdb.Tx, host, board string, id int) error {
	if _, err := tx.Delete(
		fmt.Sprintf(threadPath, escapeString(host), escapeString(board), id),
	); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func FillReplies(tx *buntdb.Tx, host string, thread *Thread) (err error) {
	thread.intReplies, err = ListReplies(tx, host, thread.Board, thread.ID)
	return
}

func ListThreads(tx *buntdb.Tx, host, board string) ([]*Thread, error) {
	var threadList = []*Thread{}
	var err error

	err = TestBoard(tx, host, board)
	if err != nil {
		return nil, err
	}

	tx.DescendKeys(
		fmt.Sprintf(
			threadSPath,
			escapeString(host),
			escapeString(board),
		),
		func(key, value string) bool {
			var thread = &Thread{}
			err = json.Unmarshal([]byte(value), thread)
			if err != nil {
				return false
			}
			thread.intReply, err = GetReply(tx, host, board, thread.ID, thread.StartReply)
			if err != nil {
				if err == buntdb.ErrNotFound {
					err = nil
					thread.intReply = &Reply{
						Board:     thread.Board,
						Thread:    thread.ID,
						ID:        -1,
						Image:     nil,
						Thumbnail: nil,
						Metadata: map[string]string{
							"deleted": "not found",
						},
						Text: "[not found]",
					}
				} else {
					return false
				}
			}

			threadList = append(threadList, thread)
			if len(threadList) >= 25 {
				return false
			}
			return true
		})
	return threadList, err
}
