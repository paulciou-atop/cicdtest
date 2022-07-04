package inventory

import (
	"encoding/json"
	"nms/lib/pgutils"
	"reflect"
	"testing"
	"time"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/go-pg/pg/v10"
	"github.com/imdario/mergo"
	lop "github.com/samber/lo/parallel"
)

var inv = Inventory{
	Name:       "test",
	DeviceType: "testtype",
	Id:         "test123",
	Location: Location{
		Path: "location-testing",
	},
	IpAddress:  "1.1.1.1",
	MacAddress: "11-11-11-11-11-11",
	HostName:   "dev1",

	CreatedAt: time.Now().Format(TIME_LAYOUT),
	LastSeen:  time.Now().Format(TIME_LAYOUT),

	SupportProtocols: []string{"snmp", "gwd"},
}

func TestDatabase(t *testing.T) {
	dbClient, err := pgutils.NewClient()
	if err != nil {
		t.Error(err)
	}
	defer dbClient.Close()

	pgdb, err := dbClient.GetDB()
	if err != nil {
		t.Error(err)
	}
	_, err = pgdb.Exec(`DROP TABLE IF EXISTS inventories`)
	if err != nil {
		t.Error(err)
	}

	err = dbClient.CreateTable(&Inventory{})
	if err != nil {
		t.Error(err)
	}

	dbClient.Insert(&inv)
	var getInv Inventory
	err = dbClient.Query(&getInv, pgutils.QueryExpr{
		Expr:  "id = ?",
		Value: "test123",
	})
	if err != nil {
		t.Error(err)
	}
	v := reflect.ValueOf(inv)
	getV := reflect.ValueOf(getInv)

	for i := 0; i < v.NumField(); i++ {
		k := reflect.TypeOf(Inventory{}).Field(i).Name
		a := v.Field(i).Interface()
		b := getV.Field(i).Interface()
		if !reflect.TypeOf(a).Comparable() {
			continue
		}
		t.Logf("[%s] a=%v b=%v", k, a, b)
		if a != b {
			t.Errorf("Should equal %v = %v", a, b)
		}
	}

	// testing where not in
	inv2 := inv
	inv2.Name = "test2"
	inv2.Id = "2"
	inv3 := inv
	inv3.Name = "test3"
	inv3.Id = "3"
	dbClient.Insert(&inv2)
	dbClient.Insert(&inv3)
	names := []string{"test", "test3"}
	var invs []Inventory
	pgdb.Model(&invs).Where("name not in (?)", pg.In(names)).
		Select()
	if len(invs) != 1 {
		t.Error("should has only one result")
	}
	if invs[0].Name != "test2" {
		t.Error("wrong result")
	}

	// Testing update new item
	inv4 := inv
	inv4.Id = "inv4"
	err = dbClient.Update(&inv4)
	if err != nil {
		t.Error(err)
	}
}

func TestSHA1(t *testing.T) {
	s, err := sha1CheckSum(inv)
	if err != nil {
		t.Error(err)
	}
	inv2 := inv

	s2, err := sha1CheckSum(inv2)
	if err != nil {
		t.Error(err)
	}
	if s != s2 {
		t.Errorf("should be equal %s %s\n", s, s2)
	}
	inv2.Name = "222"
	s3, err := sha1CheckSum(inv2)
	if err != nil {
		t.Error(err)
	}
	if s3 == s2 {
		t.Errorf("shouldn't be equal %s %s\n", s, s2)
	}
}

func TestPatch(t *testing.T) {
	inv2 := inv

	inv2.Name = "test2"
	inv2.Location.Path = "/test/v2"
	inv3 := inv2
	inv3.Name = "test3"
	inv3.FirmwareInformation.Kernel = "3.1.1"

	invJson, err := json.MarshalIndent(&inv, "", " ")
	if err != nil {
		t.Error(err)
	}
	inv2Json, err := json.MarshalIndent(&inv2, "", " ")
	if err != nil {
		t.Error(err)
	}

	inv3Json, err := json.MarshalIndent(&inv2, "", " ")
	if err != nil {
		t.Error(err)
	}
	p2 := makePatch(string(invJson), string(inv2Json))
	p3 := makePatch(string(inv2Json), string(inv3Json))

	inv2FromPatch, err := patch(string(invJson), p2)
	if err != nil {
		t.Error(err)
	}
	if inv2FromPatch != string(inv2Json) {
		t.Error("should be same")
	}

	inv3FromPatch, err := patches(string(invJson), p2, p3)
	if err != nil {
		t.Error(err)
	}
	if inv3FromPatch != string(inv3Json) {
		t.Error("should be same")
	}
}

func TestGitCommit(t *testing.T) {
	storage := memory.NewStorage()

	mfs := osfs.New("./workdir")

	// Init repo with memory filesystem and memory storage
	r, err := git.Init(storage, mfs)
	if err != nil {
		t.Error(err)
	}
	// Try to write inventory json to memory filesystem
	filename := "inventory.json"
	invFile, err := mfs.Create(filename)

	if err != nil {
		t.Error(err)
	}

	// write content to memory fs
	invJson, err := json.MarshalIndent(&inv, "", " ")
	if err != nil {
		t.Error(err)
	}
	invFile.Write(invJson)
	invFile.Close()

	// Get worktree from repo
	wt, err := r.Worktree()
	if err != nil {
		t.Error(err)
	}
	// we can check status from worktree
	status, err := wt.Status()
	if err != nil {
		t.Error(err)
	}

	// status of inventory.json should be '?' untracked either worktree or statging area
	if status[filename].Worktree != '?' || status[filename].Staging != '?' {
		t.Errorf("%s should been untracked but status is %c", filename, status[filename].Worktree)
	}
	_, err = wt.Add(filename)
	if err != nil {
		t.Error(err)
	}
	// worktree commit
	commit, err := wt.Commit("add inventory.json", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "NMS",
			Email: "none",
			When:  time.Now(),
		},
	})
	if err != nil {
		t.Error(err)
	}
	obj, err := r.CommitObject(commit)
	if err != nil {
		t.Error(err)
	}
	t.Log(obj)

	//modify
	inv2 := inv
	inv2.Name = "inv2"
	inv2Json, err := json.MarshalIndent(&inv2, "", " ")
	if err != nil {
		t.Error(err)
	}
	inv2f, err := mfs.Create(filename)
	defer inv2f.Close()
	if err != nil {
		t.Error(err)
	}
	_, err = inv2f.Write(inv2Json)
	if err != nil {
		t.Error(err)
	}

	// we can check status from worktree
	status, err = wt.Status()
	if err != nil {
		t.Error(err)
	}

	// status of inventory.json should be 'M'
	if status[filename].Worktree != 'M' {
		t.Errorf("%s should been modified but status is %c", filename, status[filename].Worktree)
	}

	commit, err = wt.Commit("update inventory.json", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "NMS",
			Email: "none",
			When:  time.Now(),
		},
	})
	if err != nil {
		t.Error(err)
	}
	obj, err = r.CommitObject(commit)
	if err != nil {
		t.Error(err)
	}
	t.Log(obj)

	//Testing open existing storage
	memf2 := memfs.New()
	r2, err := git.Open(storage, memf2)
	if err != nil {
		t.Error(err)
	}
	wk2, err := r2.Worktree()
	if err != nil {
		t.Error(err)
	}
	ref, err := r2.Head()
	if err != nil {
		t.Error(err)
	}
	wk2.Reset(&git.ResetOptions{
		Commit: ref.Hash(),
		Mode:   git.HardReset,
	})
	tree2file, err := memf2.Open(filename)
	if err != nil {
		t.Error(err)
	}
	var readInv2 = make([]byte, 1024)
	nByte, err := tree2file.Read(readInv2)
	t.Logf("%d bytes read \n\n%s\n", nByte, string(readInv2))
}

type Alan struct {
	Name string
	Age  int
}

func TestMap(t *testing.T) {
	alans := []Alan{
		{Name: "troll1", Age: 1},
		{Name: "troll2", Age: 1},
		{Name: "troll3", Age: 1},
		{Name: "troll4", Age: 1},
	}
	as := lop.Map(alans, func(a Alan, i int) Alan {
		t.Log(i)
		a.Age = 4
		return a
	})
	for _, a := range as {
		if a.Age != 4 {
			t.Errorf("Age should be 4 but %d\n", a.Age)
		}
	}
}

type Astruct struct {
	A string
	B string
	C string
	D string
}

func TestMergeStruct(t *testing.T) {
	dest := Astruct{C: "C", D: "D"}
	src := Astruct{B: "B", C: "SC"}
	err := mergo.MergeWithOverwrite(&dest, src)
	if err != nil {
		t.Error(err)
	}
	t.Log(dest)
}
