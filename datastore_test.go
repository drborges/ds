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
	Id   int64
	Name string
}

func (this Account) KeyMetadata() *ds.KeyMetadata {
	return &ds.KeyMetadata{
		Kind:  "Accounts",
		IntID: this.Id,
	}
}

type ModelMissingKind struct {
	ds.Model
}

func (this ModelMissingKind) KeyMetadata() *ds.KeyMetadata {
	return &ds.KeyMetadata{}
}

func TestLoadModel(t *testing.T) {
	c, _ := aetest.NewContext(nil)
	defer c.Close()

	Convey("ds.Datastore", t, func() {
		Convey("Load", func() {
			Convey("Given I have a model with StringID key saved in datastore", func() {
				tag := &Tag{Name: "golang", Owner: "Borges"}
				key, _ := ds.NewKey(c, tag)
				datastore.Put(c, key, tag)
				tag.SetKey(key)

				Convey("When I load it using ds.Datastore", func() {
					loadedTag := Tag{Name: "golang"}
					err := ds.Datastore{c}.Load(&loadedTag)

					Convey("Then it succeeds", func() {
						So(err, ShouldBeNil)
					})

					Convey("Then the data is loaded from datastore", func() {
						So(&loadedTag, ShouldResemble, tag)
					})

					Convey("Then the model has its key resolved", func() {
						So(loadedTag.Key(), ShouldNotBeNil)
					})
				})
			})

			Convey("Given I have a model with IntID key saved in datastore", func() {
				account := &Account{Id: 123, Name: "Borges"}
				key, _ := ds.NewKey(c, account)
				datastore.Put(c, key, account)
				account.SetKey(key)

				Convey("When I load it using ds.Datastore", func() {
					loadedAccount := Account{Id: 123}
					err := ds.Datastore{c}.Load(&loadedAccount)

					Convey("Then it succeeds", func() {
						So(err, ShouldBeNil)
					})

					Convey("Then the data is loaded from datastore", func() {
						So(&loadedAccount, ShouldResemble, account)
					})

					Convey("Then the model has its key resolved", func() {
						So(loadedAccount.Key(), ShouldNotBeNil)
					})
				})
			})

			Convey("Given I have a model with auto generated key saved in datastore", func() {
				post := &Post{Description: "This is gonna be awesome!"}
				incompleteKey, _ := ds.NewKey(c, post)
				key, _ := datastore.Put(c, incompleteKey, post)
				post.SetKey(key)

				Convey("When I load it with a key using ds.Datastore", func() {
					loadedPost := Post{}
					loadedPost.SetKey(key)
					err := ds.Datastore{c}.Load(&loadedPost)

					Convey("Then it succeeds", func() {
						So(err, ShouldBeNil)
					})

					Convey("Then the data is loaded from datastore", func() {
						So(&loadedPost, ShouldResemble, post)
					})

					Convey("Then the model has its key resolved", func() {
						So(loadedPost.Key(), ShouldNotBeNil)
					})
				})
			})

			Convey("Given I have a model missing an auto generated key", func() {
				loadedPost := Post{}

				Convey("When I try to load it from datastore", func() {
					err := ds.Datastore{c}.Load(&loadedPost)

					Convey("Then ErrMissingAutoGeneratedKey is returned", func() {
						So(err, ShouldEqual, ds.ErrUnresolvableKey)
						So(loadedPost.Key(), ShouldBeNil)
					})
				})
			})
		})

		Convey("Update", func() {
			Convey("Given I have a model with StringID key", func() {
				tag := Tag{Name: "golang", Owner: "Borges"}

				Convey("When I use ds.Datastore.Put to save it into datastore", func() {
					err := ds.Datastore{c}.Update(&tag)

					Convey("Then it succeeds", func() {
						So(err, ShouldBeNil)
					})

					Convey("Then the key is set back to the model", func() {
						So(tag.Key().String(), ShouldEqual, "/Tags,golang")
					})

					Convey("Then I am able to lookup the information from datastore", func() {
						loadedTag := Tag{}
						loadedTag.SetKey(tag.Key())
						datastore.Get(c, loadedTag.Key(), &loadedTag)

						So(loadedTag, ShouldResemble, tag)
					})
				})
			})

			Convey("Given I have a model with invalid metadata", func() {
				invalidModel := ModelMissingKind{}

				Convey("When I update it", func() {
					err := ds.Datastore{c}.Update(&invalidModel)

					Convey("Then it fails with an error", func() {
						So(err, ShouldNotBeNil)
					})

					Convey("Then the key is not set back to the model", func() {
						So(invalidModel.Key(), ShouldBeNil)
					})
				})
			})
		})

		Convey("Create", func() {
			Convey("Given I have a model with auto generated key", func() {
				post := Post{Description: "An awesome post about ds package"}

				Convey("When I create it", func() {
					err := ds.Datastore{c}.Create(&post)

					Convey("Then it succeeds", func() {
						So(err, ShouldBeNil)
					})

					Convey("Then the key is set back to the model", func() {
						So(post.Key().String(), ShouldNotBeNil)
						So(post.Key().String(), ShouldNotEqual, "/Posts,0")
					})

					Convey("Then I am able to lookup the information from datastore", func() {
						loadedPost := Post{}
						loadedPost.SetKey(post.Key())
						datastore.Get(c, loadedPost.Key(), &loadedPost)

						So(loadedPost, ShouldResemble, post)
					})
				})
			})

			Convey("Given I have a model with invalid metadata", func() {
				invalidModel := ModelMissingKind{}

				Convey("When I create it", func() {
					err := ds.Datastore{c}.Create(&invalidModel)

					Convey("Then it fails with an error", func() {
						So(err, ShouldNotBeNil)
					})

					Convey("Then the key is not set back to the model", func() {
						So(invalidModel.Key(), ShouldBeNil)
					})
				})
			})
		})

		Convey("CreateAll", func() {
			Convey("Given I have a few models", func() {
				post1 := &Post{Description: "Post 1"}
				post2 := &Post{Description: "Post 2"}
				post3 := &Post{Description: "Post 3"}
				posts := []*Post{post1, post2, post3}

				Convey("When I create all", func() {
					err := ds.Datastore{c}.CreateAll(posts)

					Convey("Then it succeeds", func() {
						So(err, ShouldBeNil)
					})

					Convey("Then the key is set back to the models", func() {
						So(post1.Key().String(), ShouldNotBeNil)
						So(post1.Key().String(), ShouldNotEqual, "/Posts,0")

						So(post2.Key().String(), ShouldNotBeNil)
						So(post2.Key().String(), ShouldNotEqual, "/Posts,0")

						So(post3.Key().String(), ShouldNotBeNil)
						So(post3.Key().String(), ShouldNotEqual, "/Posts,0")
					})
				})
			})
		})

		Convey("Delete", func() {
			Convey("Given I have a model saved in datastore", func() {
				tag := Tag{Name: "golang", Owner: "Borges"}

				Convey("When I delete it", func() {
					err := ds.Datastore{c}.Delete(&tag)

					Convey("Then it succeeds", func() {
						So(err, ShouldBeNil)
					})

					Convey("Then the model does not exist in datastore", func() {
						err = datastore.Get(c, tag.Key(), nil)
						So(err, ShouldNotBeNil)
					})
				})
			})
		})
	})
}
