package mongodb

import (
	"net/http"

	"time"

	"github.com/pushpaldev/base-api/helpers"
	"github.com/pushpaldev/base-api/helpers/params"
	"github.com/pushpaldev/base-api/models"
	"gopkg.in/mgo.v2/bson"
)

func (db *mongo) CreateUser(user *models.User) error {
	session := db.Session.Copy()
	defer session.Close()
	users := db.C(models.UsersCollection).With(session)

	user.Id = bson.NewObjectId().Hex()
	err := user.BeforeCreate()
	if err != nil {
		return err
	}

	if count, _ := users.Find(bson.M{"email": user.Email}).Count(); count > 0 {
		return helpers.NewError(http.StatusConflict, "user_already_exists", "User already exists")
	}

	err = users.Insert(user)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "user_creation_failed", "Failed to insert the user in the database")
	}

	return nil
}

func (db *mongo) FindUserById(id string) (*models.User, error) {
	session := db.Session.Copy()
	defer session.Close()
	users := db.C(models.UsersCollection).With(session)

	user := &models.User{}
	err := users.FindId(id).One(user)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "user_not_found", "User not found")
	}

	return user, err
}

func (db *mongo) FindUser(params params.M) (*models.User, error) {
	session := db.Session.Copy()
	defer session.Close()
	users := db.C(models.UsersCollection).With(session)

	user := &models.User{}

	err := users.Find(params).One(user)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "user_not_found", "User not found")
	}

	return user, err
}

func (db *mongo) ActivateUser(activationKey string, id string) error {
	session := db.Session.Copy()
	defer session.Close()
	users := db.C(models.UsersCollection).With(session)

	err := users.Update(bson.M{"$and": []bson.M{{"_id": id}, {"activationKey": activationKey}}}, bson.M{"$set": bson.M{"active": true}})
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "user_activation_failed", "Couldn't find the user to activate")
	}
	return nil
}

func (db *mongo) UpdateUser(user *models.User, params params.M) error {
	session := db.Session.Copy()
	defer session.Close()
	users := db.C(models.UsersCollection).With(session)

	err := users.UpdateId(user.Id, params)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "user_update_failed", "Failed to update the user")
	}

	return nil
}

func (db *mongo) AddLoginToken(user *models.User, ip string) (*models.LoginToken, error) {
	session := db.Session.Copy()
	defer session.Close()
	users := db.C(models.UsersCollection).With(session)

	token := &models.LoginToken{
		Id:         bson.NewObjectId().Hex(),
		Ip:         ip,
		CreatedAt:  time.Now().Unix(),
		LastAccess: time.Now().Unix(),
	}

	if err := users.UpdateId(user.Id, bson.M{"$push": bson.M{"tokens": token}}); err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "user_token_creation_failed", "Failed to create the token.")
	}

	return token, nil
}

func (db *mongo) RemoveLoginToken(user *models.User, tokenId string) error {
	session := db.Session.Copy()
	defer session.Close()
	users := db.C(models.UsersCollection).With(session)

	if err := users.UpdateId(user.Id, bson.M{"$pull": bson.M{"tokens": bson.M{"_id": tokenId}}}); err != nil {
		return helpers.NewError(http.StatusInternalServerError, "user_token_deletion_failed", "Failed to delete the token.")
	}

	return nil
}
