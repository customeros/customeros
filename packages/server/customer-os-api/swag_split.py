import copy
import sys
from pathlib import Path
from typing import Any, Dict, Set

import yaml


def find_all_refs(obj: Any, refs: Set[str]) -> None:
    """Recursively find all $ref values in an object."""
    if isinstance(obj, dict):
        for key, value in obj.items():
            if key == '$ref' and isinstance(value, str) and '#/components/schemas/' in value:
                schema_name = value.split('/')[-1]
                refs.add(schema_name)
            elif key == 'allOf' and isinstance(value, list):
                for item in value:
                    find_all_refs(item, refs)
            else:
                find_all_refs(value, refs)
    elif isinstance(obj, list):
        for item in obj:
            find_all_refs(item, refs)

def find_missing_dependencies(orig_spec: Dict[str, Any], service_prefix: str) -> Dict[str, Any]:
    """Find and collect all missing dependencies from a service spec."""
    missing_schemas = set()
    all_refs = set()

    # First, find all references in paths
    for path_obj in orig_spec.get('paths', {}).values():
        for method in path_obj.values():
            if isinstance(method, dict):
                # Check parameters
                for param in method.get('parameters', []):
                    find_all_refs(param, all_refs)
                
                # Check requestBody
                if 'requestBody' in method:
                    find_all_refs(method['requestBody'], all_refs)
                
                # Check responses
                for response in method.get('responses', {}).values():
                    find_all_refs(response, all_refs)

    # Then recursively find nested references in schemas
    existing_schemas = orig_spec.get('components', {}).get('schemas', {})
    processed = set()
    while all_refs:
        current_refs = all_refs.copy()
        all_refs.clear()
        
        for ref in current_refs:
            if ref not in processed:
                processed.add(ref)
                if ref not in existing_schemas:
                    missing_schemas.add(ref)
                else:
                    find_all_refs(existing_schemas[ref], all_refs)

    return missing_schemas

def split_openapi_spec(input_file: str, output_dir: str) -> None:
    """
    Splits an OpenAPI specification file into multiple files based on path prefixes.
    Ensures all necessary schema components are included.
    """
    # Create output directory if it doesn't exist
    Path(output_dir).mkdir(parents=True, exist_ok=True)
    
    # Read the input file
    with open(input_file, 'r') as f:
        spec = yaml.safe_load(f)
    
    # Define the path prefixes to split by
    prefixes = ['billing', 'customerbase', 'domains', 'enrich', 'outreach', 'verify']
    
    # Store all schemas
    all_schemas = spec.get('components', {}).get('schemas', {})
    
    for prefix in prefixes:
        # Create a new spec with the common elements
        new_spec = copy.deepcopy(spec)
        new_spec['paths'] = {}
        
        # Add only the paths that match the current prefix
        for path, path_item in spec['paths'].items():
            if path.startswith(f'/{prefix}'):
                new_spec['paths'][path] = path_item
        
        # Find missing dependencies
        missing_schemas = find_missing_dependencies(new_spec, prefix)
        
        # Initialize components if not present
        if 'components' not in new_spec:
            new_spec['components'] = {}
        if 'schemas' not in new_spec['components']:
            new_spec['components']['schemas'] = {}

        # Add all required schemas
        for schema_name in all_schemas:
            if (schema_name.startswith(f'{prefix}.') or 
                schema_name.startswith('rest.') or 
                schema_name in missing_schemas):
                new_spec['components']['schemas'][schema_name] = all_schemas[schema_name]
        
        # Write the new spec to a file
        output_file = Path(output_dir) / f'openapi_{prefix}.yaml'
        with open(output_file, 'w') as f:
            yaml.dump(new_spec, f, sort_keys=False)
        
        # Report what was created
        included_schemas = len(new_spec['components'].get('schemas', {}))
        print(f"Created {output_file} with {included_schemas} schemas")
        
        # Validate and report any remaining missing references
        missing = find_missing_dependencies(new_spec, prefix)
        if missing:
            print(f"Warning: {output_file} still missing schemas: {', '.join(missing)}")

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python script.py input.yaml output_dir")
        sys.exit(1)
    
    try:
        split_openapi_spec(sys.argv[1], sys.argv[2])
        print("Successfully split OpenAPI specification")
    except Exception as e:
        print(f"Error: {str(e)}")
        sys.exit(1)
