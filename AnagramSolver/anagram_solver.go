package main

import (
  "bufio"
  "fmt"
  "os"
  "strings"
  "sort"
)

//function to sort the characters of the word to make a key
func sortWord(w string) string {
  // split the string into character arrays
  s := strings.Split(w, "")
  // Sort the chracter array
  sort.Strings(s)
  // Join the characters in the array
  return strings.Join(s, "")
}

func main() {
  //input file that contains the list of words seperated by carriage return
  fileHandler, err := os.Open("wordList.txt")
  if err != nil {
    panic(err)
  }
  defer fileHandler.Close()

  //anagramsMap is a map of string keys to array[<string>] values
  anagramsMap := make(map[string][]string)
  // scanner to scan the words one by one from the input file
  scanner := bufio.NewScanner(fileHandler)
  for scanner.Scan() {
    //this is the value
    word:=scanner.Text() 
    //this is the key
    sortedWord := sortWord(word) 
    // Append the words to the map
    anagramsMap[sortedWord] = append(anagramsMap[sortedWord], word)
  }
  if err := scanner.Err(); err != nil {
    fmt.Fprintln(os.Stderr, "Error in scanning the values from file:", err)
  }

  //testing the anagramsolver by passing the word
  fmt.Println(anagramsMap[sortWord("resets")])

}