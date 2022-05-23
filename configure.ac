AC_INIT(ndstore,0.1)
AM_INIT_AUTOMAKE([foreign subdir-objects])

AC_ARG_ENABLE([fstore],
    AS_HELP_STRING([--disable-fstore], [disable fstore module]),,
    AC_CONFIG_SUBDIRS([fstore])
)

AC_ARG_ENABLE([bstore],
    AS_HELP_STRING([--disable-bstore], [disable bstore module]),,
    AC_CONFIG_SUBDIRS([bstore])
)

AC_CONFIG_FILES(Makefile)
AC_OUTPUT
