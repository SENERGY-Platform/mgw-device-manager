package storage_hdl

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	lib_model "github.com/SENERGY-Platform/mgw-device-manager/lib/model"
	"strings"
	"time"
)

var tLayout = time.RFC3339Nano

type Handler struct {
	db *sql.DB
}

func New(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) BeginTransaction(ctx context.Context) (driver.Tx, error) {
	tx, e := h.db.BeginTx(ctx, nil)
	if e != nil {
		return nil, lib_model.NewInternalError(e)
	}
	return tx, nil
}

func (h *Handler) ReadAll(ctx context.Context, filter lib_model.DevicesFilter) (map[string]lib_model.Device, error) {
	fc, val := genFilter(filter)
	q := "SELECT id, ref, name, state, type, created, updated, usr_name, usr_updated FROM devices"
	if fc != "" {
		q += fc
	}
	devRows, err := h.db.QueryContext(ctx, q, val...)
	if err != nil {
		return nil, lib_model.NewInternalError(err)
	}
	defer devRows.Close()
	attrRows, err := h.db.QueryContext(ctx, "SELECT dev_id, is_usr, key_name, value FROM device_attributes;")
	if err != nil {
		return nil, lib_model.NewInternalError(err)
	}
	defer attrRows.Close()
	devices := make(map[string]lib_model.Device)
	for devRows.Next() {
		var device lib_model.Device
		var created, updated, usrUpdated string
		if err = devRows.Scan(&device.ID, &device.Ref, &device.Name, &device.State, &device.Type, &created, &updated, &device.UserData.Name, &usrUpdated); err != nil {
			return nil, lib_model.NewInternalError(err)
		}
		device.Created, err = stringToTime(created)
		if err != nil {
			return nil, lib_model.NewInternalError(err)
		}
		device.Updated, err = stringToTime(updated)
		if err != nil {
			return nil, lib_model.NewInternalError(err)
		}
		device.UserData.Updated, err = stringToTime(usrUpdated)
		if err != nil {
			return nil, lib_model.NewInternalError(err)
		}
		devices[device.ID] = device
	}
	if err = devRows.Err(); err != nil {
		return nil, lib_model.NewInternalError(err)
	}
	for attrRows.Next() {
		var id string
		var isUsr bool
		var devAttr lib_model.DeviceAttribute
		if err = attrRows.Scan(&id, &isUsr, &devAttr.Key, &devAttr.Value); err != nil {
			return nil, lib_model.NewInternalError(err)
		}
		if dev, ok := devices[id]; ok {
			if isUsr {
				dev.UserData.Attributes = append(dev.UserData.Attributes, devAttr)
			} else {
				dev.Attributes = append(dev.Attributes, devAttr)
			}
			devices[id] = dev
		}
	}
	if err = attrRows.Err(); err != nil {
		return nil, lib_model.NewInternalError(err)
	}
	return devices, nil
}

func (h *Handler) Create(ctx context.Context, txItf driver.Tx, device lib_model.DeviceBase) error {
	var tx *sql.Tx
	if txItf != nil {
		tx = txItf.(*sql.Tx)
	} else {
		var e error
		if tx, e = h.db.BeginTx(ctx, nil); e != nil {
			return lib_model.NewInternalError(e)
		}
		defer tx.Rollback()
	}
	_, err := tx.ExecContext(ctx, "INSERT INTO devices (id, ref, name, state, type, created, updated) VALUES (?, ?, ?, ?, ?, ?, ?);", device.ID, device.Ref, device.Name, device.State, device.Type, timeToString(device.Created), timeToString(device.Updated))
	if err != nil {
		return lib_model.NewInternalError(err)
	}
	if len(device.Attributes) > 0 {
		if err = insertAttributes(ctx, tx.PrepareContext, device.ID, false, device.Attributes); err != nil {
			return err
		}
	}
	if txItf == nil {
		if err = tx.Commit(); err != nil {
			return lib_model.NewInternalError(err)
		}
	}
	return nil
}

func (h *Handler) Read(ctx context.Context, id string) (lib_model.Device, error) {
	row := h.db.QueryRowContext(ctx, "SELECT id, ref, name, state, type, created, updated, usr_name, usr_updated FROM devices WHERE id = ?;", id)
	var device lib_model.Device
	var created, updated, usrUpdated string
	err := row.Scan(&device.ID, &device.Ref, &device.Name, &device.State, &device.Type, &created, &updated, &device.UserData.Name, &usrUpdated)
	if err != nil {
		return lib_model.Device{}, lib_model.NewInternalError(err)
	}
	attrRows, err := h.db.QueryContext(ctx, "SELECT is_usr, key_name, value FROM device_attributes WHERE dev_id = ?;", id)
	if err != nil {
		return lib_model.Device{}, lib_model.NewInternalError(err)
	}
	defer attrRows.Close()
	device.Created, err = stringToTime(created)
	if err != nil {
		return lib_model.Device{}, lib_model.NewInternalError(err)
	}
	device.Updated, err = stringToTime(updated)
	if err != nil {
		return lib_model.Device{}, lib_model.NewInternalError(err)
	}
	device.UserData.Updated, err = stringToTime(usrUpdated)
	if err != nil {
		return lib_model.Device{}, lib_model.NewInternalError(err)
	}
	for attrRows.Next() {
		var isUsr bool
		var devAttr lib_model.DeviceAttribute
		if err = attrRows.Scan(&isUsr, &devAttr.Key, &devAttr.Value); err != nil {
			return lib_model.Device{}, lib_model.NewInternalError(err)
		}
		if isUsr {
			device.UserData.Attributes = append(device.UserData.Attributes, devAttr)
		} else {
			device.Attributes = append(device.Attributes, devAttr)
		}
	}
	if err = attrRows.Err(); err != nil {
		return lib_model.Device{}, lib_model.NewInternalError(err)
	}
	return device, nil
}

func (h *Handler) Update(ctx context.Context, txItf driver.Tx, deviceBase lib_model.DeviceBase) error {
	var tx *sql.Tx
	if txItf != nil {
		tx = txItf.(*sql.Tx)
	} else {
		var e error
		if tx, e = h.db.BeginTx(ctx, nil); e != nil {
			return lib_model.NewInternalError(e)
		}
		defer tx.Rollback()
	}
	res, err := tx.ExecContext(ctx, "UPDATE devices SET ref = ?, name = ?, state = ?, type = ?, created = ?, updated = ? WHERE `id` = ?", deviceBase.Ref, deviceBase.Name, deviceBase.State, deviceBase.Type, timeToString(deviceBase.Created), timeToString(deviceBase.Updated), deviceBase.ID)
	if err != nil {
		return lib_model.NewInternalError(err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return lib_model.NewInternalError(err)
	}
	if n < 1 {
		return lib_model.NewNotFoundError(errors.New("not found"))
	}
	_, err = tx.ExecContext(ctx, "DELETE FROM device_attributes WHERE dev_id = ? AND is_usr = ?", deviceBase.ID, false)
	if err != nil {
		return lib_model.NewInternalError(err)
	}
	if len(deviceBase.Attributes) > 0 {
		if err = insertAttributes(ctx, tx.PrepareContext, deviceBase.ID, false, deviceBase.Attributes); err != nil {
			return err
		}
	}
	if txItf == nil {
		if err = tx.Commit(); err != nil {
			return lib_model.NewInternalError(err)
		}
	}
	return nil
}

func (h *Handler) UpdateUserData(ctx context.Context, txItf driver.Tx, id string, userData lib_model.DeviceUserData) error {
	var tx *sql.Tx
	if txItf != nil {
		tx = txItf.(*sql.Tx)
	} else {
		var e error
		if tx, e = h.db.BeginTx(ctx, nil); e != nil {
			return lib_model.NewInternalError(e)
		}
		defer tx.Rollback()
	}
	res, err := tx.ExecContext(ctx, "UPDATE devices SET usr_name = ?, usr_updated = ? WHERE `id` = ?", userData.Name, timeToString(userData.Updated), id)
	if err != nil {
		return lib_model.NewInternalError(err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return lib_model.NewInternalError(err)
	}
	if n < 1 {
		return lib_model.NewNotFoundError(errors.New("not found"))
	}
	_, err = tx.ExecContext(ctx, "DELETE FROM device_attributes WHERE dev_id = ? AND is_usr = ?", id, true)
	if err != nil {
		return lib_model.NewInternalError(err)
	}
	if len(userData.Attributes) > 0 {
		if err = insertAttributes(ctx, tx.PrepareContext, id, true, userData.Attributes); err != nil {
			return err
		}
	}
	if txItf == nil {
		if err = tx.Commit(); err != nil {
			return lib_model.NewInternalError(err)
		}
	}
	return nil
}

func (h *Handler) UpdateStates(ctx context.Context, txItf driver.Tx, ref string, state lib_model.DeviceState, timestamp time.Time) error {
	execContext := h.db.ExecContext
	if txItf != nil {
		tx := txItf.(*sql.Tx)
		execContext = tx.ExecContext
	}
	_, err := execContext(ctx, "UPDATE devices SET state = ?, updated = ? WHERE `ref` = ?", state, timeToString(timestamp), ref)
	if err != nil {
		return lib_model.NewInternalError(err)
	}
	return nil
}

func (h *Handler) Delete(ctx context.Context, txItf driver.Tx, id string) error {
	execContext := h.db.ExecContext
	if txItf != nil {
		tx := txItf.(*sql.Tx)
		execContext = tx.ExecContext
	}
	res, err := execContext(ctx, "DELETE FROM devices WHERE id = ?", id)
	if err != nil {
		return lib_model.NewInternalError(err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return lib_model.NewInternalError(err)
	}
	if n < 1 {
		return lib_model.NewNotFoundError(errors.New("not found"))
	}
	return nil
}

func insertAttributes(ctx context.Context, pf func(ctx context.Context, query string) (*sql.Stmt, error), id string, isUser bool, attributes []lib_model.DeviceAttribute) error {
	stmt, err := pf(ctx, "INSERT INTO device_attributes (dev_id, is_usr, key_name, value) VALUES (?, ?, ?, ?);")
	if err != nil {
		return lib_model.NewInternalError(err)
	}
	defer stmt.Close()
	for _, attr := range attributes {
		if _, err = stmt.ExecContext(ctx, id, isUser, attr.Key, attr.Value); err != nil {
			return lib_model.NewInternalError(err)
		}
	}
	return nil
}

func genFilter(filter lib_model.DevicesFilter) (string, []any) {
	var fc []string
	var val []any
	if len(filter.IDs) > 0 {
		ids := removeDuplicates(filter.IDs)
		fc = append(fc, "id IN ("+strings.Repeat("?, ", len(ids)-1)+"?)")
		for _, id := range ids {
			val = append(val, id)
		}
	}
	if filter.State != "" {
		fc = append(fc, "state = ?")
		val = append(val, filter.State)
	}
	if filter.Type != "" {
		fc = append(fc, "type = ?")
		val = append(val, filter.Type)
	}
	if filter.Ref != "" {
		fc = append(fc, "ref = ?")
		val = append(val, filter.Ref)
	}
	if len(fc) > 0 {
		return " WHERE " + strings.Join(fc, " AND "), val
	}
	return "", nil
}

func removeDuplicates(sl []string) []string {
	if len(sl) < 2 {
		return sl
	}
	set := make(map[string]struct{})
	var sl2 []string
	for _, s := range sl {
		if _, ok := set[s]; !ok {
			sl2 = append(sl2, s)
		}
		set[s] = struct{}{}
	}
	return sl2
}

func timeToString(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(tLayout)
}

func stringToTime(s string) (time.Time, error) {
	if s != "" {
		return time.Parse(tLayout, s)
	}
	return time.Time{}, nil
}
