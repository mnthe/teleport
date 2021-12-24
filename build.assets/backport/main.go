/*
Copyright 2021 Gravitational, Inc.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/gravitational/trace"
)

func main() {
	input, err := parseInput()
	if err != nil {
		log.Fatal(err)
	}
	for _, baseBranch := range input.backportBranches {
		err := backport(baseBranch, input)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Backporting complete.")
}

type Input struct {
	// backportBranches is a list of branches to backport to.
	backportBranches []string

	// from is the name of the branch to pick the commits from.
	fromBranch string

	// mergeBaseCommit is the merge base commit.
	mergeBaseCommit string

	// headCommit is the HEAD of the from branch.
	headCommit string

	// startingBranch is the current branch.
	startingBranch string
}

func parseInput() (Input, error) {
	var to string
	var from string
	flag.StringVar(&to, "to", "", "List of comma-separated branch names to backport to.\n Ex: branch/v6,branch/v7\n")
	flag.StringVar(&from, "from", "", "Branch with changes to backport.")
	flag.Parse()

	if to == "" {
		return Input{}, trace.BadParameter("must supply branches to backport to.")
	}
	if from == "" {
		return Input{}, trace.BadParameter("much supply branch with changes to backport.")
	}
	// Parse branches to backport to.
	backportBranches := parseBranches(to)

	// To cherry pick all commits from a branch, the merge-base and
	// HEAD of the branch commits are needed.
	mbCommit, err := getMergeBaseCommit(from)
	if err != nil {
		return Input{}, trace.Wrap(err)
	}

	head, err := getHeadFromBranch(from)
	if err != nil {
		return Input{}, trace.Wrap(err)
	}

	currentBranchName, err := getCurrentBranch()
	if err != nil {
		return Input{}, trace.Wrap(err)
	}

	return Input{
		backportBranches: backportBranches,
		fromBranch:       from,
		mergeBaseCommit:  mbCommit,
		headCommit:       head,
		startingBranch:   currentBranchName,
	}, nil
}

func backport(baseBranch string, input Input) error {
	newBranchName, err := createBranch(input.fromBranch, baseBranch)
	if err != nil {
		return trace.BadParameter("failed to create a new branch from %s: %v", baseBranch, err)
	}
	fmt.Printf("New branch %s created.\n", newBranchName)
	// Checkout the new branch. This will fail if there are any uncommitted changes.
	// The working tree MUST be clean.
	err = checkout(newBranchName)
	if err != nil {
		err := cleanUp(newBranchName)
		if err != nil {
			fmt.Printf("*** Failed to clean up branch. please manually delete %s ***\n", newBranchName)
		}
		fmt.Println("*** Ensure your working tree is clean. ***")
		return trace.Wrap(err)
	}
	fmt.Printf("Cherry picking %s-%s \n", input.mergeBaseCommit, input.headCommit)
	err = cherryPick(input.mergeBaseCommit, input.headCommit)
	if err != nil {
		if err := cleanUp(newBranchName); err != nil {
			log.Printf("*** Failed to clean up branch. please manually delete %s ***\n", newBranchName)
		}
		return trace.BadParameter("failed to cherry-pick %s-%s: %v", input.mergeBaseCommit, input.headCommit, err)
	}
	fmt.Printf("Cherry picked %s-%s to branch %s based off of branch %s\n\n",
		input.mergeBaseCommit, input.headCommit, newBranchName, baseBranch)

	// Push new branch to remote.
	err = push(newBranchName)
	if err != nil {
		trace.BadParameter("failed to push branch %s: %v", newBranchName, err)
	}
	fmt.Println("Changes pushed successfully.")
	err = createPullRequest(baseBranch, newBranchName)
	if err != nil {
		return trace.BadParameter("failed to create a pull request for %s: %v.\n Open up a pull request on github.com.", newBranchName, err)
	}
	err = checkout(input.startingBranch)
	if err != nil {
		return trace.BadParameter("failed to checkout branch %s: %v", input.startingBranch, err)
	}
	return nil
}

// cleanUp force checks out the master branch and
// deletes the given branch.
func cleanUp(branchToDelete string) error {
	_, _, err := run("git", "checkout", "master", "--force")
	if err != nil {
		return trace.Wrap(err)
	}
	return deleteBranch(branchToDelete)
}

// deleteBranch deletes the specified branch name.
func deleteBranch(branchName string) error {
	_, _, err := run("git", "branch", "-D", branchName)
	if err != nil {
		return trace.Wrap(err)
	}
	return nil
}

// getCurrentBranch gets the current working branch. 
func getCurrentBranch() (string, error) {
	currentBranchName, _, err := run("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", trace.Wrap(err)
	}
	currentBranchName = trim(currentBranchName)
	if currentBranchName == "" {
		return "", trace.Errorf("failed to get the current branch")
	}
	return currentBranchName, nil
}

// parseBranches parses the input branches to backport to.
func parseBranches(branchesInput string) []string {
	var backportBranches []string
	branches := strings.Split(branchesInput, ",")
	for _, branch := range branches {
		backportBranches = append(backportBranches, strings.TrimSpace(branch))
	}
	return backportBranches
}

// push pushes changes to the remote repository configured
// in `.git/config` located in the project root.
func push(backportBranchName string) error {
	_, _, err := run("git", "push", "--set-upstream", "origin", backportBranchName)
	if err != nil {
		return trace.Wrap(err)
	}
	return nil
}

// createPullRequest creates a pull request with the credentials stored
// in ~/.config/gh/hosts.yaml.
func createPullRequest(baseBranch, headBranch string) error {
	_, _, err := run("gh", "pr", "create", "--fill", "--label", "backport", "--base", baseBranch, "--head", headBranch)
	if err != nil {
		return trace.Wrap(err)
	}
	return nil
}

// checkout checks out the specified branch.
func checkout(branch string) error {
	_, _, err := run("git", "checkout", branch)
	if err != nil {
		return trace.Wrap(err)
	}
	return nil
}

func getHeadFromBranch(branch string) (string, error) {
	latestCommit, _, err := run("git", "log", "-n", "1", "--pretty=format:\"%H\"", branch)
	if err != nil {
		return "", trace.Wrap(err)
	}
	latestCommit = trim(latestCommit)
	if latestCommit == "" {
		return "", trace.Errorf("failed to get the HEAD for %s: %v ", branch, err)
	}
	return latestCommit, nil
}

// trim trims the input string of new lines and quotes.
func trim(input string) string {
	input = strings.Trim(input, "\"")
	input = strings.Trim(input, "\n")
	return input
}

func getMergeBaseCommit(branchToCherryPickFrom string) (string, error) {
	mergeBaseCommit, _, err := run("git", "merge-base", "master", branchToCherryPickFrom)
	if err != nil {
		return "", trace.Wrap(err)
	}
	mergeBaseCommit = trim(mergeBaseCommit)
	if mergeBaseCommit == "" {
		return "", trace.Errorf("failed to get the merge base commit of %s", branchToCherryPickFrom)
	}
	return mergeBaseCommit, nil
}

// createBranch creates a new branch based off of the base branch.
func createBranch(fromBranchName, baseBranchName string) (string, error) {
	newBranchName := fmt.Sprintf("auto-backport/%s/%s", baseBranchName, fromBranchName)
	_, _, err := run("git", "branch", newBranchName, baseBranchName)
	if err != nil {
		return "", trace.Wrap(err)
	}
	return newBranchName, nil
}

// cherryPick cherry picks a range of commits. The first
// commit is not included.
func cherryPick(mergeBaseCommit, headCommitOfBranchWithChanges string) error {
	if mergeBaseCommit == headCommitOfBranchWithChanges {
		return trace.BadParameter("there are no changes to backport.")
	}
	commitRange := fmt.Sprintf("%s..%s", mergeBaseCommit, headCommitOfBranchWithChanges)
	_, _, err := run("git", "cherry-pick", commitRange)
	if err != nil {
		return trace.Wrap(err)
	}
	return nil
}

// run executes command on disk.
func run(ex string, command ...string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	path, err := exec.LookPath(ex)
	if err != nil {
		return "", "", err
	}
	cmd := exec.Command(path, command...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	return stdout.String(), stderr.String(), err
}
