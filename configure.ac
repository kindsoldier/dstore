AC_INIT(ndstore,0.1)
AM_INIT_AUTOMAKE([foreign subdir-objects])

AC_ARG_ENABLE([fdistr],
    AS_HELP_STRING([--disable-fdistr], [disable fdistr module]),,
    AC_CONFIG_SUBDIRS([fdistr])
)

AC_ARG_ENABLE([fstore],
    AS_HELP_STRING([--disable-fstore], [disable fstore module]),,
    AC_CONFIG_SUBDIRS([fstore])
)

AC_CONFIG_FILES(Makefile)
AC_OUTPUT
