package database

func ConnectToMongo(ctx context.Context, url string) (_ *mongo.Client, err error) {
	MongoClient, err := mongo.NewClient(options.Client().ApplyURI(url))

	if err != nil {
		return nil, err
	}

	c, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = MongoClient.Connect(c)

	if err != nil {
		return nil, err
	}

	return MongoClient, nil
}

func MongoCreateCollections(ctx context.Context, client *mongo.Client) error {
	database := client.Database(MongoDBName)

	if err := database.CreateCollection(ctx, TokensMongoCollection); err != nil {
		if _, ok := err.(mongo.CommandError); !ok {
			return err
		}
	}

	if err := database.CreateCollection(ctx, ApplicationMongoCollection); err != nil {
		if _, ok := err.(mongo.CommandError); !ok {
			return err
		}
	}

	applicationIndexes := []mongo.IndexModel{
		{
			Keys: bson.M{
				"Name": 1,
			},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := database.Collection(ApplicationMongoCollection).Indexes().CreateMany(ctx, applicationIndexes)

	if err != nil {
		return err
	}

	if err := database.CreateCollection(ctx, SecretsMongoCollection); err != nil {
		if _, ok := err.(mongo.CommandError); !ok {
			return err
		}
	}

	secretIndexes := []mongo.IndexModel{
		{
			Keys: bson.M{
				"Key": 1,
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{
				"Value": 1,
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{
				"Key":           1,
				"ApplicationId": 1,
			},
		},
	}

	_, err = database.Collection(SecretsMongoCollection).Indexes().CreateMany(ctx, secretIndexes)

	if err != nil {
		return err
	}

	return nil
}
