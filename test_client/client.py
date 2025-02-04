import requests
import time
import os
import threading
client_cert = os.path.join('..', 'certs', 'client.crt')
client_key = os.path.join('..', 'certs', 'client.key')

# Path to the server CA certificate
server_ca_cert = os.path.join('..', 'certs', 'ca.crt')

# URL of the server
url = 'https://localhost:8443/request'

# implant_info = {
#     "req_type": "InsertImplantInfo",
#     "content": {
#         "ID": "implant1",
#         "Hostname": "host1",
#         "IntIP": "192.168.1.1",
#         "ExtIP": "8.8.8.8",
#         "Os": "Linux",
#         "ProcessID": 1234,
#         "ProcessUser": "user1",
#         "ProtName": "None",
#         "LastCheckIn": int(time.time()),
#         "Active": True,
#         "KillDate": 0,
#     }
# }

implant_info = {
    "req_type": "UpdateImplantLastCheckin",
    "content": {
        "id": "implant1",
    }
}

# implant_info = {
#     "req_type": "GetAllImplants",
# }

try:
    def send_request():
        response = requests.post(url, json=implant_info, cert=(client_cert, client_key), verify=server_ca_cert)
        print(response.text)
    
    send_request()
    
    print("Done")
    
except requests.exceptions.RequestException as e:
    print(f"An error occurred: {e}")