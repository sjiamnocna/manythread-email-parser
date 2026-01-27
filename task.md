# Programming Task: Multithreaded Email Analyzer

Write a multithreaded program that takes a folder containing email files as input and produces a CSV file with statistics as output. Each line of the CSV file should contain:

- The **name of the email file**
- The **number of lines** in the text part of the email
- The **number of bytes** in the text part of the email

## Requirements

- The **input folder path** will be provided as the **first command-line argument**.
- The **output CSV file name** will be provided as the **second command-line argument**.
- The folder of email files must be processed using **multiple threads**.
- For the purpose of this task, assume each email contains **only a plain text part and an HTML part**.
- For testing, you may use the emails included in the attached archive: `test_email.zip`.
- The use of **Golang standard libraries** is preferred.
