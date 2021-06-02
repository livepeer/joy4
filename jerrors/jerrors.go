package jerrors

import "errors"

// ErrNoAudioInfoFound means that info about audio codec was not found
var ErrNoAudioInfoFound = errors.New("No audio codec info found")

// ErrNoVideoInfoFound means that info about video codec was not found
var ErrNoVideoInfoFound = errors.New("No video codec info found")
