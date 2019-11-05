package client

// Persist makes sure that the executable is run after a reboot
func Persist(path string) {
	elevated := CheckElevate()
	if elevated {
		persistAdmin(path)
	} else {
		persistUser(path)
	}
}

// persistAdmin persistence using admin privileges
func persistAdmin(path string) {

}

// persistUser persistence using user privileges
func persistUser(path string) {

}
