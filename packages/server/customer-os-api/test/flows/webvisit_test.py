import requests


def post_request(url: str, body: dict, api_key: str) -> requests.Response:
    headers = {"X-CUSTOMER-OS-API-KEY": api_key, "Content-Type": "application/json"}

    response = requests.post(url, json=body, headers=headers)

    return response


def setup(api_key: str) -> str:
    baseUrl = "http://localhost:10000/flows/v1/"

    ## Create flow
    flowBody = {
        "name": "webvisit_test_flow",
        "description": "my webvisit test flow",
        "triggerOn": "reveal.website_visit.new",
    }

    resp = post_request(baseUrl + "flows", flowBody, api_key)
    if resp.status_code == 200:
        data = resp.json()
    else:
        print(f"Error: {resp.status_code}")
        print(resp.text)
        exit(1)

    flowId = data["flow"]["id"]
    triggerNodeId = data["flow"]["triggerNodeId"]

    ## Create Node
    nodeBody = {"type": "agent", "event": "slack.notify"}

    resp = post_request(baseUrl + "flows/" + flowId + "/nodes", nodeBody, api_key)
    if resp.status_code == 200:
        data = resp.json()
    else:
        print(f"Error: {resp.status_code}")
        print(resp.text)
        exit(1)

    nodeId = data["node"]["id"]

    ## Create Edge
    edgeBody = {"fromNodeId": triggerNodeId, "toNodeId": nodeId}

    resp = post_request(baseUrl + "flows/" + flowId + "/edges", edgeBody, api_key)
    if resp.status_code != 200:
        print(f"Error: {resp.status_code}")
        print(resp.text)
        exit(1)

    ## Create Grain URL
    webhookBody = {"integration": "reveal"}

    resp = post_request(baseUrl + "hooks", webhookBody, api_key)
    if resp.status_code == 201:
        print("Flow setup complete for", flowId)
        data = resp.json()
    else:
        print(f"Error: {resp.status_code}")
        print(resp.text)
        exit(1)

    webhook = convert_to_local_url(data["hook"]["url"])
    print("Webhook url", webhook)
    print("Sending test payload...")
    return webhook


def test_payload(webhookUrl: str, api_key: str):

    payload = {
        "data": {
            "end_datetime": "2024-11-18T16:33:30Z",
            "ical_uid": "4fhPuduKu5gr2L3rfxZJxx@Cal.com",
            "id": "76270175-a8e5-4fd1-b681-b3330ec440ef",
            "intelligence_notes_md": "**Customer Budget** \n [(22:25)](https://grain.com/share/recording/76270175-a8e5-4fd1-b681-b3330ec440ef/zGaUJe4bsuFaD8C6RLSF4Q0FSMz1G80OFNn7JSLB?t=1345000) - Joseph expresses a desire to keep using the product but needs more full-featured solutions, indicating a willingness to pay for it once it's ready.\n\n**Competition** \n Joseph \"mentioned\" CustomerOS, HubSpot, Microsoft, Calendly\n\n**Sales Meeting Outcome** \n Matt Brown and Joseph Frantz discussed the current challenges with their system's data import and feature set. Joseph expresses interest in continuing to use the product while running it parallel to HubSpot, indicating a need for more comprehensive features. They agree to maintain communication via Slack and plan to connect weekly for updates on new features as they are developed. Joseph is committed to transitioning fully once the necessary enhancements are implemented.\n\n**Customer Needs** \n [(01:51)](https://grain.com/share/recording/76270175-a8e5-4fd1-b681-b3330ec440ef/zGaUJe4bsuFaD8C6RLSF4Q0FSMz1G80OFNn7JSLB?t=111000) - Matt mentions challenges with Microsoft verification, having multiple cases closed without contact.\n [(03:30)](https://grain.com/share/recording/76270175-a8e5-4fd1-b681-b3330ec440ef/zGaUJe4bsuFaD8C6RLSF4Q0FSMz1G80OFNn7JSLB?t=210000) - Joseph highlights issues with capturing random marketing and cold inbound leads.\n [(21:41)](https://grain.com/share/recording/76270175-a8e5-4fd1-b681-b3330ec440ef/zGaUJe4bsuFaD8C6RLSF4Q0FSMz1G80OFNn7JSLB?t=1301000) - Joseph expresses the need for a more full-featured solution while currently running parallel with HubSpot.\n\n**Next Steps** \n [(01:09)](https://grain.com/share/recording/76270175-a8e5-4fd1-b681-b3330ec440ef/zGaUJe4bsuFaD8C6RLSF4Q0FSMz1G80OFNn7JSLB?t=69925) - Joseph Frantz to go to customeros and share updates\n [(07:55)](https://grain.com/share/recording/76270175-a8e5-4fd1-b681-b3330ec440ef/zGaUJe4bsuFaD8C6RLSF4Q0FSMz1G80OFNn7JSLB?t=475503) - Matt Brown to check into Microsoft import logs for issues\n [(09:00)](https://grain.com/share/recording/76270175-a8e5-4fd1-b681-b3330ec440ef/zGaUJe4bsuFaD8C6RLSF4Q0FSMz1G80OFNn7JSLB?t=540340) - Joseph Frantz to add email addresses to Slack thread for search\n [(13:07)](https://grain.com/share/recording/76270175-a8e5-4fd1-b681-b3330ec440ef/zGaUJe4bsuFaD8C6RLSF4Q0FSMz1G80OFNn7JSLB?t=787043) - Matt Brown to pull in another engineer to assist with Microsoft import\n [(21:16)](https://grain.com/share/recording/76270175-a8e5-4fd1-b681-b3330ec440ef/zGaUJe4bsuFaD8C6RLSF4Q0FSMz1G80OFNn7JSLB?t=1276827) - Matt Brown to provide ETA on calendar integration\n [(23:51)](https://grain.com/share/recording/76270175-a8e5-4fd1-b681-b3330ec440ef/zGaUJe4bsuFaD8C6RLSF4Q0FSMz1G80OFNn7JSLB?t=1431721) - Joseph Frantz to schedule weekly 30-minute calls for updates\n [(27:40)](https://grain.com/share/recording/76270175-a8e5-4fd1-b681-b3330ec440ef/zGaUJe4bsuFaD8C6RLSF4Q0FSMz1G80OFNn7JSLB?t=1660463) - Matt Brown to continue development on product features discussed\n [(27:49)](https://grain.com/share/recording/76270175-a8e5-4fd1-b681-b3330ec440ef/zGaUJe4bsuFaD8C6RLSF4Q0FSMz1G80OFNn7JSLB?t=1669231) - Matt Brown to ping Joseph Frantz when there's something to show\n [(27:53)](https://grain.com/share/recording/76270175-a8e5-4fd1-b681-b3330ec440ef/zGaUJe4bsuFaD8C6RLSF4Q0FSMz1G80OFNn7JSLB?t=1673423) - Joseph Frantz to stay available on Slack for communication\n\n**Stakeholders** \n [(22:25)](https://grain.com/share/recording/76270175-a8e5-4fd1-b681-b3330ec440ef/zGaUJe4bsuFaD8C6RLSF4Q0FSMz1G80OFNn7JSLB?t=1345000) - Joseph mentioned they will run the product in parallel with HubSpot and slowly migrate over as needed.\n [(28:09)](https://grain.com/share/recording/76270175-a8e5-4fd1-b681-b3330ec440ef/zGaUJe4bsuFaD8C6RLSF4Q0FSMz1G80OFNn7JSLB?t=1689000) - Joseph stated he is the target market for the solution and wants it, emphasizing the need for completion of the product.\n [(28:31)](https://grain.com/share/recording/76270175-a8e5-4fd1-b681-b3330ec440ef/zGaUJe4bsuFaD8C6RLSF4Q0FSMz1G80OFNn7JSLB?t=1711000) - Joseph expressed interest in the solution being built and indicated a willingness to pay for it once it's fully developed. ",
            "owners": "['matt@customeros.ai']",
            "participants": "[{'name': 'Jonty Knox', 'scope': 'internal', 'email': 'jonty@customeros.ai', 'confirmed_attendee': True}, {'name': 'Matt Brown', 'scope': 'internal', 'email': 'matt@customeros.ai', 'confirmed_attendee': True}, {'name': 'Joseph Frantz', 'scope': 'unknown', 'email': None, 'confirmed_attendee': True}, {'name': 'jf', 'scope': 'external', 'email': 'jf@timesentry.ai', 'confirmed_attendee': False}]",
            "public_thumbnail_url": "https://media.grain.com/public_thumbnails/recordings/NzYyNzAxNzUtYThlNS00ZmQxLWI2ODEtYjMzMzBlYzQ0MGVmOnpHYVVKZTRic3VGYUQ4QzZSTFNGNFEwRlNNejFHODBPRk5uN0pTTEI=",
            "public_url": "https://grain.com/share/recording/76270175-a8e5-4fd1-b681-b3330ec440ef/zGaUJe4bsuFaD8C6RLSF4Q0FSMz1G80OFNn7JSLB",
            "start_datetime": "2024-11-18T16:04:44Z",
            "tags": "[]",
            "title": "Joseph L Frantz",
            "url": "https://grain.com/share/recording/76270175-a8e5-4fd1-b681-b3330ec440ef/zGaUJe4bsuFaD8C6RLSF4Q0FSMz1G80OFNn7JSLB",
        },
        "type": "recording_added",
        "user_id": "de0a2a8e-00c7-4ff0-ab7e-1cc2adf52dbe",
    }

    resp = post_request(webhookUrl, payload, api_key)
    if resp.status_code == 202:
        print("Payload accepted")
    else:
        print(f"Error: {resp.status_code}")
        print(resp.text)
        exit(1)


def convert_to_local_url(prod_url: str) -> str:
    parsed = prod_url.replace("https://api.customeros.ai", "http://localhost:10000")
    return parsed


if __name__ == "__main__":
    api_key = input("API Key: ")
    print("Creating test flow...")
    webhook = setup(api_key)
    print("Sending test payload...")
    test_payload(webhook, api_key)
