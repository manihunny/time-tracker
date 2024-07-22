// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/people": {
            "get": {
                "description": "Get people list",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "people"
                ],
                "summary": "Get people list",
                "parameters": [
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "Number of results per page",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "Number of page",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The substring that will be searched for in the 'surname' field",
                        "name": "surname",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The substring that will be searched for in the 'name' field",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The substring that will be searched for in the 'patronymic' field",
                        "name": "patronymic",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The substring that will be searched for in the 'address' field",
                        "name": "address",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The substring that will be searched for in the 'passport_number' field",
                        "name": "passport_number",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/peoplerepository.People"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "post": {
                "description": "Create people",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "people"
                ],
                "summary": "Create people",
                "parameters": [
                    {
                        "description": "User's passport number",
                        "name": "passportNumber",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/peoplerepository.People"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/people/{id}": {
            "get": {
                "description": "Get people by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "people"
                ],
                "summary": "Get people by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/peoplerepository.People"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "put": {
                "description": "Update people by ID",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "people"
                ],
                "summary": "Update people by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "User's surname",
                        "name": "surname",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "User's name",
                        "name": "name",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "User's patronymic",
                        "name": "patronymic",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "User's address",
                        "name": "address",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "User's passport number",
                        "name": "passport_number",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "delete": {
                "description": "Delete people by ID",
                "tags": [
                    "people"
                ],
                "summary": "Delete people by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "patch": {
                "description": "Partial update people by ID",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "people"
                ],
                "summary": "Partial update people by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "User's surname",
                        "name": "surname",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "User's name",
                        "name": "name",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "User's patronymic",
                        "name": "patronymic",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "User's address",
                        "name": "address",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "User's passport number",
                        "name": "passport_number",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/people/{id}/finish-task": {
            "post": {
                "description": "Complete all unfinished tasks for user",
                "tags": [
                    "people"
                ],
                "summary": "Finish task for user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/people/{id}/start-task": {
            "post": {
                "description": "Start new task for user",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "people"
                ],
                "summary": "Start new task for user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Task title",
                        "name": "title",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/people/{id}/task-statistics": {
            "post": {
                "description": "Get task statistics for user. Calculates time spent for tasks and retrieve all tasks data. If task not finished, time spent is calculated up to the current date.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "people"
                ],
                "summary": "Get user task statistics for period",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Begin of period",
                        "name": "date_from",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "End of period",
                        "name": "date_to",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/taskrepository.Task"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        }
    },
    "definitions": {
        "peoplerepository.People": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string",
                    "example": "г. Москва, ул. Ленина, д. 5, кв. 1"
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "name": {
                    "type": "string",
                    "example": "Иван"
                },
                "passport_number": {
                    "type": "string",
                    "example": "1234 567890"
                },
                "patronymic": {
                    "type": "string",
                    "example": "Иванович"
                },
                "surname": {
                    "type": "string",
                    "example": "Иванов"
                }
            }
        },
        "taskrepository.Task": {
            "type": "object",
            "properties": {
                "finished_at": {
                    "type": "string",
                    "example": "2024-07-17T00:00:00Z"
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "started_at": {
                    "type": "string",
                    "example": "2024-07-17T00:00:00Z"
                },
                "time_spent": {
                    "type": "integer",
                    "example": 48393984418000
                },
                "time_spent_formatted": {
                    "type": "string",
                    "example": "13h26m33.984418s"
                },
                "title": {
                    "type": "string",
                    "example": "Выполнить задачу 1"
                },
                "user_id": {
                    "type": "integer",
                    "example": 1
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "Time Tracker API",
	Description:      "API Server for Time Tracker application",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
