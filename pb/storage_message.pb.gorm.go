package pb

import (
	context "context"
	gorm1 "github.com/infobloxopen/atlas-app-toolkit/gorm"
	errors "github.com/infobloxopen/protoc-gen-gorm/errors"
	gorm "github.com/jinzhu/gorm"
	field_mask "google.golang.org/genproto/protobuf/field_mask"
	strings "strings"
)

type StorageORM struct {
	Driver   int32
	LaptopId *string
}

// TableName overrides the default tablename generated by GORM
func (StorageORM) TableName() string {
	return "storages"
}

// ToORM runs the BeforeToORM hook if present, converts the fields of this
// object to ORM format, runs the AfterToORM hook, then returns the ORM object
func (m *Storage) ToORM(ctx context.Context) (StorageORM, error) {
	to := StorageORM{}
	var err error
	if prehook, ok := interface{}(m).(StorageWithBeforeToORM); ok {
		if err = prehook.BeforeToORM(ctx, &to); err != nil {
			return to, err
		}
	}
	to.Driver = int32(m.Driver)
	if posthook, ok := interface{}(m).(StorageWithAfterToORM); ok {
		err = posthook.AfterToORM(ctx, &to)
	}
	return to, err
}

// ToPB runs the BeforeToPB hook if present, converts the fields of this
// object to PB format, runs the AfterToPB hook, then returns the PB object
func (m *StorageORM) ToPB(ctx context.Context) (Storage, error) {
	to := Storage{}
	var err error
	if prehook, ok := interface{}(m).(StorageWithBeforeToPB); ok {
		if err = prehook.BeforeToPB(ctx, &to); err != nil {
			return to, err
		}
	}
	to.Driver = Storage_Driver(m.Driver)
	if posthook, ok := interface{}(m).(StorageWithAfterToPB); ok {
		err = posthook.AfterToPB(ctx, &to)
	}
	return to, err
}

// The following are interfaces you can implement for special behavior during ORM/PB conversions
// of type Storage the arg will be the target, the caller the one being converted from

// StorageBeforeToORM called before default ToORM code
type StorageWithBeforeToORM interface {
	BeforeToORM(context.Context, *StorageORM) error
}

// StorageAfterToORM called after default ToORM code
type StorageWithAfterToORM interface {
	AfterToORM(context.Context, *StorageORM) error
}

// StorageBeforeToPB called before default ToPB code
type StorageWithBeforeToPB interface {
	BeforeToPB(context.Context, *Storage) error
}

// StorageAfterToPB called after default ToPB code
type StorageWithAfterToPB interface {
	AfterToPB(context.Context, *Storage) error
}

// DefaultCreateStorage executes a basic gorm create call
func DefaultCreateStorage(ctx context.Context, in *Storage, db *gorm.DB) (*Storage, error) {
	if in == nil {
		return nil, errors.NilArgumentError
	}
	ormObj, err := in.ToORM(ctx)
	if err != nil {
		return nil, err
	}
	if hook, ok := interface{}(&ormObj).(StorageORMWithBeforeCreate_); ok {
		if db, err = hook.BeforeCreate_(ctx, db); err != nil {
			return nil, err
		}
	}
	if err = db.Create(&ormObj).Error; err != nil {
		return nil, err
	}
	if hook, ok := interface{}(&ormObj).(StorageORMWithAfterCreate_); ok {
		if err = hook.AfterCreate_(ctx, db); err != nil {
			return nil, err
		}
	}
	pbResponse, err := ormObj.ToPB(ctx)
	return &pbResponse, err
}

type StorageORMWithBeforeCreate_ interface {
	BeforeCreate_(context.Context, *gorm.DB) (*gorm.DB, error)
}
type StorageORMWithAfterCreate_ interface {
	AfterCreate_(context.Context, *gorm.DB) error
}

// DefaultApplyFieldMaskStorage patches an pbObject with patcher according to a field mask.
func DefaultApplyFieldMaskStorage(ctx context.Context, patchee *Storage, patcher *Storage, updateMask *field_mask.FieldMask, prefix string, db *gorm.DB) (*Storage, error) {
	if patcher == nil {
		return nil, nil
	} else if patchee == nil {
		return nil, errors.NilArgumentError
	}
	var err error
	var updatedMemory bool
	for i, f := range updateMask.Paths {
		if f == prefix+"Driver" {
			patchee.Driver = patcher.Driver
			continue
		}
		if !updatedMemory && strings.HasPrefix(f, prefix+"Memory.") {
			if patcher.Memory == nil {
				patchee.Memory = nil
				continue
			}
			if patchee.Memory == nil {
				patchee.Memory = &Memory{}
			}
			childMask := &field_mask.FieldMask{}
			for j := i; j < len(updateMask.Paths); j++ {
				if trimPath := strings.TrimPrefix(updateMask.Paths[j], prefix+"Memory."); trimPath != updateMask.Paths[j] {
					childMask.Paths = append(childMask.Paths, trimPath)
				}
			}
			if err := gorm1.MergeWithMask(patcher.Memory, patchee.Memory, childMask); err != nil {
				return nil, nil
			}
		}
		if f == prefix+"Memory" {
			updatedMemory = true
			patchee.Memory = patcher.Memory
			continue
		}
	}
	if err != nil {
		return nil, err
	}
	return patchee, nil
}

// DefaultListStorage executes a gorm list call
func DefaultListStorage(ctx context.Context, db *gorm.DB) ([]*Storage, error) {
	in := Storage{}
	ormObj, err := in.ToORM(ctx)
	if err != nil {
		return nil, err
	}
	if hook, ok := interface{}(&ormObj).(StorageORMWithBeforeListApplyQuery); ok {
		if db, err = hook.BeforeListApplyQuery(ctx, db); err != nil {
			return nil, err
		}
	}
	db, err = gorm1.ApplyCollectionOperators(ctx, db, &StorageORM{}, &Storage{}, nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	if hook, ok := interface{}(&ormObj).(StorageORMWithBeforeListFind); ok {
		if db, err = hook.BeforeListFind(ctx, db); err != nil {
			return nil, err
		}
	}
	db = db.Where(&ormObj)
	ormResponse := []StorageORM{}
	if err := db.Find(&ormResponse).Error; err != nil {
		return nil, err
	}
	if hook, ok := interface{}(&ormObj).(StorageORMWithAfterListFind); ok {
		if err = hook.AfterListFind(ctx, db, &ormResponse); err != nil {
			return nil, err
		}
	}
	pbResponse := []*Storage{}
	for _, responseEntry := range ormResponse {
		temp, err := responseEntry.ToPB(ctx)
		if err != nil {
			return nil, err
		}
		pbResponse = append(pbResponse, &temp)
	}
	return pbResponse, nil
}

type StorageORMWithBeforeListApplyQuery interface {
	BeforeListApplyQuery(context.Context, *gorm.DB) (*gorm.DB, error)
}
type StorageORMWithBeforeListFind interface {
	BeforeListFind(context.Context, *gorm.DB) (*gorm.DB, error)
}
type StorageORMWithAfterListFind interface {
	AfterListFind(context.Context, *gorm.DB, *[]StorageORM) error
}