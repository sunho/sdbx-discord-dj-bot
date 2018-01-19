package commands

type Radio struct {
	Songs           map[string][]Song
	Categories      []string
	PlayingCategory map[string]string
}
