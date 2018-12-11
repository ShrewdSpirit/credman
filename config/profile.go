package config

type Profile struct {
	Sites []Site
}

func ProfileFromUserProfile(up *UserProfile, password string) *Profile {
	if len(up.Content) == 0 {
		return &Profile{
			Sites: make([]Site, 0),
		}
	}

	if hash(password) != up.Hash {
		return nil
	}

	return nil
}
