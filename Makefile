push:
	git tag -f $(tag) `git rev-parse HEAD`
	git push --force origin refs/tags/$(tag):refs/tags/$(tag)
