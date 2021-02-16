package model

import "time"

// User struct for our gql User type
type User struct {
	Sub        string     `json:"sub"`
	Bio        *string    `json:"bio"`
	Email      Email      `bson:"email"`
	Avatar     *string    `json:"avatar"`
	Username   *string    `json:"username"`
	JoinDate   *time.Time `bson:"join_date"`
	GivenName  *string    `bson:"given_name"`
	FamilyName *string    `bson:"family_name"`
}
