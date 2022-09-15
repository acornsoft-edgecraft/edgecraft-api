/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package utils

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

// ===== [ Constants and Variables ] =====

const (
	// Tag 처리를 위한 옵션들 (Like enum)
	tagMust    uint8 = 1 << iota // 반드시 복사되어야 할 경우 (복사 실패시는 panic 처리)
	tagNoPanic                   // 반드시 복사가 실패할 경우에 Panic 대신 error로 처리
	tagIgnore                    // 복사 무시
	hasCopied                    // 복사 여부

	// 기본 형식 변환 정의
	String  string  = ""
	Bool    bool    = false
	Int     int     = 0
	Float32 float32 = 0
	Float64 float64 = 0
)

// Errors
var (
	ErrInvalidCopyDestination        = errors.New("copy destination is invalid")
	ErrInvalidCopyFrom               = errors.New("copy from is invalid")
	ErrMapKeyNotMatch                = errors.New("map's key type doesn't match")
	ErrNotSupported                  = errors.New("not supported")
	ErrFieldNameTagStartNotUpperCase = errors.New("'Copy' field name tag must be start upper case")
)

type (
	// converterPair - 변환할 구조들에 대한 Reflect.Type 정보
	converterPair struct {
		SrcType reflect.Type
		DstType reflect.Type
	}

	// tagNameMapping - Field and Tag Mapping 정보
	tagNameMapping struct {
		FieldNameToTag map[string]string
		TagToFieldName map[string]string
	}

	// flags - Struct Tag's flag 정보
	flags struct {
		BitFlags  map[string]uint8
		SrcNames  tagNameMapping
		DestNames tagNameMapping
	}
)

// TypeConverter - 형식 변환 처리기 정보
type TypeConverter struct {
	SrcType interface{}
	DstType interface{}
	Fn      func(src interface{}) (dst interface{}, err error)
}

// Option - Copy options 정보
type Option struct {
	IgnoreEmpty bool            // true가 되면 모든 Zero value 필드들을 무시한다.
	DeepCopy    bool            // Nested Struct 포함 여부
	Converters  []TypeConverter // 형식 변환에 사용할 처리기들
}

// converters - Converter들에 대한 형식 정보들을 Map 형식으로 처리
func (opt Option) converters() map[converterPair]TypeConverter {
	var converters = map[converterPair]TypeConverter{}

	// 빠른 검색을 위해 Map으로 구성
	for i := range opt.Converters {
		pair := converterPair{
			SrcType: reflect.TypeOf(opt.Converters[i].SrcType),
			DstType: reflect.TypeOf(opt.Converters[i].DstType),
		}

		converters[pair] = opt.Converters[i]
	}
	return converters
}

// CopyTo - Struct 간 복사
func CopyTo(toValue interface{}, fromValue interface{}) (err error) {
	return copier(toValue, fromValue, Option{})
}

// CopyTo - 옵션 기준의 Struct간 복사
func CopyToWithOption(toValue interface{}, fromValue interface{}, opt Option) (err error) {
	return copier(toValue, fromValue, opt)
}

// copier - Struct간 복사
func copier(toValue interface{}, fromValue interface{}, opt Option) (err error) {
	var (
		isSlice    bool
		amount     = 1
		from       = indirect(reflect.ValueOf(fromValue))
		to         = indirect(reflect.ValueOf(toValue))
		converters = opt.converters()
	)

	if !to.CanAddr() {
		return ErrInvalidCopyDestination
	}

	if !from.IsValid() {
		return ErrInvalidCopyFrom
	}

	fromType, isPrtFrom := indirectType(from.Type())
	toType, _ := indirectType(to.Type())

	if fromType.Kind() == reflect.Interface {
		fromType = reflect.TypeOf(from.Interface())
	}

	if toType.Kind() == reflect.Interface {
		toType, _ = indirectType(reflect.TypeOf(to.Interface()))
		oldTo := to
		to := reflect.New(reflect.TypeOf(to.Interface())).Elem()
		defer func() {
			oldTo.Set(to)
		}()
	}

	// Normal to Normal
	if from.Kind() != reflect.Slice && from.Kind() != reflect.Struct && from.Kind() != reflect.Map && (from.Type().AssignableTo(to.Type()) || from.Type().ConvertibleTo(to.Type())) {
		if !isPrtFrom || !opt.DeepCopy {
			to.Set(from.Convert(to.Type()))
		} else {
			fromCopy := reflect.New(from.Type())
			fromCopy.Set(from.Elem())
			to.Set(fromCopy.Convert(to.Type()))
		}
		return
	}

	// Map to Map
	if from.Kind() != reflect.Slice && fromType.Kind() == reflect.Map && toType.Kind() == reflect.Map {
		if !fromType.Key().ConvertibleTo(toType.Key()) {
			return ErrMapKeyNotMatch
		}

		if to.IsNil() {
			to.Set(reflect.MakeMapWithSize(toType, from.Len()))
		}

		for _, k := range from.MapKeys() {
			toKey := indirect(reflect.New(toType.Key()))
			if !set(toKey, k, opt.DeepCopy, converters) {
				return fmt.Errorf("%w map, old key: %v, new key: %v", ErrNotSupported, k.Type(), toType.Key())
			}

			elemType := toType.Elem()
			if elemType.Kind() != reflect.Slice {
				elemType, _ = indirectType(elemType)
			}
			toValue := indirect(reflect.New(elemType))
			if !set(toValue, from.MapIndex(k), opt.DeepCopy, converters) {
				if err = copier(toValue.Addr().Interface(), from.MapIndex(k).Interface(), opt); err != nil {
					return err
				}
			}

			for {
				if elemType == toType.Elem() {
					to.SetMapIndex(toKey, toValue)
					break
				}
				elemType = reflect.PtrTo(elemType)
				toValue = toValue.Addr()
			}
		}
		return
	}

	// Slice to Slice (convertible)
	if from.Kind() == reflect.Slice && to.Kind() == reflect.Slice && fromType.ConvertibleTo(toType) {
		if to.IsNil() {
			slice := reflect.MakeSlice(reflect.SliceOf(to.Type().Elem()), from.Len(), from.Cap())
			to.Set(slice)
		}

		for i := 0; i < from.Len(); i++ {
			if to.Len() < i+1 {
				to.Set(reflect.Append(to, reflect.New(to.Type().Elem()).Elem()))
			}

			if !set(to.Index(i), from.Index(i), opt.DeepCopy, converters) {
				// ignore error while copy slice element
				err = copier(to.Index(i).Addr().Interface(), from.Index(i).Interface(), opt)
				if err != nil {
					continue
				}
			}
		}
		return
	}

	// not struct (least one). not suported
	if fromType.Kind() != reflect.Struct || toType.Kind() != reflect.Struct {
		return
	}

	// slice (least one)
	if from.Kind() == reflect.Slice || to.Kind() == reflect.Slice {
		isSlice = true
		if from.Kind() == reflect.Slice {
			amount = from.Len()
		}
	}

	for i := 0; i < amount; i++ {
		var dest, source reflect.Value

		if isSlice {
			// source
			if from.Kind() == reflect.Slice {
				source = indirect(from.Index(i))
			} else {
				source = indirect(from)
			}
			// dest
			dest = indirect(reflect.New(toType).Elem())
		} else {
			source = indirect(from)
			dest = indirect(to)
		}

		destKind := dest.Kind()
		initDest := false
		if destKind == reflect.Interface {
			initDest = true
			dest = indirect(reflect.New(toType))
		}

		// Get tag options
		flgs, err := getFlags(dest, source, toType, fromType)
		if err != nil {
			return err
		}

		// check source
		if source.IsValid() {
			copyUnexportedStructFields(dest, source)

			// Copy from source field to dest field or method
			fromTypeFields := deepFields(fromType)
			for _, field := range fromTypeFields {
				name := field.Name
				// Get bit flags for field
				fieldFlags := flgs.BitFlags[name]

				// Check if we should ignore copying
				if (fieldFlags & tagIgnore) != 0 {
					continue
				}

				srcFieldName, destFieldName := getFieldName(name, flgs)
				if fromField := source.FieldByName(srcFieldName); fromField.IsValid() && !shouldIgnore(fromField, opt.IgnoreEmpty) {
					// process for nested anonymous field
					destFieldNotSet := false
					if f, ok := dest.Type().FieldByName(destFieldName); ok {
						for idx := range f.Index {
							destField := dest.FieldByIndex(f.Index[:idx+1])
							if destField.Kind() != reflect.Ptr {
								continue
							}
							if !destField.IsNil() {
								continue
							}
							if !destField.CanSet() {
								destFieldNotSet = true
								break
							}
							// destField is a nil pointer that can be set
							newValue := reflect.New(destField.Type().Elem())
							destField.Set(newValue)
						}
					}

					if destFieldNotSet {
						break
					}

					toField := dest.FieldByName(destFieldName)
					if toField.IsValid() {
						if toField.CanSet() {
							if !set(toField, fromField, opt.DeepCopy, converters) {
								if err := copier(toField.Addr().Interface(), fromField.Interface(), opt); err != nil {
									return err
								}
							}
							if fieldFlags != 0 {
								// Note that a copy was made
								flgs.BitFlags[name] = fieldFlags | hasCopied
							}
						}
					} else {
						// try to set to method
						var toMethod reflect.Value
						if dest.CanAddr() {
							toMethod = dest.Addr().MethodByName(destFieldName)
						} else {
							toMethod = dest.MethodByName(destFieldName)
						}

						if toMethod.IsValid() && toMethod.Type().NumIn() == 1 && fromField.Type().AssignableTo(toMethod.Type().In(0)) {
							toMethod.Call([]reflect.Value{fromField})
						}
					}
				}
			}

			// Copy from method to dest field
			for _, field := range deepFields(toType) {
				name := field.Name
				srcFieldName, destFieldName := getFieldName(name, flgs)

				var fromMethod reflect.Value
				if source.CanAddr() {
					fromMethod = source.Addr().MethodByName(srcFieldName)
				} else {
					fromMethod = source.MethodByName(srcFieldName)
				}

				if fromMethod.IsValid() && fromMethod.Type().NumIn() == 0 && fromMethod.Type().NumOut() == 1 && !shouldIgnore(fromMethod, opt.IgnoreEmpty) {
					if toField := dest.FieldByName(destFieldName); toField.IsValid() && toField.CanSet() {
						values := fromMethod.Call([]reflect.Value{})
						if len(values) >= 1 {
							set(toField, values[0], opt.DeepCopy, converters)
						}
					}
				}
			}
		}

		if isSlice && to.Kind() == reflect.Slice {
			if dest.Addr().Type().AssignableTo(to.Type().Elem()) {
				if to.Len() < i+1 {
					to.Set(reflect.Append(to, dest.Addr()))
				} else {
					if !set(to.Index(i), dest.Addr(), opt.DeepCopy, converters) {
						// Ignore error while copyt slice element
						err = copier(to.Index(i).Addr().Interface(), dest.Addr().Interface(), opt)
						if err != nil {
							continue
						}
					}
				}
			} else if dest.Type().AssignableTo(to.Type().Elem()) {
				if to.Len() < i+1 {
					to.Set(reflect.Append(to, dest))
				} else {
					if !set(to.Index(i), dest, opt.DeepCopy, converters) {
						// Ignore error while copy
						err = copier(to.Index(i).Addr().Interface(), dest.Interface(), opt)
						if err != nil {
							continue
						}
					}
				}
			}
		} else if initDest {
			to.Set(dest)
		}

		//lint:ignore SA4006 Ignore unused variable
		err = checkBitFlags(flgs.BitFlags)
	}

	return
}

// copyUnexportedStructFields - 비공개 필드들 복사
func copyUnexportedStructFields(to, from reflect.Value) {
	if from.Kind() != reflect.Struct || to.Kind() != reflect.Struct || !from.Type().AssignableTo(to.Type()) {
		return
	}

	// 모든 필드에 대한 swallow copy 생성
	tmp := indirect(reflect.New(to.Type()))
	tmp.Set(from)

	// exported 필드 구성
	for i := 0; i < to.NumField(); i++ {
		if tmp.Field(i).CanSet() {
			tmp.Field(i).Set(to.Field(i))
		}
	}
	to.Set(tmp)
}

// shouldIgnore - 생략 가능 여부 (Zero Value)
func shouldIgnore(v reflect.Value, ignoreEmpty bool) bool {
	return ignoreEmpty && v.IsZero()
}

var deepFieldsLock sync.RWMutex
var deepFieldsMap = make(map[reflect.Type][]reflect.StructField)

// deepFields - Nested Struct Field 처리 (Caching)
func deepFields(reflectType reflect.Type) []reflect.StructField {
	deepFieldsLock.RLock()
	cache, ok := deepFieldsMap[reflectType]
	deepFieldsLock.RUnlock()
	if ok {
		return cache
	}

	var res []reflect.StructField
	if reflectType, _ = indirectType(reflectType); reflectType.Kind() == reflect.Struct {
		fields := make([]reflect.StructField, 0, reflectType.NumField())

		for i := 0; i < reflectType.NumField(); i++ {
			v := reflectType.Field(i)
			// PkgPath is the package path that qualifies a lower case (unexported)
			// field name. It is empty for upper case (exported) field names.
			// See https://golang.org/ref/spec#Uniqueness_of_identifiers
			if v.PkgPath == "" {
				fields = append(fields, v)
				if v.Anonymous {
					// also consider fields of anonymous fields as fields of the root
					fields = append(fields, deepFields(v.Type)...)
				}
			}
		}
		res = fields
	}

	deepFieldsLock.Lock()
	deepFieldsMap[reflectType] = res
	deepFieldsLock.Unlock()

	return res
}

// indirect - Ptr to Element for Value (deep)
func indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

// indirectType - Ptr or Slice to ElementType for Type (deep)
func indirectType(reflecType reflect.Type) (_ reflect.Type, isPtr bool) {
	for reflecType.Kind() == reflect.Ptr || reflecType.Kind() == reflect.Slice {
		reflecType = reflecType.Elem()
		isPtr = true
	}
	return reflecType, isPtr
}

// set - Copy value from source to dest
func set(to, from reflect.Value, deepCopy bool, converters map[converterPair]TypeConverter) bool {
	if !from.IsValid() {
		return true
	}
	if ok, err := lookupAndCopyWithConverter(to, from, converters); err != nil {
		return false
	} else if ok {
		return true
	}

	if to.Kind() == reflect.Ptr {
		// set 'to' to nil if from is nil
		if from.Kind() == reflect.Ptr && from.IsNil() {
			to.Set(reflect.Zero(to.Type()))
			return true
		} else if to.IsNil() {
			// 'from' -> 'to'
			// sql.NullString -> *string
			if fromValuer, ok := driverValuer(from); ok {
				v, err := fromValuer.Value()
				if err != nil {
					return false
				}
				// if 'from' is not valid do nothing with 'to'
				if v == nil {
					return true
				}
			}
			// allocate new 'to' variable with default value (eg. *string -> new(string))
			to.Set(reflect.New(to.Type().Elem()))
		}
		// depointer 'to'
		to = to.Elem()
	}

	if deepCopy {
		toKind := to.Kind()
		if toKind == reflect.Interface && to.IsNil() {
			if reflect.TypeOf(from.Interface()) != nil {
				to.Set(reflect.New(reflect.TypeOf(from.Interface())).Elem())
				toKind = reflect.TypeOf(to.Interface()).Kind()
			}
		}
		if from.Kind() == reflect.Ptr && from.IsNil() {
			return true
		}
		if to.Kind() == reflect.Struct || toKind == reflect.Map || toKind == reflect.Slice {
			return false
		}
	}

	if from.Type().ConvertibleTo(to.Type()) {
		to.Set(from.Convert(to.Type()))
	} else if toScanner, ok := to.Addr().Interface().(sql.Scanner); ok {
		// 'from' -> 'to'
		// *string -> sql.NullString
		if from.Kind() == reflect.Ptr {
			// if 'from' is nil do nothing with 'to'
			if from.IsNil() {
				return true
			}
			// depointer 'from'
			from = indirect(from)
		}
		// `from` -> `to`
		// string -> sql.NullString
		// set `to` by invoking method Scan(`from`)
		err := toScanner.Scan(from.Interface())
		if err != nil {
			return false
		}
	} else if fromValuer, ok := driverValuer(from); ok {
		// `from`         -> `to`
		// sql.NullString -> string
		v, err := fromValuer.Value()
		if err != nil {
			return false
		}
		// if `from` is not valid do nothing with `to`
		if v == nil {
			return true
		}
		rv := reflect.ValueOf(v)
		if rv.Type().AssignableTo(to.Type()) {
			to.Set(rv)
		} else if to.CanSet() && rv.Type().ConvertibleTo(to.Type()) {
			to.Set(rv.Convert(to.Type()))
		}
	} else if from.Kind() == reflect.Ptr {
		return set(to, from.Elem(), deepCopy, converters)
	} else {
		return false
	}

	return true
}

// lookupAndCopyWithConverter - Converter를 이용해서 복사
func lookupAndCopyWithConverter(to, from reflect.Value, converters map[converterPair]TypeConverter) (copied bool, err error) {
	pair := converterPair{
		SrcType: from.Type(),
		DstType: to.Type(),
	}

	if cnv, ok := converters[pair]; ok {
		result, err := cnv.Fn(from.Interface())
		if err != nil {
			return false, err
		}

		if result != nil {
			to.Set(reflect.ValueOf(result))
		} else {
			// in case we've got a nil value to copy
			to.Set(reflect.Zero(to.Type()))
		}
		return true, nil
	}

	return false, nil
}

// parseTags - Struct tags를 파싱해서 uint8 bit flags로 반환
func parseTags(tag string) (flg uint8, name string, err error) {
	for _, t := range strings.Split(tag, ",") {
		switch t {
		case "-":
			flg = tagIgnore
			return
		case "must":
			flg = flg | tagMust
		case "nopanic":
			flg = flg | tagNoPanic

		default:
			if unicode.IsUpper([]rune(t)[0]) {
				name = strings.TrimSpace(t)
			} else {
				err = ErrFieldNameTagStartNotUpperCase
			}
		}
	}
	return
}

// getFlags - Struct tags를 파싱해서 bit flags와 필드 명을 반환
func getFlags(dest, src reflect.Value, toType, fromType reflect.Type) (flags, error) {
	flgs := flags{
		BitFlags: map[string]uint8{},
		SrcNames: tagNameMapping{
			FieldNameToTag: map[string]string{},
			TagToFieldName: map[string]string{},
		},
		DestNames: tagNameMapping{
			FieldNameToTag: map[string]string{},
			TagToFieldName: map[string]string{},
		},
	}
	var toTypeFields, fromTypeFields []reflect.StructField
	if dest.IsValid() {
		toTypeFields = deepFields(toType)
	}
	if src.IsValid() {
		fromTypeFields = deepFields(fromType)
	}

	// Get a list dest of tags
	for _, field := range toTypeFields {
		tags := field.Tag.Get("copier")
		if tags != "" {
			var name string
			var err error
			if flgs.BitFlags[field.Name], name, err = parseTags(tags); err != nil {
				return flags{}, err
			} else if name != "" {
				flgs.DestNames.FieldNameToTag[field.Name] = name
				flgs.DestNames.TagToFieldName[name] = field.Name
			}
		}
	}

	// Get a list source of tags
	for _, field := range fromTypeFields {
		tags := field.Tag.Get("copier")
		if tags != "" {
			var name string
			var err error
			if _, name, err = parseTags(tags); err != nil {
				return flags{}, err
			} else if name != "" {
				flgs.SrcNames.FieldNameToTag[field.Name] = name
				flgs.SrcNames.TagToFieldName[name] = field.Name
			}
		}
	}

	return flgs, nil
}

// checkBitFlags - Error 또는 Panic 조건에 대한 Flasg 검증
func checkBitFlags(flagsList map[string]uint8) (err error) {
	// Check flag condition were met
	for name, flags := range flagsList {
		if flags&hasCopied == 0 {
			switch {
			case flags&tagMust != 0 && flags&tagNoPanic != 0:
				err = fmt.Errorf("field %s has must tag but was not copied", name)
				return
			case flags&(tagMust) != 0:
				panic(fmt.Sprintf("Field %s has must tag but was not copied", name))
			}
		}
	}
	return
}

// getFieldName - Field와 Tag를 기준으로 필드명 반환
func getFieldName(fieldName string, flgs flags) (srcFieldName string, destFieldName string) {
	// get dest field name
	if srcTagName, ok := flgs.SrcNames.FieldNameToTag[fieldName]; ok {
		destFieldName = srcTagName
		if destTagName, ok := flgs.DestNames.TagToFieldName[srcTagName]; ok {
			destFieldName = destTagName
		}
	} else {
		if destTagName, ok := flgs.DestNames.TagToFieldName[fieldName]; ok {
			destFieldName = destTagName
		}
	}

	if destFieldName == "" {
		destFieldName = fieldName
	}

	// get source field name
	if destTagName, ok := flgs.DestNames.FieldNameToTag[fieldName]; ok {
		srcFieldName = destTagName
		if srcField, ok := flgs.SrcNames.TagToFieldName[destTagName]; ok {
			srcFieldName = srcField
		}
	} else {
		if srcField, ok := flgs.SrcNames.TagToFieldName[fieldName]; ok {
			destFieldName = srcField
		}
	}

	if srcFieldName == "" {
		srcFieldName = fieldName
	}

	return
}

// driverValuer - DB 값 처리기 (database/sql 사용)
func driverValuer(v reflect.Value) (i driver.Valuer, ok bool) {
	if !v.CanAddr() {
		i, ok = v.Interface().(driver.Valuer)
		return
	}
	i, ok = v.Addr().Interface().(driver.Valuer)
	return
}
