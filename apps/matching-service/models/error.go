package models

import "errors"

var ExistingUserError = errors.New("already has an existing user in matchmaking")
var NoMatchFound = errors.New("no matches found for current user")
