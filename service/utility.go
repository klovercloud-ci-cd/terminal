package service

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt"
	"log"
)

func checkAuthority(userResourcePermission UserResourcePermission, resourceName, role, permission string) error {
	var resourceWiseRoles ResourceWiseRoles
	for _, resource := range userResourcePermission.Resources {
		if resource.Name == resourceName {
			resourceWiseRoles = resource
			break
		}
	}
	if role != "" {
		for _, each := range resourceWiseRoles.Roles {
			if each.Name == role {
				return nil
			}
		}
	} else if permission != "" {
		for _, each := range resourceWiseRoles.Roles {
			for _, perm := range each.Permissions {
				if perm.Name == permission {
					return nil
				}
			}

		}
	}
	return errors.New("[ERROR]: Insufficient permission")
}

func GetUserResourcePermissionFromBearerToken(token string) (UserResourcePermission, error) {
	res, _ := ValidateToken(token)
	if !res {
		return UserResourcePermission{}, errors.New("[ERROR]: Token is expired")
	}
	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})
	jsonBody, err := json.Marshal(claims["data"])
	if err != nil {
		log.Println(err)
	}
	userResourcePermission := UserResourcePermission{}
	if err := json.Unmarshal(jsonBody, &userResourcePermission); err != nil {
		return UserResourcePermission{}, errors.New("[ERROR]: No resource permissions")
	}
	return userResourcePermission, nil
}
