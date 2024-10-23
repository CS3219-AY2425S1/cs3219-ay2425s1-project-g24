package models

import "errors"

var ExistingUserError = errors.New("already has an existing user in matchmaking")
