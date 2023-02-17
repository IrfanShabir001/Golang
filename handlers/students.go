package handlers

import (
	"awesomeTestProject/models"
	"awesomeTestProject/services"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	. "awesomeTestProject/shared"
)

func PostStudentHandler(req HttpWebRequest, ctx context.Context, config *MapPropertySource) *ResponseOut {
	ri := &ResponseOut{}
	Info(ctx, "Parse request")

	Info(ctx, "Parse request body")
	postRequestPayload, err := service.ParseStudent(ctx, req)
	ErrorCheck(err)
	ErrorCheckNilThrowInvalidParam(postRequestPayload)

	Info(ctx, "Parsing completed, Create Student")
	student, err := service.SaveStudent(ctx, postRequestPayload, config)
	ErrorCheck(err)

	Info(ctx, "Student created.")
	res, _ := json.Marshal(student)
	ri.Body(res)
	ri.Status(http.StatusCreated)
	return ri
}

func PutStudentHandler(req HttpWebRequest, ctx context.Context, config *MapPropertySource) *ResponseOut {
	ri := &ResponseOut{}
	Info(ctx, "Parse request")

	Info(ctx, "Parse request body")
	postRequestPayload, err := service.ParseStudent(ctx, req)
	ErrorCheck(err)
	ErrorCheckNilThrowInvalidParam(postRequestPayload)

	Info(ctx, fmt.Sprintf("Parsing completed, Update Student(%s)", postRequestPayload.Id))
	student, err := service.PutStudent(ctx, postRequestPayload, config)
	ErrorCheck(err)

	Info(ctx, fmt.Sprintf("Student(%s) updated.", student.Id))
	res, _ := json.Marshal(student)
	ri.Body(res)
	ri.Status(http.StatusOK)
	return ri
}

func PatchStudentHandler(req HttpWebRequest, ctx context.Context, config *MapPropertySource) *ResponseOut {
	ri := &ResponseOut{}
	Info(ctx, "Parse request")

	Info(ctx, "Parse request body")
	patchRequestPayload := service.ParsePatchStudent(ctx, req)
	ErrorCheckNilThrowInvalidParam(patchRequestPayload)

	Info(ctx, fmt.Sprintf("Parsing completed, Patch Student(%s)", patchRequestPayload.Id))
	student, err := service.PatchStudent(ctx, patchRequestPayload, config)
	ErrorCheck(err)

	Info(ctx, fmt.Sprintf("Student(%s) patched.", student.Id))
	res, _ := json.Marshal(student)
	ri.Body(res)
	ri.Status(http.StatusOK)
	return ri
}

func GetStudentsHandler(req HttpWebRequest, ctx context.Context, config *MapPropertySource) *ResponseOut {
	ri := &ResponseOut{}
	Info(ctx, "Parse request")

	Info(ctx, "Parse request body")
	params, err := service.ParseGetRequest(ctx, req)
	ErrorCheck(err)

	Info(ctx, "Parsing completed, Getting Students")
	item, err := service.GetStudents(ctx, params, config)
	ErrorCheck(err)

	Info(ctx, "Got Students.")
	ri.Body(item)
	ri.Status(http.StatusOK)
	return ri
}

func GetStudentByIdHandler(req HttpWebRequest, ctx context.Context, config *MapPropertySource) *ResponseOut {
	ri := &ResponseOut{}
	Info(ctx, "Parse request")

	var student models.Student
	student.Id = req.Param("id")
	Info(ctx, fmt.Sprintf("Get Student(%s)", student.Id))
	err := service.GetStudent(ctx, &student, config)
	if err != nil && err.Error() == "mongo: no documents in result" {
		ri.Body([]byte{})
		ri.Status(http.StatusOK)
		return ri
	}
	ErrorCheck(err)

	Info(ctx, fmt.Sprintf("Student(%s).", student.Id))
	res, _ := json.Marshal(student)
	ri.Body(res)
	ri.Status(http.StatusOK)
	return ri
}

func DeleteStudentHandler(req HttpWebRequest, ctx context.Context, config *MapPropertySource) *ResponseOut {
	ri := &ResponseOut{}
	Info(ctx, "Parse request")

	id := req.Param("id")
	Info(ctx, fmt.Sprintf("Delete Student(%s)", id))
	err := service.DeleteStudent(ctx, id, config)
	ErrorCheck(err)

	Info(ctx, fmt.Sprintf("Deleted Student(%s).",id))
	ri.Body([]byte{})
	ri.Status(http.StatusOK)
	return ri
}
