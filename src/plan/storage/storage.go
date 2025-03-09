package plan_storage

import (
	"context"
	p "uller/src/permission"
	pl "uller/src/plan"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PlanStorage struct {
  client *mongo.Client
  collection *mongo.Collection
  databaseName, collectionName string
}

func New() *PlanStorage {
  return &PlanStorage{}
}

func (ps *PlanStorage) Configure(client *mongo.Client, database string) {
  ps.client = client
  ps.databaseName = database
  ps.collectionName = "plan"
  collection := client.Database(ps.databaseName).Collection(ps.collectionName)
  ps.collection = collection

  ps.EnsureFreePlanCreated()
}

func (ps *PlanStorage) EnsureFreePlanCreated() {
  filter := bson.M{"name": "free"}
  result := ps.collection.FindOne(context.Background(), filter)

  if result.Err() != nil {
    ps.collection.InsertOne(context.Background(), &pl.Plan{
      Name: "free",
      Permissions: []*p.Permission{
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
        

        

        // {
        //   Name: "ent_search_name",
        //   IsLimited: false,
        //   RemainingUses: 0,
        // },
        // {
        //   Name: "ent_search_fantasy_name",
        //   IsLimited: false,
        //   RemainingUses: 0,
        // },
        // {
        //   Name: "ent_search_document",
        //   IsLimited: false,
        //   RemainingUses: 0,
        // },
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
          Name: "ent_search_variationPercentage",
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
        // {
        //   Name: "ent_search_partners",
        //   IsLimited: false,
        //   RemainingUses: 0,
        // },
        {
          Name: "ent_details",
          IsLimited: true,
          RemainingUses: 20,
        },
        {
          Name: "ent_history",
          IsLimited: true,
          RemainingUses: 20,
        },
      },
    })
  }
}

func (ps *PlanStorage) GetFreePlan() *pl.Plan {
  filter := bson.M{"name": "free"}
  result := ps.collection.FindOne(context.Background(), filter)
  var plan pl.Plan
  result.Decode(&plan)

  return &plan
}