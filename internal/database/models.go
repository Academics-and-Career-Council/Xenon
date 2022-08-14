package database

type User struct {
	ID         string `bson:"_id,omitempty" json:"id" csv:"-"`
	Username   string `bson:"username" json:"username" csv:"username"`
	Name       string `bson:"name" json:"name" csv:"name"`
	Banned     bool   `bson:"banned" json:"banned" csv:"-"`
	RollNo     string `bson:"rollno" json:"rollno" csv:"rollno"`
	Department string `bson:"branch" json:"department" csv:"department"`
	KratosID   string `bson:"kid" json:"registration_id" csv:"-"`
	EmailID    string `bson:"email" json:"email" csv:"emailid"`
	Role       string `bson:"role" json:"role" csv:"role"`
}
