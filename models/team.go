package models

import (
	"time"
)

type AccessMode int

const (
	ACCESS_MODE_NONE AccessMode = iota
	ACCESS_MODE_READ
	ACCESS_MODE_WRITE
	ACCESS_MODE_ADMIN
	ACCESS_MODE_OWNER
)

type Team struct {
	Id         int64
	LowerName  string  `xorm:"varchar(30) unique INDEX NOT NULL"`
	Name       string  `xorm:"varchar(30) unique INDEX NOT NULL"`
	Desc       string  `xorm:"varchar(255)"`
	Email      string  `xorm:"varchar(80) unique"`
	Members    []*User `xorm:"-"`
	AdminTeam  []*User `xorm:"-"`
	NumMembers int
	Access     int
	Created    time.Time `xorm:"created"`
	Updated    time.Time `xorm:"updated"`
}

func AddTeam(t *Team, owner *User) (err error) {
	err = IsTeamExistByName(t.LowerName)
	if err != nil {
		return err
	}

	err = IsTeamExistByEmail(t.Email)
	if err != nil {
		return err
	}

	err = Insert(t)
	if err != nil {
		return err
	}

	err = GetByExample(t)
	if err != nil {
		return err
	}

	return AddTeamUser(t.Id, owner.Id, ACCESS_MODE_OWNER)
}

func GetTeamByName(name string) (*Team, error) {
	var t = Team{
		Name: name,
	}
	err := GetByExample(&t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func IsTeamExistByName(name string) error {
	if has := IsExist(&Team{LowerName: name}); has {
		return ErrNameAlreadyExist
	}
	return nil
}

func IsTeamExistByEmail(email string) error {
	if has := IsExist(&Team{Email: email}); has {
		return ErrEmailAlreadyUsed
	}
	return nil
}

// OrgUser represents an organization-user relation.
type TeamUser struct {
	Id     int64
	Uid    int64 `xorm:"INDEX UNIQUE(s)"`
	TeamId int64 `xorm:"INDEX UNIQUE(s)"`
	Access AccessMode
}

func GetTeamUsersByTeamId(tid int64) ([]*TeamUser, error) {
	tus := make([]*TeamUser, 0, 10)
	err := orm.Where("team_id=?", tid).Find(&tus)
	return tus, err
}

func AddTeamUser(tid, uid int64, role AccessMode) error {
	if IsTeamMember(tid, uid) {
		return nil
	}

	tu := &TeamUser{
		Uid:    uid,
		TeamId: tid,
		Access: role,
	}

	return Insert(tu)
}

// IsTeamMember returns true if given user is member of organization.
func IsTeamMember(tid, uid int64) bool {
	has, _ := orm.Where("uid=?", uid).And("team_id=?", tid).Get(new(TeamUser))
	return has
}
