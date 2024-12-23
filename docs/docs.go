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
        "/shortener/v1/link": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Links"
                ],
                "summary": "Create a short link",
                "parameters": [
                    {
                        "description": "New link",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.CreateLinkInput"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/dto.CreateLinkOutput"
                        }
                    },
                    "302": {
                        "description": "Found",
                        "schema": {
                            "$ref": "#/definitions/dto.CreateLinkOutput"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/http.ErrHTTP"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/http.ErrHTTP"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.ErrHTTP"
                        }
                    }
                }
            }
        },
        "/shortener/v1/link/{alias}": {
            "get": {
                "consumes": [
                    "text/plain"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Links"
                ],
                "summary": "Get a short link by alias",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Link alias",
                        "name": "alias",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.GetLinkOutput"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/http.ErrHTTP"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/http.ErrHTTP"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.ErrHTTP"
                        }
                    }
                }
            }
        },
        "/shortener/v1/link/{alias}/redirect": {
            "get": {
                "consumes": [
                    "text/plain"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Links"
                ],
                "summary": "Redirect to URL by alias",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Link alias",
                        "name": "alias",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "302": {
                        "description": "redirect to the original url"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/http.ErrHTTP"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/http.ErrHTTP"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.ErrHTTP"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.CreateLinkInput": {
            "type": "object",
            "properties": {
                "url": {
                    "type": "string"
                }
            }
        },
        "dto.CreateLinkOutput": {
            "type": "object",
            "properties": {
                "alias": {
                    "type": "string"
                },
                "expired_at": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "dto.GetLinkOutput": {
            "type": "object",
            "properties": {
                "alias": {
                    "type": "string"
                },
                "expired_at": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "http.ErrHTTP": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.0.0",
	Host:             "",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Title",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
