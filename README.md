## EasyStorage

An program allows you to easily self-host a storage server that displays all the files with allowed extensions, allowing users to download the files. It provides a simple and secure way to share files with others without relying on third-party storage servers.

### How It Works

1. **File Storage**: The program serves as a web server that hosts files in a designated directory (`files_download`). Users can upload files to this directory to make them accessible for download.

2. **Allowed Extensions**: The program reads a JSON file (`allowed_extensions.json`) to determine which file extensions are allowed for download. Users can configure this file to specify the allowed extensions.

3. **Download Tracking**: The program tracks the download count for each file category (based on extension) and records the counts in a JSON file (`download_counts.json`). This feature allows you to monitor the popularity of different file types.

### Why I Made this Program?

- **Security**: By hosting your own storage server, you have full control over the security and privacy of your files. You don't need to rely on third-party servers, reducing the risk of data breaches or unauthorized access.

- **Customization**: You can customize the allowed file extensions and configure the server according to your specific needs. This flexibility allows you to tailor the server to meet your requirements.

- **Ease of Use**: The program is designed to be easy to set up and use. With just a few configuration steps, you can have your own storage server up and running, ready to share files with your friends!

### Getting Started

1. **Installation**: Clone the repository and install any dependencies required by the program.

2. **Configuration**: Customize the allowed file extensions by editing the `allowed_extensions.json` file. You can specify which file types are allowed for download.

3. **Run the Server**: Start the server by running the `main.go` file. The server will start listening on port 8080 by default.

4. **Upload Files**: Place your files in the `files_download` directory. These files will be accessible for download through the web interface.

5. **Access the Web Interface**: Open a web browser and navigate to yourwebsite.com/list to access the web interface. You can now view and download the files hosted on your storage server.

6. **Share Your Server**: Share the URL of your storage server with others to allow them to access and download the files. You can send the URL to your friends or colleagues to share files securely.

### Contributing

Contributions to this project are welcome! If you have any ideas for improvements or new features, feel free to submit a pull request or open an issue on GitHub.

### License

This project is licensed under the [MIT License](LICENSE). Feel free to use, modify, and distribute the code for your own purposes.
