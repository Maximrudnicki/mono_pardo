package response

type GroupResponse struct {
	UserId   int    `json:"user_id"`
	GroupId  string `json:"group_id"`
	Title    string `json:"title"`
	Students []int  `json:"students"`
}

type StudentResponse struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

type TeacherResponse struct {
	TeacherId int    `json:"teacher_id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
}

type StatisticsResponse struct {
	StatId    string `json:"statistics_id"`
	GroupId   string `json:"group_id"`
	TeacherId int    `json:"teacher_id"`
	StudentId int    `json:"student_id"`
	Words     []int  `json:"words"`
}

type AddWordToUserResponse struct {
	WordId int `json:"word_id"`
}

type StudentInfo struct {
	StudentId int    `json:"student_id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
}

type StudentInformation struct {
	StudentId int             `json:"student_id"`
	Email     string          `json:"email"`
	Username  string          `json:"username"`
	Words     []VocabResponse `json:"words"`
}
