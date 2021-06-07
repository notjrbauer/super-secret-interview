#!/bin/bash
# generate test purpose certificates for workerservice communication

# this also works for MAC. For linux we can use "readlink -f"
realpath() {
  [[ $1 = /* ]] && echo "$1" || echo "$PWD/${1#./}"
}
SCRIPT_PATH=$(realpath "$0")
SCRIPT_DIR=$(dirname $SCRIPT_PATH)
PROJ_DIR=$SCRIPT_DIR/../

mkdir -p $PROJ_DIR/certs/

# test purpose, generate a self signed CA cert and generate two certificates for two PF instances.

gen_key() {
  openssl ecparam -name prime256v1 -out $1_param.pem
  openssl ecparam -name prime256v1 -in $1_param.pem -genkey -noout -out $1_ec.key
  # convert the EC key to PKCS key
  openssl pkcs8 -topk8 -nocrypt -in $1_ec.key -outform PEM -out $1.key
  rm -f $1_param.pem $1_ec.key
}

# TODO: configs may be merged into a common CA config file, so the shell cmd could look elegant.
gen_cert() {
  openssl req -new -key $1.key -out $1.csr -subj "/C=US/ST=IL/L=Champaign/O=Pyre Inc./OU=Interview/CN=$1" -config <(cat $SCRIPT_DIR/openssl.cnf <(printf "\nDNS.0 = $1"))
  openssl x509 -req -days 1000 -in $1.csr -CA $2.crt -CAkey $2.key -set_serial 0101 -out $1.crt -sha256 -extensions 'v3_req' -extfile <(cat $SCRIPT_DIR/openssl.cnf <(printf "\nDNS.0 = $1"))
}

gen_key root_ca
openssl req -new -x509 -days 365 -key root_ca.key -out root_ca.crt -subj "/C=US/ST=IL/L=Champaign/O=Pyre Inc./OU=Interview/CN=workerservice"

gen_key server
gen_cert server root_ca

gen_key client_read
gen_key client_write
gen_key client_unauthorized
gen_key client_unauthorized_permissions

gen_cert client_read root_ca
gen_cert client_write root_ca
gen_cert client_unauthorized root_ca
gen_cert client_unauthorized_permissions root_ca

mv *.key *.csr *.crt $PROJ_DIR/certs/
