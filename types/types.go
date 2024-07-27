package types

type UserInfo struct {
	Name              string
	Email             string
	Company           string
	Phone             string
	Position          string
	G7Company         string
	ExistingG7Project string
	Range             string
	UpcomingG7Project string
	Description       string
}

type User struct {
	Email string
}

type Tender struct {
	ID            int
	Name          string
	Link          string
	CreatedAt     string
	IsNotified    bool
	TarikhLawatan string
	KodBidang     string
	KebenaranKhas string
	TarikhIklan   string
	Taraf         string
	Index         int
}
