package dao

type RelationInfo struct {
	ID       int64 `gorm:"primaryKey;autoIncrement"`
	Follower int64 `gorm:"uniqueIndex:idx_followee_follower,priority:2;index:follower_status,priority:1"`
	Followee int64 `gorm:"uniqueIndex:idx_followee_follower,priority:1;index:followee_status,priority:1"`

	Status uint8 `gorm:"index:follower_status,priority:2;index:followee_status,priority:2'"`

	Ctime int64
	Utime int64
}
type FollowCount struct {
	ID  int64 `gorm:"primaryKey;autoIncrement"`
	UID int64 `gorm:"uniqueIndex"`

	FollowerCount uint32
	FolloweeCount uint32

	Ctime int64
	Utime int64
}
