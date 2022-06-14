package apis

type collaboration struct {
	Members        []collaborationMembership `json:"members"`
	PendingMembers []collaborationMembership `json:"pendingMembers"`
	Teams          []collaborationMembership `json:"teams"`
}

type collaborationMembership struct {
	Name  string   `json:"name"`
	Roles []string `json:"roles"`
}
