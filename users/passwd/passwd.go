// Copyright (c) 2015 Deepin Ltd. All rights reserved.
// Use of this source is govered by General Public License that can be found
// in the LICENSE file.
package passwd

// #include <stdlib.h>
// #include <pwd.h>
import "C"

import (
	"fmt"
	"unsafe"
)

// Passwd wraps up `passwd` struct used in kernel
type Passwd struct {
	Name    string
	Passwd  string // Password is encrypted
	Uid     uint32
	Gid     uint32
	Comment string
	Home    string
	Shell   string
}

// passwdC2Go converts `passwd` struct from C to golang native struct
func passwdC2Go(passwdC *C.struct_passwd) *Passwd {
	return &Passwd{
		Name:    C.GoString(passwdC.pw_name),
		Passwd:  C.GoString(passwdC.pw_passwd),
		Uid:     uint32(passwdC.pw_uid),
		Gid:     uint32(passwdC.pw_gid),
		Comment: C.GoString(passwdC.pw_gecos),
		Home:    C.GoString(passwdC.pw_dir),
		Shell:   C.GoString(passwdC.pw_shell),
	}
}

type UserNotFoundError struct {
	Name string
	Uid  uint32
}

func (err *UserNotFoundError) Error() string {
	if len(err.Name) > 0 {
		return fmt.Sprintf("User with name `%s` not found!", err.Name)
	} else {
		return fmt.Sprintf("User with uid `%d` not found!", err.Uid)
	}
}

// GetPasswdByName wraps up `getpwnam` system call
func GetPasswdByName(name string) (*Passwd, error) {
	nameC := C.CString(name)
	defer C.free(unsafe.Pointer(nameC))
	passwdC, err := C.getpwnam(nameC)
	if passwdC == nil {
		if err == nil {
			return nil, &UserNotFoundError{Name: name}
		} else {
			return nil, err
		}
	} else {
		return passwdC2Go(passwdC), nil
	}
}

// GetPasswdByUid wraps up `getpwuid` system call
func GetPasswdByUid(uid uint32) (*Passwd, error) {
	uidC := C.__uid_t(uid)
	passwdC, err := C.getpwuid(uidC)
	if passwdC == nil {
		if err == nil {
			return nil, &UserNotFoundError{Uid: uid}
		} else {
			return nil, err
		}
	} else {
		return passwdC2Go(passwdC), nil
	}
}