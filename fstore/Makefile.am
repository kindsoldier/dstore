
AUTOMAKE_OPTIONS = foreign no-dependencies no-installinfo

SUBDIRS = rc

SUFFIXES = .go
OBJEXT= none

sbin_PROGRAMS = fstored fstorecli 
fstored_SOURCES = fssrv/fsserv.go fssrv/fsconf.go
fstorecli_SOURCES = fscli/fscli.go

EXTRA_fstored_SOURCES = 

GOFLAGS = -ldflags="-s -w"

fstored$(EXEEXT): $(fstored_SOURCES) $(EXTRA_fstored_SOURCES)
	$(GO) build $(GOFLAGS) -o fstored$(EXEEXT) $(fstored_SOURCES)

fstorecli$(EXEEXT): $(fstorecli_SOURCES) $(EXTRA_fstorecli_SOURCES)
	$(GO) build $(GOFLAGS) -o fstorecli$(EXEEXT) $(fstorecli_SOURCES)

EXTRA_DIST = \
	fstore.conf

distclean-local:
	rm -rf autom4te.cache
	rm -rf tmp.log/ tmp.run/ tmp.data/
	rm -f fssrv/fssrv
	rm -f fscli/fscli

install-data-local:
	test -z $(DESTDIR)$(SRV_CONFDIR) || $(MKDIR_P) $(DESTDIR)$(SRV_CONFDIR)
	test -z $(DESTDIR)$(SRV_LOGDIR) || $(MKDIR_P) $(DESTDIR)$(SRV_LOGDIR)
	test -z $(DESTDIR)$(SRV_RUNDIR) || $(MKDIR_P) $(DESTDIR)$(SRV_RUNDIR)
	test -z $(DESTDIR)$(SRV_DATADIR) || $(MKDIR_P) $(DESTDIR)$(SRV_DATADIR)

