#!/bin/bash

if [ $# -ne 2 ]; then
  echo "Uso: $0 <nombre_del_archivo_de_salida> <cantidad_de_clientes>"
  exit 1
fi

output_file=$1
num_clients=$2

python3 generador.py $output_file $num_clients

