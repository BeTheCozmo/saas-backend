package role_storage

import (
	"context"
	"fmt"
	"log"
	"time"

	p "uller/src/permission"
	r "uller/src/role"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type RoleStorage struct {
  client *mongo.Client
  collection *mongo.Collection
  databaseName, collectionName string
}

func New() *RoleStorage {
  return &RoleStorage{}
}

func (rs *RoleStorage) Configure(client *mongo.Client, database string) {
  rs.client = client
  rs.databaseName = database
  rs.collectionName = "role"
  collection := client.Database(rs.databaseName).Collection(rs.collectionName)
  rs.collection = collection

  rs.EnsureUserRoleCreated()
  rs.EnsureAdminRoleCreated()

  rs.ensureIndexes()
}

func (rs *RoleStorage) ensureIndexes() {
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()

  indexModel := []mongo.IndexModel{
    {
      Keys: bson.M{"name": 1},
      Options: options.Index().SetUnique(true),
    },
  }

  _, err := rs.collection.Indexes().CreateMany(ctx, indexModel)
  if err != nil {
    log.Fatalf("Error creating unique index: %v", err)
  }
}

func (rs *RoleStorage) EnsureUserRoleCreated() {
  filter := bson.M{"name": "user"}
  result := rs.collection.FindOne(context.Background(), filter)

  if result.Err() != nil {
    _, err := rs.collection.InsertOne(context.Background(), &r.Role{
      Name: "user",
      Permissions: []*p.Permission{
        {
          Name: "ent_get_name",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "ent_get_fantasyName",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "ent_get_document",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "ent_get_debt",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "ent_get_variation",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "ent_get_variationPercentage",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "ent_get_state",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "ent_get_city",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "ent_get_partners",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "ent_details",
          IsLimited: false,
          RemainingUses: 0,
        },
      },
    })

    if err != nil {
      fmt.Printf("Erro ao inserir o role user: %v", err)
    }
  }
}

func (rs *RoleStorage) EnsureAdminRoleCreated() {
  filter := bson.M{"name": "admin"}
  result := rs.collection.FindOne(context.Background(), filter)

  if result.Err() != nil {
    _, err := rs.collection.InsertOne(context.Background(), &r.Role{
      Name: "admin",
      Permissions: []*p.Permission{
        {
          Name: "user_login",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "user_get",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "user_change_name",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "user_change_profile",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "user_req_change_email",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "user_change_email",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "user_req_change_phone",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "user_change_phone",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "plan_get",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "plan_create",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "plan_change",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "plan_delete",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "plan_buy",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "plan_cancel",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "ent_get_name",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "ent_get_fantasyName",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "ent_get_document",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "ent_get_debt",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "ent_get_variation",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "ent_get_variatioPercentage",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "ent_get_state",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "ent_get_city",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "ent_get_partners",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "ent_search_partners",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "ent_get_details",
          IsLimited: false,
          RemainingUses: 0,
        },

        {
          Name: "ent_search_name",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "ent_search_fantasyName",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "ent_search_document",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "ent_search_debt",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "ent_search_variation",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "ent_search_variatioPercentage",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "ent_search_state",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "ent_search_city",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "ent_search_partners",
          IsLimited: false,
          RemainingUses: 0,
        },
        {
          Name: "ent_search_details",
          IsLimited: false,
          RemainingUses: 0,
        },
      },
    })
    if err != nil {
      fmt.Printf("Erro ao inserir a role admin: %v", err)
    }
  }
}

func (rs *RoleStorage) GetUserRole() *r.Role {
  return rs.GetRoleByName("user")
}

func (rs *RoleStorage) GetRoleByName(name string) *r.Role {
  filter := bson.M{"name": name}
  result := rs.collection.FindOne(context.Background(), filter)
  var role r.Role
  result.Decode(&role)
  return &role
}