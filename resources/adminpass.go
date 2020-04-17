package resources

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/hlandau/passlib.v1"
	"github.com/tidwall/buntdb"
)

type AdminPass struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

func (a *AdminPass) HashLogin(pass string) error {
	var err error
	a.Password, err = passlib.Hash(pass)
	return err
}

func (a *AdminPass) VerifyLogin(pass string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("pkg: %v", r)
			}
		}
	}()
	if a == nil {
		return errors.New("no login")
	}
	err = passlib.VerifyNoUpgrade(pass, a.Password)
	return err
}

func NewAdmin(tx *buntdb.Tx, in *AdminPass) error {
	dat, err := json.Marshal(in)
	if err != nil {
		return err
	}
	_, replaced, err := tx.Set(
		fmt.Sprintf(adminPassPath, escapeString(in.ID)),
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

func GetAdmin(tx *buntdb.Tx, id string) (*AdminPass, error) {
	var ret = &AdminPass{}
	dat, err := tx.Get(
		fmt.Sprintf(adminPassPath, escapeString(id)),
	)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal([]byte(dat), ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func DelAdmin(tx *buntdb.Tx, id string) error {
	if _, err := tx.Delete(
		fmt.Sprintf(adminPassPath, escapeString(id)),
	); err != nil {
		return err
	}
	return nil
}
