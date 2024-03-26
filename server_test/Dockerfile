FROM python:3.9.7-slim

# Copiar los archivos al contenedor
COPY server_checker.py /server_checker.py
COPY config.ini /config.ini

# Instalar netcat en la imagen
RUN apt-get update && apt-get install -y netcat

# Definir el punto de entrada para el contenedor
CMD ["python", "/server_checker.py"]
