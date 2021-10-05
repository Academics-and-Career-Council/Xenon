package database

type User struct {
	ID         string `bson:"_id,omitempty"`
	Username   string `bson:"username"`
	Name       string `bson:"name"`
	Banned     bool   `bson:"banned"`
	RollNo     string `bson:"rollno"`
	Department string `bson:"department"`
	KratosID   string `bson:"kid"`
}
