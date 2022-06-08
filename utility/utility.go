package utility

import (
	"errors"
	"github.com/klovercloud-ci-cd/terminal/adapter/enum/Authority"
	"github.com/klovercloud-ci-cd/terminal/adapter/universal"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

const (
	CURRENT_USER   = "current-user"
	SUCCESS_STATUS = "success"
)

func GetHeaderData(r *http.Request) (universal.User, bson.ObjectId, error) {
	currentUserString := r.Header.Get(CURRENT_USER)
	log.Println("current-user:", currentUserString)
	user, err := universal.GenerateUserFromJsonStr(currentUserString)
	if err != nil {
		return user, "", err
	}
	return user, "", nil
}

func CheckPermission(user universal.User) error {
	if user.HasAnyAuthority(Authority.ROLE_AGENT_SERVICE, Authority.ROLE_GIT_AGENT, Authority.ANONYMOUS, Authority.ROLE_SNAPSHOTTER) {
		return errors.New("permission denied")
	}
	return nil
}
