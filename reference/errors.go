package reference // import "github.com/DevanshMathur19/docker-v23/reference"

type notFoundError string

func (e notFoundError) Error() string {
	return string(e)
}

func (notFoundError) NotFound() {}

type invalidTagError string

func (e invalidTagError) Error() string {
	return string(e)
}

func (invalidTagError) InvalidParameter() {}

type conflictingTagError string

func (e conflictingTagError) Error() string {
	return string(e)
}

func (conflictingTagError) Conflict() {}
