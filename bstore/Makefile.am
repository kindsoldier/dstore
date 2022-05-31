
AUTOMAKE_OPTIONS = foreign no-dependencies no-installinfo

SUFFIXES = .go
OBJEXT= none

sbin_PROGRAMS = bstored bstorecli
bstored_SOURCES = bssrv/bsmain.go

EXTRA_bstored_SOURCES = \
	bssrv/bsbreg/block.go \
	bssrv/bsbreg/comm.go \
	bssrv/bsbreg/block_test.go \
	bssrv/bsureg/common.go \
	bssrv/bsureg/users.go \
	bssrv/bsureg/users_test.go \
	bssrv/bsbcont/comm.go \
	bssrv/bsbcont/block.go \
	bssrv/bsbcont/block_test.go \
	bssrv/bsuser/bsuser.go \
	bssrv/bsucont/comm.go \
	bssrv/bsucont/user_test.go \
	bssrv/bsucont/user.go \
	bssrv/bsblock/bsrec.go


bstorecli_SOURCES = bscli/bscli.go


GOFLAGS = -ldflags="-s -w"

bstored$(EXEEXT): $(bstored_SOURCES) $(EXTRA_bstored_SOURCES)
	$(GO) build $(GOFLAGS) -o bstored$(EXEEXT) $(bstored_SOURCES)

bstorecli$(EXEEXT): $(bstorecli_SOURCES) $(EXTRA_bstorecli_SOURCES)
	$(GO) build $(GOFLAGS) -o bstorecli$(EXEEXT) $(bstorecli_SOURCES)

EXTRA_DIST = \
	bstore.conf

clean-local:
	rm -rf autom4te.cache
	rm -rf log/ run/ data/
	rm -f bssrv/bssrv
	rm -f bscli/bscli


install-data-local:
	test -z $(DESTDIR)$(SRV_CONFDIR) || $(MKDIR_P) $(DESTDIR)$(SRV_CONFDIR)
	test -z $(DESTDIR)$(SRV_LOGDIR) || $(MKDIR_P) $(DESTDIR)$(SRV_LOGDIR)
	test -z $(DESTDIR)$(SRV_RUNDIR) || $(MKDIR_P) $(DESTDIR)$(SRV_RUNDIR)
	test -z $(DESTDIR)$(SRV_DATADIR) || $(MKDIR_P) $(DESTDIR)$(SRV_DATADIR)

PATH := $(PATH):$(HOME)/go/bin

instvis:
	$(GO) get github.com/ofabry/go-callvis
	$(GO) install github.com/ofabry/go-callvis

CVCMD = go-callvis -nostd  -skipbrowser -graphviz
CVIGNORE = nsstore/dslog
CVDIR = .
GOMAIN = ndstore/bstore/bssrv
GOROOT = ndstore

graph:
	$(CVCMD) -ignore $(CVIGNORE) -focus $(GOROOT)/bstore/ndsrv -file $(CVDIR)/bstore $(GOMAIN)
	$(CVCMD) -ignore $(CVIGNORE) -focus $(GOROOT)/dsrpc -file $(CVDIR)/dcrpc $(GOMAIN)

EXTRA_DIST += \
	bstore.svg \
	dsrpc.svg
