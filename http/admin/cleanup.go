package admin

import (
	"github.com/tidwall/buntdb"
	"go.rls.moe/nyx/http/errw"
	"go.rls.moe/nyx/http/middle"
	"net/http"
)

func handleCleanup(w http.ResponseWriter, r *http.Request) {
	db := middle.GetDB(r)
	err := db.Update(func(tx *buntdb.Tx) error {
		/* Insert cleanup codes here */
		return nil
	})
	err = db.Shrink()
	if err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}
}
