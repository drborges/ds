package ds_test

import (
	"appengine/aetest"
	"appengine/datastore"
	"github.com/drborges/ds"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type Tag struct {
	ds.Model
	Name  string
	Owner string
}

func (this Tag) KeyMetadata() *ds.KeyMetadata {
	return &ds.KeyMetadata{
		Kind:     "Tags",
		StringID: this.Name,
	}
}

type Post struct {
	ds.Model
	Description string
}

func (this Post) KeyMetadata() *ds.KeyMetadata {
	return &ds.KeyMetadata{
		Kind: "Posts",
	}
}

type Account struct {
	ds.Model
	Id int64
	Name string
}

func (this Account) KeyMetadata() *ds.KeyMetadata {
	return &ds.KeyMetadata{
		Kind: "Accounts",
		IntID: this.Id,
	}
}

func TestLoadModel(t *testing.T) {
	c, _ := aetest.NewContext(nil)
	defer c.Close()

	Convey("ds.Datastore", t, func() {
		Convey("Given I have a model with StringID key saved in datastore", func() {
			tag := Tag{Name: "golang", Owner: "Borges"}
			key, _ := ds.NewKey(c, tag.KeyMetadata())
			datastore.Put(c, key, &tag)
			tag.SetKey(key)

			Convey("Then I can load its information along with its key using ds.Datastore.Load", func() {
				loadedTag := Tag{Name: "golang"}
				err := ds.Datastore{c}.Load(&loadedTag)

				So(err, ShouldBeNil)
				So(loadedTag, ShouldResemble, tag)
				So(loadedTag.Key(), ShouldNotBeNil)
			})
		})

		Convey("Given I have a model with IntID key saved in datastore", func() {
			account := Account{Id: 123, Name: "Borges"}
			key, _ := ds.NewKey(c, account.KeyMetadata())
			datastore.Put(c, key, &account)
			account.SetKey(key)

			Convey("Then I can load its information along with its key using ds.Datastore.Load", func() {
				loadedAccount := Account{Id: 123}
				err := ds.Datastore{c}.Load(&loadedAccount)

				So(err, ShouldBeNil)
				So(loadedAccount, ShouldResemble, account)
				So(loadedAccount.Key(), ShouldNotBeNil)
			})
		})

		Convey("Given I have a model with auto generated key saved in datastore", func() {
			post := Post{Description: "This is gonna be awesome!"}
			incompleteKey, _ := ds.NewKey(c, post.KeyMetadata())
			key, _ := datastore.Put(c, incompleteKey, &post)
			post.SetKey(key)

			Convey("Then I can load its information with its key using ds.Datastore.Load", func() {
				loadedPost := Post{}
				loadedPost.SetKey(key)
				err := ds.Datastore{c}.Load(&loadedPost)

				So(err, ShouldBeNil)
				So(loadedPost, ShouldResemble, post)
				So(loadedPost.Key(), ShouldNotBeNil)
			})
		})
	})
}
