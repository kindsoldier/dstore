
AUTOMAKE_OPTIONS = foreign no-dependencies no-installinfo

SUBDIRS = initrc

SUFFIXES = .go
OBJEXT= none

sbin_PROGRAMS = bstored bstorecli
bstored_SOURCES = bssrv/bsmain.go

EXTRA_bstored_SOURCES = \
	bsapi/bdeleteblock.go \
	bsapi/xcheckuser.go \
	bsapi/bsaveblock.go \
	bsapi/bblockexists.go \
	bsapi/xadduser.go \
	bsapi/bloadblock.go \
	bsapi/xlistusers.go \
	bsapi/blistblocks.go \
	bsapi/bcheckblock.go \
	bsapi/bgethello.go \
	bsapi/xupdateuser.go \
	bsapi/xdeleteuser.go \
	bsapi/bpurgeall.go \
	bscom/descr.go \
	bsfunc/funcblock.go \
	bssrv/bsbreg/bregblock.go \
	bssrv/bsbreg/bregcomm.go \
	bssrv/bsureg/uregcommon.go \
	bssrv/bsureg/uregusers.go \
	bssrv/bsbcont/bcontcomm.go \
	bssrv/bsbcont/bcontblock.go \
	bssrv/bsmain.go \
	bssrv/bsuser/bsuser.go \
	bssrv/bsucont/ucontuser.go \
	bssrv/bsucont/ucontcomm.go \
	bssrv/bsucont/ucontauth.go \
	bssrv/bsblock/bsrec.go

EXTRA_DIST = \
	bscli/bscli_test.go \
	bssrv/bsbreg/bregblock_test.go \
	bssrv/bsureg/uregusers_test.go \
	bssrv/bsbcont/bcontblock_test.go \
	bssrv/bsucont/ucontuser_test.go


bstorecli_SOURCES = bscli/bscli.go


GOFLAGS = -ldflags="-s -w"

bstored$(EXEEXT): $(bstored_SOURCES) $(EXTRA_bstored_SOURCES)
	$(GO) build $(GOFLAGS) -o bstored$(EXEEXT) $(bstored_SOURCES)

bstorecli$(EXEEXT): $(bstorecli_SOURCES) $(EXTRA_bstorecli_SOURCES)
	$(GO) build $(GOFLAGS) -o bstorecli$(EXEEXT) $(bstorecli_SOURCES)

EXTRA_DIST += \
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



