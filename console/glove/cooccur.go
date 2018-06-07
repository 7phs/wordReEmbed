package main

import (
	"fmt"
	"os"

	"github.com/7phs/wordReEmbed/console/options"
	"github.com/7phs/wordReEmbed/console/storages"
	"github.com/7phs/wordReEmbed/lib/model/glove"
	storages2 "github.com/7phs/wordReEmbed/lib/storages"
)

func showOptions(op *options.GloveOptions) {
	maxProduct, overflowLimit := options.Op.GetMemoryLimit()

	os.Stdout.WriteString(fmt.Sprintln("window size: ", options.Op.WindowSize))
	os.Stdout.WriteString(fmt.Sprintln("max product: ", maxProduct))
	os.Stdout.WriteString(fmt.Sprintln("overflow length: ", overflowLimit))
}

func main() {
	os.Stdout.WriteString("BUILDING VOCABULARY\n")

	showOptions(&options.Op)

	os.Stdout.WriteString("Reading vocab from file \"" + options.VocabFileName + "\"...")
	vocabCount, err := storages.ReadVocab(options.VocabFileName)
	if err != nil {
		os.Stderr.WriteString("\nFailed to read vocab file '" + options.VocabFileName + "': " + err.Error())
		return
	}

	os.Stdout.WriteString(fmt.Sprintln("loaded ", vocabCount.Len(), " words."))

	corpus, err := storages.OpenTextReader(options.CorpusFileName, ' ')
	if err != nil {
		os.Stderr.WriteString("Error open file: '" + options.CorpusFileName + "': " + err.Error() + "\n")
		return
	}
	defer corpus.Close()

	//bigramStorage := lib.NewLevelDB(options.CooccurenceFileName)
	//_, err = bigramStorage.Init()
	//if err!=nil {
	//	os.Stderr.WriteString("Error open db file: '" + options.CooccurenceFileName + "': " + err.Error() + "\n")
	//	return
	//}
	os.Stdout.WriteString("Building lookup table...")
	bigramStorage := storages2.NewFileBigramTable(options.CooccurenceFileName, vocabCount.Len(), &options.Op)
	//os.Stdout.WriteString(fmt.Sprintln("table contains ", bigramStorage.LookupTableSize(), " elements."))

	var (
		cooccur    = glove.NewCooccur(vocabCount, bigramStorage, &options.Op)
		totalCount = 0
	)

	err = nil

	os.Stdout.WriteString(fmt.Sprint("Processed ", totalCount, " tokens"))

	corpus.ForEach(func(word []byte) bool {
		if len(word) == 0 {
			return true
		}

		cooccur.Push(string(word))

		totalCount++
		if (totalCount % 100000) == 0 {
			os.Stdout.WriteString(fmt.Sprint("\033[11G ", totalCount, " tokens"))
		}

		return true
	})

	os.Stdout.WriteString(fmt.Sprint("\033[0GProcessed ", totalCount, " tokens\n"))

	bigramStorage.Flush()
}
