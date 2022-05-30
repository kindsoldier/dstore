
AUTOMAKE_OPTIONS = foreign no-dependencies no-installinfo

SUFFIXES = .go
OBJEXT= none

sbin_PROGRAMS = fstored fstorecli
fstored_SOURCES = fssrv/fsmain.go

EXTRA_fstored_SOURCES = \
	fssrv/fscont/user.go \
	fssrv/fscont/comm_test.go \
	fssrv/fscont/file_test.go \
	fssrv/fscont/file.go \
	fssrv/fscont/comm.go \
	fssrv/fsrec/user.go \
	fssrv/fsrec/file.go \
	fssrv/fsrec/comm.go \
	fssrv/fsrec/comm_test.go \
	fssrv/fsrec/file_test.go \
	fssrv/fsconf/fsconf.go.in \
	fssrv/fsfile/file.go \
	fssrv/fsfile/block.go \
	fssrv/fsfile/batch_test.go \
	fssrv/fsfile/file_test.go \
	fssrv/fsfile/block_test.go \
	fssrv/fsfile/batch.go \
	fssrv/fsreg/users.go \
	fssrv/fsreg/comm.go \
	fssrv/fsreg/entries_test.go \
	fssrv/fsreg/bstores_test.go \
	fssrv/fsreg/entries.go \
	fssrv/fsreg/bstores.go \
	fssrv/fsreg/files_test.go \
	fssrv/fsreg/users_test.go \
	fssrv/fsreg/files.go


fstorecli_SOURCES = fscli/fscli.go


GOFLAGS = -ldflags="-s -w"

fstored$(EXEEXT): $(fstored_SOURCES) $(EXTRA_fstored_SOURCES)
	$(GO) build $(GOFLAGS) -o fstored$(EXEEXT) $(fstored_SOURCES)

fstorecli$(EXEEXT): $(fstorecli_SOURCES) $(EXTRA_fstorecli_SOURCES)
	$(GO) build $(GOFLAGS) -o fstorecli$(EXEEXT) $(fstorecli_SOURCES)

EXTRA_DIST = \
	fstore.conf

clean-local:
	rm -rf autom4te.cache
	rm -rf log/ run/ data/

install-data-local:
	test -z $(DESTDIR)$(SRV_CONFDIR) || $(MKDIR_P) $(DESTDIR)$(SRV_CONFDIR)
	test -z $(DESTDIR)$(SRV_LOGDIR) || $(MKDIR_P) $(DESTDIR)$(SRV_LOGDIR)
	test -z $(DESTDIR)$(SRV_RUNDIR) || $(MKDIR_P) $(DESTDIR)$(SRV_RUNDIR)
	test -z $(DESTDIR)$(SRV_DATADIR) || $(MKDIR_P) $(DESTDIR)$(SRV_DATADIR)

PATH := $(PATH):$(HOME)/go/bin

instvis:
	$(GO) get github.com/ofabry/go-callvis
	$(GO) install github.com/ofabry/go-callvis

CVCMD = go-callvis -nostd -skipbrowser
CVIGNORE = nsstore/dslog
CVDIR = ./docs
GOMAIN = ndstore/fstore/fssrv
GOROOT = ndstore

graph:
	$(CVCMD) -ignore $(CVIGNORE) -focus $(GOROOT)/fstore/fssrv -file $(CVDIR)/fstore $(GOMAIN)
	$(CVCMD) -ignore $(CVIGNORE) -focus $(GOROOT)/fstore/fssrv/fsrec -file $(CVDIR)/fscont $(GOMAIN)
	$(CVCMD) -ignore $(CVIGNORE) -focus $(GOROOT)/fstore/fssrv/fsrec -file $(CVDIR)/fsrec $(GOMAIN)
	$(CVCMD) -ignore $(CVIGNORE) -focus $(GOROOT)/fstore/fssrv/fsreg -file $(CVDIR)/fsreg $(GOMAIN)
	rm -f $(CVDIR)/*.gv

EXTRA_DIST += \
	$(CVDIR)/fstore.svg
