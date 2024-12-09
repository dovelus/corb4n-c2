import requests

# Paths to the client certificate and key
client_cert = '../certs/client.crt'
client_key = '../certs/client.key'

# Path to the server CA certificate
server_ca_cert = '../certs/server.crt'

# URL of the server
url = 'https://localhost:8443/'

# Make a GET request to the server with client certificate and key
response = requests.get(url, cert=(client_cert, client_key), verify=server_ca_cert)

# Print the response from the server
print(response.text)