import requests


class CustomerOsApiCient:

    def __init__(self, base_url, api_key, openline_username):
        self.base_url = base_url
        self.api_key = api_key
        self.openline_username = openline_username

    def get_user(self, user_id):
        url = f"{self.base_url}/query/"
        headers = {
            "X-Openline-API-KEY": self.api_key,
            "X-Openline-USERNAME": self.openline_username,
        }
        query = '''
            query ($id: ID!) {
  				user(id: $id){
    				firstName
                    lastName
  				}
  			}
        '''
        variables = {
            'id': user_id
        }
        response = requests.post(url, json={'query': query, 'variables': variables}, headers=headers)
        result = response.json()
        # Check for errors in the response
        if 'errors' in result:
            print(result)
            return {'firstName': '', 'lastName': ''}
        return {'firstName': result['data']['user']['firstName'], 'lastName': result['data']['user']['lastName']}

    def get_contact(self, contact_id):
        url = f"{self.base_url}/query/"

        headers = {
            "X-Openline-API-KEY": self.api_key,
            "X-Openline-USERNAME": self.openline_username,
        }

        query = '''
			query ($id: ID!) {
  				contact(id: $id){
    				firstName
                    lastName
                    organizations{
                        content {
                            id
                            name
                            industry
                            description
                        }
                    }
                    jobRoles {
                        organization {
                            id
                            name
                            industry
                            description
                        }
                    }
  				}
            }
        '''

        variables = {
            'id': contact_id
        }

        response = requests.post(url, json={'query': query, 'variables': variables}, headers=headers)
        result = response.json()

        # Check for errors in the response
        if 'errors' in result:
            print(result)
            return {'firstName': '', 'lastName': ''}
        response_obj =  {'firstName': result['data']['contact']['firstName'], 'lastName': result['data']['contact']['lastName']}

        response_obj['organizations'] = []
        if len(result['data']['contact']['organizations']['content']) > 0:
            response_obj['organizations'] = result['data']['contact']['organizations']['content']

        if len(result['data']['contact']['jobRoles']) > 0:
            for job_role in result['data']['contact']['jobRoles']:
                response_obj['organizations'].append(job_role['organization'])


        return response_obj