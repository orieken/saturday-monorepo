---
description: How to self-heal a test failure using the Saturday MCP
---
Step 1: Execute the tests to capture the failure.
Use the `run_tests` tool.
args:
  projectPath: /path/to/project
  filter: "failed test name" (optional)

// turbo
Step 2: Parse the failure output to find the exact file and error.
Use the `parse_test_failure` tool.
args:
  output: "{{output from Step 1}}"

Step 3: Analyze the impact of the failing file.
Use `analyze_impact` to see what else might be affected.
args:
  projectPath: /path/to/project
  targetFile: "{{file from Step 2}}"

Step 4: Read the content of the failing file.
Use `read_resource` (or `view_file`) to get the code context.
args:
  url: "file://{{projectPath}}/{{file from Step 2}}"

Step 5: Heal the code.
Use the `replace_file_content` tool (or `generate_page` if it's a new element needed) to fix the error identified in Step 2.
Common fixes:
- Updating selectors.
- Adding `await`.
- Fixing typos.

Step 6: Verify the fix.
Run `run_tests` again with the same filter to ensure it passes.
