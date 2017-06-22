package main

import (
	"io"
	"os"

	"fmt"

	"github.com/hajimehoshi/oto"
	"github.com/korandiz/mpa"
	"github.com/korandiz/mpseek"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
	}
}

func run() error {
	f, err := os.Open("Keyscape.mp3")
	if err != nil {
		return err
	}

	seekTable, err := mpseek.CreateTable(f, 1)
	if err != nil {
		return err
	}

	r := seekTable.FindTime(120)

	fmt.Println(r.Offset)

	m := &mpa.Reader{Decoder: &mpa.Decoder{Input: f}}

	fmt.Println(r.Time)

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	for i := 0; i < r.WarmUp; i++ {
		err = m.Decoder.DecodeFrame()
		if err != nil {
			return err
		}
	}

	p, err := oto.NewPlayer(m.Decoder.SamplingFrequency(), 2, 2, 65536)
	if err != nil {
		return err
	}

	_, err = io.Copy(p, m)
	return err
}
