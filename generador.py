import sys

def generar_docker_compose(nombre_archivo, cantidad_clientes):

    servicios = """name: tp0
services:
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
    networks:
      - testing_net
    volumes:
      - type: bind
        source: ./server/config.ini
        target: /config.ini
  
"""

    for i in range(1, cantidad_clientes + 1):
        cliente_id = f'{i}'
        CLIENTE = f"""  client{cliente_id}:
    container_name: client{cliente_id}
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID={cliente_id}
    networks:
      - testing_net
    depends_on:
      - server
    volumes:
      - type: bind
        source: ./client/config.yaml
        target: /config.yaml
        
      - type: bind
        source: ./.data/agency-{cliente_id}.csv
        target: /agency.csv
"""
      
        servicios += CLIENTE + '\n'
        
    NETWORKS = """networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24 
    """

    with open(nombre_archivo, 'w') as archivo:
        texto_final = servicios + NETWORKS
        archivo.write(texto_final)

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Uso: python3 generador.py <nombre_del_archivo_de_salida> <cantidad_de_clientes>")
        sys.exit(1)

    nombre_archivo = sys.argv[1]
    cantidad_clientes = int(sys.argv[2])

    generar_docker_compose(nombre_archivo, cantidad_clientes)
    print(f"Archivo {nombre_archivo} generado con {cantidad_clientes} clientes.")

