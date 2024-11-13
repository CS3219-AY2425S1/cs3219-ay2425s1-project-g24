package utils

import (
	"cloud.google.com/go/firestore"
	"context"
	"execution-service/models"
	"fmt"
	"google.golang.org/api/iterator"
	"log"
	"strings"
)

// RepopulateTests Populate deletes all testcases and repopulates testcases
func RepopulateTests(ctx context.Context, client *firestore.Client,
	questionTitleToDocRefIdMap map[string]string) error {

	// delete all existing document in the collection
	iter := client.Collection("tests").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to iterate over tests: %v", err)
		}

		if _, err := doc.Ref.Delete(ctx); err != nil {
			return fmt.Errorf("failed to delete test: %v", err)
		}
	}
	defer iter.Stop()

	var tests = []models.Test{
		{
			QuestionTitle: "Reverse a String",
			VisibleTestCases: `
2
hello
olleh
Hannah
hannaH
`,
			HiddenTestCases: `
2
Hannah
hannaH
abcdefg
gfedcba
`,
			InputValidation: getPackagesAndFunction([]string{}, `
return len(inputOrOutput) > 0
`),
			OutputValidation: getPackagesAndFunction([]string{}, `
return len(inputOrOutput) > 0
`),
		},
		{
			QuestionTitle: "Linked List Cycle Detection",
			VisibleTestCases: `
2
[3,2,0,-4] -> pos = 1
true
[1]
false
`,
			HiddenTestCases: `
2
[1,2] -> pos = 0
true
[1]
false
`,
			InputValidation: getPackagesAndFunction([]string{"strconv", "strings"}, `
hasCycle := false

// Step 1: Split by " -> pos = "
parts := strings.Split(inputOrOutput, " -> pos = ")
if len(parts) == 2 {
	hasCycle = true
} else if len(parts) != 1 {
	return false
}

listPart := strings.TrimSpace(parts[0]) // Get the list part

//Validate the list format
if len(listPart) < 2 || listPart[0] != '[' || listPart[len(listPart)-1] != ']' {
	return false
}

listContent := listPart[1 : len(listPart)-1] // Remove brackets
if listContent == "" {                       // Check for empty list
	return false
}

// Split the list by commas and validate each element
elements := strings.Split(listContent, ",")
for _, elem := range elements {
	elem = strings.TrimSpace(elem)                // Trim whitespace
	if _, err := strconv.Atoi(elem); err != nil { // Check if it's an integer
		return false
	}
}

if !hasCycle {
	return true
}

posPart := strings.TrimSpace(parts[1]) // Get the position part

//Validate the position
posValue, err := strconv.Atoi(posPart)
if err != nil || posValue < 0 || posValue >= len(elements) {
	return false
}

return true
`),
			OutputValidation: getPackagesAndFunction([]string{}, `
return inputOrOutput == "true" || inputOrOutput == "false"
`),
		},
		{
			QuestionTitle: "Roman to Integer",
			VisibleTestCases: `
2
III
3
IV
4
`,
			HiddenTestCases: `
2
IV
4
MCMXCIV
1994
`,
			InputValidation: getPackagesAndFunction([]string{"strings"}, `
valid := "IVXLCDM"
for _, char := range inputOrOutput {
	if !strings.ContainsRune(valid, char) {
		return false
	}
}
return true
`),
			OutputValidation: getPackagesAndFunction([]string{"strconv"}, `
_, err := strconv.Atoi(inputOrOutput)
return err == nil
`),
		},
		{
			QuestionTitle: "Add Binary",
			VisibleTestCases: `
2
"11", "1"
"100"
"1010", "1011"
"10101"
`,
			HiddenTestCases: `
2
"1010", "1011"
"10101"
"111", "111"
"1110"
`,
			InputValidation: getPackagesAndFunction([]string{"regexp"}, `
binaryRegex := regexp.MustCompile("^\"([01]+)\",\\s*\"([01]+)\"$")
return binaryRegex.MatchString(inputOrOutput)
`),
			OutputValidation: getPackagesAndFunction([]string{"regexp"}, `
binaryRegex := regexp.MustCompile("^\"([01]+)\"$")
return binaryRegex.MatchString(inputOrOutput)
`),
		},
		{
			QuestionTitle: "Fibonacci Number",
			VisibleTestCases: `
2
0
0
10
55
`,
			HiddenTestCases: `
2
1
1
10
55
`,
			InputValidation: getPackagesAndFunction([]string{"strconv"}, `
num, err := strconv.Atoi(inputOrOutput)
return err == nil && num >= 0
`),
			OutputValidation: getPackagesAndFunction([]string{"strconv"}, `
num, err := strconv.Atoi(inputOrOutput)
return err == nil && num >= 0
`),
		},
		{
			QuestionTitle: "Implement Stack using Queues",
			VisibleTestCases: `
2
push(1), push(2), top()
2
push(1), empty()
false
`,
			HiddenTestCases: `
2
push(1), push(2), pop(), top()
1
push(1), empty()
false
`,
			InputValidation: getPackagesAndFunction([]string{"strconv", "strings"}, `
// Split the line by commas to handle multiple operations
operations := strings.Split(inputOrOutput, ",")
for _, op := range operations {
	op = strings.TrimSpace(op) // Trim whitespace
	// Check if the operation is valid
	if strings.HasPrefix(op, "push(") && strings.HasSuffix(op, ")") {
		// Check if it's a push operation with a valid integer
		numStr := op[5 : len(op)-1] // Extract the number string
		if _, err := strconv.Atoi(numStr); err != nil {
			return false
		}
	} else if op != "pop()" && op != "top()" && op != "empty()" {
		// If the operation is not one of the valid ones
		return false
	}
}
return true
`),
			OutputValidation: getPackagesAndFunction([]string{"strconv"}, `
if inputOrOutput == "true" || inputOrOutput == "false" || inputOrOutput == "null" {
	return true
}
_, err := strconv.Atoi(inputOrOutput)
return err == nil
`),
		}, {
			QuestionTitle: "Combine Two Tables",
			VisibleTestCases: `
2
Person: [(1, "Smith", "John"), (2, "Doe", "Jane")], Address: [(1, 1, "NYC", "NY"), (2, 3, "LA", "CA")]
[("John", "Smith", "NYC", "NY"), ("Jane", "Doe", null, null)]
Person: [(1, "White", "Mary")], Address: []
[("Mary", "White", null, null)]
`,
			HiddenTestCases: `
2
Person: [(1, "Black", "Jim")], Address: [(1, 1, "Miami", "FL")]
[("Jim", "Black", "Miami", "FL")]
Person: [(1, "White", "Mary")], Address: []
[("Mary", "White", null, null)]
`,
			InputValidation: getPackagesAndFunction([]string{}, `
return len(inputOrOutput) > 0
`),
			OutputValidation: getPackagesAndFunction([]string{}, `
return len(inputOrOutput) > 0
`),
		},
		{
			QuestionTitle: "Repeated DNA Sequences",
			VisibleTestCases: `
2
AAAAACCCCCAAAAACCCCCCAAAAAGGGTTT
["AAAAACCCCC", "CCCCCAAAAA"]
ACGTACGTACGT
[]
`,
			HiddenTestCases: `
2
AAAAAAAAAAAAA
["AAAAAAAAAA"]
ACGTACGTACGT
[]
`,
			InputValidation: getPackagesAndFunction([]string{}, `
input := inputOrOutput

// Check that input length is at least 10
if len(input) < 10 {
	return false
}

// Check that input contains only 'A', 'C', 'G', and 'T'
for _, ch := range input {
	if ch != 'A' && ch != 'C' && ch != 'G' && ch != 'T' {
		return false
	}
}
return true
`),
			OutputValidation: getPackagesAndFunction([]string{"strings"}, `
output := inputOrOutput

if output == "[]" {
	return true
}

// Check that the output is enclosed in square brackets
if len(output) < 2 || output[0] != '[' || output[len(output)-1] != ']' {
	return false
}

// Extract the content between square brackets
content := output[1 : len(output)-1]

// Split by commas without trimming spaces
sequences := strings.Split(content, ", ")
for _, seq := range sequences {
	// Check if each sequence is properly enclosed in double quotes and is exactly 10 characters
	if len(seq) != 12 || seq[0] != '"' || seq[11] != '"' {
		return false
	}

	// Check that the sequence only contains valid DNA characters between the quotes
	for i := 1; i < 11; i++ {
		ch := seq[i]
		if ch != 'A' && ch != 'C' && ch != 'G' && ch != 'T' {
			return false
		}
	}
}
return true
`),
		},
		{
			QuestionTitle: "Course Schedule",
			VisibleTestCases: `
2
2, [[1,0]]
true
2, [[1,0],[0,1]]
false
`,
			HiddenTestCases: `
2
2, [[1,0],[0,1]]
false
4, [[1,0],[2,0],[3,1],[3,2]]
true
`,
			InputValidation: getPackagesAndFunction([]string{}, `
return len(inputOrOutput) > 0
`),
			OutputValidation: getPackagesAndFunction([]string{}, `
return len(inputOrOutput) > 0
`),
		},
		{
			QuestionTitle: "LRU Cache Design",
			VisibleTestCases: `
2
put(1, 1), put(2, 2), get(1)
1
put(1, 1), put(2, 2), put(3, 3), get(2)
-1
`,
			HiddenTestCases: `
2
put(1, 1), put(2, 2), put(3, 3), get(2)
-1
put(1, 1), put(2, 2), put(1, 10), get(1)
10
`,
			InputValidation: getPackagesAndFunction([]string{}, `
return len(inputOrOutput) > 0
`),
			OutputValidation: getPackagesAndFunction([]string{}, `
return len(inputOrOutput) > 0
`),
		},
		{
			QuestionTitle: "Longest Common Subsequence",
			VisibleTestCases: `
2
"abcde", "ace"
3
"abc", "def"
0
`,
			HiddenTestCases: `
2
"abc", "abc"
3
"abc", "def"
0
`,
			InputValidation: getPackagesAndFunction([]string{}, `
return len(inputOrOutput) > 0
`),
			OutputValidation: getPackagesAndFunction([]string{}, `
return len(inputOrOutput) > 0
`),
		},
		{
			QuestionTitle: "Rotate Image",
			VisibleTestCases: `
2
[[1,2,3],[4,5,6],[7,8,9]]
[[7,4,1],[8,5,2],[9,6,3]]
[[1]]
[[1]]
`,
			HiddenTestCases: `
2
[[5,1,9,11],[2,4,8,10],[13,3,6,7],[15,14,12,16]]
[[15,13,2,5],[14,3,4,1],[12,6,8,9],[16,7,10,11]]
[[1]]
[[1]]
`,
			InputValidation: getPackagesAndFunction([]string{}, `
return len(inputOrOutput) > 0
`),
			OutputValidation: getPackagesAndFunction([]string{}, `
return len(inputOrOutput) > 0
`),
		},
		{
			QuestionTitle: "Airplane Seat Assignment Probability",
			VisibleTestCases: `
2
1
1.00000
3
0.50000
`,
			HiddenTestCases: `
2
2
0.50000
3
0.50000
`,
			InputValidation: getPackagesAndFunction([]string{}, `
return len(inputOrOutput) > 0
`),
			OutputValidation: getPackagesAndFunction([]string{}, `
return len(inputOrOutput) > 0
`),
		},
		{
			QuestionTitle: "Validate Binary Search Tree",
			VisibleTestCases: `
2
[2,1,3]
true
[5,1,4,null,null,3,6]
false
`,
			HiddenTestCases: `
2
[5,1,4,null,null,3,6]
false
[1]
true
`,
			InputValidation: getPackagesAndFunction([]string{}, `
return len(inputOrOutput) > 0
`),
			OutputValidation: getPackagesAndFunction([]string{}, `
return len(inputOrOutput) > 0
`),
		},
		{
			QuestionTitle: "Sliding Window Maximum",
			VisibleTestCases: `
2
[1,3,-1,-3,5,3,6,7], k=3
[3,3,5,5,6,7]
[9, 11], k=2
[11]
`,
			HiddenTestCases: `
2
[1, -1], k=1
[1, -1]
[9, 11], k=2
[11]
`,
			InputValidation: getPackagesAndFunction([]string{}, `
return len(inputOrOutput) > 0
`),
			OutputValidation: getPackagesAndFunction([]string{}, `
return len(inputOrOutput) > 0
`),
		},
		{
			QuestionTitle: "N-Queen Problem",
			VisibleTestCases: `
2
4
[[".Q..","...Q","Q...","..Q."],["..Q.","Q...","...Q",".Q.."]]
2
[]
`,
			HiddenTestCases: `
2
1
[["Q"]]
2
[]
`,
			InputValidation: getPackagesAndFunction([]string{}, `
return len(inputOrOutput) > 0
`),
			OutputValidation: getPackagesAndFunction([]string{}, `
return len(inputOrOutput) > 0
`),
		},
		{
			QuestionTitle: "Serialize and Deserialize a Binary Tree",
			VisibleTestCases: `
2
[1,2,3,null,null,4,5]
"1 2 null null 3 4 null null 5 null null"
[]
"null"
`,
			HiddenTestCases: `
2
[]
"null"
[1]
"1 null null"
`,
			InputValidation: getPackagesAndFunction([]string{}, `
return len(inputOrOutput) > 0
`),
			OutputValidation: getPackagesAndFunction([]string{}, `
return len(inputOrOutput) > 0
`),
		},
		{
			QuestionTitle: "Wildcard Matching",
			VisibleTestCases: `
2
"aa", "a"
false
"aa", "*"
true
`,
			HiddenTestCases: `
2
"aa", "*"
true
"cb", "?a"
false
`,
			InputValidation: getPackagesAndFunction([]string{}, `
return len(inputOrOutput) > 0
`),
			OutputValidation: getPackagesAndFunction([]string{}, `
return len(inputOrOutput) > 0
`),
		},
		{
			QuestionTitle: "Chalkboard XOR Game",
			VisibleTestCases: `
2
[1,1,2]
false
[1,2,3]
true
`,
			HiddenTestCases: `
2
[1,2,3]
true
[0,0,0]
true
`,
			InputValidation: getPackagesAndFunction([]string{}, `
return len(inputOrOutput) > 0
`),
			OutputValidation: getPackagesAndFunction([]string{}, `
return len(inputOrOutput) > 0
`),
		},
		{
			QuestionTitle: "Trips and Users",
			VisibleTestCases: `
2
Trips: [(1, 1, 10, 'NYC', 'completed', '2013-10-01'), (2, 2, 11, 'NYC', 'cancelled_by_driver', '2013-10-01')],Users: [(10, 'No', 'client'), (11, 'No', 'driver')]
0.50
Trips: [(1, 1, 10, 'NYC', 'completed', '2013-10-03'), (2, 2, 11, 'NYC', 'cancelled_by_client', '2013-10-03')],Users: [(10, 'No', 'client'), (11, 'Yes', 'driver')]
0.00
`,
			HiddenTestCases: `
2
Trips: [(1, 1, 10, 'NYC', 'completed', '2013-10-02')],Users: [(10, 'No', 'client'), (11, 'No', 'driver')]
0.00
Trips: [(1, 1, 10, 'NYC', 'completed', '2013-10-03'), (2, 2, 11, 'NYC', 'cancelled_by_client', '2013-10-03')],Users: [(10, 'No', 'client'), (11, 'Yes', 'driver')]
0.00
`,
			InputValidation: getPackagesAndFunction([]string{}, `
return len(inputOrOutput) > 0
`),
			OutputValidation: getPackagesAndFunction([]string{}, `
return len(inputOrOutput) > 0
`),
		},
	}

	for _, test := range tests {
		if _, exists := questionTitleToDocRefIdMap[test.QuestionTitle]; !exists {
			return fmt.Errorf("testcase's question title %s not found in questionTitleToDocRefIdMap",
				test.QuestionTitle)
		}
	}

	for _, test := range tests {
		_, err := ValidateTestCaseFormat(test.VisibleTestCases, test.InputValidation, test.OutputValidation)
		if err != nil {
			return fmt.Errorf("failed to validate visible test case format: %v", err)
		}
		_, err = ValidateTestCaseFormat(test.HiddenTestCases, test.InputValidation, test.OutputValidation)
		if err != nil {
			return fmt.Errorf("failed to validate hidden test case format: %v", err)
		}
	}

	for _, test := range tests {
		_, _, err := client.Collection("tests").Add(ctx, map[string]interface{}{
			"questionDocRefId": questionTitleToDocRefIdMap[test.QuestionTitle],
			"visibleTestCases": test.VisibleTestCases,
			"hiddenTestCases":  test.HiddenTestCases,
			"createdAt":        firestore.ServerTimestamp,
			"updatedAt":        firestore.ServerTimestamp,
		})
		if err != nil {
			return fmt.Errorf("failed to add test: %v", err)
		}
	}

	log.Println("Cleaned tests collection and repopulated tests.")
	return nil
}

func getPackagesAndFunction(imports []string, functionBody string) string {
	// Use a string builder for efficient string concatenation
	var importCode strings.Builder

	if len(imports) > 0 {
		// Start the import block
		importCode.WriteString("import (\n")

		// Iterate over the imports and add them to the builder
		for _, imp := range imports {
			importCode.WriteString(fmt.Sprintf("\t%q\n", imp))
		}

		// Close the import block
		importCode.WriteString(")")
	}

	// add tab before every line in functionBody including first line
	functionBody = strings.ReplaceAll(functionBody, "\n", "\n\t")

	return fmt.Sprintf(`
%s

func validateInputOrOutput(inputOrOutput string) bool {
%s
}
`, importCode.String(), functionBody)
}

func GetDefaultValidation() string {
	return getPackagesAndFunction([]string{}, `
return len(inputOrOutput) > 0
`)
}
