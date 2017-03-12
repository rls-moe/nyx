package resources

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hlandau/passlib"
	"github.com/tidwall/buntdb"
)

type ModPass struct {
	ID       string `json:"id"`
	Password string `json:"password"`
	Board    string `json:"board"`
}

func (m *ModPass) HashLogin(pass string) error {
	var err error
	m.Password, err = passlib.Hash(pass)
	return err
}

func (m *ModPass) VerifyLogin(pass string) error {
	var err error
	err = passlib.VerifyNoUpgrade(pass, m.Password)
	return err
}

func NewMod(tx *buntdb.Tx, host string, in *ModPass) error {
	dat, err := json.Marshal(in)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, replaced, err := tx.Set(
		fmt.Sprintf(modPassPath, escapeString(host), escapeString(in.ID)),
		string(dat),
		nil)
	if err != nil {
		tx.Rollback()
		return err
	}
	if replaced {
		tx.Rollback()
		return errors.New("Admin already exists")
	}
	return nil
}

func GetMod(tx *buntdb.Tx, host, id string) (*ModPass, error) {
	var ret = &ModPass{}
	dat, err := tx.Get(
		fmt.Sprintf(modPassPath, escapeString(host), escapeString(id)),
	)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal([]byte(dat), ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func DelMod(tx *buntdb.Tx, host, id string) error {
	if _, err := tx.Delete(
		fmt.Sprintf(modPassPath, escapeString(host), escapeString(id)),
	); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
