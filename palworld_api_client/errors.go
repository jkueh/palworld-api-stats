package api_client

import "fmt"

type Non200Error struct {
	error
	statusCode int
}

func (e *Non200Error) String() string {
	return fmt.Sprintf("The server responded with a non-200 error code: %d", e.statusCode)
}

type BodyReadError struct{ error }

func (e *BodyReadError) String() string {
	return fmt.Sprintf("Unable to read the response body: %v", e.error)
}
