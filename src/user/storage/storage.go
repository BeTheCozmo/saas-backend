package user_storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
	"uller/src/logger"
	r "uller/src/role"
	user "uller/src/user"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UserStorage struct {
  logger *logger.Logger
  client *mongo.Client
  collection *mongo.Collection
  databaseName, collectionName string
}

func New() *UserStorage {
  return &UserStorage{}
}

func (us *UserStorage) Configure(client *mongo.Client, database string, logger *logger.Logger) {
  us.client = client
  us.databaseName = database
  us.collectionName = "user"
  collection := client.Database(us.databaseName).Collection(us.collectionName)
  us.collection = collection
  us.logger = logger

  us.ensureIndexes()
}

func (us *UserStorage) ensureIndexes() {
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()

  indexModel := []mongo.IndexModel{
    {
      Keys: bson.M{"email": 1},
      Options: options.Index().SetUnique(true),
    },
  }

  _, err := us.collection.Indexes().CreateMany(ctx, indexModel)
  if err != nil {
    log.Fatalf("Error creating unique index: %v", err)
  }
}

func (us *UserStorage) Create(user *user.User) (error) {
  _, err := us.collection.InsertOne(context.Background(), user)
  if err != nil {
    return err
  }
  return nil
}

func (us *UserStorage) GetByEmail(email string) (*user.User, error) {
  filter := bson.M{"email": email}
  var userSearched user.User
  err := us.collection.FindOne(context.Background(), filter).Decode(&userSearched)
  
  if err != nil {
    return nil, err
  }
  return &userSearched, nil
}

func (us *UserStorage) GetByPhone(phone string) (*user.User, error) {
  filter := bson.M{"phone": phone}
  var userSearched user.User
  err := us.collection.FindOne(context.Background(), filter).Decode(&userSearched)
  if err != nil {
    return nil, err
  }
  return &userSearched, nil
}

func (us *UserStorage) ChangeUserRoleTo(user *user.User, role *r.Role) error {
  filter := bson.M{"email": user.Email}

  update := bson.M{
    "$set": bson.M{
      "role": role,
    },
  }

  result, err := us.collection.UpdateOne(context.Background(), filter, update)
  if err != nil {
    us.logger.Warn(err)
    return err
  }

  if result.MatchedCount == 0 {
    message := fmt.Sprintf("no user found with email %s", user.Email)
    us.logger.Warn(message)
    return errors.New(message)
  }

  if result.ModifiedCount == 0 {
    message := fmt.Sprintf("failed to update role '%s' from user '%s'", role.Name, user.Email)
    us.logger.Warn(message)
    return errors.New(message)
  }

  return nil
}

func (us *UserStorage) SaveUser(user *user.User) error {
  filter := bson.M{"email": user.Email}
  if user.Email == "" {
    return errors.New("user email is required to save")
  }

  update := bson.M{
    "$set": bson.M{
      "permissions": user.Permissions,
      "role": user.Role,
      "plan": user.Plan,
    },
  }

  result, err := us.collection.UpdateOne(context.Background(), filter, update)
  if err != nil {
    us.logger.Error("failed to save user:", err)
    return fmt.Errorf("failed to save user: %v", err)
  }

  if result.MatchedCount == 0 {
    us.logger.Warn("no user found with email:", user.Email)
    return fmt.Errorf("user not found")
  }
  if result.ModifiedCount == 0 {
    us.logger.Info("No changes made to user with email:", user.Email)
  }
  us.logger.Debug("User", user.Email,  "saved successfully")
  return nil
}