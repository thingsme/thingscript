package stdlib

import (
	"github.com/thingsme/thingscript/object"
	"github.com/thingsme/thingscript/stdlib/arrays"
	"github.com/thingsme/thingscript/stdlib/fmt"
	"github.com/thingsme/thingscript/stdlib/hashmap"
	"github.com/thingsme/thingscript/stdlib/strings"
)

func Packages() []object.Package {
	return []object.Package{
		fmt.New(),
		strings.New(),
		arrays.New(),
		hashmap.New(),
	}
}
