package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Definindo um tipo personalizado para os métodos HTTP
type HTTPMethod string

// Definindo constantes para os métodos HTTP válidos
const (
	GET    HTTPMethod = "GET"
	POST   HTTPMethod = "POST"
	PUT    HTTPMethod = "PUT"
	DELETE HTTPMethod = "DELETE"
)

type Headers struct {
	Authorization string            `json:"Authorization,omitempty"` // Header de autenticação
	ContentType   string            `json:"Content-Type,omitempty"`  // Tipo do conteúdo
	Custom        map[string]string // Headers personalizados adicionais
}

// HttpRequestParams encapsula os parâmetros para a requisição
type HttpRequestParams struct {
	Method  HTTPMethod  // Método HTTP: GET, POST, PUT, DELETE.
	URL     string      // URL da API
	Headers Headers     // Headers personalizados
	Body    interface{} // Corpo da requisição (opcional)
}

func SendHttpRequest(params HttpRequestParams) (*http.Response, error) {
	// Validação do método HTTP
	switch params.Method {
	case GET, POST, PUT, DELETE:
		// O método é válido
	default:
		return nil, fmt.Errorf("método HTTP inválido: %s", params.Method)
	}
	if params.URL == "" {
		return nil, fmt.Errorf("URL é obrigatória")
	}
	// Converter o corpo para JSON, se fornecido
	var bodyData []byte
	if params.Body != nil {
		var err error
		bodyData, err = json.Marshal(params.Body)
		if err != nil {
			return nil, fmt.Errorf("erro ao converter o corpo para JSON: %v", err)
		}
	}

	// Criar a requisição
	req, err := http.NewRequest(string(params.Method), params.URL, bytes.NewBuffer(bodyData))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar a requisição: %v", err)
	}

	// Adicionar headers estruturados
	if params.Headers.Authorization != "" {
		req.Header.Set("Authorization", params.Headers.Authorization)
	}
	if params.Headers.ContentType != "" {
		req.Header.Set("Content-Type", params.Headers.ContentType)
	}
	for key, value := range params.Headers.Custom {
		req.Header.Set(key, value)
	}
	// Criar o cliente HTTP e enviar a requisição
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao enviar a requisição: %v", err)
	}

	return resp, nil
}
