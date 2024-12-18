import requests


def post_request(url: str, body: dict, api_key: str) -> requests.Response:
    headers = {"X-CUSTOMER-OS-API-KEY": api_key, "Content-Type": "application/json"}

    response = requests.post(url, json=body, headers=headers)

    return response


def run(api_key):
    baseUrl = "http://localhost:10000/flows/v1/flows"

    ## Create flow
    flowBody = {
        "name": "test_flow",
        "description": "my test flow",
        "triggerOn": "grain.meeting_summary.created",
    }

    resp = post_request(baseUrl, flowBody, api_key)
    if resp.status_code == 200:
        data = resp.json()
    else:
        print(f"Error: {resp.status_code}")
        print(resp.text)
        return

    flowId = data["flow"]["id"]
    triggerNodeId = data["flow"]["triggerNodeId"]

    ## Create Node
    nodeBody = {"type": "agent", "event": "timeline_event.create"}

    resp = post_request(baseUrl + "/" + flowId + "/nodes", nodeBody, api_key)
    if resp.status_code == 200:
        data = resp.json()
    else:
        print(f"Error: {resp.status_code}")
        print(resp.text)
        return

    nodeId = data["node"]["id"]

    ## Create Edge
    edgeBody = {"fromNodeId": triggerNodeId, "toNodeId": nodeId}

    resp = post_request(baseUrl + "/" + flowId + "/edges", edgeBody, api_key)
    if resp.status_code == 200:
        print("Flow setup complete for", flowId)
        return
    else:
        print(f"Error: {resp.status_code}")
        print(resp.text)
        return


if __name__ == "__main__":
    api_key = input("API Key: ")
    run(api_key)
