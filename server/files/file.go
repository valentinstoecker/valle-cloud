package files

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/nfnt/resize"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/valentinstoecker/valle-cloud/server/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const THUMBNAIL_SIZE = 512

var DATA_PATH = "data/"
var Collection *mongo.Collection

func init() {
	env, is_set := os.LookupEnv("DATA_PATH")
	if is_set {
		DATA_PATH = env
	}
	Collection = db.DB.Collection("files")
	unique := true
	_, err := Collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.M{"hash": 1},
		Options: &options.IndexOptions{
			Unique: &unique,
		},
	})
	if err != nil {
		panic(err)
	}
	os.MkdirAll(DATA_PATH, 0755)
}

type file struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
	Type string             `bson:"type" json:"type"`
	Time primitive.DateTime `bson:"time" json:"time"`
	Hash string             `bson:"hash" json:"hash"`
}

func (f *file) Save(ctx context.Context) error {
	upsert := true
	_, err := Collection.UpdateOne(ctx, bson.M{"_id": f.ID}, bson.M{"$set": f}, &options.UpdateOptions{
		Upsert: &upsert,
	})
	if err != nil {
		return err
	}
	return nil
}

func decExif(done chan *exif.Exif, err_c chan error, r io.Reader) {
	ex, err := exif.Decode(r)
	ioutil.ReadAll(r)
	if err != nil {
		err_c <- err
		return
	}
	done <- ex
}

func multiTee(r io.Reader, w ...io.Writer) io.Reader {
	out_r := r
	for _, w := range w {
		out_r = io.TeeReader(out_r, w)
	}
	return out_r
}

func renameImage(oldName string, newName string) error {
	err := os.Rename(DATA_PATH+oldName, DATA_PATH+newName)
	if err != nil {
		return err
	}
	err = os.Rename(DATA_PATH+oldName+".thumb", DATA_PATH+newName+".thumb")
	if err != nil {
		return err
	}
	return nil
}

func makeThumbnail(name string, img image.Image) error {
	thumb := resize.Thumbnail(THUMBNAIL_SIZE, THUMBNAIL_SIZE, img, resize.NearestNeighbor)
	thumb_f, err := os.Create(DATA_PATH + name + ".thumb")
	if err != nil {
		return err
	}
	defer thumb_f.Close()
	err = jpeg.Encode(thumb_f, thumb, nil)
	if err != nil {
		return err
	}
	return nil
}

func NewImage(ctx context.Context, name string, reader io.Reader) (*file, error) {
	obj_id := primitive.NewObjectID()

	f, err := os.Create(DATA_PATH + obj_id.Hex())
	if err != nil {
		return nil, err
	}
	defer f.Close()

	hash := sha256.New()
	r, w := io.Pipe()
	tr := multiTee(reader, f, hash, w)

	done := make(chan *exif.Exif)
	err_c := make(chan error)
	go decExif(done, err_c, r)

	img, f_type, err := image.Decode(tr)
	if err != nil {
		return nil, err
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	makeThumbnail(obj_id.Hex(), img)

	h := hex.EncodeToString(hash.Sum(nil))
	renameImage(obj_id.Hex(), h)

	var old file
	err = Collection.FindOne(ctx, bson.M{"hash": h}).Decode(&old)
	if err == nil {
		return &old, nil
	}

	select {
	case ex := <-done:
		t, err := ex.DateTime()
		if err != nil {
			return nil, err
		}
		return &file{ID: obj_id, Name: name, Type: f_type, Time: primitive.NewDateTimeFromTime(t), Hash: h}, nil
	case <-err_c:
		return &file{
			ID:   obj_id,
			Name: name,
			Type: f_type,
			Time: primitive.NewDateTimeFromTime(time.Now()),
			Hash: h,
		}, nil
	}
}

func GetImages(ctx context.Context) ([]*file, error) {
	files := make([]*file, 0)
	cur, err := Collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		f := &file{}
		err := cur.Decode(f)
		files = append(files, f)
		if err != nil {
			return nil, err
		}
	}
	return files, nil
}

func GetImageFile(ctx context.Context, hash string) (*file, error) {
	var file file
	err := Collection.FindOne(ctx, bson.M{"hash": hash}).Decode(&file)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (f *file) Thumbnail() (io.Reader, error) {
	fr, err := os.Open(DATA_PATH + f.Hash + ".thumb")
	if err != nil {
		return nil, err
	}
	return fr, nil
}

func (f *file) Image() (io.Reader, error) {
	fr, err := os.Open(DATA_PATH + f.Hash)
	if err != nil {
		return nil, err
	}
	return fr, nil
}
