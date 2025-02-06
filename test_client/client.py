import requests
import os
import json

client_cert = os.path.join('..', 'certs', 'client.crt')
client_key = os.path.join('..', 'certs', 'client.key')
server_ca_cert = os.path.join('..', 'certs', 'ca.crt')
url = 'https://localhost:8443/request'

headers = { 'Content-Type': 'application/json' }

# data = {
#     "req_type": "UploadTaskResults",
#     "content": {
#         "id": "implant1",
#         "task_id": "cuid3g3a6s20ad4c1nf0",
#         "result": {
#             "status": "success",
#             "output": "SKYNUT\\Dovelus"
#         }
#     }
# }

data = {
    "req_type": "UploadTaskResults",
    "content": json.dumps({
        "id": "implant1",
        "task_id": "cuid4aja6s20ad4c1nfg",
        "result": {
            "status": "success",
            "output": {
                "file_name": "example.txt",
                "file_type": "plain/text" #MIME type
            }
        }
    })
}

# File to upload
file_path = "example.txt"
files = {
    "file": open(file_path, "rb")
}
print(json.dumps(data, indent=4))
try:
    #response = requests.post(url, headers=headers, json=data, cert=(client_cert, client_key), verify=server_ca_cert)
    response = requests.post(url, data=data, files=files, cert=(client_cert, client_key), verify=server_ca_cert)
    print(response.status_code)
    print(response.text)
except Exception as e:
    print(f"An error occurred: {e}")