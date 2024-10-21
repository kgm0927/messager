package serverclient

type commandID int

const (
	CMD_NICK  commandID = iota // 0
	CMD_JOIN                   // 1
	CMD_ROOMS                  // 2
	CMD_MSG                    // 3
	CMD_QUIT                   // 4
)

type command struct {
	id     commandID
	client *client
	args   []string
}
