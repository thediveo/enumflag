// Copyright 2020 Harald Albrecht.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy
// of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package enumflag

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// Flag represents a CLI (enumeration) flag which can take on only a
// single enumeration value out of a fixed set of enumeration values.
// Applications using the enumflag package might want to derive their
// enumeration flags from Flag, such as "type MyFoo enumflag.Flag"
type Flag uint

// Mapper returns the mapping between enumeration values and their
// corresponding textual representations. If multiple textual representations
// exist for the same enumeration value, then the first textual representation
// is considered to be canonical and the one returned by String() when applied
// on an enumeration flag. Enumeration flags must implement this interface in
// order to be managed by the enumflag package.
type Mapper interface {
	Enums() (interface{}, EnumCaseSensitivity)
}

// EnumCaseSensitivity specifies whether the textual representations of enum
// values are considered to be case sensitive, or not.
type EnumCaseSensitivity bool

// Controls whether the textual representations for enum values are case
// sensitive, or not. If they are case insensitive, then the textual
// representation must be registered in lower case.
const (
	EnumCaseInsensitive EnumCaseSensitivity = false // case insensitive textual representations
	EnumCaseSensitive   EnumCaseSensitivity = true  // case sensitive textual representations
)

// String returns the textual representation of an enumeration (flag) value.
// In case multiple textual representations (identifiers) exist for the same
// enumeration value, then only the first one is returned, which is considered
// to be the canonical textual representation.
func String(flag interface{}) string {
	idmap := cachedMapping(flag)
	flagval := reflect.ValueOf(flag)
	if flagval.Kind() == reflect.Ptr {
		flagval = flagval.Elem()
	}
	return idmap.EnumIdentifiers[flagval.Convert(enumFlagType).Interface().(Flag)][0]
}

// Set sets an enumeration flag to the value corresponding with the given
// textual representation. It returns an error if the textual representation
// doesn't match any registered one for the given enumeration flag.
func Set(flag interface{}, s string) error {
	idmap := cachedMapping(flag)
	if idmap.EnumCaseSensitivity == EnumCaseInsensitive {
		s = strings.ToLower(s)
	}
	flagval := reflect.ValueOf(flag).Elem()
	for enumval, ids := range idmap.EnumIdentifiers {
		for _, id := range ids {
			if s == id {
				flagval.Set(reflect.ValueOf(enumval).Convert(flagval.Type()))
				return nil
			}
		}
	}
	// Oh no! An invalid textual enum value was specified, so let's generate
	// some useful error explaining which textual representations are valid.
	// We're ordering values by their canonical names in order to achieve a
	// stable error message.
	allids := []string{}
	for _, ids := range idmap.EnumIdentifiers {
		s := []string{}
		for _, id := range ids {
			s = append(s, "'"+id+"'")
		}
		allids = append(allids, strings.Join(s, "/"))
	}
	sort.Strings(allids)
	return fmt.Errorf("must be %s", strings.Join(allids, ", "))
}

// EnumIdentifiers maps enumeration values to their corresponding textual
// representations. This mapping is a one-to-many mapping in that the same
// enumeration value may have more than only one associated textual
// representation.
type EnumIdentifiers map[Flag][]string

// enumTypeCache maps individual enumeration types to their enumeration
// mappings, where these mappings associate enumeration values with textual
// representations. Additionally, the cache also stores whether the textual
// representations are case sensitive or not.
var enumTypeCache = map[string]enumTypeCacheEntry{}

// enumTypeCache stores the mapping between enum values and textual
// representations, as well as the case sensitivity, for a single specific
// enumeration type.
type enumTypeCacheEntry struct {
	EnumIdentifiers
	EnumCaseSensitivity
}

// Reflection types used in this package.
var enumFlagType = reflect.TypeOf(Flag(0))
var uintType = reflect.TypeOf(uint(0))

// cachedMapping returns the enumeration value-to-textual identifiers mapping
// for the enumeration type of the specified value. For better performance, it
// will automatically cache these mappings after it has converted them into
// our internal canonical format.
func cachedMapping(flag interface{}) enumTypeCacheEntry {
	// Do we already have this enumeration flag type in our cache? If so, then
	// return the cached mapping, based on the enum flag type.
	flagval := reflect.ValueOf(flag)
	if flagval.Kind() == reflect.Ptr { // ...be forgiving in the receiver type.
		flagval = flagval.Elem()
	}
	flagtype := flagval.Type()
	flagtypename := flagtype.Name()
	if cacheentry, ok := enumTypeCache[flagtypename]; ok {
		return cacheentry
	}
	// Nothing was cached so far, so we now need to ask this new type for
	// its enum mapping, convert that into our canonical representation,
	// and then cache it. First, make sure that the enumeration type is
	// compatible with our enumeration flag mechanism.
	if !flagtype.ConvertibleTo(uintType) {
		panic(fmt.Sprintf("incompatible enumeration type %s", flagtypename))
	}
	// Next, ensure that the enumeration value-to-textual ids map actually
	// is a map. Please note that we must not use "flagval" here, as it
	// might have been dereferences in case of a pointer to an enum flag:
	// this would then forbid us to convert the enum flag into a Mapper
	// interface.
	enumids, sensitivity := reflect.ValueOf(flag).Interface().(Mapper).Enums()
	enumidsrval := reflect.ValueOf(enumids)
	if enumidsrval.Kind() != reflect.Map {
		panic(fmt.Sprintf("incompatible enumeration identifiers map type %s",
			enumidsrval.Type().Name()))
	}
	// Oh, the magic of Golang reflexions makes us put on our beer googles,
	// erm, goggles: we now convert the specified enum values to enum textual
	// representations mapping into our "canonical" mapping. While we can keep
	// the map values (which are string slices), we have to convert the map
	// keys into our canonical EnumFlag key type (=enum values type).
	cacheentry := enumTypeCacheEntry{
		EnumIdentifiers:     EnumIdentifiers{},
		EnumCaseSensitivity: sensitivity,
	}
	for _, key := range enumidsrval.MapKeys() {
		cacheentry.EnumIdentifiers[key.Convert(enumFlagType).Interface().(Flag)] =
			enumidsrval.MapIndex(key).Interface().([]string)
	}
	enumTypeCache[flagtypename] = cacheentry
	return cacheentry
}
