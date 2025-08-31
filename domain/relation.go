package domain

const (
	StatusUnknown uint8 = iota
	StatusFollowing
	StatusUnFollowing
)

type RelationInfo struct {
	ID       int64
	Follower int64
	Followee int64

	Status uint8

	Ctime int64
	Utime int64
}

type FollowCount struct {
	ID  int64
	UID int64

	FollowerCount uint32
	FolloweeCount uint32

	Ctime int64
	Utime int64
}
