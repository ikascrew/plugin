package telop

import (
	"log"
	"os"
	"path/filepath"

	"github.com/ikascrew/core"
	"github.com/ikascrew/plugin/image"

	"github.com/sclevine/agouti"
	"golang.org/x/xerrors"
)

var tempFileFormat = "%s/.temp_%s.png"

type Telop struct {
	*image.Image
}

func New(args ...string) (core.Video, error) {

	if args == nil || len(args) <= 0 {
		return nil, xerrors.Errorf("arguments error need html file path.")
	}

	fn := args[0]

	info, err := os.Stat(fn)
	if err != nil {
		return nil, xerrors.Errorf("file path error[%s]: %w", fn, err)
	}

	abs, err := filepath.Abs(info.Name())
	if err != nil {
		return nil, xerrors.Errorf("file absolute path error[%s]: %w", info.Name(), err)
	}

	//PATH

	err = convert("file:///"+abs, ".temp-telop.png")
	if err != nil {
		return nil, xerrors.Errorf("convert error[%s]: %w", args[0], err)
	}

	return image.New(".temp-telop.png")
}

func convert(in string, out string) error {

	log.Println(in)

	options := agouti.ChromeOptions(
		"args", []string{
			"--headless",
			"--disable-gpu",
			"--window-size=1280,720",
			"--no-sandbox",
		})

	driver := agouti.ChromeDriver(options)
	defer driver.Stop()
	driver.Start()

	page, err := driver.NewPage()
	if err != nil {
		return xerrors.Errorf("driver new page: %w", err)
	}

	err = page.Navigate(in)
	if err != nil {
		return xerrors.Errorf("page navigate: %w", err)
	}

	err = page.Screenshot(out)
	if err != nil {
		return xerrors.Errorf("page screenshot: %w", err)
	}
	return nil
}
