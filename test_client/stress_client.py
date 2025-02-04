import requests
import threading
import os
import time

client_cert = os.path.join('..', 'certs', 'client.crt')
client_key = os.path.join('..', 'certs', 'client.key')
server_ca_cert = os.path.join('..', 'certs', 'ca.crt')
url = 'https://localhost:8443/request'

implant_info = {
    "req_type": "GetAllImplants",
}

def send_request():
    try:
        response = requests.post(url, json=implant_info, cert=(client_cert, client_key), verify=server_ca_cert)
        print(response.status_code, response.text)
    except requests.exceptions.RequestException as e:
        print(f"An error occurred: {e}")

def stress_test(num_requests):
    threads = []
    for _ in range(num_requests):
        thread = threading.Thread(target=send_request)
        threads.append(thread)
        thread.start()

    for thread in threads:
        thread.join()

if __name__ == "__main__":
    num_requests = 10000  # Number of parallel requests
    start_time = time.time()
    stress_test(num_requests)
    end_time = time.time()
    print(f"Completed {num_requests} requests in {end_time - start_time} seconds")