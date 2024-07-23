package request

type AddStudentRequest struct {
	Token   string `json:"token"`
	GroupId string `json:"group_id"`
}

type AddWordToUserRequest struct {
	Word       string `json:"word"`
	Definition string `json:"definition"`
	GroupId    string `json:"group_id"`
	UserId     int    `json:"user_id"`
	Token      string `json:"token"`
}

type CreateGroupRequest struct {
	Token string `json:"token"`
	Title string `json:"title"`
}

type DeleteGroupRequest struct {
	Token   string `json:"token"`
	GroupId string `json:"group_id"`
}

type FindGroupRequest struct {
	Token   string `json:"token"`
	GroupId string `json:"group_id"`
}

type FindStudentRequest struct {
	Token     string `json:"token"`
	StudentId int    `json:"student_id"`
	GroupId   string `json:"group_id"`
}

type FindTeacherRequest struct {
	Token   string `json:"token"`
	GroupId string `json:"group_id"`
}

type FindGroupsTeacherRequest struct {
	Token string `json:"token"`
}

type FindGroupsStudentRequest struct {
	Token string `json:"token"`
}

type GetStatisticsRequest struct {
	Token     string `json:"token"`
	GroupId   string `json:"group_id"`
	StudentId int    `json:"student_id"`
}

type RemoveStudentRequest struct {
	Token   string `json:"token"`
	GroupId string `json:"group_id"`
	UserId  int    `json:"user_id"`
}
