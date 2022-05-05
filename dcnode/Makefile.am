
AUTOMAKE_OPTIONS = foreign no-dependencies no-installinfo

SUFFIXES = .go
OBJEXT= none

sbin_PROGRAMS = dcnode
dcnode_SOURCES = dcnode.go

GOFLAGS = -ldflags="-s -w"

EXTRA_dcnode_SOURCES = \
	dcnode.go \
	config/config.go

#config/config.go: config/config.go.in
#	./configure.status $<

dcnode$(EXEEXT): $(dcnode_SOURCES) $(EXTRA_dcnode_SOURCES)
	$(GO) build $(GOFLAGS) -o dcnode$(EXEEXT) $(dcnode_SOURCES)

clean-local:
	rm -rf autom4te.cache
	rm -rf log/ run/ data/

install-data-local:
	test -z $(DESTDIR)$(SRV_CONFDIR) || $(MKDIR_P) $(DESTDIR)$(SRV_CONFDIR)
	test -z $(DESTDIR)$(SRV_LOGDIR) || $(MKDIR_P) $(DESTDIR)$(SRV_LOGDIR)
	test -z $(DESTDIR)$(SRV_RUNDIR) || $(MKDIR_P) $(DESTDIR)$(SRV_RUNDIR)
	test -z $(DESTDIR)$(SRV_DATADIR) || $(MKDIR_P) $(DESTDIR)$(SRV_DATADIR)
