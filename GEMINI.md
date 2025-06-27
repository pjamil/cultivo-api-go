# Gemini Project Overview: cultivo-api-go

This document provides a brief overview of the `cultivo-api-go` project and key commands for development and testing, intended for use by the Gemini AI assistant.

## Project Description

`cultivo-api-go` is a Go-based API for managing cultivation-related data. Based on the project structure, it appears to handle entities such as plants (`planta`), environments (`ambiente`), genetics (`genetica`), and users (`usuario`). The API follows a standard layered architecture with controllers, services, repositories, and domain models.

## Key Commands

The following commands are available in the `Makefile`:

- **Build the application:**

  ```bash
  make build
  ```

  This command compiles the project and creates a binary at `bin/cultivo-api`.

- **Run tests:**

  ```bash
  make test
  ```

  This command executes all the tests in the project.

- **Build Docker image:**

  ```bash
  make docker-build
  ```

  This command builds a Docker image for the application with the tag `cultivo-api`.
