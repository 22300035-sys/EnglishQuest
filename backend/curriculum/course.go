// Package curriculum embeds private authoring metadata; it is never served to learners.
package curriculum

import _ "embed"

//go:embed course.json
var Data []byte
