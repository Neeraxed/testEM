// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://github.com/Neeraxed",
        "contact": {
            "name": "Alina Kuznetsova",
            "email": "Neeraxed@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/songs": {
            "get": {
                "description": "get string by filters",
                "produces": [
                    "application/json"
                ],
                "summary": "Get songs",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entities.SongsWrapper"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/delivery.HttpError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/delivery.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/delivery.HttpError"
                        }
                    }
                }
            },
            "post": {
                "description": "add song",
                "produces": [
                    "application/json"
                ],
                "summary": "Add song",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entities.Song"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/delivery.HttpError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/delivery.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/delivery.HttpError"
                        }
                    }
                }
            },
            "delete": {
                "description": "delete song with specified id",
                "produces": [
                    "application/json"
                ],
                "summary": "Delete song",
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/delivery.HttpError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/delivery.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/delivery.HttpError"
                        }
                    }
                }
            },
            "patch": {
                "description": "update song with specified id",
                "produces": [
                    "application/json"
                ],
                "summary": "Patch song",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entities.Song"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/delivery.HttpError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/delivery.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/delivery.HttpError"
                        }
                    }
                }
            }
        },
        "/songs/{id}/verses": {
            "get": {
                "description": "get verses for song",
                "produces": [
                    "application/json"
                ],
                "summary": "Get verses",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entities.VersesWrapper"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/delivery.HttpError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/delivery.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/delivery.HttpError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "delivery.HttpError": {
            "type": "object"
        },
        "entities.Song": {
            "type": "object",
            "properties": {
                "group": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "link": {
                    "type": "string"
                },
                "releaseDate": {
                    "type": "string"
                },
                "song": {
                    "type": "string"
                }
            }
        },
        "entities.SongsWrapper": {
            "type": "object",
            "properties": {
                "songs": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entities.Song"
                    }
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "entities.Verse": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "num": {
                    "type": "integer"
                },
                "song_id": {
                    "type": "string"
                }
            }
        },
        "entities.VersesWrapper": {
            "type": "object",
            "properties": {
                "total": {
                    "type": "integer"
                },
                "verses": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entities.Verse"
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.0.1",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "TestEM API",
	Description:      "Test app for EM from Alina Kuznetsova",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
