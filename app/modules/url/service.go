package url

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func redirect(dto dto.RedirectDto) string {
	var reportUrlModel models.ReportUrl

	queryOptions := options.FindOneOptions{}
	queryOptions.SetSort(bson.D{{Key: "createdAt", Value: -1}})
	reportUrlCollection.FindOne(ctx, bson.M{}, &queryOptions).Decode(&reportUrlModel)

	if dto.Gid != "" && dto.Range != "" {
		return reportUrlModel.Url + "/edit#gid=" + dto.Gid + "&range=" + dto.Range
	} else if dto.Gid != "" {
		return reportUrlModel.Url + "/edit#gid=" + dto.Gid
	}

	return reportUrlModel.Url
}
