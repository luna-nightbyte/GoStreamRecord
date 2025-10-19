package settings

// type Logins struct {
// 	Users []Login
// }

// type Login struct {
// 	Name string `json:"name"`
// 	Key  string `json:"key"`
// }

// func (l *Logins) GetKey(name string) string {
// 	for _, user := range l.Users {
// 		if user.Name == name {
// 			return user.Key
// 		}
// 	}
// 	return ""
// }
// func (l *Logins) Remove(name string) {
// 	newList := []Login{}
// 	for _, user := range l.Users {
// 		if user.Name == name {
// 			continue
// 		}
// 		newList = append(newList, user)
// 	}
// 	l.Users = newList
// }
// func (l *Logins) Add(name, key string) bool {
// 	for _, user := range l.Users {
// 		if user.Name == name {
// 			return false
// 		}
// 	}
// 	l.Users = append(l.Users, Login{Name: name, Key: key})
// 	return true
// }
// func (l *Logins) Modify(oldname, newname, key string) bool {

// 	for i, user := range l.Users {
// 		if user.Name == oldname {
// 			l.Users[i].Name = newname
// 			l.Users[i].Key = key
// 			return true
// 		}
// 	}
// 	return false
// }
// func (l *Logins) Update(name, key string) {
// 	for i, user := range l.Users {
// 		if user.Name == name {
// 			l.Users[i].Key = key
// 			return
// 		}
// 	}
// }
// func (l *Logins) Exists(name string) bool {
// 	for _, user := range l.Users {
// 		if user.Name == name {
// 			return true
// 		}
// 	}
// 	return false
// }
