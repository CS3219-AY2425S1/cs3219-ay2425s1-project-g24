package processes

import (
	"context"
	"log"
	"matching-service/databases"
	"matching-service/models"
	"matching-service/utils"
	"strings"

	"github.com/redis/go-redis/v9"
)

// Find the first matching user from the front on the queue, based on topics, then difficulty.
func findMatchingUser(tx *redis.Tx, currentUsername string, ctx context.Context) (*models.MatchFound, error) {
	currentUser, err := databases.GetUserDetails(tx, currentUsername, ctx)
	if err != nil {
		return nil, err
	}

	// If currentUser has no topics, it is treated as user accepts any topic.
	shouldDoTopicMatching := len(currentUser.Topics) != 0

	// If currentUser has no difficulty, it is treated as user accepts any difficulty
	shouldDoDifficultyMatching := len(currentUser.Difficulties) != 0

	queuedUsers, err := databases.GetAllQueuedUsers(tx, ctx)
	if err != nil {
		return nil, err
	}

	// For each step, filter the number of potential matches
	potentialMatches := append([]string(nil), queuedUsers...)

	if shouldDoTopicMatching {
		potentialMatches, err = doTopicMatching(tx, ctx, &currentUser, potentialMatches)
		if err != nil {
			return nil, err
		}
	}

	if shouldDoDifficultyMatching {
		potentialMatches, err = doDifficultyMatching(tx, ctx, &currentUser, potentialMatches)
		if err != nil {
			return nil, err
		}
	}

	// Pick the first user from potential matches, which is earlier in the queue.
	var foundUser *string
	for _, otherUser := range potentialMatches {
		if otherUser == currentUsername {
			continue
		}
		foundUser = &otherUser
	}

	if foundUser == nil {
		return nil, models.NoMatchFound
	}
	return foundMatch(tx, ctx, &currentUser, foundUser)
}

func foundMatch(tx *redis.Tx, ctx context.Context, currentUser *models.MatchRequest, matchedUsername *string) (*models.MatchFound, error) {
	// Generate a random match ID
	matchId, err := utils.GenerateMatchID()
	if err != nil {
		log.Println("Unable to randomly generate matchID")
	}

	matchedUser, err := databases.GetUserDetails(tx, *matchedUsername, ctx)
	if err != nil {
		return nil, err
	}

	var matchedTopics []string
	for _, topic := range currentUser.Topics {
		for _, otherTopic := range matchedUser.Topics {
			if topic == otherTopic {
				matchedTopics = append(matchedTopics, topic)
			}
		}
	}
	var matchedDifficulties []string
	for _, topic := range currentUser.Difficulties {
		for _, otherTopic := range matchedUser.Difficulties {
			if topic == otherTopic {
				matchedDifficulties = append(matchedDifficulties, topic)
			}
		}
	}

	matchFound := models.MatchFound{
		Type:                "match_found",
		MatchID:             matchId,
		User:                currentUser.Username,
		MatchedUser:         *matchedUsername,
		MatchedTopics:       matchedTopics,
		MatchedDifficulties: matchedDifficulties,
	}

	return &matchFound, nil
}

func findSameTopicsUsers(tx *redis.Tx, ctx context.Context, currentUser *models.MatchRequest) (map[string]struct{}, error) {
	sameTopicUsersSet := make(map[string]struct{})
	for _, topic := range currentUser.Topics {
		topicUsers, err := tx.SMembers(ctx, strings.ToLower(topic)).Result()
		if err != nil {
			return nil, err
		}

		for _, potentialMatchUser := range topicUsers {
			sameTopicUsersSet[potentialMatchUser] = struct{}{}
		}
	}
	return sameTopicUsersSet, nil
}

func findSameDifficultiesUsers(tx *redis.Tx, ctx context.Context, currentUser *models.MatchRequest) (map[string]struct{}, error) {
	sameDifficultyUsersSet := make(map[string]struct{})
	queuedUsers, err := databases.GetAllQueuedUsers(tx, ctx)
	if err != nil {
		return nil, err
	}

	for _, potentialMatch := range queuedUsers {
		if potentialMatch == currentUser.Username {
			continue
		}

		potentialMatchUser, err := databases.GetUserDetails(tx, potentialMatch, ctx)
		if err != nil {
			return nil, err
		}

		for _, a := range potentialMatchUser.Difficulties {
			for _, b := range currentUser.Difficulties {
				if a == b {
					sameDifficultyUsersSet[potentialMatch] = struct{}{}
				}
			}
		}
	}
	return sameDifficultyUsersSet, nil
}

func doTopicMatching(tx *redis.Tx, ctx context.Context, currentUser *models.MatchRequest, potentialMatches []string) ([]string, error) {
	sameTopicUsers, err := findSameTopicsUsers(tx, ctx, currentUser)
	if err != nil {
		return nil, err
	}

	// Iterate through the queue to find a match, so a user in the queue the longest is more likely to be matched.
	var foundUsers []string
	for _, otherUsername := range potentialMatches {
		if otherUsername == currentUser.Username {
			continue
		}

		// Include users without any topic selected
		otherUser, err := databases.GetUserDetails(tx, otherUsername, ctx)
		if err != nil {
			return nil, err
		}
		if len(otherUser.Topics) == 0 {
			foundUsers = append(foundUsers, otherUsername)
		}

		// other user has matching topic
		if _, ok := sameTopicUsers[otherUsername]; ok {
			foundUsers = append(foundUsers, otherUsername)
			break
		}
	}

	// If no other user with same topics, then is not valid match
	if len(foundUsers) == 0 {
		return nil, models.NoMatchFound
	}

	return foundUsers, nil
}

func doDifficultyMatching(tx *redis.Tx, ctx context.Context, currentUser *models.MatchRequest, potentialMatches []string) ([]string, error) {
	sameDifficultyUsers, err := findSameDifficultiesUsers(tx, ctx, currentUser)
	if err != nil {
		return nil, err
	}

	// Iterate through the queue to find a match, so a user in the queue the longest is more likely to be matched.
	var foundUsers []string
	for _, otherUsername := range potentialMatches {
		if otherUsername == currentUser.Username {
			continue
		}

		// Include users without any difficulty selected
		otherUser, err := databases.GetUserDetails(tx, otherUsername, ctx)
		if err != nil {
			return nil, err
		}
		if len(otherUser.Topics) == 0 {
			foundUsers = append(foundUsers, otherUsername)
		}

		// other user has matching difficulty or has no difficulty
		if _, ok := sameDifficultyUsers[otherUsername]; ok {
			foundUsers = append(foundUsers, otherUsername)
			break
		}
	}

	// If no other user with same difficulties, then skip this criteria
	// so that at least a match is found
	if len(foundUsers) == 0 {
		return nil, models.NoMatchFound
	}

	return foundUsers, nil
}
