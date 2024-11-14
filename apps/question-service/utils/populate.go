package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"question-service/models"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// PopulateSampleQuestionsInTransaction deletes all existing questions and then adds new ones in a single transaction
func populateSampleQuestionsInTransaction(ctx context.Context, client *firestore.Client) error {

	// Sample questions to be added after deletion
	var sampleQuestions = []models.Question{
		{
			Title:      "Reverse a String",
			Categories: []string{"Strings", "Algorithms"},
			Complexity: models.Easy,
			CreatedAt:  time.Now().AddDate(-1, -5, -12),
			Description: `Write a function that reverses a string. The input string is given as an array of characters s.

You must do this by modifying the input array in-place with O(1) extra memory.

Example 1:
Input: s = ["h","e","l","l","o"]
Output: ["o","l","l","e","h"]

Example 2:
Input: s = ["H","a","n","n","a"," h"]
Output: ["h","a","n","n","a"," H"]

Constraints:
1 <= s.length <= 105 s[i] is a printable ascii character.`,
		},
		{
			Title:       "Linked List Cycle Detection",
			Categories:  []string{"Data Structures", "Algorithms"},
			Complexity:  models.Easy,
			CreatedAt:   time.Now().AddDate(-2, -3, -25),
			Description: `Implement a function to detect if a linked list contains a cycle.`,
		},
		{
			Title:       "Roman to Integer",
			Categories:  []string{"Algorithms"},
			Complexity:  models.Easy,
			CreatedAt:   time.Now().AddDate(-3, -2, -7),
			Description: `Given a roman numeral, convert it to an integer.`,
		},
		{
			Title:       "Add Binary",
			Categories:  []string{"Bit Manipulation", "Algorithms"},
			Complexity:  models.Easy,
			CreatedAt:   time.Now().AddDate(-4, -6, -18),
			Description: `Given two binary strings a and b, return their sum as a binary string.`,
		},
		{
			Title:      "Fibonacci Number",
			Categories: []string{"Recursion", "Algorithms"},
			Complexity: models.Easy,
			CreatedAt:  time.Now().AddDate(-1, -11, -30),
			Description: `The Fibonacci numbers, commonly denoted F(n) form a sequence, called the Fibonacci sequence, such that each number is the sum of the two preceding ones, starting from 0 and 1. That is,

F(0) = 0, F(1) = 1 F(n) = F(n - 1) + F(n - 2), for n > 1.
Given n, calculate F(n).`,
		},
		{
			Title:       "Implement Stack using Queues",
			Categories:  []string{"Data Structures"},
			Complexity:  models.Easy,
			CreatedAt:   time.Now().AddDate(-2, -4, -15),
			Description: `Implement a last-in- first-out (LIFO) stack using only two queues. The implemented stack should support all the functions of a normal stack (push, top, pop, and empty).`,
		},
		{
			Title:      "Combine Two Tables",
			Categories: []string{"Databases"},
			Complexity: models.Easy,
			CreatedAt:  time.Now().AddDate(-3, -5, -21),
			Description: `Given table Person with the following columns:
1. personId (int)
2. lastName (varchar)
3. firstName (varchar)
personId is the primary key.

And table Address with the following columns:
1. addressId (int)
2. personId (int)
3. city (varchar)
4. state (varchar)
addressId is the primary key.

Write a solution to report the first name, last name, city, and state of each person in the Person table.
If the address of a personId is not present in the Address table,
report null instead. Return the result table in any order.`,
		},
		{
			Title:      "Repeated DNA Sequences",
			Categories: []string{"Algorithms", "Bit Manipulation"},
			Complexity: models.Medium,
			CreatedAt:  time.Now().AddDate(-1, -7, -28),
			Description: `The DNA sequence is composed of a series of nucleotides abbreviated as 'A', 'C', 'G', and 'T'.

For example, "ACGAATTCCG" is a DNA sequence. When studying DNA, it is useful to identify repeated sequences within the DNA.

Given a string s that represents a DNA sequence, return all the 10-letter-long sequences (substrings) that occur more than once in a DNA molecule. You may return the answer in any order.`,
		},
		{
			Title:      "Course Schedule",
			Categories: []string{"Data Structures", "Algorithms"},
			Complexity: models.Medium,
			CreatedAt:  time.Now().AddDate(-2, -6, -9),
			Description: `There are a total of numCourses courses you have to take, labeled from 0 to numCourses - 1. You are given an array prerequisites where prerequisites[i] = [ai, bi] indicates that you must take course bi first if you want to take course ai.

For example, the pair [0, 1], indicates that to take course 0 you have to first take course 1.
Return true if you can finish all courses. Otherwise, return false.`,
		},
		{
			Title:       "LRU Cache Design",
			Categories:  []string{"Data Structures"},
			Complexity:  models.Medium,
			CreatedAt:   time.Now().AddDate(-1, -8, -17),
			Description: `Design and implement an LRU (Least Recently Used) cache.`,
		},
		{
			Title:      "Longest Common Subsequence",
			Categories: []string{"Strings", "Algorithms"},
			Complexity: models.Medium,
			CreatedAt:  time.Now().AddDate(-3, -9, -3),
			Description: `Given two strings text1 and text2, return the length of their longest common subsequence. If there is no common subsequence, return 0.

A subsequence of a string is a new string generated from the original string with some characters (can be none) deleted without changing the relative order of the remaining characters.

For example, "ace" is a subsequence of "abcde".
A common subsequence of two strings is a subsequence that is common to both strings.`,
		},
		{
			Title:       "Rotate Image",
			Categories:  []string{"Arrays", "Algorithms"},
			Complexity:  models.Medium,
			CreatedAt:   time.Now().AddDate(-4, -11, -13),
			Description: `You are given an n x n 2D matrix representing an image, rotate the image by 90 degrees (clockwise).`,
		},
		{
			Title:      "Airplane Seat Assignment Probability",
			Categories: []string{"Brainteaser"},
			Complexity: models.Medium,
			CreatedAt:  time.Now().AddDate(-1, -3, -29),
			Description: `n passengers board an airplane with exactly n seats. The first passenger has lost the ticket and picks a seat randomly. But after that, the rest of the passengers will:

Take their own seat if it is still available, and Pick other seats randomly when they find their seat occupied

Return the probability that the nth person gets his own seat.`,
		},
		{
			Title:       "Validate Binary Search Tree",
			Categories:  []string{"Data Structures", "Algorithms"},
			Complexity:  models.Medium,
			CreatedAt:   time.Now().AddDate(-2, -12, -6),
			Description: `Given the root of a binary tree, determine if it is a valid binary search tree (BST).`,
		},
		{
			Title:      "Sliding Window Maximum",
			Categories: []string{"Arrays", "Algorithms"},
			Complexity: models.Hard,
			CreatedAt:  time.Now().AddDate(-3, -1, -22),
			Description: `You are given an array of integers nums, there is a sliding window of size k which is moving from the very left of the array to the very right. You can only see the k numbers in the window. Each time the sliding window moves right by one position.

Return the max sliding window.`,
		},
		{
			Title:      "N-Queen Problem",
			Categories: []string{"Algorithms"},
			Complexity: models.Hard,
			CreatedAt:  time.Now().AddDate(-2, -7, -10),
			Description: `The n-queens puzzle is the problem of placing n queens on an n x n chessboard such that no two queens attack each other.

Given an integer n, return all distinct solutions to the n- queens puzzle. You may return the answer in any order.

Each solution contains a distinct board configuration of the n-queens' placement, where 'Q' and '.' both indicate a queen and an empty space, respectively.`,
		},
		{
			Title:      "Serialize and Deserialize a Binary Tree",
			Categories: []string{"Data Structures", "Algorithms"},
			Complexity: models.Hard,
			CreatedAt:  time.Now().AddDate(-1, -4, -19),
			Description: `Serialization is the process of converting a data structure or object into a sequence of bits so that it can be stored in a file or memory buffer, or transmitted across a network connection link to be reconstructed later in the same or another computer environment.

Design an algorithm to serialize and deserialize a binary tree. There is no restriction on how your serialization/deserializ ation algorithm should work. You just need to ensure that a binary tree can be serialized to a string and this string can be deserialized to the original tree structure.`,
		},
		{
			Title:      "Wildcard Matching",
			Categories: []string{"Strings", "Algorithms"},
			Complexity: models.Hard,
			CreatedAt:  time.Now().AddDate(-4, -2, -27),
			Description: `Given an input string (s) and a pattern (p), implement wildcard pattern matching with support for '?' and '*' where:

'?' Matches any single character.
'*' Matches any sequence of characters (including the empty sequence). The matching should cover the entire input string (not partial).`,
		},
		{
			Title:      "Chalkboard XOR Game",
			Categories: []string{"Brainteaser"},
			Complexity: models.Hard,
			CreatedAt:  time.Now().AddDate(-4, -2, -27),
			Description: `You are given an array of integers nums represents the numbers written on a chalkboard.

Alice and Bob take turns erasing exactly one number from the chalkboard, with Alice starting first. If erasing a number causes the bitwise XOR of all the elements of the chalkboard to become 0, then that player loses. The bitwise XOR of one element is that element itself, and the bitwise XOR of no elements is 0.

Also, if any player starts their turn with the bitwise XOR of all the elements of the chalkboard equal to 0, then that player wins.

Return true if and only if Alice wins the game, assuming both players play optimally.`,
		},
		{
			Title:      "Trips and Users",
			Categories: []string{"Databases"},
			Complexity: models.Hard,
			CreatedAt:  time.Now().AddDate(-1, -6, -5),
			Description: `Given table Trips:
1. id (int)
2. client_id (int)
3. driver_id (int)
4. city_id (int)
5. status (enum)
6. request_at(date)
id is the primary key.
The table holds all taxi trips. Each trip has a unique id, while client_id and driver_id are foreign keys to the users_id at the Users table.
Status is an ENUM (category) type of ('completed', 'cancelled_by_driver', 'cancelled_by_client').

And table Users:
1. users_id (int)
2. banned (enum)
3. role (enum)
users_id is the primary key (column with unique values) for this table.
The table holds all users. Each user has a unique users_id, and role is an ENUM type of ('client', 'driver', 'partner').
banned is an ENUM (category) type of ('Yes', 'No').
The cancellation rate is computed by dividing the number of canceled (by client or driver) requests with unbanned users by the total number of requests with unbanned users on that day.

Write a solution to find the cancellation rate of requests with unbanned users (both client and driver must not be banned) each day between "2013- 10-01" and "2013-10- 03". Round Cancellation Rate to two decimal points.

Return the result table in any order.`,
		},
	}

	// Start a Firestore transaction
	return client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		counterDocRef := client.Collection("counters").Doc("questions")

		collection := client.Collection("questions")

		// Step 1: Delete all existing documents in the collection
		iter := collection.Documents(ctx)
		for {
			doc, err := iter.Next()
			if err != nil {
				if err == iterator.Done {
					break
				}
				return err // Abort transaction on error
			}
			// Delete each document in the transaction
			if err := tx.Delete(doc.Ref); err != nil {
				return err // Abort transaction if delete fails
			}
		}
		defer iter.Stop()

		// Step 2: Add new questions
		currentCounter := 0
		for _, question := range sampleQuestions {
			// Use the currentCounter as the ID for the new document
			currentCounter = currentCounter + 1
			_, _, err := client.Collection("questions").Add(ctx, map[string]interface{}{
				"id":          currentCounter,
				"title":       question.Title,
				"description": question.Description,
				"complexity":  question.Complexity,
				"categories":  question.Categories,
				"createdAt":   question.CreatedAt,
			})
			if err != nil {
				return err
			}
		}

		// Update the counter in Firestore
		if err := tx.Set(counterDocRef, map[string]interface{}{
			"count": 20,
		}); err != nil {
			return err
		}

		return nil
	})
}

func repopulateTests(dbClient *firestore.Client) error {
	ctx := context.Background()

	executionServiceUrl := os.Getenv("EXECUTION_SERVICE_URL")
	url := executionServiceUrl + "tests/populate"

	// get title and docRefId of all questions
	var questions []map[string]string
	iter := dbClient.Collection("questions").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to fetch question: %v", err)
		}

		questions = append(questions, map[string]string{
			"title":    doc.Data()["title"].(string),
			"docRefId": doc.Ref.ID,
		})
	}
	defer iter.Stop()

	jsonData, err := json.Marshal(questions)
	if err != nil {
		return fmt.Errorf("failed to marshal question titles: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request to populate tests: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to populate tests: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to populate tests: %v", resp.Status)
	}

	return nil
}

func Populate(client *firestore.Client, populateTests bool) {
	ctx := context.Background()

	// Run the transaction to delete all questions and add new ones
	err := populateSampleQuestionsInTransaction(ctx, client)
	if err != nil {
		log.Fatalf("Failed to populate sample questions in transaction: %v", err)
	}

	if !populateTests {
		log.Println("Counter reset, all questions deleted and sample questions added successfully in a transaction.")
		return
	}

	// Populate testcases in the execution-service
	err = repopulateTests(client)
	if err != nil {
		log.Fatalf("Failed to populate testcases: %v", err)
	}

	log.Println("Counter reset, " +
		"all questions deleted and sample questions added successfully in a transaction. Testcases repopulated in execution-service.")
}
