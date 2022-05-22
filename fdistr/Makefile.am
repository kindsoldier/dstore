
AUTOMAKE_OPTIONS = foreign no-dependencies no-installinfo

SUFFIXES = .go
OBJEXT= none

sbin_PROGRAMS = fdistrd fdistrcli
fdistrd_SOURCES = fdsrv/fdmain.go

EXTRA_fdistrd_SOURCES = \
	fdsrv/fdconf/fdconf.go \
	fdsrv/fdconf/fdconf.go.in \
	fdsrv/fdcont/fdcont.go

fdistrcli_SOURCES = fdcli/fdcli.go


GOFLAGS = -ldflags="-s -w"

fdistrd$(EXEEXT): $(fdistrd_SOURCES) $(EXTRA_fdistrd_SOURCES)
	$(GO) build $(GOFLAGS) -o fdistrd$(EXEEXT) $(fdistrd_SOURCES)

fdistrcli$(EXEEXT): $(fdistrcli_SOURCES) $(EXTRA_fdistrcli_SOURCES)
	$(GO) build $(GOFLAGS) -o fdistrcli$(EXEEXT) $(fdistrcli_SOURCES)

EXTRA_DIST = \
	fdistr.conf

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
CVIGNORE = nsstore/dslog
CVDIR = .
GOMAIN = ndstore/fdistr/fdsrv
GOROOT = ndstore

graph:
	$(CVCMD) -ignore $(CVIGNORE) -focus $(GOROOT)/fdistr/ndsrv -file $(CVDIR)/fdistr $(GOMAIN)
	$(CVCMD) -ignore $(CVIGNORE) -focus $(GOROOT)/dsrpc -file $(CVDIR)/dcrpc $(GOMAIN)

EXTRA_DIST += \
	fdistr.svg \
	dsrpc.svg
