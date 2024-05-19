# Concurrency-pattern Pipeline for Secure Archiving

## Table of Contents
1. [Introduction](#introduction)
2. [Project Overview](#project-overview)
3. [Architecture](#architecture)
4. [Concurrency-pattern Pipeline](#concurrency-pattern-pipeline)
    - [Reading Files](#reading-files)
    - [Encrypting Content](#encrypting-content)
    - [Multiplexing](#multiplexing)
    - [Archiving](#archiving)
5. [Sequence Diagram](#sequence-diagram)
6. [Installation](#installation)
7. [Usage](#usage)
8. [Conclusion](#conclusion)
9. [References](#references)

## Introduction
This project demonstrates the use of a concurrency-pattern pipeline for securely archiving files. The pipeline improves performance and scalability by concurrently processing files through multiple stages: reading, encrypting, multiplexing, and archiving.

## Project Overview
The goal of this project is to showcase the benefits and implementation details of using a concurrency-pattern pipeline for secure file archiving. The process involves reading files from a source directory, encrypting them using AES encryption, multiplexing the encrypted data, and finally archiving the result.

## Architecture
The architecture of this project follows a pipeline pattern where each stage of the process is handled concurrently. This design leverages Go's concurrency features to enhance performance and scalability.

## Concurrency-pattern Pipeline

### Reading Files
The first stage in the pipeline is reading files from the source directory. This is done by the `readFile` function, which reads the contents of the directory and sends file data through a channel.

### Encrypting Content
The next stage involves encrypting the content of the files. This is achieved by three concurrent encryptor functions (`encryptContent`), each operating on the same input channel (`chanFile`). These functions encrypt the file data using AES encryption and send the encrypted data through separate channels.

### Multiplexing
The multiplexing stage combines the output from the three encryption channels into a single channel. This is done by the `multiPlexerEncrypt` function, which merges the encrypted data streams for further processing.

### Archiving
The final stage of the pipeline is archiving the encrypted data. The `archive` function takes the multiplexed encrypted data and creates a ZIP archive, storing the result in the output directory.

## Sequence Diagram
![Concurrency Pattern](https://firebasestorage.googleapis.com/v0/b/personal-website-1d263.appspot.com/o/project-pict%2Fconcurrent.png?alt=media&token=f802c366-678f-456b-a6de-d2fa5f03868a)

## Installation
To install and run this project, follow these steps:

1. Clone the repository:
    ```sh
    git clone https://github.com/vizucode/go-concurrency-pipeline
    cd go-concurrency-pipeline
    cd concurrent-pattern-archive
    ```

2. Build the project:
    ```sh
    go build -o secure-archive
    ```

3. Run the executable:
    ```sh
    ./secure-archive
    ```

## Usage
To use this project, ensure you have a directory of files you want to archive securely. Modify the source directory path in the `main` function and run the program. The encrypted and archived files will be saved in the output directory.

## Conclusion
The concurrency-pattern pipeline significantly enhances the performance and scalability of secure file archiving. By leveraging concurrent processing, the pipeline can handle large datasets efficiently, making it suitable for high-performance applications.

## References
- [Concurrency in Go](https://go.dev/blog/pipelines)

---

*Feel free to reach out if you have any questions or need further assistance!*

Thank You!
