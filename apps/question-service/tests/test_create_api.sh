#!/bin/bash

# Output file
output_file="test_create_api_results.txt"

cd .. || { echo "Failed to navigate to parent directory."; exit 1; }
go run populate.go || { echo "Failed to populate database."; exit 1; }
cd - || { echo "Failed to navigate back."; exit 1; }

> "$output_file"

# List of curl commands (formatted properly)
requests=(
# Successful case
"curl -X POST http://localhost:8080/questions \
-H 'Content-Type: application/json' \
-d '{\"title\": \"Sample Question\", \"description\": \"This is a sample description.\", \"complexity\": \"medium\", \"categories\": [\"Data Structures\", \"Algorithms\"]}'"

# Missing question
"curl -X POST http://localhost:8080/questions \
-H 'Content-Type: application/json' \
-d '{\"description\": \"This is a sample description.\", \"complexity\": \"medium\", \"categories\": [\"Data Structures\", \"Algorithms\"]}'"

# Missing description
"curl -X POST http://localhost:8080/questions \
-H 'Content-Type: application/json' \
-d '{\"title\": \"Sample Question 1\", \"complexity\": \"medium\", \"categories\": [\"Data Structures\", \"Algorithms\"]}'"

# Missing complexity
"curl -X POST http://localhost:8080/questions \
-H 'Content-Type: application/json' \
-d '{\"title\": \"Sample Question 2\", \"description\": \"This is a sample description.\", \"categories\": [\"Data Structures\", \"Algorithms\"]}'"

# Missing categories
"curl -X POST http://localhost:8080/questions \
-H 'Content-Type: application/json' \
-d '{\"title\": \"Sample Question 3\", \"description\": \"This is a sample description.\", \"complexity\": \"medium\"}'"

# Invalid complexity
"curl -X POST http://localhost:8080/questions \
-H 'Content-Type: application/json' \
-d '{\"title\": \"Sample Question 4\", \"description\": \"This is a sample description.\", \"complexity\": \"extreme\", \"categories\": [\"Data Structures\", \"Algorithms\"]}'"

# Duplicate question title
"curl -X POST http://localhost:8080/questions \
-H 'Content-Type: application/json' \
-d '{\"title\": \"Sample Question\", \"description\": \"This is a sample description.\", \"complexity\": \"medium\", \"categories\": [\"Data Structures\", \"Algorithms\"]}'"

# Incorrect JSON
"curl -X POST http://localhost:8080/questions \
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
