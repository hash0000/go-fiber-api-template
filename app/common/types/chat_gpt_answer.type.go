package types

type StoryChapter struct {
	Title              string    `json:"title"`
	Chapters           []Chapter `json:"chapters"`
	QuestionsAboutTale []string  `json:"questions_about_tale"`
}

type Chapter struct {
	Title         string `json:"title"`
	Text          string `json:"text"`
	PicGeneration string `json:"pic_generation"`
}
