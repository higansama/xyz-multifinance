package mongodb

type MongoDoc interface {
	Decode(val any) error
	Err() error
}
