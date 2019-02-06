# Boggle-Challenge
Boggle Challenge - given a set of boggle dice, and a time limit, come up with the best possible boggle board and its corresponding score.

It starts out with a random configuration of the boggle board and for every iteration, checks to see if the new board is better than the previous one. If it is, then the new board is accepted. If not, it is accepted with a probability that is high initially and decreases later on. This is based on the principle of Simulated Annealing https://en.wikipedia.org/wiki/Simulated_annealing, which in turn is based on the Metropolis-Hastings algorithm https://en.wikipedia.org/wiki/Metropolis–Hastings_algorithm.

Build it using: go build Boggle-Challenge

Run it using: go run Boggle-Challenge

It can be also be run using these flags

-a=false for not accepting Boggle boards that have a worst score than the current one

-l=true for logging the progress and printing the results at the end

-p=3 for perturbation count of dice. In this case, 3 dice will be perturbed to show another face and generate a different board.

-t=2000 for the initial temperature

-c=0.95 for the cooling rate

-m=1 for number of minutes to run the search for the best board

Here are the default values for the flags: -a=true -l=false -p=1 -t=1000 -c=0.99 -m=2

Here is an sample run command for runnng the board search for 5 minutes

go run Boggle-Challenge -m=5
