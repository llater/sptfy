package sptfy

type Sptfy struct{}

// Returns whether a Spotify user exists with this id, then the user, and an error.
// Use to check if a user exists with a given Spotify user id.
func (s *Sptfy) GetUserById(spotifyId string) (bool, *SptfyUser, error) {
   
}