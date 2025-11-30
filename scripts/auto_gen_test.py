import os
import sys
import google.generativeai as genai
from git import Repo

def get_changed_go_files():
    """
    Detects changed .go files in the current PR/commit range.
    In GitHub Actions, we usually compare HEAD against the base branch.
    """
    repo = Repo(".")
    changed_files = []
    try:
        files = repo.git.diff('--name-only', 'origin/main...HEAD').splitlines()
        for f in files:
            if f.endswith('.go') and not f.endswith('_test.go'):
                changed_files.append(f)
    except Exception as e:
        print(f"Error getting git diff: {e}")
        return []

    return changed_files

def generate_test(file_path):
    print(f"Generating test for: {file_path}")
    
    with open(file_path, 'r') as f:
        content = f.read()

    model = genai.GenerativeModel('gemini-2.0-flash')
    
    prompt = f"""You are an expert Go developer. Generate comprehensive unit tests for the following Go code using the standard 'testing' package. 
Output ONLY the code for the test file, including package declaration and imports. 
Do not include markdown code blocks or any other text.

Code:
{content}"""

    response = model.generate_content(prompt)
    test_content = response.text
    
    test_content = test_content.replace("```go", "").replace("```", "")
    dir_name = os.path.dirname(file_path)
    base_name = os.path.basename(file_path)
    test_file_name = base_name.replace(".go", "_test.go")
    test_file_path = os.path.join(dir_name, test_file_name)
    
    with open(test_file_path, 'w') as f:
        f.write(test_content)
    
    print(f"Created: {test_file_path}")

def main():
    api_key = os.environ.get("GEMINI_API_KEY")
    if not api_key:
        print("GEMINI_API_KEY not set")
        sys.exit(1)
        
    genai.configure(api_key=api_key)
    
    changed_files = get_changed_go_files()
    if not changed_files:
        print("No .go files changed.")
        return

    for f in changed_files:
        if os.path.exists(f):
            generate_test(f)
        else:
            print(f"File not found (maybe deleted): {f}")

if __name__ == "__main__":
    main()
