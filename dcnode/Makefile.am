
AUTOMAKE_OPTIONS = foreign no-dependencies no-installinfo

SUFFIXES = .go
OBJEXT= none

sbin_PROGRAMS = dcnode
dcnode_SOURCES = ndsrv/ndmain.go

EXTRA_dcnode_SOURCES = \
	ndsrv/ndconf/ndconf.go \
	ndsrv/ndconf/ndconf.go.in \
	ndsrv/ndstore/ndstore.go \
	ndsrv/ndcontr/ndcontr.go \
	ndsrv/ndreg/ndreg.go

GOFLAGS = -ldflags="-s -w"

dcnode$(EXEEXT): $(dcnode_SOURCES) $(EXTRA_dcnode_SOURCES)
	$(GO) build $(GOFLAGS) -o dcnode$(EXEEXT) $(dcnode_SOURCES)

EXTRA_DIST = \
	dcnode.conf

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

CVCMD = go-callvis -nostd  -skipbrowser -graphviz
CVIGNORE = dcstore/dclog
CVDIR = .
GOMAIN = dcstore/dcnode/ndsrv
GOROOT = dcstore

graph:
	$(CVCMD) -ignore $(CVIGNORE) -focus $(GOROOT)/dcnode/ndsrv -file $(CVDIR)/dcnode $(GOMAIN)
	$(CVCMD) -ignore $(CVIGNORE) -focus $(GOROOT)/dcrpc -file $(CVDIR)/dcrpc $(GOMAIN)

EXTRA_DIST += \
	dcnode.svg \
	dcrpc.svg
	rm -f $(CVDIR)/*.gv
