package db

import (
	"fmt"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type ProducerLevel struct {
	Pid   string `json:"pid" form:"pid" query:"pid"`
	Level int    `json:"level" form:"level" query:"level"`
}

type ProducerLevelInfo struct {
	Aid string `json:"aid" form:"aid" query:"aid"`
	//ProducerLevels []*ProducerLevel `json:"producerlevels" form:"producerlevels" query:"producerlevels"`
	ProducerLevels []*ProducerLevel `json:"producerlevels" form:"producerlevels" query:"producerlevels"`
}

type AwardInfo struct {
	Aid         string `json:"aid" form:"aid" query:"aid"`
	StartTime   int64  `json:"start_time" form:"start_time" query:"start_time"`
	EndTime     int64  `json:"end_time" form:"end_time" query:"end_time"`
	TotalAmount int64  `json:"total_amount" form:"total_amount" query:"total_amount"`
	CountTime   int64  `json:"count_time" form:"count_time" query:"count_time"`
}

type UserAward struct {
	Aid      string  `json:"aid" form:"aid" query:"aid"`
	Username string  `json:"username" form:"username" query:"username"`
	Pid      string  `json:"pid" form:"pid" query:"pid"`
	Vote     float64 `json:"vote" form:"vote" query:"vote"`
	Award    float64 `json:"award" form:"award" query:"award"`
}

type ProducerAward struct {
	Aid   string  `json:"aid" form:"aid" query:"aid"`
	Pid   string  `json:"pid" form:"pid" query:"pid"`
	Vote  float64 `json:"vote" form:"vote" query:"vote"`
	Award float64 `json:"award" form:"award" query:"award"`
}

func retryWriteMongo(b *mgo.Bulk) {
	var retryTime int
	for {
		if _, err := b.Run(); err != nil {
			log.Println("fail to write data to mongo ", err)
			time.Sleep(time.Second)
			retryTime++
			if retryTime > 10 {
				log.Fatalln("fail to write data to mongo, retry time exceeds")
			}
			continue
		}
		return
	}
}

func SaveUserContributionAward(userAwards []UserAward) error {
	UAC := GetCollection(CollectionUserContributionAward)
	UAB := UAC.Bulk()
	for k, ua := range userAwards {
		UAB.Upsert(bson.M{"aid": ua.Aid, "pid": ua.Pid, "username": ua.Username}, ua)
		if k % 100 == 0 {
			retryWriteMongo(UAB)
			UAB = UAC.Bulk()
			fmt.Println("finished:", k)
		}
	}
	retryWriteMongo(UAB)
	return nil
}

func SaveProducerContributionAward(producerAwards []ProducerAward) error {
	UAC := GetCollection(CollectionProducerContributionAward)
	UAB := UAC.Bulk()
	for _, ua := range producerAwards {
		UAB.Upsert(bson.M{"aid": ua.Aid, "pid": ua.Pid}, ua)
	}
	retryWriteMongo(UAB)
	return nil
}

func SaveUserAward(userAwards []UserAward) error {
	UAC := GetCollection(CollectionUserAward)
	UAB := UAC.Bulk()
	for _, ua := range userAwards {
		UAB.Insert(ua)
	}
	retryWriteMongo(UAB)
	return nil
}

func SaveProducerAward(producerAwards []ProducerAward) error {
	UAC := GetCollection(CollectionProducerAward)
	UAB := UAC.Bulk()
	for _, ua := range producerAwards {
		UAB.Insert(ua)
	}
	retryWriteMongo(UAB)
	return nil
}

func GetAwardInfo(aid string) (ainfo AwardInfo, err error) {
	BPC := GetCollection(CollectionAwardInfo)
	err = BPC.Find(bson.M{"aid": aid}).One(&ainfo)
	return
}

func SaveProducerLevelInfo(aInfo ProducerLevelInfo) error {
	BPC := GetCollection(CollectionProducerLevelInfo)
	_, err := BPC.Upsert(bson.M{"aid": aInfo.Aid}, aInfo)
	return err
}

func GetProducerLevelInfo(aid string) (pInfo ProducerLevelInfo, err error) {
	BPC := GetCollection(CollectionProducerLevelInfo)
	err = BPC.Find(bson.M{"aid": aid}).One(&pInfo)
	return
}

func SaveAwardInfo(aInfo AwardInfo) error {
	BPC := GetCollection(CollectionAwardInfo)
	count, err := BPC.Find(bson.M{"aid": aInfo.Aid}).Count()
	if err != nil {
		return err
	}
	if count != 0 {
		return errors.New("Aid Exist!")
	}
	err = BPC.Insert(aInfo)
	return err
}

func GetVoteAwardList() ([]*AwardInfo, error) {
	BPC := GetCollection(CollectionAwardInfo)
	var awards []*AwardInfo
	err := BPC.Find(bson.M{}).Sort("-_id").All(&awards)
	if err != nil {
		return nil, err
	}
	return awards, err
}

func GetVoteAwardInfo(id string) ([]*AwardInfo, error) {
	BPC := GetCollection(CollectionAwardInfo)
	var awards []*AwardInfo
	err := BPC.Find(bson.M{"aid": id}).All(&awards)
	if err != nil {
		return nil, err
	}
	return awards, err
}

func GetUserContributionAward(id string) ([]*UserAward, error) {
	BPC := GetCollection(CollectionUserContributionAward)
	var awards []*UserAward
	err := BPC.Find(bson.M{"aid": id}).All(&awards)
	if err != nil {
		return nil, err
	}
	return awards, err
}

func GetProducerContributionAward(id string) ([]*ProducerAward, error) {
	BPC := GetCollection(CollectionProducerContributionAward)
	var awards []*ProducerAward
	err := BPC.Find(bson.M{"aid": id}).All(&awards)
	if err != nil {
		return nil, err
	}
	return awards, err
}

func GetUserAward(id string) ([]*UserAward, error) {
	BPC := GetCollection(CollectionUserAward)
	var awards []*UserAward
	err := BPC.Find(bson.M{"aid": id}).All(&awards)
	if err != nil {
		return nil, err
	}
	return awards, err
}

func GetProducerAward(id string) ([]*ProducerAward, error) {
	BPC := GetCollection(CollectionProducerAward)
	var awards []*ProducerAward
	err := BPC.Find(bson.M{"aid": id}).All(&awards)
	if err != nil {
		return nil, err
	}
	return awards, err
}

func GetVoteTxs(endBlock int64) (voteTx []*VoteTx, err error) {
	BPC := GetCollection(CollectionVoteTx)
	fmt.Println("vote dbb:", BPC.Database.Name, endBlock)
	err = BPC.Find(bson.M{"blockNumber": bson.M{"$lte": endBlock}}).Sort("blockNumber").All(&voteTx)
	if err != nil {
		return nil, err
	}
	return
}
