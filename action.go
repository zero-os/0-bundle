package main

import (
	"fmt"
	"os"
	"path"
	"github.com/codegangsta/cli"
	"sync"
)

var (
	//BaseFSDir where we keep the cache and the working place of fuse
	BaseFSDir    = path.Join(os.TempDir(), "zbundle.db")
	BaseMountDir = path.Join(os.TempDir(), "zbundle")
)

func action(ctx *cli.Context) error {
	if ctx.NArg() < 1 {
		return fmt.Errorf("flist is missing")
	}

	if !isRoot() {
		return fmt.Errorf("please run as root")
	}

	updateChan := make(chan bool)
	// Check for updates
	flist := ctx.Args().First()
	checkFlistUpdate(flist, updateChan)

	// Start the sandbox
	chroot := Chroot{
		ID:      ctx.GlobalString("id"),
		Flist:   ctx.Args().First(),
		Storage: ctx.GlobalString("storage"),
	}

	sandbox := Sandbox{
		Root:       chroot.Root(),
		UserEnv:    ctx.GlobalStringSlice("env"),
		EntryPoint: ctx.GlobalString("entry-point"),
		Args:       ctx.Args().Tail(),
	}
	bundle := Bundle{
		chroot: &chroot,
		sandbox: &sandbox,
	}
	var wg sync.WaitGroup
	wg.Add(1)
	bundle.Run(ctx, updateChan, &wg)
	wg.Wait()

	return nil
}
