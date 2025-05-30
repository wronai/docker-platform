# Media Vault Project Explorer

A web-based file explorer for the Media Vault project that renders markdown files and displays project structure.

## Features

- Browse project files and directories
- View and render markdown files in real-time
- Syntax highlighting for code blocks
- Responsive design that works on desktop and mobile
- File size and modification date information

## Getting Started

1. Make sure you have Node.js installed (v14 or later)
2. Install the dependencies:

```bash
npm install
```

3. Start the development server:

```bash
npm start
```

4. Open your browser and navigate to:

```
http://localhost:3000
```

## Development

For development with auto-reload:

```bash
npm run dev
```

## Security Notes

- The server is configured to only serve files within the project directory
- File size is limited to 1MB for display
- Hidden files and node_modules are excluded from the file listing

## Dependencies

- Express.js - Web server framework
- Marked - Markdown parser
- Highlight.js - Syntax highlighting
- Bootstrap 5 - Styling

## License

This project is part of the Media Vault platform.
