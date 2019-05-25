package bodies

type JSColors struct {
	fillStyle   string
	strokeStyle string
}

type Body struct {
	JSColors
	sticky bool
}
