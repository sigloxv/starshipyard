## Usage
A very simple usage example:

```Go
t := trie.New()
// Add can take in meta information which can be stored with the key.
// i.e. you could store any information you would like to associate with
// this particular key.
t.Add("foobar", 1)

node, ok := t.Find("foobar")
meta := node.Meta()
// use meta with meta.(type)
t.Remove("foobar")

t.PrefixSearch("foo")
t.HasKeysWithPrefix("foo")

t.FuzzySearch("fb")
```
