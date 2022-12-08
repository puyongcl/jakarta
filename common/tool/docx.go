package tool

import "github.com/lukasjarosch/go-docx"

func GenNewDocxFromTemplate(templateFile, outFile string, r docx.PlaceholderMap) (err error) {
	// replaceMap is a key-value map whereas the keys
	// represent the placeholders without the delimiters
	//replaceMap := docx.PlaceholderMap{
	//	"NAME":  "张三",
	//	"NAME2": "张三",
	//	"PHONE": "13800009078",
	//	"IDNO":  "8909786789876564567",
	//}

	// read and parse the template docx
	doc, err := docx.Open(templateFile)
	if err != nil {
		return
	}

	// replace the keys with values from replaceMap
	err = doc.ReplaceAll(r)
	if err != nil {
		return
	}

	// write out a new file
	err = doc.WriteToFile(outFile)
	if err != nil {
		return
	}
	return
}
