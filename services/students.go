package service

import (
	"awesomeTestProject/datastore"
	"awesomeTestProject/models"
	. "awesomeTestProject/shared"
	"context"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func ParseStudent(ctx context.Context, req HttpWebRequest) (*models.Student, error) {
	var student models.Student

	b, _ := req.Body()
	err := json.Unmarshal([]byte(b), &student)
	if err != nil {
		Fatal(ctx, "Unable to deserialize the request body")
		return nil, err
	}

	id := req.Param("id")
	if id != "" {
		student.Id = id
	}
	return &student, nil
}

func ParsePatchStudent(ctx context.Context, req HttpWebRequest) *models.PatchRequestPayload {
	var patchPayload models.PatchRequestPayload

	b, _ := req.Body()
	err := json.Unmarshal([]byte(b), &patchPayload)
	if err != nil {
		Fatal(ctx, "Unable to deserialize the request body")
		return nil
	}
	patchPayload.Id = req.Param("id")
	return &patchPayload
}

func ParseGetRequest(ctx context.Context, req HttpWebRequest) (map[string]interface{}, error) {
	params := make(map[string]interface{}, 0)

	id := req.Param("id")
	name := req.Param("name")
	order := req.Param("sortorder")
	field := req.Param("sortfield")

	if id != "" {
		params["id"] = id
	}
	if name != "" {
		params["name"] = name
	}
	if order != "" {
		params["sortorder"] = order
	}
	if field != "" {
		params["sortfield"] = field
	}
	return params, nil
}

func SaveStudent(ctx context.Context, student *models.Student, config *MapPropertySource) (*models.Student, error) {
	student.Id = uuid.NewV4().String()
	student.Meta.Created = time.Now()
	student.Meta.LastModified = time.Now()

	err := datastore.GetDatastore().Save(ctx, config.GetString("students-collection"), student)
	if err != nil {
		return nil, err
	}
	return student, nil
}

func GetStudents(ctx context.Context, params map[string]interface{}, config *MapPropertySource) ([]byte, error) {
	var student models.Student
	filter, opt, err := GetFilter(params)
	if err != nil {
		return nil, err
	}

	return datastore.GetDatastore().GetByFilter(ctx, config.GetString("students-collection"), filter, opt, &student)
}

func GetStudent(ctx context.Context, student *models.Student, config *MapPropertySource) error {
	filter := bson.D{{"id", student.Id}}
	return datastore.GetDatastore().GetById(ctx, config.GetString("students-collection"), filter, student)
}

func PatchStudent(ctx context.Context, patchPayload *models.PatchRequestPayload, config *MapPropertySource) (*models.Student, error) {
	var setElements bson.D
	for _, v := range patchPayload.Operations {
		setElements = append(setElements, bson.E{Key: v.Path, Value: v.Value})
	}
	setElements = append(setElements, bson.E{Key: "meta.lastModified", Value: time.Now()})

	query := bson.D{{"$set", setElements}}

	filter := bson.D{{"id", patchPayload.Id}}
	err := datastore.GetDatastore().Update(ctx, config.GetString("students-collection"), filter, query)
	if err != nil {
		return nil, err
	}

	var student models.Student
	student.Id = patchPayload.Id
	err = GetStudent(ctx, &student, config)
	if err != nil {
		return nil, err
	}

	return &student, err
}

func PutStudent(ctx context.Context, student *models.Student, config *MapPropertySource) (*models.Student, error) {
	var setElements bson.D

	setElements = append(setElements, bson.E{Key: "name", Value: student.Name})
	setElements = append(setElements, bson.E{Key: "meta.lastModified", Value: time.Now()})

	query := bson.D{{"$set", setElements}}

	filter := bson.D{{"id", student.Id}}
	err := datastore.GetDatastore().Update(ctx, config.GetString("students-collection"), filter, query)
	if err != nil {
		return nil, err
	}

	var studentDB models.Student
	studentDB.Id = student.Id
	err = GetStudent(ctx, &studentDB, config)
	if err != nil {
		return nil, err
	}

	return &studentDB, err
}

func DeleteStudent(ctx context.Context, id string, config *MapPropertySource) error {
	filter := bson.D{{"id", id}}
	return datastore.GetDatastore().Delete(ctx, config.GetString("students-collection"), filter)
}