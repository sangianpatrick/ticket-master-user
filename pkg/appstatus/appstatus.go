package appstatus

type Status string

var (
	OK                  Status = Status("OK")
	Created             Status = Status("Created")
	NotFound            Status = Status("NotFound")
	InternalServerError Status = Status("InternalServerError")
	Conflict            Status = Status("Conflict")
	UnprocessableEntity Status = Status("UnprocessableEntity")
	BadRequest          Status = Status("BadRequest")
)
