#!/bin/bash

echo "This may take a few minutes. Please wait..."

GCLOUD_CMD="$RGX_PACKAGES_DIR/google-cloud-sdk/gcloudsdk-$RGX_PACKAGE_VERSION/google-cloud-sdk/bin/gcloud"

echo "Installing skaffold and kubectl components..."
$GCLOUD_CMD components install skaffold kubectl --quiet

echo "export PATH=$RGX_PACKAGES_DIR/google-cloud-sdk/bin:\$PATH" > "${RGX_RCFILE_DIR}/.gcloud-$RGX_PACKAGE_MAJORVERSION-rc"

echo "Script completed. RC file created at: ${RGX_RCFILE_DIR}/.gcloud-$RGX_PACKAGE_MAJORVERSION-rc"