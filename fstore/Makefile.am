
AUTOMAKE_OPTIONS = foreign no-dependencies no-installinfo

SUBDIRS = initrc

SUFFIXES = .go
OBJEXT= none

sbin_PROGRAMS = fstored fstorecli
fstored_SOURCES = fssrv/fsmain.go

EXTRA_fstored_SOURCES = \
	    fsapi/fdeletefile.go \
	    fsapi/fgethello.go \
	    fsapi/flistfiles.go \
	    fsapi/bupdatebstore.go \
	    fsapi/ucheckuser.go \
	    fsapi/blistbstores.go \
	    fsapi/fsavefile.go \
	    fsapi/uupdateuser.go \
	    fsapi/udeleteuser.go \
	    fsapi/bdeletebstore.go \
	    fsapi/floadfile.go \
	    fsapi/uadduser.go \
	    fsapi/baddbstore.go \
	    fsapi/ulistusers.go \
	    fssrv/fscont/contruser.go \
	    fssrv/fscont/contrbstore.go \
	    fssrv/fscont/contrfile.go \
	    fssrv/fscont/contrcomm.go \
	    fssrv/fsmain.go \
	    fssrv/fsrec/recfile.go \
	    fssrv/fsrec/recuser.go \
	    fssrv/fsrec/recfiledistr.go \
	    fssrv/fsrec/reccomm.go \
	    fssrv/fsrec/recbstore.go \
	    fssrv/fsconf/fsconf.go \
	    fssrv/fsfile/fsbatch.go \
	    fssrv/fsfile/fsfile.go \
	    fssrv/fsfile/fsblock.go \
	    fssrv/fsreg/regbatch.go \
	    fssrv/fsreg/regfile.go \
	    fssrv/fsreg/reguser.go \
	    fssrv/fsreg/regentry.go \
	    fssrv/fsreg/regbstore.go \
	    fssrv/fsreg/regblock.go \
	    fssrv/fsreg/regcomm.go

fstorecli_SOURCES = fscli/fscli.go


GOFLAGS = -ldflags="-s -w"

fstored$(EXEEXT): $(fstored_SOURCES) $(EXTRA_fstored_SOURCES)
	$(GO) build $(GOFLAGS) -o fstored$(EXEEXT) $(fstored_SOURCES)

fstorecli$(EXEEXT): $(fstorecli_SOURCES) $(EXTRA_fstorecli_SOURCES)
	$(GO) build $(GOFLAGS) -o fstorecli$(EXEEXT) $(fstorecli_SOURCES)

EXTRA_DIST = \
	fstore.conf

dictclean-local:
	rm -rf autom4te.cache
	rm -rf log/ run/ data/
	rm -f fssrv/fssrv
	rm -f fscli/fscli

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
