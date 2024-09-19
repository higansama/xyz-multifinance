package mongodb

//func Paginate(ctx context.Context, listQuery *utils.ListQuery, collection *mongo.Collection, filter any) (*utils.ListResult, error) {
//	if filter == nil {
//		filter = bson.D{}
//	}
//
//	count, err := collection.CountDocuments(ctx, filter)
//	if err != nil {
//		return nil, err
//	}
//
//	limit := int64(listQuery.GetLimit())
//	skip := int64(listQuery.GetOffset())
//
//	cursor, err := collection.Find(ctx, filter, &options.FindOptions{
//		Limit: &limit,
//		Skip:  &skip,
//	})
//	if err != nil {
//		return nil, err
//	}
//	defer cursor.Close(ctx) // nolint: errcheck
//
//	data := make([]T, 0, listQuery.GetSize())
//
//	for cursor.Next(ctx) {
//		var prod T
//		if err := cursor.Decode(&prod); err != nil {
//			return nil, err
//		}
//		data = append(data, prod)
//	}
//
//	if err := cursor.Err(); err != nil {
//		return nil, err
//	}
//
//	return utils.NewListResult(data, listQuery.GetSize(), listQuery.GetPage(), count), nil
//}
