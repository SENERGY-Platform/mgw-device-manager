package storage_hdl

import (
	"context"
	"database/sql"
	sql_db_hdl "github.com/SENERGY-Platform/go-service-base/sql-db-hdl"
	lib_model "github.com/SENERGY-Platform/mgw-device-manager/lib/model"
	"github.com/SENERGY-Platform/mgw-device-manager/util/db"
	"reflect"
	"testing"
	"time"
)

func initDB(t *testing.T) (*sql.DB, error) {
	testDB, err := db.New(t.TempDir())
	if err != nil {
		return nil, err
	}
	if err = sql_db_hdl.InitDB(context.Background(), testDB, "../../include/storage_schema.sql", time.Second*5, time.Second); err != nil {
		return nil, err
	}
	return testDB, nil
}

func TestHandler(t *testing.T) {
	testDB, err := initDB(t)
	if err != nil {
		t.Fatal(err)
	}
	h := New(testDB)
	t.Run("list devices db empty", func(t *testing.T) {
		devices, err := h.ReadAll(context.Background(), lib_model.DevicesFilter{})
		if err != nil {
			t.Error(err)
		}
		if len(devices) != 0 {
			t.Error("expected empty map")
		}
	})
	id := "1"
	a := lib_model.Device{
		DeviceBase: lib_model.DeviceBase{
			DeviceData: lib_model.DeviceData{
				ID:    id,
				Ref:   "test",
				Name:  "test",
				State: "test",
				Type:  "test",
				Attributes: []lib_model.DeviceAttribute{
					{
						Key:   "test",
						Value: "test",
					},
				},
			},
			Created: time.Now().Round(0),
			Updated: time.Now().Round(0),
		},
		UserData: lib_model.DeviceUserData{
			DeviceUserDataBase: lib_model.DeviceUserDataBase{
				Name: "test",
				Attributes: []lib_model.DeviceAttribute{
					{
						Key:   "test",
						Value: "test",
					},
				},
			},
			Updated: time.Now().Round(0),
		},
	}
	t.Run("create device", func(t *testing.T) {
		err = h.Create(context.Background(), nil, a.DeviceBase)
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("create device exists", func(t *testing.T) {
		err = h.Create(context.Background(), nil, a.DeviceBase)
		if err == nil {
			t.Error("expected error")
		}
	})
	t.Run("read device", func(t *testing.T) {
		b, err := h.Read(context.Background(), id)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(a.DeviceBase, b.DeviceBase) {
			t.Error("expected\n", a.DeviceBase, "got\n", b.DeviceBase)
		}
	})
	t.Run("read device does not exist", func(t *testing.T) {
		_, err = h.Read(context.Background(), "2")
		if err == nil {
			t.Error("expected error")
		}
	})
	t.Run("update device", func(t *testing.T) {
		a.Ref = "test2"
		a.Name = "test2"
		a.State = "test2"
		a.Type = "test2"
		a.Created = time.Now().Round(0)
		a.Updated = time.Now().Round(0)
		a.Attributes = append(a.Attributes, lib_model.DeviceAttribute{
			Key:   "test2",
			Value: "test2",
		})
		err = h.Update(context.Background(), nil, a.DeviceBase)
		if err != nil {
			t.Error(err)
		}
		b, err := h.Read(context.Background(), id)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(a.DeviceBase, b.DeviceBase) {
			t.Error("expected\n", a.DeviceBase, "got\n", b.DeviceBase)
		}
	})
	t.Run("update device does not exist", func(t *testing.T) {
		err = h.Update(context.Background(), nil, lib_model.DeviceBase{DeviceData: lib_model.DeviceData{ID: "2"}})
		if err == nil {
			t.Error("expected error")
		}
	})
	t.Run("update user data", func(t *testing.T) {
		a.UserData.Name = "test2"
		a.UserData.Updated = time.Now().Round(0)
		a.UserData.Attributes = []lib_model.DeviceAttribute{
			{
				Key:   "test2",
				Value: "test2",
			},
		}
		err = h.UpdateUserData(context.Background(), nil, id, a.UserData)
		if err != nil {
			t.Error(err)
		}
		b, err := h.Read(context.Background(), id)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(a, b) {
			t.Error("expected\n", a, "got\n", b)
		}
	})
	t.Run("update user data does not exist", func(t *testing.T) {
		err = h.UpdateUserData(context.Background(), nil, "2", lib_model.DeviceUserData{})
		if err == nil {
			t.Error("expected error")
		}
	})
	t.Run("update states", func(t *testing.T) {
		ts := time.Now().Round(0)
		err = h.UpdateStates(context.Background(), nil, "test2", lib_model.Offline, ts)
		if err != nil {
			t.Error(err)
		}
		b, err := h.Read(context.Background(), id)
		if err != nil {
			t.Error(err)
		}
		if b.State != lib_model.Offline {
			t.Error("expected", lib_model.Offline, "got", b.State)
		}
		if b.Updated != ts {
			t.Error("expected", ts, "got", b.Updated)
		}
	})
	t.Run("delete device", func(t *testing.T) {
		err = h.Delete(context.Background(), nil, id)
		if err != nil {
			t.Error(err)
		}
		_, err = h.Read(context.Background(), id)
		if err == nil {
			t.Error("expected error")
		}
	})
	t.Run("delete device does not exist", func(t *testing.T) {
		err = h.Delete(context.Background(), nil, id)
		if err == nil {
			t.Error("expected error")
		}
	})
	t.Run("attributes table empty", func(t *testing.T) {
		rows, err := testDB.Query("SELECT is_usr, key_name, value FROM device_attributes;")
		if err != nil {
			t.Fatal(err)
		}
		defer rows.Close()
		c := 0
		for rows.Next() {
			c++
		}
		if c != 0 {
			t.Error("expected empty table")
		}
		if err = rows.Err(); err != nil {
			t.Fatal(err)
		}
	})
}
