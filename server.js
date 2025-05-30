const express = require('express');
const fs = require('fs');
const path = require('path');
const cors = require('cors');
const morgan = require('morgan');

const app = express();
const PORT = 3000;

// Middleware
app.use(cors());
app.use(morgan('dev'));
app.use(express.static('.'));

// Helper function to check if path is within project directory
function isSafePath(basePath, targetPath) {
    const resolvedPath = path.resolve(basePath);
    const resolvedTarget = path.resolve(targetPath);
    return resolvedTarget.startsWith(resolvedPath + path.sep) || resolvedTarget === resolvedPath;
}

// Get list of files and directories
app.get('/api/project-files', (req, res) => {
    try {
        const basePath = __dirname;
        const relativePath = req.query.path || '';
        const fullPath = path.join(basePath, relativePath);

        // Security check
        if (!isSafePath(basePath, fullPath)) {
            return res.status(403).json({ error: 'Access denied' });
        }

        // Check if path exists and is a directory
        if (!fs.existsSync(fullPath) || !fs.statSync(fullPath).isDirectory()) {
            return res.status(404).json({ error: 'Directory not found' });
        }

        const items = fs.readdirSync(fullPath, { withFileTypes: true });
        const result = {
            path: relativePath,
            dirs: [],
            files: []
        };

        items.forEach(item => {
            const itemPath = path.join(relativePath, item.name);
            const fullItemPath = path.join(fullPath, item.name);
            
            // Skip hidden files and node_modules
            if (item.name.startsWith('.') || item.name === 'node_modules') {
                return;
            }

            if (item.isDirectory()) {
                result.dirs.push({
                    name: item.name,
                    path: itemPath
                });
            } else {
                result.files.push({
                    name: item.name,
                    path: itemPath,
                    size: fs.statSync(fullItemPath).size
                });
            }
        });

        // Sort directories and files
        result.dirs.sort((a, b) => a.name.localeCompare(b.name));
        result.files.sort((a, b) => a.name.localeCompare(b.name));

        res.json(result);
    } catch (error) {
        console.error('Error reading directory:', error);
        res.status(500).json({ error: 'Failed to read directory' });
    }
});

// Get file content
app.get('/api/file-content', (req, res) => {
    try {
        const basePath = __dirname;
        const filePath = req.query.path;
        
        if (!filePath) {
            return res.status(400).json({ error: 'File path is required' });
        }

        const fullPath = path.join(basePath, filePath);

        // Security check
        if (!isSafePath(basePath, fullPath)) {
            return res.status(403).json({ error: 'Access denied' });
        }

        // Check if file exists and is a file
        if (!fs.existsSync(fullPath) || !fs.statSync(fullPath).isFile()) {
            return res.status(404).json({ error: 'File not found' });
        }

        // Check file size (limit to 1MB)
        const stats = fs.statSync(fullPath);
        if (stats.size > 1024 * 1024) {
            return res.status(413).json({ error: 'File is too large to display' });
        }

        const content = fs.readFileSync(fullPath, 'utf8');
        res.json({
            path: filePath,
            content: content,
            size: stats.size,
            lastModified: stats.mtime
        });
    } catch (error) {
        console.error('Error reading file:', error);
        res.status(500).json({ error: 'Failed to read file' });
    }
});

// Serve index.html for all other routes
app.get('*', (req, res) => {
    res.sendFile(path.join(__dirname, 'index.html'));
});

// Start server
app.listen(PORT, () => {
    console.log(`Server is running on http://localhost:${PORT}`);
    console.log('Project explorer is available at http://localhost:3000');
});

// Handle uncaught exceptions
process.on('uncaughtException', (err) => {
    console.error('Uncaught Exception:', err);
    process.exit(1);
});
