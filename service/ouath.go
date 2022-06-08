package service

// Permission roles permission
type Permission struct {
	Name string `json:"name" bson:"name"`
}

// Role users roles
type Role struct {
	Name        string       `json:"name" bson:"name"`
	Permissions []Permission `json:"permissions" bson:"permissions"`
}

// UserResourcePermission user and resources wise role
type UserResourcePermission struct {
	Metadata  UserMetadata        `json:"metadata" bson:"-"`
	UserId    string              `json:"user_id" bson:"user_id"`
	Resources []ResourceWiseRoles `json:"resources" bson:"resources"`
}

type AgentData struct {
	Name string `json:"name" bson:"name"`
}

// UserMetadata users metadata
type UserMetadata struct {
	CompanyId string `json:"company_id" bson:"company_id"`
}

// ResourceWiseRoles resource wise roles
type ResourceWiseRoles struct {
	Name  string `json:"name" bson:"name"`
	Roles []Role `json:"roles" bson:"roles"`
}
