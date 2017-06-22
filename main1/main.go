package main

import (
	"errors"
	"io"
	"log"
	"os"

	"github.com/hajimehoshi/oto"

	"fmt"

	"github.com/hajimehoshi/go-mp3"
)

func run() error {
	f, err := os.Open("Horizon.mp3")
	if err != nil {
		return err
	}
	defer f.Close()
	d, err := mp3.Decode(f)
	if err != nil {
		return err
	}
	defer d.Close()
	fmt.Println(d.SampleRate())

	p, err := oto.NewPlayer(d.SampleRate(), 2, 2, 65536)
	if err != nil {
		return err
	}

	if err := Play(p, d, nil); err != nil && err != io.EOF {
		return err
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func Play(dst io.Writer, src io.Reader, buf []byte) (err error) {
	written := int64(0)
	if buf == nil {
		buf = make([]byte, 32*1024)
	}
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = errors.New("Short write")
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	fmt.Println(written)
	return
}
