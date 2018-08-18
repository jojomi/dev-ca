# dev-ca

This tool is designed to easily create a Certificate Authority (CA) and TLS/SSL Certificates
matching any domain you use for your local development.

There are instructions included for getting the CA to be locally trusted and thus make the certificates feel like for any other with a public domain and a known CA.

Linux and MacOS X are supported, `openssl` is required to be installed on the machine generating the CA and certificates.

## How to Use

Create a new Dev Certificate Authority (CA) first:

    dev-ca gen-ca

Then make certificates as needed:

    dev-ca gen-cert --domains "mydomain.test,*.local.dev"
    dev-ca gen-cert --domains "*.project1.test"

## Source

https://stackoverflow.com/a/43666288/4021739