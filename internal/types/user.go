package types

// User represents the user entity from the API
type User struct {
	ID               string
	Email            string
	Name             string
	ProfilePicture   string
	ActiveWorkspace  string
	DefaultWorkspace string
}
