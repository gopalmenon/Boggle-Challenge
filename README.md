# Boggle-Challenge
(as described [here](https://codegolf.stackexchange.com/questions/5654/best-scoring-boggle-board))

Given a set of boggle dice, and a time limit, come up with the best possible boggle board, its corresponding score and the words found.

It starts out with a random configuration of the boggle board and for every iteration, checks to see if the new board (generated by perturbing dice) is better than the previous one. If it is, then the new board is accepted. If not, it is accepted with a probability that is high initially and decreases later on. This is based on the principle of [Simulated Annealing](https://en.wikipedia.org/wiki/Simulated_annealing), which is a variation of the [Metropolis-Hastings algorithm](https://en.wikipedia.org/wiki/Metropolis–Hastings_algorithm).

This project uses the trie data structure developed by [Derek Parker](https://github.com/derekparker). Thanks Derek. The list of valid words is loaded into the trie from [here](http://coursera.cs.princeton.edu/algs4/testing/boggle/dictionary-yawl.txt). When starting a word search at every die, if the prefix does not exist in the trie, then the rest of the search is ended, thus making the search more efficient.

Get the prefix tree or trie data structure using: ```go get github.com/derekparker/trie```

Build it using: ```go build Boggle-Challenge```

Run it using: ```go run Boggle-Challenge```

It can be also be run using these flags

* ```-a=false``` for not accepting Boggle boards that have a worst score than the current one

* ```-l=true``` for logging the progress and printing the results at the end

* ```-p=3``` for perturbation count of dice. In this case, 3 dice will be perturbed to show another face and generate a different board.

* ```-t=2000``` for the initial temperature

* ```-c=0.95``` for the cooling rate

* ```-m=1``` for number of minutes to run the search for the best board

Here are the default values for the flags: ```-a=true -l=false -p=1 -t=1000 -c=0.99 -m=2```

Here is a sample run command for runnng the board search for 5 minutes

```go run Boggle-Challenge -m=5```

The current best score I get is 1346. There are some improvements for which I have logged issues and will attempt soon.
