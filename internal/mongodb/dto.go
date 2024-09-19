package mongodb

type AggregateCountRes struct {
	Count int64 `bson:"count_docs"`
}
