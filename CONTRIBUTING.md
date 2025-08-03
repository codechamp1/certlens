# Contributing to certlens

### First off, thank you for considering contributing to `certlens`. It's people like you that make it such a great tool.


### Following these guidelines helps to communicate that you respect the time of the developers managing and developing this open source project. In return, they should reciprocate that respect in addressing your issue, assessing changes, and helping you finalize your pull requests.


### What We're Looking For:

We welcome any contributions, whether it's fixing bugs, improving performance, or enhancing the UI. If you encounter a bug, feel free to submit a fix, but **please open an issue first** so we can track and discuss it. For new features or larger changes, it's important to **start by opening an issue** to propose and discuss the idea. This helps ensure we're aligned before development begins and prevents duplicated or misaligned work.

### What We're Not Looking For:
Please don’t use the issue tracker for general Kubernetes support or Go questions.
Don’t submit large PRs without prior discussion; we prefer iterative, focused improvements. All such issues will be closed without comment.


## Commit Message Guidelines

To keep the commit history clean and useful, please ensure your commit messages are:

- Clear and specific
- Focused on a single change or purpose
- Written in the imperative mood (e.g., `fix: crash when parsing invalid cert`)

We **require** all commits that are merged into the `main` branch to follow the **[Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/)** format. This is important because our release pipeline uses these messages to automatically generate changelogs and version bumps.

## Issue Labels

We use [Standard Issue Labels](https://github.com/wagenet/StandardIssueLabels) to keep our issue tracker consistent and understandable.

These labels help classify and prioritize issues, making collaboration smoother for everyone involved.

If you're opening a new issue, please familiarize yourself with the label definitions. You can apply labels if you have permission, or a maintainer will do so after your issue is reviewed.


## Pull Request Requirements

Before submitting a pull request, please ensure that you have completed **all items** in the PR checklist. This helps us maintain high-quality contributions and speeds up the review process.

## Continuous Integration Requirements

Every pull request must successfully pass all GitHub Actions workflows before it can be merged.

To ensure quality control, the workflows require manual approval by a maintainer before running. Please wait for approval and verify that all checks complete successfully prior to requesting a review.
