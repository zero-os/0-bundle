#!/bin/bash
set -ex

ncpu=$(grep 'model name' /proc/cpuinfo | wc -l)

apt-get update
apt-get install -y curl git build-essential
apt-get install -y libgflags-dev libsnappy-dev zlib1g-dev libbz2-dev liblz4-dev # libzstd-dev

# install go 1.8 (needed by fuse)
curl https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz > /tmp/go1.8.linux-amd64.tar.gz
tar -C /usr/local -xzf /tmp/go1.8.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
mkdir -p /gopath
export GOPATH=/gopath

# compiling rocksdb
rversion="439855a7743a10bc036c7bc05563521500b83068"
curl -L https://github.com/facebook/rocksdb/archive/${rversion}.tar.gz > /tmp/rocksdb.tar.gz
tar -xf /tmp/rocksdb.tar.gz -C /tmp/

pushd /tmp/rocksdb-${rversion}
PORTABLE=1 make -j ${ncpu} static_lib
PORTABLE=1 make install
popd

# (patched) gorocksdb
go get -v github.com/gigforks/gorocksdb

# 0-bundle
mkdir -p $GOPATH/src/github.com/zero-os/0-bundle
cp -ar /0-bundle $GOPATH/src/github.com/zero-os/

pushd $GOPATH/src/github.com/zero-os/0-bundle
go get -v ./...

# Fix Makefile to produce static build
sed -i 's/-lrocks/-static -lrocks/' Makefile
make build

# reduce binary size
strip -s zbundle

# print shared libs
ldd zbundle || true
popd

mkdir -p /tmp/root/bin
cp $GOPATH/src/github.com/zero-os/0-bundle/zbundle /tmp/root/bin

mkdir -p /tmp/archives/
tar -czf "/tmp/archives/0-bundle.tar.gz" -C /tmp/root .
