package board

import (
	"bytes"
	"github.com/nfnt/resize"
	"go.rls.moe/nyx/http/errw"
	"go.rls.moe/nyx/resources"
	"image"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

func parseImage(reply *resources.Reply, file multipart.File, hdr *multipart.FileHeader, err error) error {
	if err != nil && err != http.ErrMissingFile {
		return err
	}
	if err == http.ErrMissingFile {
		return nil
	}
	cfg, _, err := image.DecodeConfig(file)
	if err != nil {
		return err
	}
	if cfg.Height > 8000 || cfg.Width > 8000 {
		log.Println("Somebody tried to detonate the memory!")
		return errw.MakeErrorWithTitle("Too large", "Your upload was too large")
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}
	if img.Bounds().Dx() > 8000 || img.Bounds().Dy() > 8000 {
		log.Println("Somebody tried to detonate the memory!")
		return errw.MakeErrorWithTitle("Too large", "Your upload was too large")
	}
	thumb := resize.Thumbnail(128, 128, img, resize.Lanczos3)
	imgBuf := bytes.NewBuffer([]byte{})
	err = png.Encode(imgBuf, thumb)
	if err != nil {
		return err
	}
	log.Printf("Thumb has size %d KiB", imgBuf.Len()/1024)
	reply.Thumbnail = make([]byte, imgBuf.Len())
	copy(reply.Thumbnail, imgBuf.Bytes())
	imgBuf.Reset()
	err = png.Encode(imgBuf, img)
	if err != nil {
		return err
	}
	log.Printf("Image has size %d KiB", imgBuf.Len()/1024)
	reply.Image = make([]byte, imgBuf.Len())
	copy(reply.Image, imgBuf.Bytes())
	return nil
}
