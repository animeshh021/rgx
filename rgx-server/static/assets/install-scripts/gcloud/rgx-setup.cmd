@echo off

echo This may take a few minutes . Please wait...

xcopy /E /I /Y %RGX_PACKAGES_DIR%\google-cloud-sdk\gcloudsdk-%RGX_PACKAGE_VERSION%\google-cloud-sdk\platform\bundledpython %RGX_PACKAGES_DIR%\google-cloud-sdk\google-cloud-sdk-python-%RGX_PACKAGE_VERSION% > temp_output.txt

for /f "tokens=*" %%i in (temp_output.txt) do set lastline=%%i
echo.
echo %lastline%

del temp_output.txt

set CLOUDSDK_PYTHON=%RGX_PACKAGES_DIR%\google-cloud-sdk\google-cloud-sdk-python-%RGX_PACKAGE_VERSION%\python.exe
set CLOUDSDK_PYTHON %CLOUDSDK_PYTHON%

echo Copied Bundledpython successfully.

set GCLOUD_CMD=%RGX_PACKAGES_DIR%\google-cloud-sdk\gcloudsdk-%RGX_PACKAGE_VERSION%\google-cloud-sdk\bin\gcloud

for /f "delims=" %%i in ('""%GCLOUD_CMD%" components copy-bundled-python"') do (
    set CLOUDSDK_PYTHON=%%i
)

call %GCLOUD_CMD% components install skaffold kubectl --quiet

echo set CLOUDSDK_PYTHON=%RGX_PACKAGES_DIR%\google-cloud-sdk\google-cloud-sdk-python-%RGX_PACKAGE_VERSION%\python.exe > %RGX_RCFILE_DIR%\use-gcloud-%RGX_PACKAGE_VERSION%.cmd
echo PATH %RGX_PACKAGE_SCRIPTDIR%\google-cloud-sdk\bin;%PATH% >> %RGX_RCFILE_DIR%\use-gcloud-%RGX_PACKAGE_VERSION%.cmd
echo To use: source %RGX_RCFILE_DIR%\use-gcloud-%RGX_PACKAGE_VERSION%.cmd

echo export PATH=%RGX_PACKAGE_SCRIPTDIR%/google-cloud-sdk/bin:%PATH% > %RGX_RCFILE_DIR%/use-gcloud-%RGX_PACKAGE_VERSION%.rc
echo To use: source %RGX_RCFILE_DIR%/use-gcloud-%RGX_PACKAGE_VERSION%.rc

echo Google Cloud SDK setup completed successfully.