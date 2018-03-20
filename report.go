package main

import (
	"crypto/tls"
	"net/url"
	"time"

	"encoding/json"

	"github.com/codegangsta/cli"
	"github.com/go-redis/redis"
)

//Report structure
type Report struct {
	//ID of the sandbox
	ID string `json:"-"`

	//Stdout captured stdout of the sandbox process
	Stdout string `json:"stdout"`

	//Stderr captured stderr of the sandbox process
	Stderr string `json:"stderr"`

	//Error containers the str representation of the exit error (ususally exit code)
	Error string `json:"error"`

	ExitTime time.Time `json:"time"`
}

//Reporter defines a reporter callback
type Reporter func(u *url.URL, report *Report) error

var (
	reporters = map[string]Reporter{
		"redis":     redisReporter,
		"redis+tls": redisReporter,
	}
)

func report(ctx *cli.Context, stdout, stderr []byte, result error) error {
	recievers := ctx.GlobalStringSlice("report")

	report := Report{
		ID:       ctx.GlobalString("id"),
		ExitTime: time.Now(),

		//we convert both stdout and stderr to string, to make it humand readable
		//in the json representation
		Stdout: string(stdout),
		Stderr: string(stderr),
		Error:  result.Error(),
	}

	for _, receiver := range recievers {
		u, err := url.Parse(receiver)
		if err != nil {
			log.Errorf("invalid report url: %v", err)
			continue
		}

		reporter, ok := reporters[u.Scheme]
		if !ok {
			log.Errorf("no reporter handler for scheme '%v'", u.Scheme)
			continue
		}

		if err := reporter(u, &report); err != nil {
			log.Errorf("failed to report to '%v': %v", receiver, err)
		}
	}

	return nil
}

//redis reporter, report error to redis
//accepts the following URL
// redis[+tls]://[password@]host:port
func redisReporter(u *url.URL, report *Report) error {
	var config *tls.Config
	if u.Scheme == "redis+tls" {
		config = &tls.Config{
			ServerName:         u.Hostname(),
			InsecureSkipVerify: true, // TODO: this should be removed in production
		}
	}

	var password string
	if u.User != nil {
		password = u.User.String()
	}

	cl := redis.NewClient(&redis.Options{
		Addr:      u.Host,
		Password:  password,
		TLSConfig: config,
	})

	defer cl.Close()
	data, err := json.Marshal(report)
	if err != nil {
		return err
	}

	return cl.LPush(report.ID, string(data)).Err()
}
