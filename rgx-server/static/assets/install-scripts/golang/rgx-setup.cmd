@echo off

ren go go-%RGX_PACKAGE_VERSION%
echo Go %RGX_PACKAGE_VERSION% is now present in rgx-packages\golang\go-%RGX_PACKAGE_VERSION%

PATH go-%RGX_PACKAGE_VERSION%\go\bin\;%PATH%

echo Writing use-golang and .golang-rc files...
echo @echo off> %RGX_RCFILE_DIR%\use-golang-%RGX_PACKAGE_MAJORVERSION%.cmd
echo set GOLANG_HOME=%RGX_PACKAGE_SCRIPTDIR%\go-%RGX_PACKAGE_VERSION%\go>> %RGX_RCFILE_DIR%\use-golang-%RGX_PACKAGE_MAJORVERSION%.cmd
echo PATH %%GOLANG_HOME%%\bin;%%PATH%%>> %RGX_RCFILE_DIR%\use-golang-%RGX_PACKAGE_MAJORVERSION%.cmd
echo echo Go %RGX_PACKAGE_MAJORVERSION% added to PATH>> %RGX_RCFILE_DIR%\use-golang-%RGX_PACKAGE_MAJORVERSION%.cmd

echo export GOLANG_HOME=%RGX_PACKAGE_SCRIPTDIR_MSYS%/go-%RGX_PACKAGE_VERSION%/go> %RGX_RCFILE_DIR%\.golang-%RGX_PACKAGE_MAJORVERSION%-rc
echo export PATH=$GOLANG_HOME/bin:$PATH>> %RGX_RCFILE_DIR%\.golang-%RGX_PACKAGE_MAJORVERSION%-rc
echo echo 'Go %RGX_PACKAGE_MAJORVERSION% added to PATH'>> %RGX_RCFILE_DIR%\.golang-%RGX_PACKAGE_MAJORVERSION%-rc

echo.
echo To use Go,
echo in cmd.exe, type %RGX_RCFILE_DIR%\use-golang-%RGX_PACKAGE_MAJORVERSION%.cmd
echo in git bash, type %RGX_RCFILE_DIR%\.golang-%RGX_PACKAGE_MAJORVERSION%-rc