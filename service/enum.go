package service

// PERMISSIONS permission type
type PERMISSIONS string

const (
	// CREATE permission
	CREATE = PERMISSIONS("CREATE")
	// READ permission
	READ = PERMISSIONS("READ")
	// DELETE permission
	DELETE = PERMISSIONS("DELETE")
	// UPDATE permission
	UPDATE = PERMISSIONS("UPDATE")
)

// RESOURCE name
type RESOURCE string

const (
	// PIPELINE refers to pipeline resource
	PIPELINE = RESOURCE("pipeline")
	// PROCESS refers to process resource
	PROCESS = RESOURCE("process")
	// COMPANY refers to company resource
	COMPANY = RESOURCE("company")
	// REPOSITORY refers to repository resource
	REPOSITORY = RESOURCE("repository")
	// APPLICATION refers to application resource
	APPLICATION = RESOURCE("application")
)
