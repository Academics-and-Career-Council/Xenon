package database

type User struct {
	ID         string `bson:"_id,omitempty" json:"id"`
	Username   string `bson:"username" json:"username"`
	Name       string `bson:"name" json:"name"`
	Banned     bool   `bson:"banned" json:"banned"`
	RollNo     string `bson:"rollno" json:"rollno"`
	Department string `bson:"branch" json:"department"`
	KratosID   string `bson:"kid" json:"registration_id"`
	EmailID    string `bson:"email" json:"email"`
}
