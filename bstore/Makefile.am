
AUTOMAKE_OPTIONS = foreign no-dependencies no-installinfo

SUBDIRS = rc

SUFFIXES = .go
OBJEXT= none

sbin_PROGRAMS = bstored bstorecli 
bstored_SOURCES = bssrv/bsserv.go bssrv/bsconf.go
bstorecli_SOURCES = bscli/bscli.go

EXTRA_bstorecli_SOURCES = 
EXTRA_bstored_SOURCES = 

GOFLAGS = -ldflags="-s -w"

bstored$(EXEEXT): $(bstored_SOURCES) $(EXTRA_bstored_SOURCES)
	$(GO) build $(GOFLAGS) -o bstored$(EXEEXT) $(bstored_SOURCES)

bstorecli$(EXEEXT): $(bstorecli_SOURCES) $(EXTRA_bstorecli_SOURCES)
	$(GO) build $(GOFLAGS) -o bstorecli$(EXEEXT) $(bstorecli_SOURCES)

EXTRA_DIST = \
	bstore.conf

distclean-local:
	rm -rf autom4te.cache
	rm -rf tmp.log/ tmp.run/ tmp.data/
	rm -f bssrv/bssrv
	rm -f bscli/bscli

install-data-local:
	test -z $(DESTDIR)$(SRV_CONFDIR) || $(MKDIR_P) $(DESTDIR)$(SRV_CONFDIR)
	test -z $(DESTDIR)$(SRV_LOGDIR) || $(MKDIR_P) $(DESTDIR)$(SRV_LOGDIR)
	test -z $(DESTDIR)$(SRV_RUNDIR) || $(MKDIR_P) $(DESTDIR)$(SRV_RUNDIR)
	test -z $(DESTDIR)$(SRV_DATADIR) || $(MKDIR_P) $(DESTDIR)$(SRV_DATADIR)

