package main2

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
	f, err := os.Open("Horizon.mp3")
	if err != nil {
		return err
	}

	seekTable, err := mpseek.CreateTable(f, 0.2)
	if err != nil {
		return err
	}
	fmt.Println(seekTable.Length())
	// change this to be able to find a time to seek to based on frame offsets - seek to that time before decoding
	// otherwise, seeking to the decoded stream is impossible
	r := seekTable.FindTime(0)

	fmt.Println(r.Offset)

	m := &mpa.Reader{Decoder: &mpa.Decoder{Input: f}}

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	// decode enough frames to hit the time we wish to seek to
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
