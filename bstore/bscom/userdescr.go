package bscom

type UserDescr struct {
    Id      int64       `json:"id"      db:"id"`
    Login   string      `json:"login"   db:"login"`
    Pass    string      `json:"pass"    db:"pass"`
    State   string      `json:"state"   db:"state"`
}
