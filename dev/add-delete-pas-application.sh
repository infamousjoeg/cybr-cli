#!/bin/bash
set -e

# This script will create an application named "$APP_NAME".
# This application will only be able to authenticate from machine '$APP_ADDRESS' using the linux OS user '$APP_OS_USER'
# This application will have retrieve access to all accounts inside of the '$SAFE_NAME' safe
# In this example "$APP_NAME" application will have access to the MYSQL account specified below in 'Account Information'

# Application information
APP_NAME="FirstApplication"
APP_OS_USER="firstapp"
APP_ADDRESS="10.0.1.10"
SAFE_NAME="TEST_${APP_NAME}"

# Account information
DATABASE_ACCOUNT_NAME="mysql"
DATABASE_ADDRESS="10.0.1.12"
DATABASE_USERNAME="firstapp"
DATABASE_DEFAULT_PASSWORD="thisIsTheDefaultPassword"
DATABASE_PLATFORM="MySQL"

cybr logon -a ldap -b "$PAS_HOSTNAME" -u "$PAS_USERNAME"

# create the application
cybr applications add --app-id "$APP_NAME" --location "\\"
cybr application add-authn --app-id "$APP_NAME" -t OSUser -v "$APP_OS_USER"
cybr application add-authn --app-id "$APP_NAME" -t machineAddress -v "$APP_ADDRESS"

# create the safe
cybr safes add -s "$SAFE_NAME" --days 0 --desc "Safe for application "$APP_NAME""
cybr safes add-member -m "$APP_NAME" -s "$SAFE_NAME" --access-content-without-confirmation --retrieve-accounts

# add account to safe
account=$(cybr accounts add --safe "$SAFE_NAME" -n "$DATABASE_ACCOUNT_NAME" -a "$DATABASE_ADDRESS" -u "$DATABASE_USERNAME" -t password -c "$DATABASE_DEFAULT_PASSWORD" -p "$DATABASE_PLATFORM" --automatic-management)
id=$(echo "$account" | jq -r .id)

# clean up
cybr app delete -a "$APP_NAME"
cybr accounts delete -i "$id"
cybr safe delete -s "$SAFE_NAME"

cybr logoff
