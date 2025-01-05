package mail_sending

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/aws/aws-sdk-go/service/ses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (mh *MailHandler) Handler(ctx context.Context, event eventbridge.EventBridge) {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection(os.Getenv("MONGO_COLLECTION"))
	distinctOrg, err := collection.Distinct(ctx, "org", bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	var result *mongo.Cursor
	for _, org := range distinctOrg {
		result, err = collection.Aggregate(ctx, bson.A{
			bson.D{{Key: "$match", Value: bson.D{{Key: "org", Value: org}}}},
			bson.D{
				{Key: "$project",
					Value: bson.D{
						{Key: "_id", Value: 1},
						{Key: "assignedTokenLimit", Value: "$assignedTokenLimit"},
						{Key: "org", Value: "$org"},
						{Key: "provider", Value: "$provider"},
						{Key: "remainingTokenLimit", Value: "$remainingTokenLimit"},
						{Key: "usedTokenLimit",
							Value: bson.D{
								{Key: "$subtract",
									Value: bson.A{
										"$assignedTokenLimit",
										"$remainingTokenLimit",
									},
								},
							},
						},
					},
				},
			},
		})
		if err != nil {
			log.Fatal(err)
		}
		var tokenResponse []TokenLimit
		if err = result.All(ctx, &tokenResponse); err != nil {
			panic(err)
		}

		mh.SendEmail(ctx, fmt.Sprint(org), tokenResponse)
	}
}

func (mh *MailHandler) SendEmail(ctx context.Context, org string, details []TokenLimit) {
	emailTemplate, err := mh.GenerateSESTemplate(ctx, org, details)
	if err != nil {
		log.Fatal(err.Error())
	}

	service := ses.New(mh.SessionInstance)

	// Attempt to send the email.
	_, err = service.SendEmail(emailTemplate)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Fatal(aerr.Error())
		} else {
			log.Fatal(err)
		}
	}
}
