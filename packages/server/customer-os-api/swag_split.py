import copy
from pathlib import Path
from typing import Any, Dict

import yaml


def split_openapi_spec(input_file: str, output_dir: str) -> None:
    """
    Splits an OpenAPI specification file into multiple files based on path prefixes.
    
    Args:
        input_file (str): Path to the input OpenAPI specification file
        output_dir (str): Directory where the split specifications will be saved
    """
    # Create output directory if it doesn't exist
    Path(output_dir).mkdir(parents=True, exist_ok=True)
    
    # Read the input file
    with open(input_file, 'r') as f:
        spec = yaml.safe_load(f)
    
    # Define the path prefixes to split by
    prefixes = ['billing', 'customerbase', 'domains', 'enrich', 'outreach', 'verify']
    
    # Track all used schema references for each split spec
    used_schemas: Dict[str, set] = {prefix: set() for prefix in prefixes}
    
    # Create a separate spec for each prefix
    for prefix in prefixes:
        # Create a new spec with the common elements
        new_spec = copy.deepcopy(spec)
        new_spec['paths'] = {}
        
        # Add only the paths that match the current prefix
        for path, path_item in spec['paths'].items():
            if path.startswith(f'/{prefix}'):
                new_spec['paths'][path] = path_item
                
                # Track schema references in this path
                for method in path_item.values():
                    # Check request body schemas
                    if 'requestBody' in method:
                        content = method['requestBody'].get('content', {})
                        for content_type in content.values():
                            if 'schema' in content_type:
                                if '$ref' in content_type['schema']:
                                    ref = content_type['schema']['$ref'].split('/')[-1]
                                    used_schemas[prefix].add(ref)
                    
                    # Check response schemas
                    if 'responses' in method:
                        for response in method['responses'].values():
                            if 'content' in response:
                                for content_type in response['content'].values():
                                    if 'schema' in content_type:
                                        if '$ref' in content_type['schema']:
                                            ref = content_type['schema']['$ref'].split('/')[-1]
                                            used_schemas[prefix].add(ref)
        
        # Only include referenced schemas
        if 'components' in new_spec and 'schemas' in new_spec['components']:
            referenced_schemas = {}
            for schema_name, schema in new_spec['components']['schemas'].items():
                # Include the schema if it's used or if it's prefixed with the current service
                if (schema_name in used_schemas[prefix] or 
                    schema_name.startswith(f'{prefix}.') or 
                    schema_name.startswith('rest.')):
                    referenced_schemas[schema_name] = schema
                    
                    # Also include any nested references
                    def find_nested_refs(obj):
                        if isinstance(obj, dict):
                            if '$ref' in obj:
                                ref = obj['$ref'].split('/')[-1]
                                referenced_schemas[ref] = spec['components']['schemas'][ref]
                            for value in obj.values():
                                find_nested_refs(value)
                        elif isinstance(obj, list):
                            for item in obj:
                                find_nested_refs(item)
                    
                    find_nested_refs(schema)
            
            new_spec['components']['schemas'] = referenced_schemas
        
        # Write the new spec to a file
        output_file = Path(output_dir) / f'openapi_{prefix}.yaml'
        with open(output_file, 'w') as f:
            yaml.dump(new_spec, f, sort_keys=False)
        
        print(f"Created {output_file}")

if __name__ == "__main__":
    import sys
    
    if len(sys.argv) != 3:
        print("Usage: python script.py input.yaml output_dir")
        sys.exit(1)
    
    input_file = sys.argv[1]
    output_dir = sys.argv[2]
    
    try:
        split_openapi_spec(input_file, output_dir)
        print("Successfully split OpenAPI specification")
    except Exception as e:
        print(f"Error: {str(e)}")
        sys.exit(1)
