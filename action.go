package main

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/codegangsta/cli"
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

	chroot := Chroot{
		ID:      ctx.GlobalString("id"),
		Flist:   ctx.Args().First(),
		Storage: ctx.GlobalString("storage"),
	}

	if err := chroot.Start(); err != nil {
		return err
	}

	defer chroot.Stop()

	sandbox := Sandbox{
		Root:       chroot.Root(),
		UserEnv:    ctx.GlobalStringSlice("env"),
		EntryPoint: ctx.GlobalString("entry-point"),
		Args:       ctx.Args().Tail(),
	}

	//handle termination signals
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		//wait for termination signal to forward
		sig := <-ch
		if err := sandbox.Signal(sig); err != nil {
			log.Infof("process has already exited, Ctrl+C again to terminate the sandbox")
		}
	}()

	stdout, stderr, err := sandbox.Run()

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
		log.Infof("the sandbox is mounted under: %s", chroot.Root())
		log.Infof("Ctrl+C to terminate the sandbox")
		go func() {
			//wait for termination signal to terminate the sandbox
			<-ch
			log.Infof("terminating ...")
			chroot.Stop()
		}()

		chroot.Wait()
	}

	return err
}
