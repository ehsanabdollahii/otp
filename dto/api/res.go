package api

type Error struct {
	Code    string
	Message string
}

type Response struct {
	OK       bool
	Response any
	Errors   []Error
}
