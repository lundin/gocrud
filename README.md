# gocrud
Go framework to simplify creating, reading, updating, and deleting arbitrary depth structured data — to make building REST services fast and easy.

 * [Gocrud Presentation](https://go-talks.appspot.com/github.com/manishrjain/gocrud/presentation/gomeetup.slide)
 * [Link to blog posts](https://mrjn.xyz/categories/gocrud/)
 * [Example Usage](example.md)

[![wercker status](https://app.wercker.com/status/08406d64ea74cab2dd8155944e56b87d/m "wercker status")](https://app.wercker.com/project/bykey/08406d64ea74cab2dd8155944e56b87d)
[![GoDoc](https://godoc.org/github.com/manishrjain/gocrud?status.svg)](https://godoc.org/github.com/manishrjain/gocrud)

## Releases

Note that the master branch always refers to latest version of Gocrud, which would contain breaking changes.
**To use stable version of Gocrud APIs, please use the packages released via gopkg.in.**

Gocrud version | Install instructions | Godoc | Source
:---: | --- | --- | ---
**v1** (stable) | `go get -v gopkg.in/manishrjain/gocrud.v1/...` |  [godoc](http://godoc.org/gopkg.in/manishrjain/gocrud.v1) | [source](https://github.com/manishrjain/gocrud/tree/v1)
**master** (dev) | `go get -v github.com/manishrjain/gocrud/...` | [godoc](https://godoc.org/github.com/manishrjain/gocrud) | [source](https://github.com/manishrjain/gocrud)


## Questions / Support
I hang out at the gophers.slack.com, at `#gocrud` channel.
You can get an invitation to join Slack via this link:
[Gopher Slack Signup](http://bit.ly/go-slack-signup).
You can also direct message me, my user id is `@manishrjain`.

## Why?
![](allthings.jpg)

<sub>Courtesy: Monish, co-founder karma.wiki</sub>

Having built over 3 different startup backends, I think a lot of time is wasted figuring out and coding CRUD for data structures. In addition, the choice of database has to be made up front, which causes a lot of headache for startup founders. Gocrud was written with the aim to make CRUD easy, and provide the flexibility to switch out both the underlying storage and search engines at any stage of development.

#### Data stores
Datastore | Driver Available | Status
--- | :---: | ---
LevelDB | Yes | Ready
MySQL | Yes | Needs to implement `Iterate` func
PostgreSQL | Yes | Needs to implement `Iterate` func
Cassandra | Yes | Ready
MongoDB | Yes | Needs to implement `Iterate` func
Google Datastore | Yes | Needs to implement `Iterate` func
RethinkDB | Yes | Needs to implement `Iterate` func
Amazon DynamoDB | No | Needs work
**[Datastore usage](datastore.md)** shows how to use and initialize various datastores. One can add support for more by implementing this interface:
```go
type Store interface {
  Init(args ...string)
  Commit(its []*x.Instruction) error
  IsNew(subject string) bool
  GetEntity(subject string) ([]x.Instruction, error)
  Iterate(fromId string, num int, ch chan x.Entity) (int, error)
}
```

#### Search engines
Search Engine | Drive Available
--- | :---:
Elastic Search | Yes
Solr | No

Can be added by implementing these interfaces:
```go
type Engine interface {
	Init(args ...string)
	Update(x.Doc) error
	NewQuery(kind string) Query
}

type Query interface {
	Limit(num int) Query
	Order(field string) Query
	Run() ([]x.Doc, error)
  // and few others
}
```

## Framework
This framework is built to follow these principles:

1. **Versioning**: Keep track of all edits to the data, including deletion operations.
1. **Authorship**: Be able to track who edited (/deleted) what.
1. **Retention**: On deletion, only mark it as deleted. Never actually delete any data.

The framework makes it easy to have *Parent-Child* relationships, quite common in today’s CRUD operations. For e.g.
```
- Posts created by User (User -> Post)
- Comments on Posts (Post -> Comment)
- Likes on Posts (Post -> Like)
- Likes on Comments (Comment -> Like)
```
And be able to traverse these relationships and retrieve all of the children, grandchildren etc. For e.g. `(User -> Post -> [(Comment -> Like), Like])`

The framework does this by utilizing Graph operations, but without using a Graph database. This means the framework can be used to quickly build a Go backend to serve arbitrarily complex data, while still using your database of choice. See [example usage](example.md)

## Dependency management
Users who import Gocrud into their packages are responsible to organize
and maintain all of their dependencies to ensure code compatibility and build
reproducibility. Gocrud makes no direct use of dependency management tools like
[Godep](https://github.com/tools/godep).

## Performance considerations
For the [example](example.md), this is what gets stored in the database:
```
mysql> select * from instructions;
+------------+--------------+-----------+--------------------------------------+-----------+---------------------+---------+----+
| subject_id | subject_type | predicate | object                               | object_id | nano_ts             | source  | id |
+------------+--------------+-----------+--------------------------------------+-----------+---------------------+---------+----+
| uid_oNM    | User         | Post      | NULL                                 | wClGp     | 1435408916326573229 | uid_oNM |  1 |
| wClGp      | Post         | body      | "You can search for cat videos here" |           | 1435408916326573229 | uid_oNM |  2 |
| wClGp      | Post         | tags      | ["search","cat","videos"]            |           | 1435408916326573229 | uid_oNM |  3 |
| wClGp      | Post         | url       | "www.google.com"                     |           | 1435408916326573229 | uid_oNM |  4 |
| wClGp      | Post         | Like      | NULL                                 | kStx9     | 1435408916341828408 | uid_qB3 |  5 |
| kStx9      | Like         | thumb     | 1                                    |           | 1435408916341828408 | uid_qB3 |  6 |
| wClGp      | Post         | Comment   | NULL                                 | 8f78r     | 1435408916341828408 | uid_qB3 |  7 |
| 8f78r      | Comment      | body      | "Comment by on the post"             |           | 1435408916341828408 | uid_qB3 |  8 |
| wClGp      | Post         | Like      | NULL                                 | Gyd7G     | 1435408916352622582 | uid_a30 |  9 |
| Gyd7G      | Like         | thumb     | 1                                    |           | 1435408916352622582 | uid_a30 | 10 |
| 8f78r      | Comment      | Like      | NULL                                 | q2IKK     | 1435408916357443075 | uid_I5u | 11 |
| q2IKK      | Like         | thumb     | 1                                    |           | 1435408916357443075 | uid_I5u | 12 |
| 8f78r      | Comment      | Comment   | NULL                                 | g8llL     | 1435408916357443075 | uid_I5u | 13 |
| g8llL      | Comment      | body      | "Comment xv on comment"              |           | 1435408916357443075 | uid_I5u | 14 |
| q2IKK      | Like         | Comment   | NULL                                 | oaztb     | 1435408916368908590 | uid_SPX | 15 |
| oaztb      | Comment      | body      | "Comment kL on Like"                 |           | 1435408916368908590 | uid_SPX | 16 |
| 8f78r      | Comment      | censored  | true                                 |           | 1435408916377065650 | uid_D2g | 17 |
| kStx9      | Like         | _delete_  | true                                 |           | 1435408916384422689 | uid_2a5 | 18 |
+------------+--------------+-----------+--------------------------------------+-----------+---------------------+---------+----+
18 rows in set (0.00 sec)
```

The writes are in constant time, where each (entity,predicate) constitutes one row. As the properties per entity grow, more rows need to be read (1 row = 1 edge/predicate) to get the entity, it's predicates and it's children. This however, shouldn't be much of a concern for any standard data, which has limited number of predicates/properties per entity. Gocrud in addition, retrieves all children in parallel via `goroutines`, instead of retrieving them one by one.

Property value filtering, sorting, full and partial text matching are now being made available via various search engines. Gocrud provides a search interface, which provides the most common search functionality right out of the box. Thus, there's a clear distinction between data store and search right from the beginning.

## Reserved keywords
The following predicates are reserved by the framework, and shouldn't be used by the caller. Currently, this guideline isn't being hardly enforced by the framework.

Predicate | Meaning
--- | ---
`_parent_` | Stores an edge from child -> parent entity.
`_delete_` | Marks a particular entity as deleted.

## Contact
Feel free to [contact me](https://twitter.com/manishrjain) at my Twitter handle **@manishrjain** for any discussions related to this framework. Also, feel free to send pull requests, they're welcome!
