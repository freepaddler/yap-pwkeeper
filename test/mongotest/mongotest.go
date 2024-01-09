package main

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	StateActive  = "Active"
	StateDeleted = "Deleted"
)

type Entity struct {
	CreatedAt  time.Time `bson:"created_at"`
	ModifiedAt time.Time `bson:"modified_at"`
	State      string    `bson:"state"`
}

type User struct {
	Id       string `bson:"_id,omitempty"`
	Login    string `bson:"login"`
	Password string `bson:"password"`
	Entity   `bson:"inline"`
}

type Note struct {
	Id     string `bson:"_id,omitempty"`
	UserId string `bson:"user_id"`
	Name   string `bson:"name"`
	Note   string `bson:"note"`
	Entity `bson:"inline"`
}

type File struct {
	Id       string  `bson:"_id,omitempty"`
	Name     string  `bson:"name"`
	FileName string  `bson:"filename"`
	FileSize int64   `bson:"filesize"`
	Sha256   string  `bson:"sha256"`
	Data     *[]byte `bson:"data"`
}

var (
	orig       = "/tmp/ui1.dat"
	clientfile = "/tmp/ui1.dat"
	serverfile = "/tmp/ui2.dat"
)

func readchunk() {
	//st, err := os.Stat(orig)
	//if err != nil {
	//	log.Fatal(err)
	//}
	f1, err := os.Open(orig)
	if err != nil {
		log.Fatal(err)
	}
	defer f1.Close()
	h := sha256.New()
	_, err = io.Copy(h, f1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%x", h.Sum(nil))

	return

	//defer f1.Close()
	//fc, err := os.Create(clientfile)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer fc.Close()
	//fs, err := os.Create(serverfile)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer fs.Close()
	//
	//b := make([]byte, 0, st.Size())
	//buf := bytes.NewBuffer(b)
	//hash := sha256.New()
	//mw := io.MultiWriter(fc, hash, buf)
	//
	//for {
	//	tmp := make([]byte, 1<<8)
	//	n, err := f1.Read(tmp)
	//	log.Printf("read %d", n)
	//	if err != nil {
	//		if errors.Is(io.EOF, err) {
	//			log.Println("end of file")
	//			break
	//		}
	//		log.Fatal(err)
	//	}
	//	_, err = mw.Write(tmp)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}
	//_, err = fs.Write(buf.Bytes())
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Println(hex.EncodeToString(hash.Sum(nil)))
}

func main() {
	f1, err := os.Open(orig)
	if err != nil {
		log.Fatal(err)
	}
	defer f1.Close()
	h := sha256.New()
	_, err = io.Copy(h, f1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%x", h.Sum(nil))
	return

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.Background())
	coll := client.Database("test").Collection("files")

	//fname := "/tmp/ui.log"
	//f2name := "/tmp/ui2.log"
	//
	//stat, err := os.Stat(fname)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//b := make([]byte, 0, stat.Size())
	//buf := bytes.NewBuffer(b)
	//hash := sha256.New()
	//mw := io.MultiWriter(buf, hash)
	//
	//f, err := os.Open(fname)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer f.Close()
	////_, err := io.
	//
	//buf := bytes
	//b, err := io.ReadAll(f)
	//st, err := os.Stat(fname)
	//h := sha256.New()
	//if _, err := io.Copy(h, f); err != nil {
	//	log.Fatal(err)
	//}
	//log.Println("sha is " + hex.EncodeToString(h.Sum(nil)) + " " + st.Name())
	//
	//fst := File{
	//	Name:     "file",
	//	FileName: st.Name(),
	//	FileSize: st.Size(),
	//	Sha256:   hex.EncodeToString(h.Sum(nil)),
	//	Data:     &b,
	//}
	//
	//res, err := coll.InsertOne(context.Background(), &fst)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//log.Println("document_id " + res.InsertedID.(primitive.ObjectID).Hex())
	//
	//f2 := File{}
	//file2, err := os.Create(f2name)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer file2.Close()
	//
	//err = coll.FindOne(context.Background(), bson.D{{"_id", res.InsertedID}}).Decode(&f2)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//file2.Write(*f2.Data)
	//
	//log.Println("Done")
	//
	//return

	userNew := User{
		Login:    "login",
		Password: "password",
		Entity: Entity{
			CreatedAt:  time.Now(),
			ModifiedAt: time.Now(),
			State:      StateActive,
		},
	}

	fmt.Printf("User orig: %+v\n", userNew)

	resUserNew, err := coll.InsertOne(context.TODO(), userNew)

	userNew.Id, err = objectId2string(resUserNew.InsertedID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("User added: %+v\n", userNew)

	oid, err := primitive.ObjectIDFromHex(userNew.Id)
	if err != nil {
		panic(err)
	}
	userFilter := bson.D{{"_id", oid}}

	var foundUser User
	err = coll.FindOne(context.Background(), userFilter).Decode(&foundUser)

	if err != nil {
		panic(err)
	}

	fmt.Printf("User found: %+v\n", foundUser)

	noteNew := Note{
		UserId: foundUser.Id,
		Name:   "note1",
		Note:   "this is a note text",
		Entity: Entity{
			CreatedAt:  time.Now(),
			ModifiedAt: time.Now(),
			State:      StateActive,
		},
	}
	noteNew2 := Note{
		UserId: foundUser.Id,
		Name:   "note2",
		Note:   "this is a second note text",
		Entity: Entity{
			CreatedAt:  time.Now(),
			ModifiedAt: time.Now(),
			State:      StateActive,
		},
	}
	fmt.Printf("Note orig: %+v\n", noteNew)

	coll = client.Database("test").Collection("note")
	notesRes, err := coll.InsertMany(context.Background(), []interface{}{noteNew, noteNew2})
	if err != nil {
		panic(err)
	}

	noteEditId, _ := objectId2string(notesRes.InsertedIDs[0])
	fmt.Println("Note edit id ", noteEditId)
	noteNew.Id = noteEditId
	noteNew.Note = "this is corrected note"
	<-time.After(time.Second)
	noteNew.ModifiedAt = time.Now()
	fmt.Printf("Edited Note: %+v\n", noteNew)

	updateId, _ := primitive.ObjectIDFromHex(noteEditId)
	noteReplace := struct {
		Id   primitive.ObjectID `bson:"_id"`
		Note `bson:"inline"`
	}{
		Id:   updateId,
		Note: noteNew,
	}
	updateFilter := bson.D{{"_id", updateId}, {"state", StateActive}}
	resUpdate, err := coll.ReplaceOne(context.Background(), updateFilter, noteReplace)
	if err != nil {
		panic(err)
	}
	fmt.Printf("modified: %d\n", resUpdate.ModifiedCount)

	filterNotes := bson.D{{"user_id", foundUser.Id}}
	cursor, err := coll.Find(context.Background(), filterNotes)
	var notes []Note
	for cursor.Next(context.TODO()) {
		var note Note
		if err := cursor.Decode(&note); err != nil {
			cursor.Close(context.Background())
			panic(err)
		}
		notes = append(notes, note)
	}
	cursor.Close(context.Background())
	fmt.Println("notes found")
	for _, f := range notes {
		fmt.Printf("%+v\n", f)
	}

}

func objectId2string(oid interface{}) (string, error) {
	if id, ok := oid.(primitive.ObjectID); !ok {
		return "", errors.New("failed to get oid from response")
	} else {
		return id.Hex(), nil
	}
}
