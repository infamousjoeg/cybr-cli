#!/bin/bash

main() {
    case "$1" in
        ldap)
            ldap
            ;;
        cyberark)
            cyberark
            ;;
        conjur)
            conjur
            ;;
        conjur-non-interactive)
            conjur-non-interactive
            ;;
        conjur-authn-iam)
            conjur-authn-iam
            ;;
        *)
            echo "Usage: $0 {ldap|cyberark|conjur|conjur-non-interactive|conjur-authn-iam}"
            exit 1
            ;;
    esac
}

ldap() {
    cybr logon -a ldap -b https://cyberark.joegarcia.dev -u jgarcia --concurrent
}

cyberark() {
    cybr logon -a cyberark -b https://cyberark.joegarcia.dev -u Administrator --concurrent
}

conjur() {
    cybr conjur logon -b https://conjur.joegarcia.dev -l admin -a cyberarkdemo
}

conjur-non-interactive() {
    CONJUR_ACCOUNT="cyberarkdemo"
    CONJUR_AUTHN_LOGIN="admin"
    CONJUR_APPLIANCE_URL="https://conjur.joegarcia.dev"
    CONJUR_AUTHN_API_KEY=$(summon -p ring.py --yaml 'API_KEY: !var conjur/admin' bash -c 'echo "$API_KEY"')
    export CONJUR_ACCOUNT CONJUR_AUTHN_LOGIN CONJUR_APPLIANCE_URL CONJUR_AUTHN_API_KEY
    cybr conjur list
    unset CONJUR_ACCOUNT CONJUR_AUTHN_LOGIN CONJUR_APPLIANCE_URL CONJUR_AUTHN_API_KEY
}

conjur-authn-iam() {
    CONJUR_ACCOUNT="cyberarkdemo"
    CONJUR_AUTHN_LOGIN="host/cloud/aws/ec2/735280068473/ConjurAWSRoleEC2"
    CONJUR_APPLIANCE_URL="https://conjur.joegarcia.dev"
    CONJUR_AUTHENTICATOR="authn-iam"
    CONJUR_AUTHN_SERVICE_ID="prod"
    CONJUR_AWS_TYPE="ec2"
    export CONJUR_ACCOUNT CONJUR_AUTHN_LOGIN CONJUR_APPLIANCE_URL CONJUR_AUTHENTICATOR CONJUR_AUTHN_SERVICE_ID CONJUR_AWS_TYPE
    cybr conjur list
    unset CONJUR_ACCOUNT CONJUR_AUTHN_LOGIN CONJUR_APPLIANCE_URL CONJUR_AUTHENTICATOR CONJUR_AUTHN_SERVICE_ID CONJUR_AWS_TYPE
}
main "$@"