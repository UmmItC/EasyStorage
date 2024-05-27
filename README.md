# EasyStorage

EasyStorage is a simple program that allows you to easily self-host a storage server for sharing files. It provides a convenient way to share files securely without relying on third-party storage servers.

## How It Works

1. **File Storage**: EasyStorage serves as a web server that hosts files in a designated directory (`files_download`). Users can upload files to this directory to make them accessible for download.

2. **Allowed Extensions**: EasyStorage reads a JSON file (`allowed_extensions.json`) to determine which file extensions are allowed for download. Users can customize this file to specify the allowed extensions.

3. **Download Tracking**: EasyStorage tracks the download count for each file category (based on extension) and records the counts in a JSON file (`download_counts.json`). This feature allows you to monitor the popularity of different file types.

## Why Use EasyStorage?

- **Security**: Hosting your own storage server gives you full control over the security and privacy of your files. You're not reliant on third-party servers, reducing the risk of data breaches or unauthorized access.

- **Customization**: EasyStorage allows you to customize the allowed file extensions and configure the server according to your specific needs. This flexibility enables you to tailor the server to meet your requirements.

- **Ease of Use**: EasyStorage is designed to be simple to set up and use. With just a few configuration steps, you can have your own storage server up and running, ready to share files with others!

## Getting Started

### Installation

Make sure you have Go and Nginx installed on your system:

```bash
sudo apt install golang nginx
```

Clone the EasyStorage repository:

```bash
git clone https://codeberg.org/UmmIt/EasyStorage.git
```

### Configuration

Customize the allowed file extensions by editing the `allowed_extensions.json` file. Specify which file types are allowed for download.

### Running the Server

#### Manually

1. Install screen if you haven't already:

```bash
sudo apt install screen
```

2. Start a new screen session:

```bash
screen
```

3. Navigate to the EasyStorage directory and build the program:

```bash
cd EasyStorage
go build
```

4. Run EasyStorage:

```bash
./EasyStorage
```

Press `Ctrl+A` followed by `Ctrl+D` to detach from the screen session while leaving EasyStorage running.

#### Systemd Service

1. Create a systemd service file for EasyStorage:

```bash
sudo vim /lib/systemd/system/easy-storage.service
```

2. Copy and paste the following configuration into the file:

```ini
[Unit]
Description=EasyStorage

[Service]
Type=simple
Restart=always
RestartSec=5s
WorkingDirectory=/path/to/your/EasyStorage/directory
ExecStart=/path/to/EasyStorage/EasyStorage

[Install]
WantedBy=multi-user.target
```

3. Reload systemd to apply changes:

```bash
sudo systemctl daemon-reload
```

4. Start and enable the EasyStorage service:

```bash
sudo systemctl start easy-storage
sudo systemctl enable easy-storage
```

5. Check the status of the service:

```bash
sudo systemctl status easy-storage
```

### Nginx Configuration

Edit your Nginx configuration file (e.g., `/etc/nginx/sites-available/default` or a custom configuration file):

```nginx
server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;

    # SSL Certificate
    ssl_certificate /etc/ssl/cert.pem;
    ssl_certificate_key /etc/ssl/key.pem;

    # Serve Listing Page
    location / {
        try_files $uri $uri/ =404;
        proxy_pass http://localhost:8080/list;
    }
    
    # Download Files
    location /download {
        proxy_pass http://localhost:8080;
    }

    # Serve Files for Download
    location /files_download {
        root /var/www/storage;
    }

    # Domain Configuration
    server_name storage storage.domain.com;
}
```

Make sure to replace `/etc/ssl/cert.pem`, `/etc/ssl/key.pem`, and `/var/www/storage` with the paths relevant to your system.

## Usage

1. **Upload Files**: Place your files in the `files_download` directory. These files will be accessible for download through the web interface.

2. **Access the Web Interface**: Open a web browser and navigate to your website's URL. You can now view and download the files hosted on your storage server.

3. **Share Your Server**: Share the URL of your storage server with others to allow them to access and download the files securely.

## Contributing

Contributions to this project are welcome! If you have any ideas for improvements or new features, feel free to submit a pull request or open an issue on GitHub.

## License

This project is licensed under the [MIT License](./LICENSE.md). Feel free to use, modify, and distribute the code for your own purposes.
