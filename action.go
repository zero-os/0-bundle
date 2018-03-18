package main

import (
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/codegangsta/cli"
	"github.com/pborman/uuid"
	g8ufs "github.com/zero-os/0-fs"
	"github.com/zero-os/0-fs/meta"
	"github.com/zero-os/0-fs/storage"
)

func action(ctx *cli.Context) error {
	if ctx.NArg() != 2 {
		return fmt.Errorf("invalid number of arguments")
	}

	if !isRoot() {
		return fmt.Errorf("please run as root")
	}

	u, err := url.Parse(ctx.GlobalString("storage"))
	if err != nil {
		return fmt.Errorf("invalid storage url: %s", err)
	}

	flist := ctx.Args()[0]
	mount := ctx.Args()[1]

	os.MkdirAll(mount, 0755)
	// should we do this under temp?
	id := uuid.New()
	namespace := path.Join(os.TempDir(), "zbundle", id)

	metaPath, err := getMetaDB(namespace, flist)
	if err != nil {
		return err
	}
	log.Info("db path: %s", metaPath)

	metaStore, err := meta.NewRocksMeta("", metaPath)
	if err != nil {
		return err
	}

	stor, err := storage.NewARDBStorage(u)
	if err != nil {
		return err
	}

	opt := g8ufs.Options{
		Backend:   namespace,
		Target:    mount,
		MetaStore: metaStore,
		Storage:   stor,
	}

	fs, err := g8ufs.Mount(&opt)

	return fs.Wait()
	//validation
	//return nil
}
