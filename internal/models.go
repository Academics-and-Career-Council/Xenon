package internal

type User struct {
	ID              string   `bson:"_id,omitempty"`
	Index           int      `bson:"index"`
	Username        string   `bson:"username"`
	Name            string   `bson:"name"`
	Banned          bool     `bson:"banned"`
	RollNo          string   `bson:"rollno"`
	Role            string   `bson:"role"`
	KratosID        string   `bson:"kid"`
	RavenID         string   `bson:"rvid"`
	Credits         int      `bson:"credits"`
	CoursesUnlocked []string `bson:"courses_unlocked"`
	CoursesReviewed []string `bson:"courses_reviewed"`
}

type UserMetaData struct {
	Banned          bool     `bson:"banned"`
	Credits         int      `bson:"credits"`
	CoursesUnlocked []string `bson:"courses_unlocked"`
	Role            string   `bson:"role"`
}
