#!/bin/bash

# Functions
ok() { echo -e '\e[32m'$1'\e[m'; } # Green

MYSQL=`which mysql`

Q1="CREATE DATABASE IF NOT EXISTS hydra;"
Q2="CREATE DATABASE IF NOT EXISTS whisper;"

SQL="${Q1}${Q2}"

$MYSQL -uroot -p$MYSQL_ROOT_PASSWORD -e "${SQL}"

ok "Created databases hydra and whisper"