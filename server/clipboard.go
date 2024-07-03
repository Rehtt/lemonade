package server

import (
	"github.com/Rehtt/lemonade/lemon"
	"github.com/atotto/clipboard"
)

type Clipboard struct{}

func (_ *Clipboard) Copy(text string, _ *struct{}) error {
	<-connCh
	// Logger instance needs to be passed here somehow?
	return clipboard.WriteAll(lemon.ConvertLineEnding(text, LineEndingOpt))
}

func (_ *Clipboard) Paste(_ struct{}, resp *string) error {
	<-connCh
	t, err := clipboard.ReadAll()
	*resp = t
	return err
}
