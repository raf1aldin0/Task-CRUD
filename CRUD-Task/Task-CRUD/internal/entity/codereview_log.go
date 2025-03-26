package entity

import "time"

// CodeReviewLog merepresentasikan log review kode oleh AI
type CodeReviewLog struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	RepoID    int       `json:"repo_id" bson:"repo_id"`
	Review    string    `json:"review" bson:"review"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
