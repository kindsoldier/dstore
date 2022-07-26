
AUTOMAKE_OPTIONS = foreign no-dependencies no-installinfo

SUBDIRS = rc

SUFFIXES = .go
OBJEXT= none

sbin_PROGRAMS = bstored bstorecli 
bstored_SOURCES = bssrv/bsserv.go 
nodist_bstored_SOURCES = bssrv/bsconf.go
bstorecli_SOURCES = bscli/bscli.go

EXTRA_bstorecli_SOURCES = \
	bsapi/blockapi.go \
	bsapi/servapi.go \
	bsapi/userapi.go

EXTRA_bstored_SOURCES = \
	bsapi/blockapi.go \
	bsapi/servapi.go \
	bsapi/userapi.go \
	\
	bssrv/bsblock/fsblock.go \
	bssrv/bsblock/fscrate.go \
	bssrv/bsblock/fstool.go \
	bssrv/bsconf.go.in \
	\
	bssrv/bscont/contauth.go \
	bssrv/bscont/contblock.go \
	bssrv/bscont/contcomm.go \
	bssrv/bscont/contserv.go \
	bssrv/bscont/contuser.go \
	\
	bssrv/bsreg/regblock.go \
	bssrv/bsreg/regcomm.go \
	bssrv/bsreg/reguser.go \
	\
	bssrv/bstore/storeblock.go \
	bssrv/bstore/storecomm.go \
	bssrv/bstore/storeuser.go


GOFLAGS = -ldflags="-s -w"

bstored$(EXEEXT): $(bstored_SOURCES) $(EXTRA_bstored_SOURCES)
	$(GO) build $(GOFLAGS) -o bstored$(EXEEXT) $(bstored_SOURCES) \
	   $(nodist_bstored_SOURCES)

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

	test -z $(DESTDIR)$(SRV_LOGDIR) || $(CHOWN) $(SRV_USER):$(SRV_GROUP) $(DESTDIR)$(SRV_LOGDIR)
	test -z $(DESTDIR)$(SRV_RUNDIR) || $(CHOWN) $(SRV_USER):$(SRV_GROUP) $(DESTDIR)$(SRV_RUNDIR)
	test -z $(DESTDIR)$(SRV_DATADIR) || $(CHOWN) $(SRV_USER):$(SRV_GROUP) $(DESTDIR)$(SRV_DATADIR)

