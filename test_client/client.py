import requests
import time
import os

client_cert = os.path.join('..', 'certs', 'client.crt')
client_key = os.path.join('..', 'certs', 'client.key')

# Path to the server CA certificate
server_ca_cert = os.path.join('..', 'certs', 'ca.crt')

# URL of the server
url = 'https://localhost:8443/request'

implant_info = {
    "req_type": "ImplantInfo",
    "content": {
        "ID": "implant1",
        "Hostname": "host1",
        "IntIP": "192.168.1.1",
        "ExtIP": "8.8.8.8",
        "Os": "Linux",
        "ProcessID": 1234,
        "ProcessUser": "user1",
        "ProtName": "None",
        "LastCheckIn": int(time.time()),
        "Active": True,
        "KillDate": 0,
    }
}

# implant_info = {
#     "req_type": "RemoveImplant",
#     "content": {
#         "id": "implant1",
#     }
# }

try:
    # Make a POST request to the server with client certificate and key
    response = requests.post(
        url,
        json=implant_info,
        cert=(client_cert, client_key),
        verify=server_ca_cert
    )

    # Print the response from the server
    print(response.status_code)
    print(response.text)
except requests.exceptions.RequestException as e:
    print(f"An error occurred: {e}")