package mongodb

import (
	"context"
	"reflect"
	"strings"

	ierrors "github.com/higansama/xyz-multi-finance/internal/errors"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func ConvertMongoError(err error, entity string) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		return ierrors.NewEntityNotFoundError(entity)
	}

	return err
}

func MapStringToObjectId(ids []string) []any {
	var pids []any
	for _, v := range ids {
		id, err := primitive.ObjectIDFromHex(v)
		if err == nil {
			pids = append(pids, id)
		} else {
			pids = append(pids, v)
		}
	}
	return pids
}

func StringToObjectId(id string) primitive.ObjectID {
	oid, _ := primitive.ObjectIDFromHex(id)
	return oid
}

func OidToHexOrEmpty(id primitive.ObjectID) string {
	if id.IsZero() {
		return ""
	}
	return id.Hex()
}

// ConvertStructToBsonM convert the given struct to bson.M
//
//	toObjIdKeys: keys that needs to transform to primitive.ObjectId
func ConvertStructToBsonM(obj any, toObjIdKeys []string) (bson.M, error) {
	if reflect.ValueOf(obj).Kind() != reflect.Struct &&
		!(reflect.ValueOf(obj).Kind() == reflect.Ptr && reflect.ValueOf(obj).Elem().Kind() == reflect.Struct) {
		return bson.M{}, errors.New("obj should be a struct")
	}

	data, err := bson.Marshal(obj)
	if err != nil {
		return bson.M{}, errors.WithStack(err)
	}

	var res bson.M
	err = bson.Unmarshal(data, &res)
	if err != nil {
		return bson.M{}, errors.WithStack(err)
	}

	toObjIdKeys = append(toObjIdKeys, "_id") // always convert _id to ObjectId
	for _, item := range toObjIdKeys {
		if vl, ok := res[item]; ok {
			if v, ok := vl.(string); ok && v != "" {
				res[item] = StringToObjectId(v)
			}
		} else if strings.Contains(item, ".") {
			itemPs := strings.Split(item, ".")
			if len(itemPs) > 1 {
				prVal := res[itemPs[0]]
				if prVal == nil {
					continue
				}
				prKnd := reflect.TypeOf(prVal).Kind()
				switch prKnd {
				case reflect.Slice:
					prValS := prVal.(bson.A)
					if len(prValS) > 0 {
						if _, ok := prValS[0].(bson.M); ok {
							nPrValS := make([]bson.M, 0)
							for _, itemJ := range prValS {
								itemJM := itemJ.(bson.M)
								if v, ok := itemJM[itemPs[1]].(string); ok && v != "" {
									itemJM[itemPs[1]] = StringToObjectId(v)
								}
								nPrValS = append(nPrValS, itemJM)
							}
							prVal = nPrValS
						}
					}
				case reflect.Map:
					prValM := prVal.(bson.M)
					if v, ok := prValM[itemPs[1]].(string); ok && v != "" {
						prValM[itemPs[1]] = StringToObjectId(v)
					}
					prVal = prValM
				default:
				}
			}
		}
	}

	return res, nil
}

func AppendFilterToArrayQuery(queries []bson.M, filter ...bson.M) []bson.M {
	for _, f := range filter {
		if len(f) > 0 {
			queries = append(queries, f)
		}
	}

	return queries
}

func AppendFilterToAndQuery(query bson.M, filter ...bson.M) bson.M {
	if _, ok := query["$and"]; !ok {
		query["$and"] = make([]bson.M, 0)
	}

	query["$and"] = append(query["$and"].([]bson.M), filter...)
	return query
}

func AppendFilterToOrQuery(query bson.M, filter ...bson.M) bson.M {
	if query == nil {
		query = make(bson.M)
	}
	if _, ok := query["$or"]; !ok {
		query["$or"] = make([]bson.M, 0)
	}

	query["$or"] = append(query["$or"].([]bson.M), filter...)
	return query
}

func NormalizeMatchQuery(query bson.M) bson.M {
	query = removeOperatorIfEmpty(query, "$or")
	query = removeOperatorIfEmpty(query, "$and")
	query = removeOperatorIfEmpty(query, "$nor")
	return query
}

func removeOperatorIfEmpty(query bson.M, key string) bson.M {
	if v, ok := query[key]; ok {
		if fv, ok := v.([]bson.M); ok {
			if len(fv) == 0 {
				delete(query, key)
			}
		}
	}
	return query
}

func AddStagesToPipeline(pipeline mongo.Pipeline, stages ...bson.D) mongo.Pipeline {
	for _, s := range stages {
		pipeline = append(pipeline, s)
	}
	return pipeline
}

func AddLimitOffsetToPipeline(pipeline mongo.Pipeline, offset int64, limit int64) mongo.Pipeline {
	if offset > 0 {
		pipeline = append(pipeline, bson.D{{"$skip", offset}})
	}
	if limit > 0 {
		pipeline = append(pipeline, bson.D{{"$limit", limit}})
	}
	return pipeline
}

func GetAggregateCount(ctx context.Context, col *mongo.Collection, pipeline mongo.Pipeline) (int64, error) {
	hasMatch := false
	for _, x := range pipeline {
		if v, ok := x.Map()["$match"]; ok {
			v := reflect.ValueOf(v)
			switch v.Kind() {
			case reflect.Map, reflect.Slice:
				hasMatch = v.Len() > 0
			}
		}

		if hasMatch {
			break
		}
	}

	if hasMatch {
		cursor, err := col.Aggregate(ctx, AddStagesToPipeline(
			pipeline,
			bson.D{{"$count", "count_docs"}}))
		if err != nil {
			return 0, errors.WithStack(err)
		}

		var countRes []AggregateCountRes
		err = cursor.All(ctx, &countRes)
		if err != nil {
			return 0, errors.WithStack(err)
		}

		if len(countRes) > 0 {
			return countRes[0].Count, nil
		}
		return 0, nil
	}

	count, err := col.CountDocuments(ctx, bson.D{})
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}

func ConvertValueToOtherType(obj any, target any) error {
	data, err := bson.Marshal(obj)
	if err != nil {
		return errors.WithStack(err)
	}

	err = bson.Unmarshal(data, target)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
