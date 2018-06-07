# Word embedding models reimplementation in Go/Golang

Re-implementation of well-known word embedding algorithms to clearly understand how they training.
One of the goals is to separate the functionality of storing the model from learning.

Plan:
- [ ] [Glove](https://github.com/stanfordnlp/GloVe) (WIP)
- [ ] [FastText](https://github.com/facebookresearch/fastText/)
- [ ] [Word2Vec](https://code.google.com/archive/p/word2vec/)

Production ready Go library - [word-embedding](https://github.com/ynqa/word-embedding)

# Glove

Build console tools:
```
go build -o ./bin/glove/vocab_count ./console/glove/vocab_count.go
go build -o ./bin/glove/cooccur ./console/glove/cooccur.go  
```

Using a text corpus of the library: [text8](http://mattmahoney.net/dc/text8.zip).
