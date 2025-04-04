# NLP Suite Runner

The NLP Suite Runner is the new, official tool to install and run the NLP Suite.
The NLP Suite Runner is a binary that automatically installs 3 main packages:

- NLP Suite UI: the user interface to interact with the NLP Suite Agent
- NLP Suite Agent: the core algorithms and tools provided by the NLP Suite
- Stanford CoreNLP: a dependency needed to run the NLP Suite

**The NLP Suite Runner requires to download and install** [Docker Desktop](https://www.docker.com/products/docker-desktop/).

## ⚙️ Requirements

To install the NLP Suite, Docker desktop requires at least 4GB of RAM.
The NLP Suite Runner strongly encourages that you have at least 10 GB of storage available.

## 📦 Installation

To install the NLP Suite Runner, first ensure that you have downloaded and installed [Docker Desktop](https://www.docker.com/products/docker-desktop/).
You can download the proper NLP Suite via Docker executable for your operating system in the latest release [here](https://github.com/NLP-Suite/nlp-suite-runner/releases/latest).
The executable will be downloaded and placed to a location of your choice. The executable (nlp-suite-runner.exec) will look like this: <img width="131" alt="Screenshot 2025-03-31 at 4 22 37 PM" src="https://github.com/user-attachments/assets/df863362-df60-4bcc-b266-e56932ff9717" />


To run GIS (Geographic Information System) mapping via **Google Earth Pro**, rather than the default Python Folium, you need to download and install Google Earth Pro; click [here](https://www.google.com/earth/download/gep/agree.html?hl=en-GB) to download Google Earth Pro.

To display network graphs via **Gephi**, you need to download and install Gephi; click [here](https://gephi.org/users/download/) to download Gephi. 

-# Mac

Mac and Linux systems can easily extract `.tar.gz` files without additional tools.

#### M-Chip Series

If you have a mac running the ARM M-chip series, ensure that you install the `darwin-arm64` version of the NLP Suite Runner.

#### Intel Chip Series

If you have the Intel chip, you will need to install the `darwin-amd64` version of the NLP Suite Runner.

### Windows

To extract the runner, you will need [7-zip](https://www.7-zip.org/), a file archiver tool for Windows.

You most likely will need to install the `windows-amd64` version of the NLP Suite Runner.
If your device has an ARM chip, ensure you have `windows-arm64` installed. To check which version you have, check [here](https://www.tenforums.com/tutorials/176966-how-check-if-processor-32-bit-64-bit-arm-windows-10-a.html).

## ⚡️ Running the NLP Suite Runner

To run the NLP Suite Runner, double click the binary/exe.
**You must first start Docker Desktop by clicking on the Docker app wherever you installed on your computer.
![image](https://github.com/user-attachments/assets/0b881596-c8b8-47e7-a63d-522475ffee4e)
It is required to run the NLP Suite Runner.**

When the runner starts, cpy the address on an open brower to visualize the NLP Suite 
![image](https://github.com/user-attachments/assets/bea8b5ab-10b6-472d-b039-fc882b877900)

![image](https://github.com/user-attachments/assets/fe12e55b-8af0-4b9b-8392-652eed9ad1ef)

Click on the NLP tool you wish to run, e.g., Sentiment analysis

![image](https://github.com/user-attachments/assets/a1d157e1-87bd-42f5-8cc5-98a107f0e508)

The input and output directories are disabled. You cannot change the values.

![image](https://github.com/user-attachments/assets/7c492712-d65c-4152-b00c-1d3931396a8a)

So... how can yyou change those directories? See, below the section **Input / Output Directories**


### Mac Systems

You may see an error indicating that the NLP Suite Runner cannot be opened because it is not coming from a trusted source.
This is expected to occur the first time you run the NLP Suite Runner. To start the runner, follow [these instructions to allow the NLP Suite Runner to open.](https://www.macworld.com/article/672947/how-to-open-a-mac-app-from-an-unidentified-developer.html#:~:text=Open%20System%20Settings.-,Go%20to%20Privacy%20%26%20Security.,Click%20the%20Open%20Anyway%20button.)

### Windows / Linux Systems

There are no known issues for starting and running the NLP Suite Runner at this time.

### Input / Output Directories
Currently, the input and output directories are hard coded. DO NOT TRY ALTER THE INPUT AND OUTPUT PATHS.

For **Mac users**, the folder is located in your home folder. The home folder may not directly appear in your Finder. To convinently add your home folder to Finder, use these [instructions](https://www.tomsguide.com/how-to/how-to-find-the-home-folder-on-mac-and-add-it-to-finder).

For **Windows users**, the two Input & Output folders are automatically installed for you in your Users: folder.
![image](https://github.com/user-attachments/assets/f0738bb2-5903-4aee-8e02-05a4a16004f6)
![image](https://github.com/user-attachments/assets/e55dc6f7-0929-4a7f-9f1d-99cde4bd0022)  
For instance rfranzo and within it nlp-suite  
![image](https://github.com/user-attachments/assets/f8b58041-1683-45a7-a073-6d2a3895abc7)  
The path would be **C:\Users\rfranzo\nlp-suite**

**Make sure that your corpus files are uploaded into the input directory (e.g., C:\Users\rfranzo\nlp-suite\input) and not inside a sub-folder of this directory.**

### Updating the NLP Suite Runner

While the NLP Suite Runner automatically updates the installed packages mentioned above, the NLP Suite Runner does not automatically update itself.
Periodically check this page to ensure you have the latest version installed.

## 📣 Reporting issues

If there are any issues you are encountering, header over to [Issues](https://github.com/NLP-Suite/nlp-suite-runner/issues) to open a ticket.

## 💻 Contributing

If you would like to contribute to the NLP Suite, please contact [Roberto Franzosi](https://sociology.emory.edu/people/bios/Franzosi-Roberto.html)
