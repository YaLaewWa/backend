// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag/v2"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},"swagger":"2.0","info":{"description":"{{escape .Description}}","title":"{{.Title}}","contact":{"name":"API Support"},"version":"{{.Version}}"},"host":"{{.Host}}","basePath":"{{.BasePath}}","paths":{"/auth/login":{"post":{"description":"Login to the system","consumes":["application/json"],"produces":["application/json"],"tags":["Auth"],"summary":"Login to the system","parameters":[{"description":"Description of the authentication body","name":"user","in":"body","required":true,"schema":{"$ref":"#/definitions/socket_internal_dto.AuthBody"}}],"responses":{"200":{"description":"successfully login","schema":{"$ref":"#/definitions/socket_internal_dto.SuccessResponse-socket_internal_dto_UserWithTokenResponse"}},"400":{"description":"your request is invalid or your request body is incorrect or cannot save user","schema":{"$ref":"#/definitions/socket_internal_dto.ErrorResponse"}},"500":{"description":"cannot use this password","schema":{"$ref":"#/definitions/socket_internal_dto.ErrorResponse"}}}}},"/auth/register":{"post":{"description":"Register a user","consumes":["application/json"],"produces":["application/json"],"tags":["Auth"],"summary":"Register a user","parameters":[{"description":"Description of the authentication body","name":"user","in":"body","required":true,"schema":{"$ref":"#/definitions/socket_internal_dto.AuthBody"}}],"responses":{"201":{"description":"register successfully","schema":{"$ref":"#/definitions/socket_internal_dto.SuccessResponse-socket_internal_dto_UserResponse"}},"400":{"description":"your request is invalid or your request body is incorrect or cannot save user","schema":{"$ref":"#/definitions/socket_internal_dto.ErrorResponse"}},"409":{"description":"username already exists","schema":{"$ref":"#/definitions/socket_internal_dto.ErrorResponse"}},"500":{"description":"cannot use this password","schema":{"$ref":"#/definitions/socket_internal_dto.ErrorResponse"}}}}}},"definitions":{"socket_internal_dto.AuthBody":{"type":"object","required":["password","username"],"properties":{"password":{"type":"string","minLength":9},"username":{"type":"string","minLength":4}}},"socket_internal_dto.ErrorResponse":{"type":"object","properties":{"error":{"type":"string"}}},"socket_internal_dto.SuccessResponse-socket_internal_dto_UserResponse":{"type":"object","properties":{"data":{"$ref":"#/definitions/socket_internal_dto.UserResponse"}}},"socket_internal_dto.SuccessResponse-socket_internal_dto_UserWithTokenResponse":{"type":"object","properties":{"data":{"$ref":"#/definitions/socket_internal_dto.UserWithTokenResponse"}}},"socket_internal_dto.UserResponse":{"type":"object","properties":{"username":{"type":"string"}}},"socket_internal_dto.UserWithTokenResponse":{"type":"object","properties":{"accessToken":{"type":"string"},"user":{"$ref":"#/definitions/socket_internal_dto.UserResponse"}}}}}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "chadChat API",
	Description:      "This is a swagger to show all RestAPI of chadChat project",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
