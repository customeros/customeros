import json

import requests
import yaml


def convert_swagger_to_openapi(input_file: str, output_file: str) -> None:
    """
    Convert Swagger/OpenAPI 2.0 specification to OpenAPI 3.0
    using Swagger Editor's conversion API.
    """
    # Read the input swagger file
    with open(input_file, 'r') as f:
        swagger_dict = yaml.safe_load(f)
    
    # Convert to JSON for the API request
    swagger_json = json.dumps(swagger_dict)
    
    # Call Swagger Editor's conversion API
    response = requests.post(
        'https://converter.swagger.io/api/convert',
        data=swagger_json,
        headers={'Content-Type': 'application/json'}
    )
    
    if response.status_code != 200:
        raise Exception(f"Conversion failed: {response.text}")
    
    # Parse the response and write to file
    openapi_dict = response.json()
    with open(output_file, 'w') as f:
        yaml.dump(openapi_dict, f, sort_keys=False)

if __name__ == "__main__":
    import sys
    
    if len(sys.argv) != 3:
        print("Usage: python converter.py input.yaml output.yaml")
        sys.exit(1)
        
    try:
        convert_swagger_to_openapi(sys.argv[1], sys.argv[2])
        print(f"Successfully converted {sys.argv[1]} to OpenAPI 3.0 format")
    except Exception as e:
        print(f"Error: {str(e)}")
        sys.exit(1)
