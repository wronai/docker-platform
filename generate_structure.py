#!/usr/bin/env python3
import os
import json
from pathlib import Path
import re

def get_file_info(file_path, base_path):
    """Get file information in the required format."""
    rel_path = os.path.relpath(file_path, base_path)
    if os.path.isdir(file_path):
        return {
            'name': os.path.basename(file_path),
            'type': 'directory',
            'path': rel_path.replace('\\', '/'),
            'children': scan_directory(file_path, base_path)
        }
    else:
        try:
            file_size = os.path.getsize(file_path)
            content = ''
            # Only read text files for content
            if file_size < 1024 * 1024:  # Only read files smaller than 1MB
                try:
                    with open(file_path, 'r', encoding='utf-8') as f:
                        content = f.read()
                except (UnicodeDecodeError, PermissionError):
                    content = f"[Binary file: {os.path.basename(file_path)}]"
            
            return {
                'name': os.path.basename(file_path),
                'type': 'file',
                'path': rel_path.replace('\\', '/'),
                'size': file_size,
                'content': content
            }
        except (PermissionError, OSError) as e:
            print(f"Skipping {file_path}: {e}")
            return None

def should_skip_directory(dirname):
    """Check if directory should be skipped."""
    skip_dirs = {
        'node_modules',
        '.git',
        '.github',
        '__pycache__',
        'venv',
        'env',
        '.venv',
        'dist',
        'build',
        'coverage',
        '.vscode',
        '.idea'
    }
    return dirname in skip_dirs

def scan_directory(directory, base_path):
    """Recursively scan a directory and return its structure."""
    result = []
    try:
        entries = sorted(os.listdir(directory), key=lambda x: (not os.path.isdir(os.path.join(directory, x)), x.lower()))
        for entry in entries:
            if entry.startswith('.') or should_skip_directory(entry):
                continue
                
            full_path = os.path.join(directory, entry)
            file_info = get_file_info(full_path, base_path)
            if file_info:
                result.append(file_info)
    except PermissionError as e:
        print(f"Permission denied: {directory}")
    return result

def generate_json_structure(root_dir, output_file):
    """Generate JSON structure for the given directory."""
    root_dir = os.path.abspath(root_dir)
    structure = {
        'name': os.path.basename(root_dir),
        'type': 'directory',
        'path': '.',
        'children': scan_directory(root_dir, root_dir)
    }
    
    with open(output_file, 'w', encoding='utf-8') as f:
        json.dump(structure, f, indent=2, ensure_ascii=False)
    
    print(f"Generated structure saved to {output_file}")
    return structure

def update_html_template(html_file, json_file):
    """Update HTML template with file system data if template exists."""
    try:
        # Check if template file exists
        if not os.path.exists(html_file):
            print(f"Note: HTML template {html_file} not found, skipping update")
            return
            
        # Read the template file
        with open(html_file, 'r', encoding='utf-8') as f:
            content = f.read()
        
        # Read the JSON file
        with open(json_file, 'r', encoding='utf-8') as f:
            json_data = f.read()
        
        # Create a script tag with the file system data
        script_tag = f'<script>\n        // This is automatically generated - do not edit manually\n        const FILE_SYSTEM = {json_data};\n    </script>'
        
        # Replace the placeholder script tag
        content = re.sub(
            r'<script>\s*// This will be replaced by the actual file system data\s*const FILE_SYSTEM = \{.*?\};?\s*</script>',
            script_tag,
            content,
            flags=re.DOTALL
        )
        
        # Write the updated content back to the HTML file
        with open(html_file, 'w', encoding='utf-8') as f:
            f.write(content)
        
        print(f"Updated {html_file} to include the file system data")
    except Exception as e:
        print(f"Warning: Could not update HTML template: {e}")

def main():
    # Configuration
    root_dir = os.path.dirname(os.path.abspath(__file__))
    data_dir = os.path.join(root_dir, 'data')
    output_json = os.path.join(data_dir, 'file-system.json')
    
    # Create data directory if it doesn't exist
    os.makedirs(data_dir, exist_ok=True)
    
    # Generate the JSON structure
    print("Generating file system structure...")
    structure = generate_json_structure(root_dir, output_json)
    
    # Update the HTML file if it exists
    html_file = os.path.join(root_dir, 'index.html')
    update_html_template(html_file, output_json)
    
    print("\nFile structure generation complete!")
    
    # Count directories and files in the generated structure
    def count_items(node):
        if not node or 'children' not in node:
            return 0, 0
            
        dirs = 0
        files = 0
        
        for child in node.get('children', []):
            if child['type'] == 'directory':
                dirs += 1
                child_dirs, child_files = count_items(child)
                dirs += child_dirs
                files += child_files
            else:
                files += 1
                
        return dirs, files
    
    dirs, files = count_items(structure)
    print(f"Total directories: {dirs}")
    print(f"Total files: {files}")
    print(f"Output file: {output_json}")

if __name__ == "__main__":
    main()
