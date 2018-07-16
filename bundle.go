package main

import (
	"github.com/codegangsta/cli"
	"os"
	"sync"
	"syscall"
	"os/signal"
)

type Bundle struct {
	chroot    *Chroot
	sandbox   *Sandbox
}


func (bundle *Bundle) Run(ctx *cli.Context, updateCh chan bool, wg *sync.WaitGroup){
	defer wg.Done()
	defer bundle.chroot.Stop()
	signalChan := make(chan os.Signal)
	//listen for termination signals
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)

	for {
		if err := bundle.chroot.Start(); err != nil {
			log.Error(err.Error())
			return
		}
		exitChan := make(chan struct{}, 1)
		go func (ch chan struct{}){
			stdout, stderr, err := bundle.sandbox.Run()
			if err != nil {
				if err := report(ctx, stdout, stderr, err); err != nil {
					log.Errorf("report: %s", err)
				}
			}
			ch <- struct {}{}
		}(exitChan)

		select {
		case <-updateCh:
			bundle.sandbox.Signal(syscall.SIGTERM)
			bundle.chroot.Stop()
			bundle.chroot.Wait()
			continue
		case <-exitChan:
			if ctx.GlobalBool("no-exit") {
				bundle.sandBoxNoExit(signalChan)
			}
			return
		case <- signalChan:
			if ctx.GlobalBool("no-exit") {
				bundle.sandBoxNoExit(signalChan)
			}
			return
		}
	}
	return
}

func (bundle *Bundle) Stop(ctx *cli.Context, signal os.Signal){
	if err := bundle.sandbox.Signal(signal); err != nil {
		log.Infof("process has already exited, Ctrl+C again to terminate the sandbox")
	}
}

func (bundle *Bundle) sandBoxNoExit(ch chan os.Signal) {
	log.Infof("flist exited, waiting for unmount (--no-exit was set)")
	log.Infof("the sandbox is mounted under: %s", bundle.chroot.Root())
	log.Infof("Ctrl+C to terminate the sandbox")
	go func() {
		//wait for termination signal to terminate the sandbox
		<-ch
		log.Infof("terminating ...")
		bundle.chroot.Stop()
	}()
	bundle.chroot.Wait()
}
