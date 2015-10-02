package couchbase

// To test this couchbase integration, you need to set up a data,query and index service
// in couchbase. Then add a index on your selected bucket from the cbq CLI
// Definition: CREATE PRIMARY INDEX primary ON activity
// Definition: CREATE INDEX subjectsid ON activity(subject_id)
//

import (
	"fmt"
	couchbase "gopkg.in/couchbase/gocb.v1"
	"gopkg.in/manishrjain/gocrud.v1/store"
	"gopkg.in/manishrjain/gocrud.v1/x"
)

var log = x.Log("couchbase")

type Container struct {
	Data struct {
		Id string `json:"id"`
	}
}

// CouchbaseDB store backed by Couchbase
type CouchbaseDB struct {
	bucket     *couchbase.Bucket
	bucketname string
}

// Init setup a new collection using the name provided
func (cb *CouchbaseDB) Init(args ...string) {
	if len(args) != 2 {
		log.WithField("args", args).Fatal("Invalid arguments")
		return
	}

	ipaddr := args[0]
	bucket := args[1]

	cluster, err := couchbase.Connect("couchbase://" + ipaddr)
	if err != nil {
		x.LogErr(log, err).Fatal("Error connecting:  %v", err)
	}

	cb.bucket, err = cluster.OpenBucket(bucket, "")
	if err != nil {
		x.LogErr(log, err).Fatal("Error open bucket:  %v", err)
	}
	cb.bucketname = bucket
	fmt.Print("Couchbase registered")

}

// Commit inserts the instructions into the collection as documents
func (cb *CouchbaseDB) Commit(its []*x.Instruction) error {
	c := cb.bucket
	for _, it := range its {
		var key string
		key = fmt.Sprintf("%s_%s", it.SubjectId, x.UniqueString(5))
		cas, err := c.Insert(key, &it, 0)
		if cas == 0 {
			x.LogErr(log, err).Error("While executing batch")
			return nil
		}

	}

	log.WithField("inserted", len(its)).Debug("Stored instructions")

	return nil
}

// IsNew checks if the supplied subject identifier exists in the collection
func (cb *CouchbaseDB) IsNew(subject string) bool {
	c := cb.bucket

	query := couchbase.NewN1qlQuery("SELECT subject_id FROM " + cb.bucketname + " WHERE subject_id='" + subject + "'")
	var row x.Instruction
	rows, err := c.ExecuteN1qlQuery(query, nil)

	if rows.Next(&row) != true {
		fmt.Print(err)
		//x.LogErr(log, err).Error("While running query")
		return true
	} else {
		return false
	}

	return true
}

// GetEntity retrieves all documents matching the subject identifier
func (cb *CouchbaseDB) GetEntity(subject string) (result []x.Instruction, err error) {
	c := cb.bucket

	query := couchbase.NewN1qlQuery("SELECT META() as data FROM " + cb.bucketname + " WHERE subject_id='" + subject + "'")

	rows, err := c.ExecuteN1qlQuery(query, nil)
	if err != nil {
		fmt.Print(err)
	}

	var rowx x.Instruction
	var row Container
	for rows.Next(&row) {
		fmt.Print(row.Data.Id)
		c.Get(row.Data.Id, &rowx)
		result = append(result, rowx)
	}
	rows.Close()

	return result, err
}

func (c *CouchbaseDB) Iterate(fromId string, num int, ch chan x.Entity) (found int, last x.Entity, err error) {
	log.Fatal("Not implemented")
	return
}

func init() {
	log.Info("Registering Couchbase")
	store.Register("couchbase", new(CouchbaseDB))
}
