// Current path stack
let currentPath = [];
let fileSystem = null;

// DOM Elements
const fileExplorer = document.getElementById('file-explorer');
const fileContent = document.getElementById('file-content');
const fileTitle = document.getElementById('file-title');
const breadcrumb = document.getElementById('breadcrumb');
const refreshBtn = document.getElementById('refresh-btn');
const copyBtn = document.getElementById('copy-btn');
const downloadBtn = document.getElementById('download-btn');
const fileActions = document.getElementById('file-actions');

// File type icons
const FILE_ICONS = {
    'directory': 'folder',
    'file': 'file-earmark',
    'json': 'filetype-json',
    'js': 'filetype-js',
    'css': 'filetype-css',
    'html': 'filetype-html',
    'md': 'file-earmark-text',
    'py': 'filetype-py',
    'jpg': 'file-image',
    'jpeg': 'file-image',
    'png': 'file-image',
    'gif': 'file-image',
    'svg': 'file-image'
};

// Get file icon based on file extension
function getFileIcon(name, type) {
    if (type === 'directory') return 'folder';
    const ext = name.split('.').pop().toLowerCase();
    return FILE_ICONS[ext] || 'file-earmark';
}

// Format file size
function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// Render breadcrumb
function renderBreadcrumb() {
    breadcrumb.innerHTML = '';
    
    // Add root/home
    const homeItem = document.createElement('li');
    homeItem.className = 'breadcrumb-item' + (currentPath.length === 0 ? ' active' : '');
    const homeLink = document.createElement('a');
    homeLink.href = '#';
    homeLink.textContent = 'Home';
    homeLink.addEventListener('click', (e) => {
        e.preventDefault();
        navigateTo([]);
    });
    homeItem.appendChild(homeLink);
    breadcrumb.appendChild(homeItem);
    
    // Add path segments
    let currentPathStr = '';
    currentPath.forEach((segment, index) => {
        currentPathStr += (currentPathStr ? '/' : '') + segment;
        const isLast = index === currentPath.length - 1;
        
        const item = document.createElement('li');
        item.className = 'breadcrumb-item' + (isLast ? ' active' : '');
        
        if (!isLast) {
            const link = document.createElement('a');
            link.href = '#';
            link.textContent = segment;
            const pathToNavigate = currentPath.slice(0, index + 1);
            link.addEventListener('click', (e) => {
                e.preventDefault();
                navigateTo(pathToNavigate);
            });
            item.appendChild(link);
        } else {
            item.textContent = segment;
        }
        
        breadcrumb.appendChild(item);
    });
}

// Find item by path
function findItemByPath(pathArray, root = fileSystem) {
    let current = root;
    
    for (const segment of pathArray) {
        if (!current.children) return null;
        const found = current.children.find(item => item.name === segment);
        if (!found) return null;
        current = found;
    }
    
    return current;
}

// Render file explorer
function renderFileExplorer() {
    fileExplorer.innerHTML = '';
    
    // Get current directory
    const currentDir = currentPath.length === 0 ? fileSystem : findItemByPath(currentPath);
    if (!currentDir || !currentDir.children) return;
    
    // Sort: directories first, then files, both alphabetically
    const sortedItems = [...currentDir.children].sort((a, b) => {
        if (a.type === b.type) {
            return a.name.localeCompare(b.name);
        }
        return a.type === 'directory' ? -1 : 1;
    });
    
    // Render each item
    sortedItems.forEach(item => {
        const itemElement = document.createElement('div');
        itemElement.className = `file-item ${item.type === 'directory' ? 'dir-item' : ''}`;
        
        const icon = document.createElement('i');
        icon.className = `bi bi-${getFileIcon(item.name, item.type)}`;
        
        const name = document.createElement('span');
        name.textContent = item.name;
        
        const size = document.createElement('span');
        size.className = 'file-size';
        size.textContent = item.type === 'file' ? formatFileSize(item.size || 0) : '';
        
        itemElement.appendChild(icon);
        itemElement.appendChild(name);
        itemElement.appendChild(size);
        
        itemElement.addEventListener('click', (e) => {
            e.stopPropagation();
            if (item.type === 'directory') {
                navigateTo([...currentPath, item.name]);
            } else {
                showFileContent(item);
            }
        });
        
        fileExplorer.appendChild(itemElement);
    });
    
    // Show empty state if no items
    if (sortedItems.length === 0) {
        const emptyState = document.createElement('div');
        emptyState.className = 'empty-state';
        emptyState.innerHTML = `
            <i class="bi bi-folder-x"></i>
            <p class="mb-0">This directory is empty</p>
        `;
        fileExplorer.appendChild(emptyState);
    }
}

// Show file content
function showFileContent(file) {
    fileTitle.textContent = file.name;
    fileActions.style.display = 'block';
    
    // Set up copy button
    copyBtn.onclick = () => {
        navigator.clipboard.writeText(file.content || '');
        const originalText = copyBtn.innerHTML;
        copyBtn.innerHTML = '<i class="bi bi-check"></i>';
        setTimeout(() => {
            copyBtn.innerHTML = originalText;
        }, 2000);
    };
    
    // Set up download button
    downloadBtn.onclick = () => {
        const blob = new Blob([file.content || ''], { type: 'text/plain' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = file.name;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
    };
    
    // Determine content type and render accordingly
    let content = '';
    const ext = file.name.split('.').pop().toLowerCase();
    
    if (['jpg', 'jpeg', 'png', 'gif', 'svg'].includes(ext)) {
        // For images, create an image element
        content = `<div class="text-center"><img src="${file.path}" class="img-fluid" alt="${file.name}" onerror="this.parentElement.innerHTML='<i class=\'bi bi-file-earmark-image\'></i> <span>Cannot display image</span>';"></div>`;
    } else if (file.content) {
        // For text files, use highlight.js
        let codeClass = '';
        let codeContent = file.content;
        
        switch(ext) {
            case 'js':
                codeClass = 'javascript';
                break;
            case 'json':
                codeClass = 'json';
                try { codeContent = JSON.stringify(JSON.parse(file.content), null, 2); } catch(e) {}
                break;
            case 'html':
                codeClass = 'html';
                break;
            case 'css':
                codeClass = 'css';
                break;
            case 'md':
                // Render markdown
                content = `<div class="markdown-body">${window.marked.parse(file.content)}</div>`;
                break;
            default:
                codeClass = '';
        }
        
        if (content === '') {
            content = `<pre><code class="${codeClass}">${escapeHtml(codeContent)}</code></pre>`;
        }
    } else {
        content = '<div class="empty-state"><i class="bi bi-file-earmark"></i><p class="mb-0">No content to display</p></div>';
    }
    
    fileContent.innerHTML = content;
    
    // Apply syntax highlighting
    if (file.content && !['md', 'jpg', 'jpeg', 'png', 'gif', 'svg'].includes(ext)) {
        document.querySelectorAll('pre code').forEach((block) => {
            hljs.highlightElement(block);
        });
    }
    
    // Update active state
    document.querySelectorAll('.file-item').forEach(el => el.classList.remove('active'));
    const activeItem = Array.from(document.querySelectorAll('.file-item')).find(el => 
        el.textContent.trim().startsWith(file.name)
    );
    if (activeItem) activeItem.classList.add('active');
}

// Helper to escape HTML
function escapeHtml(unsafe) {
    return unsafe
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#039;");
}

// Navigate to path
function navigateTo(pathArray) {
    currentPath = pathArray;
    renderBreadcrumb();
    renderFileExplorer();
    
    // Hide file content when navigating directories
    fileContent.innerHTML = '<div class="empty-state"><i class="bi bi-file-earmark-text"></i><p class="mb-0">Select a file to view its contents</p></div>';
    fileTitle.textContent = 'Select a file';
    fileActions.style.display = 'none';
    
    // Update URL hash
    const hash = currentPath.length > 0 ? '#' + currentPath.join('/') : '';
    window.history.pushState(null, '', window.location.pathname + hash);
}

// Handle back/forward navigation
window.addEventListener('popstate', () => {
    const hash = window.location.hash.slice(1);
    currentPath = hash ? hash.split('/') : [];
    renderBreadcrumb();
    renderFileExplorer();
});

// Initialize
async function init() {
    try {
        // Load file system from inline script or fallback to fetch
        if (typeof FILE_SYSTEM !== 'undefined') {
            fileSystem = FILE_SYSTEM;
        } else {
            const response = await fetch('data/file-system.json');
            if (!response.ok) throw new Error('Failed to load file system');
            fileSystem = await response.json();
        }
        
        // Handle initial navigation based on URL hash
        const hash = window.location.hash.slice(1);
        if (hash) {
            currentPath = hash.split('/').filter(Boolean);
            // Check if the path exists
            const item = findItemByPath(currentPath);
            if (!item) {
                currentPath = [];
                window.location.hash = '';
            } else if (item.type === 'file') {
                const dirPath = currentPath.slice(0, -1);
                currentPath = dirPath;
                renderBreadcrumb();
                renderFileExplorer();
                showFileContent(item);
                return;
            }
        }
        
        renderBreadcrumb();
        renderFileExplorer();
        
        // Set up refresh button
        refreshBtn.addEventListener('click', () => {
            window.location.reload();
        });
        
    } catch (error) {
        console.error('Error initializing file explorer:', error);
        fileExplorer.innerHTML = `
            <div class="alert alert-danger m-3">
                <i class="bi bi-exclamation-triangle-fill"></i>
                Failed to load file system: ${error.message}
            </div>
        `;
    }
}

// Start the application
document.addEventListener('DOMContentLoaded', init);
