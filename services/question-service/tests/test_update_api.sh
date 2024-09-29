#!/bin/bash

# Output file
output_file="test_update_api_results.txt"

> "$output_file"

docRefID="VlnxEsVxduJPhkCU9zQq"

# List of curl commands (formatted properly)
requests=(
# Successful case
"curl -X PUT http://localhost:8080/questions/$docRefID \
-H 'Content-Type: application/json' \
-d '{\"title\": \"Sample Question\", \"description\": \"This is a sample description.\", \"complexity\": \"medium\", \"categories\": [\"Data Structures\", \"Algorithms\"]}'"

# Missing question
"curl -X PUT http://localhost:8080/questions/$docRefID \
-H 'Content-Type: application/json' \
-d '{\"description\": \"This is a sample description.\", \"complexity\": \"medium\", \"categories\": [\"Data Structures\", \"Algorithms\"]}'"

# Missing description
"curl -X PUT http://localhost:8080/questions/$docRefID \
-H 'Content-Type: application/json' \
-d '{\"title\": \"Sample Question 1\", \"complexity\": \"medium\", \"categories\": [\"Data Structures\", \"Algorithms\"]}'"

# Missing complexity
"curl -X PUT http://localhost:8080/questions/$docRefID \
-H 'Content-Type: application/json' \
-d '{\"title\": \"Sample Question 2\", \"description\": \"This is a sample description.\", \"categories\": [\"Data Structures\", \"Algorithms\"]}'"

# Missing categories
"curl -X PUT http://localhost:8080/questions/$docRefID \
-H 'Content-Type: application/json' \
-d '{\"title\": \"Sample Question 3\", \"description\": \"This is a sample description.\", \"complexity\": \"medium\"}'"

# Invalid complexity
"curl -X PUT http://localhost:8080/questions/$docRefID \
-H 'Content-Type: application/json' \
-d '{\"title\": \"Sample Question 4\", \"description\": \"This is a sample description.\", \"complexity\": \"extreme\", \"categories\": [\"Data Structures\", \"Algorithms\"]}'"

# Duplicate question title
"curl -X PUT http://localhost:8080/questions/$docRefID \
-H 'Content-Type: application/json' \
-d '{\"title\": \"Sample Question\", \"description\": \"This is a sample description.\", \"complexity\": \"medium\", \"categories\": [\"Data Structures\", \"Algorithms\"]}'"

# Incorrect JSON
"curl -X PUT http://localhost:8080/questions/$docRefID \
-H 'Content-Type: application/json' \
-d '{\"title\": \"Sample Question\", \"description\": \"This is a sample description.\", \"complexity\": \"medium\", \"categories\": [\"Data Structures\", \"Algorithms\"}'"
)

# Execute each request and save the output to the file
for request in "${requests[@]}"; do
    echo "Executing: $request" >> "$output_file"
    echo "Response:" >> "$output_file"
    eval $request >> "$output_file"
    echo -e "\n" >> "$output_file"
done

echo "All requests executed. Results saved in $output_file."
