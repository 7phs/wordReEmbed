package glove

import (
	"fmt"
	"os"

	"github.com/7phs/wordReEmed/console/options"
	"github.com/7phs/wordReEmed/console/storages"
	"github.com/7phs/wordReEmed/lib/model/glove"
)

func main() {
	os.Stdout.WriteString("BUILDING VOCABULARY\n")

	corpus, err := storages.OpenTextReader(options.CorpusFileName, ' ')
	if err != nil {
		os.Stderr.WriteString("Error open file: '" + options.CorpusFileName + "': " + err.Error() + "\n")
		return
	}
	defer corpus.Close()

	var (
		vocab = glove.NewVocabCount()
		count = 0
	)

	err = nil

	os.Stdout.WriteString(fmt.Sprint("Processing ", count, " tokens"))

	corpus.ForEach(func(word []byte) bool {
		if len(word) == 0 {
			return true
		}

		vocab.Add(string(word))

		count++
		if (count % 100000) == 0 {
			os.Stdout.WriteString(fmt.Sprint("\033[11G ", count, " tokens"))
		}

		return true
	})

	os.Stdout.WriteString(fmt.Sprint("\033[0GProcessed ", count, " tokens\n"))

	os.Stdout.WriteString(fmt.Sprint("Counted ", len(vocab), " unique words\n"))

	lenV, err := storages.StoreVocab(options.VocabFileName, vocab, &options.Op)
	if err != nil {
		os.Stderr.WriteString("\nFailed to write vocab file '" + options.VocabFileName + "': " + err.Error())
		return
	}

	os.Stdout.WriteString(fmt.Sprint("Using vocabulary of size ", lenV, "\n"))
}
