package resources

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/buntdb"
)

type Board struct {
	ShortName string   `json:"short"`
	LongName  string   `json:"long"`
	Metadata  Metadata `json:"meta"`
}

func NewBoard(tx *buntdb.Tx, hostname string, in *Board) error {
	dat, err := json.Marshal(in)
	if err != nil {
		return err
	}
	_, replaced, err := tx.Set(
		fmt.Sprintf(boardPath, escapeString(hostname), escapeString(in.ShortName)),
		string(dat),
		nil)
	if err != nil {
		return err
	}
	if replaced {
		return errors.New("Board " + escapeString(in.ShortName) + " already exists")
	}
	return nil
}

func TestBoard(tx *buntdb.Tx, hostname, shortname string) (error) {
	_, err := tx.Get(
		fmt.Sprintf(boardPath, escapeString(hostname), escapeString(shortname)),
	)
	return err
}

func UpdateBoard(tx *buntdb.Tx, hostname string, b *Board) error {
	if err := TestBoard(tx, hostname, b.ShortName); err != nil {
		return err
	}

	dat, err := json.Marshal(b)
	if err != nil {
		return err
	}
	_, replaced, err := tx.Set(
		fmt.Sprintf(boardPath, escapeString(hostname), escapeString(b.ShortName)),
		string(dat),
		nil)
	if err != nil {
		return err
	}
	if !replaced {
		return errors.New("Board " + escapeString(b.ShortName) + " does not exist")
	}
	return nil
}

func GetBoard(tx *buntdb.Tx, hostname, shortname string) (*Board, error) {
	var ret = &Board{}
	dat, err := tx.Get(
		fmt.Sprintf(boardPath, escapeString(hostname), escapeString(shortname)),
	)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal([]byte(dat), ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func DelBoard(tx *buntdb.Tx, hostname, shortname string) error {
	if _, err := tx.Delete(
		fmt.Sprintf(boardPath, escapeString(hostname), escapeString(shortname)),
	); err != nil {
		return err
	}
	return nil
}

func ListBoards(tx *buntdb.Tx, hostname string) ([]*Board, error) {
	var boardList = []*Board{}
	var err error
	tx.AscendKeys(fmt.Sprintf(boardPath, escapeString(hostname), "*"),
		func(key, value string) bool {
			var board = &Board{}
			err = json.Unmarshal([]byte(value), board)
			if err != nil {
				return false
			}
			boardList = append(boardList, board)
			return true
		})
	return boardList, err
}
