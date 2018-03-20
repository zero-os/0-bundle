package main

import (
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/codegangsta/cli"
	g8ufs "github.com/zero-os/0-fs"
	"github.com/zero-os/0-fs/meta"
	"github.com/zero-os/0-fs/storage"
)

var (
	//BaseFSDir where we keep the cache and the working place of fuse
	BaseFSDir    = path.Join(os.TempDir(), "zbundle.db")
	BaseMountDir = path.Join(os.TempDir(), "zbundle")
)

func action(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
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
	id := ctx.GlobalString("id")
	root := path.Join(BaseMountDir, id)

	os.MkdirAll(root, 0755)
	if g8ufs.IsMount(root) {
		return fmt.Errorf("a sandbox is running with the same id")
	}
	// should we do this under temp?
	namespace := path.Join(BaseFSDir, id)

	metaPath, err := getMetaDB(namespace, flist)
	if err != nil {
		return err
	}

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
		Target:    root,
		MetaStore: metaStore,
		Storage:   stor,
		Cache:     path.Join(BaseFSDir, "cache"),
	}

	fs, err := g8ufs.Mount(&opt)

	defer os.RemoveAll(namespace)
	defer os.RemoveAll(fmt.Sprintf("%s.db", namespace))
	defer fs.Unmount()

	stdout, stderr, err := sandbox(root, ctx.GlobalStringSlice("env"))

	if err != nil {
		if err := report(ctx, stdout, stderr, err); err != nil {
			log.Errorf("report: %s", err)
		}
	}

	if ctx.GlobalBool("no-exit") {
		if err != nil {
			log.Errorf("%v", err)
		}
		log.Infof("flist exited, waiting for unmount (--no-exit was set)")
		log.Infof("the sandbox is mounted under: %s", root)
		fs.Wait()
	}

	return err
}
