PWD=$(shell pwd)

init:
	brew install pre-commit
	pre-commit install